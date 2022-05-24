package rpc_test

import (
	"context"
	"testing"

	"capnproto.org/go/capnp/v3/rpc"
	testcp "capnproto.org/go/capnp/v3/rpc/internal/testcapnp"
	"capnproto.org/go/capnp/v3/rpc/transport"
)

func BenchmarkPingPong(b *testing.B) {
	p1, p2 := transport.NewPipe(1)
	srv := testcp.PingPong_ServerToClient(pingPongServer{}, nil)
	conn1 := rpc.NewConn(p2, &rpc.Options{
		ErrorReporter:   testErrorReporter{tb: b},
		BootstrapClient: srv.Client,
	})
	defer func() {
		<-conn1.Done()
		if err := conn1.Close(); err != nil {
			b.Error("conn1.Close:", err)
		}
	}()
	conn2 := rpc.NewConn(p1, &rpc.Options{
		ErrorReporter: testErrorReporter{tb: b},
	})
	defer func() {
		if err := conn2.Close(); err != nil {
			b.Error("conn2.Close:", err)
		}
	}()

	ctx := context.Background()
	client := testcp.PingPong{Client: conn2.Bootstrap(ctx)}
	defer client.Client.Release()
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
