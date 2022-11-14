// test-tool is a useful command line tool for testing flow limiters.
package main

import (
	context "context"
	"encoding/json"
	"flag"
	"math"
	"net"
	"os"
	"sync"
	"time"

	capnp "capnproto.org/go/capnp/v3"
	"capnproto.org/go/capnp/v3/flowcontrol"
	"capnproto.org/go/capnp/v3/flowcontrol/tracing"
	"capnproto.org/go/capnp/v3/rpc"
)

var (
	addr      = flag.String("addr", ":2323", "Address to listen on/connect to")
	dial      = flag.Bool("client", false, "Should we be the client?")
	limiter   = flag.String("limiter", "", "What limiter should we use?")
	totaldata = flag.Int("totaldata", 512*1024*1024, "How much data should we send?")
	bandwidth = flag.Uint64("bandwidth", math.MaxInt, "Maximum bandwidth the server will permit (B/s)")
)

type Report struct {
	StartTime time.Time
	EndTime   time.Time
	Duration  time.Duration
	Sent      int
	Bandwidth float64

	Records []tracing.TraceRecord
}

func main() {
	flag.Parse()
	if *dial {
		doClient(context.Background())
	} else {
		doServer()
	}
}

func doClient(ctx context.Context) {
	netConn, err := net.Dial("tcp", *addr)
	chkfatal(err)
	rpcConn := rpc.NewConn(rpc.NewStreamTransport(netConn), nil)
	defer rpcConn.Close()
	w := Writer(rpcConn.Bootstrap(ctx))
	l := capnp.Client(w).GetFlowLimiter()
	switch *limiter {
	case "":
	case "fixed":
		l = flowcontrol.NewFixedLimiter(1024 * 1024)
	case "bbr":
		l = flowcontrol.NewBBR()
	default:
		panic("Unknown limiter type: " + *limiter)
	}
	tl := &tracing.TraceLimiter{Underlying: l}
	capnp.Client(w).SetFlowLimiter(tl)
	startTime := time.Now()
	sent := 0

	wg := &sync.WaitGroup{}
	for sent < *totaldata && ctx.Err() == nil {
		fut, rel := w.Write(ctx, func(p Writer_write_Params) error {
			chkfatal(p.SetData(make([]byte, 8192)))
			sz, _ := p.Message().TotalSize()
			sent += int(sz)
			return nil
		})
		waitAsync(wg, fut, rel)
	}

	wg.Wait()

	endTime := time.Now()
	chkfatal(ctx.Err())

	duration := endTime.Sub(startTime)
	bandwidth := float64(sent) / (float64(duration) / float64(time.Second))

	report := Report{
		StartTime: startTime,
		EndTime:   endTime,
		Duration:  duration,
		Sent:      sent,
		Bandwidth: bandwidth,
		Records:   tl.Records(),
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	chkfatal(enc.Encode(report))
}

func waitAsync(wg *sync.WaitGroup, fut Writer_write_Results_Future, rel capnp.ReleaseFunc) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer rel()
		_, err := fut.Struct()
		chkfatal(err)
	}()
}

func doServer() {
	l, err := net.Listen("tcp", *addr)
	chkfatal(err)
	for {
		netConn, err := l.Accept()
		if err != nil {
			continue
		}
		go func() {
			rpcConn := rpc.NewConn(rpc.NewStreamTransport(netConn),
				&rpc.Options{
					BootstrapClient: capnp.Client(Writer_ServerToClient(writerImpl{})),
				})
			<-rpcConn.Done()
		}()
	}
}

func chkfatal(err error) {
	if err != nil {
		panic(err)
	}
}

type writerImpl struct {
}

func (writerImpl) Write(ctx context.Context, p Writer_write) error {
	data, err := p.Args().Data()
	if err != nil {
		return err
	}
	time.Sleep(time.Duration(float64(len(data)) / (float64(*bandwidth) * float64(time.Second))))
	return nil
}
