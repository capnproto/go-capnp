package rpc_test

import (
	"context"
	"fmt"
	"net"
	"testing"

	"capnproto.org/go/capnp/v3"
	"capnproto.org/go/capnp/v3/flowcontrol"
	"capnproto.org/go/capnp/v3/rpc"
	testcp "capnproto.org/go/capnp/v3/rpc/internal/testcapnp"
)

type benchmarkStreamingConfig struct {
	FlowLimit    int64
	MessageCount int
	MessageSize  int
}

func BenchmarkStreaming(b *testing.B) {
	cfg := benchmarkStreamingConfig{
		MessageSize: 32,
	}
	for i := 0; i < 10; i++ {
		cfg.MessageSize *= 2
		cfg.MessageCount = 1 << 16
		cfg.FlowLimit = int64(cfg.MessageSize) * (1 << 12)
		b.Run(fmt.Sprintf("MessageSize=0x%x,MessageCount=0x%x,FlowLimit=0x%x",
			cfg.MessageSize, cfg.MessageCount, cfg.FlowLimit),
			func(b *testing.B) {
				benchmarkStreaming(b, &cfg)
			})
	}
}

func benchmarkStreaming(b *testing.B, cfg *benchmarkStreamingConfig) {
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
	bootstrap.SetFlowLimiter(flowcontrol.NewFixedLimiter(cfg.FlowLimit))
	data := make([]byte, cfg.MessageSize)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for j := 0; j < cfg.MessageCount; j++ {
			err := bootstrap.Push(ctx, func(p testcp.StreamTest_push_Params) error {
				return p.SetData(data)
			})
			if err != nil {
				b.Fatalf("Streaming call #%v failed: %v", j, err)
			}
		}
	}
	if err := bootstrap.WaitStreaming(); err != nil {
		b.Errorf("Error waiting on streaming calls: %v", err)
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
		Logger:          testErrorReporter{tb: b},
		BootstrapClient: capnp.Client(srv),
	})
	defer func() {
		<-conn1.Done()
		if err := conn1.Close(); err != nil {
			b.Error("conn1.Close:", err)
		}
	}()
	conn2 := rpc.NewConn(rpc.NewStreamTransport(p1), &rpc.Options{
		Logger: testErrorReporter{tb: b},
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
