package rpc_test

import (
	"context"
	"net"
	"testing"

	"capnproto.org/go/capnp/v3"
	"capnproto.org/go/capnp/v3/flowcontrol"
	"capnproto.org/go/capnp/v3/rpc"
	testcp "capnproto.org/go/capnp/v3/rpc/internal/testcapnp"
	"capnproto.org/go/capnp/v3/std/capnp/stream"
)

func BenchmarkStreaming(b *testing.B) {
	ctx := context.Background()
	p1, p2 := net.Pipe()
	srv := testcp.StreamTest_ServerToClient(nullStream{})
	conn1 := rpc.NewConn(rpc.NewStreamTransport(p1), &rpc.Options{
		BootstrapClient: capnp.Client(srv),
	})
	defer conn1.Close()
	conn2 := rpc.NewConn(rpc.NewStreamTransport(p2), nil)
	defer conn2.Close()
	bootstrap := testcp.StreamTest(conn2.Bootstrap(ctx))
	defer bootstrap.Release()
	var (
		futures      []stream.StreamResult_Future
		releaseFuncs []capnp.ReleaseFunc
	)
	bootstrap.SetFlowLimiter(flowcontrol.NewFixedLimiter(1 << 9))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for j := 0; j < 1<<16; j++ {
			fut, rel := bootstrap.Push(ctx, nil)
			futures = append(futures, fut)
			releaseFuncs = append(releaseFuncs, rel)
		}
	}
	for i, fut := range futures {
		_, err := fut.Struct()
		if err != nil {
			b.Errorf("Error waiting on future #%v: %v", i, err)
		}
	}
	for _, rel := range releaseFuncs {
		rel()
	}
}

// nullStream implements testcp.StreamTest, ignoring the data it is sent.
type nullStream struct {
}

func (nullStream) Push(context.Context, testcp.StreamTest_push) error {
	return nil
}

func BenchmarkPingPong(b *testing.B) {
	p1, p2 := net.Pipe()
	srv := testcp.PingPong_ServerToClient(pingPongServer{})
	conn1 := rpc.NewConn(rpc.NewStreamTransport(p2), &rpc.Options{
		ErrorReporter:   testErrorReporter{tb: b},
		BootstrapClient: capnp.Client(srv),
	})
	defer func() {
		<-conn1.Done()
		if err := conn1.Close(); err != nil {
			b.Error("conn1.Close:", err)
		}
	}()
	conn2 := rpc.NewConn(rpc.NewStreamTransport(p1), &rpc.Options{
		ErrorReporter: testErrorReporter{tb: b},
	})
	defer func() {
		if err := conn2.Close(); err != nil {
			b.Error("conn2.Close:", err)
		}
	}()

	ctx := context.Background()
	client := testcp.PingPong(conn2.Bootstrap(ctx))
	defer client.Release()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		ans, release := client.EchoNum(ctx, func(args testcp.PingPong_echoNum_Params) error {
			args.SetN(42)
			return nil
		})
		result, err := ans.Struct()
		if err != nil {
			release()
			b.Errorf("call failed on iteration %d: %v", i, err)
			break
		}
		n := result.N()
		release()
		if n != 42 {
			b.Errorf("n = %d; want 42", n)
			break
		}
	}
}

type pingPongServer struct{}

func (pingPongServer) EchoNum(ctx context.Context, call testcp.PingPong_echoNum) error {
	out, err := call.AllocResults()
	if err != nil {
		return err
	}
	out.SetN(call.Args().N())
	return nil
}
