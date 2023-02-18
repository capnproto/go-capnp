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
	"capnproto.org/go/capnp/v3/flowcontrol/bbr"
	"capnproto.org/go/capnp/v3/flowcontrol/tracing"
	"capnproto.org/go/capnp/v3/internal/syncutil"
	"capnproto.org/go/capnp/v3/rpc"
	"zenhack.net/go/util"
)

var (
	addr       = flag.String("addr", ":2323", "Address to listen on/connect to")
	dial       = flag.Bool("client", false, "Should we be the client?")
	limiter    = flag.String("limiter", "", "What limiter should we use?")
	packetsize = flag.Int("packetsize", 8192, "Size of individual packets")
	totaldata  = flag.Int("totaldata", 512*1024*1024, "How much data should we send?")
	bandwidth  = flag.Uint64("bandwidth", math.MaxInt, "Maximum bandwidth the server will permit (B/s)")
	noTrace    = flag.Bool("no-trace", false, "don't capture a trace")
)

type Report struct {
	StartTime time.Time
	EndTime   time.Time
	Duration  time.Duration
	Sent      int
	Bandwidth float64

	Records []tracing.TraceRecord

	Snapshots []any
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
	util.Chkfatal(err)
	rpcConn := rpc.NewConn(rpc.NewStreamTransport(netConn), nil)
	defer rpcConn.Close()
	w := Writer(rpcConn.Bootstrap(ctx))
	l := capnp.Client(w).GetFlowLimiter()
	var (
		snapshots  []bbr.Snapshot
		snapshotMu sync.Mutex
	)
	switch *limiter {
	case "":
	case "fixed":
		l = flowcontrol.NewFixedLimiter(1024 * 1024)
	case "bbr":
		l = bbr.NewLimiter(nil)
	case "bbr-snapshot":
		l = bbr.SnapshottingLimiter{
			Limiter: bbr.NewLimiter(nil),
			RecordSnapshot: func(s bbr.Snapshot) {
				syncutil.With(&snapshotMu, func() {
					snapshots = append(snapshots, s)
				})
			},
		}
	default:
		panic("Unknown limiter type: " + *limiter)
	}
	var (
		records   []tracing.TraceRecord
		recordsMu sync.Mutex
	)
	tl := tracing.New(l, func(record tracing.TraceRecord) {
		syncutil.With(&recordsMu, func() {
			records = append(records, record)
		})
	})
	if *noTrace {
		capnp.Client(w).SetFlowLimiter(l)
	} else {
		capnp.Client(w).SetFlowLimiter(tl)
	}
	startTime := time.Now()
	sent := 0

	wg := &sync.WaitGroup{}
	for sent < *totaldata && ctx.Err() == nil {
		fut, rel := w.Write(ctx, func(p Writer_write_Params) error {
			util.Chkfatal(p.SetData(make([]byte, *packetsize)))
			sz, _ := p.Message().TotalSize()
			sent += int(sz)
			return nil
		})
		waitAsync(wg, fut, rel)
	}

	wg.Wait()

	endTime := time.Now()
	util.Chkfatal(ctx.Err())

	duration := endTime.Sub(startTime)
	bandwidth := float64(sent) / (float64(duration) / float64(time.Second))

	report := Report{
		StartTime: startTime,
		EndTime:   endTime,
		Duration:  duration,
		Sent:      sent,
		Bandwidth: bandwidth,
		Records:   records,
	}
	for _, s := range snapshots {
		report.Snapshots = append(report.Snapshots, s.Json())
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	util.Chkfatal(enc.Encode(report))
}

func waitAsync(wg *sync.WaitGroup, fut Writer_write_Results_Future, rel capnp.ReleaseFunc) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer rel()
		_, err := fut.Struct()
		util.Chkfatal(err)
	}()
}

func doServer() {
	l, err := net.Listen("tcp", *addr)
	util.Chkfatal(err)
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
