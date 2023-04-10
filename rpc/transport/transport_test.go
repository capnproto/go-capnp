package transport

import (
	"errors"
	"io"
	"net"
	"testing"
	"time"

	capnp "capnproto.org/go/capnp/v3"
	rpccp "capnproto.org/go/capnp/v3/std/capnp/rpc"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testTransport(t *testing.T, makePipe func() (t1, t2 Transport, err error)) {
	t.Run("Close", func(t *testing.T) {
		t1, t2, err := makePipe()
		if err != nil {
			t.Fatal("makePipe:", err)
		}
		if err := t1.Close(); err != nil {
			t.Error("t1.Close:", err)
		}
		if err := t2.Close(); err != nil {
			t.Error("t2.Close:", err)
		}
	})
	t.Run("Send", func(t *testing.T) {
		t1, t2, err := makePipe()
		if err != nil {
			t.Fatal("makePipe:", err)
		}
		defer func() {
			if err := t1.Close(); err != nil {
				t.Error("t1.Close:", err)
			}
			if err := t2.Close(); err != nil {
				t.Error("t2.Close:", err)
			}
		}()

		// Create messages out of sending order
		callMsg, err := t1.NewMessage()
		if err != nil {
			t.Fatal("t1.NewMessage #1:", err)
		}
		defer callMsg.Release()
		bootMsg, err := t1.NewMessage()
		if err != nil {
			t.Fatal("t1.NewMessage #2:", err)
		}
		defer bootMsg.Release()

		// Fill in bootstrap message
		boot, err := bootMsg.Message().NewBootstrap()
		if err != nil {
			t.Fatal("NewBootstrap:", err)
		}
		boot.SetQuestionId(42)

		// Fill in call message
		call, err := callMsg.Message().NewCall()
		if err != nil {
			t.Fatal("NewCall:", err)
		}
		call.SetQuestionId(123)
		call.SetInterfaceId(456)
		call.SetMethodId(7)
		tgt, err := call.NewTarget()
		if err != nil {
			t.Fatal("NewTarget:", err)
		}
		pa, err := tgt.NewPromisedAnswer()
		if err != nil {
			t.Fatal("NewPromisedAnswer:", err)
		}
		pa.SetQuestionId(42)
		params, err := call.NewParams()
		if err != nil {
			t.Fatal("NewParams:", err)
		}
		// simulate mutating CapTable
		callMsg.Message().Message().CapTable().Add(capnp.ErrorClient(errors.New("foo")))
		callMsg.Message().Message().CapTable().Reset()
		capPtr := capnp.NewInterface(params.Segment(), 0).ToPtr()
		if err := params.SetContent(capPtr); err != nil {
			t.Fatal("SetContent:", err)
		}
		capTable, err := params.NewCapTable(1)
		if err != nil {
			t.Fatal("NewCapTable:", err)
		}
		capTable.At(0).SetSenderHosted(777)

		// Send/receive first message (bootstrap)
		if err := bootMsg.Send(); err != nil {
			t.Fatal("sendBoot():", err)
		}
		bootMsg.Release()
		r1, err := t2.RecvMessage()
		if err != nil {
			t.Fatal("t2.RecvMessage:", err)
		}
		if r1.Message().Message().CapTable().Len() != 0 {
			t.Error("t2.RecvMessage(ctx).Message().CapTable() is not empty")
		}
		if r1.Message().Which() != rpccp.Message_Which_bootstrap {
			t.Errorf("t2.RecvMessage(ctx).Which = %v; want bootstrap", r1.Message().Which())
		} else {
			rboot, _ := r1.Message().Bootstrap()
			if rboot.QuestionId() != 42 {
				t.Errorf("t2.RecvMessage(ctx).Bootstrap.QuestionID = %d; want 42", rboot.QuestionId())
			}
		}
		r1.Release()

		// Send/receive second message (call)
		if err := callMsg.Send(); err != nil {
			t.Fatal("sendCall():", err)
		}
		callMsg.Release()
		r2, err := t2.RecvMessage()
		if err != nil {
			t.Fatal("t2.RecvMessage:", err)
		}
		if r2.Message().Message().CapTable().Len() != 0 {
			t.Error("t2.RecvMessage(ctx).Message().CapTable() is not empty")
		}
		if r2.Message().Which() != rpccp.Message_Which_call {
			t.Errorf("t2.RecvMessage(ctx).Which = %v; want call", r2.Message().Which())
		} else {
			rcall, _ := r2.Message().Call()
			if rcall.QuestionId() != 123 {
				t.Errorf("t2.RecvMessage(ctx).Call.QuestionID = %d; want 123", rcall.QuestionId())
			}
			if rcall.InterfaceId() != 456 {
				t.Errorf("t2.RecvMessage(ctx).Call.InterfaceID = %d; want 456", rcall.InterfaceId())
			}
			if rcall.MethodId() != 7 {
				t.Errorf("t2.RecvMessage(ctx).Call.MethodID = %d; want 7", rcall.InterfaceId())
			}
			rparams, _ := rcall.Params()
			rctab, _ := rparams.CapTable()
			if rctab.Len() != 1 {
				t.Errorf("len(t2.RecvMessage(ctx).Call.Params.CapTable) = %d; want 1", rctab.Len())
			} else if rctab.At(0).Which() != rpccp.CapDescriptor_Which_senderHosted {
				t.Errorf("t2.RecvMessage(ctx).Call.Params.CapTable.Which = %v; want senderHosted", rctab.At(0).Which())
			} else if rctab.At(0).SenderHosted() != 777 {
				t.Errorf("t2.RecvMessage(ctx).Call.Params.CapTable.SenderHosted = %d; want 777", rctab.At(0).SenderHosted())
			}
		}
		r2.Release()
	})
	t.Run("InterruptRecv", func(t *testing.T) {
		t1, t2, err := makePipe()
		require.NoError(t, err, "makePipe should not fail")

		go func() {
			time.Sleep(100 * time.Millisecond)
			t1.Close()
		}()
		_, err = t1.RecvMessage() // hangs here if doesn't work
		assert.Error(t, err,
			"RecvMessage() should return error when transport is closed")

		err = t2.Close()
		assert.NoError(t, err,
			"Close() should not return error after remote side closes")

	})
}

func TestTCPStreamTransport(t *testing.T) {
	t.Run("Unpacked", func(t *testing.T) {
		t.Parallel()

		testTCPStreamTransport(t, NewStream)
	})

	t.Run("Packed", func(t *testing.T) {
		t.Parallel()

		testTCPStreamTransport(t, NewPackedStream)
	})
}

func testTCPStreamTransport(t *testing.T, newTransport func(io.ReadWriteCloser) Transport) {
	type listenCall struct {
		c   *net.TCPConn
		err error
	}

	makePipe := func() (t1, t2 Transport, err error) {
		host, err := net.LookupIP("localhost")
		if err != nil {
			return nil, nil, err
		}
		l, err := net.ListenTCP("tcp", &net.TCPAddr{IP: host[0]})
		if err != nil {
			return nil, nil, err
		}
		ch := make(chan listenCall)
		abort := make(chan struct{})
		go func() {
			c, err := l.AcceptTCP()
			select {
			case ch <- listenCall{c, err}:
			case <-abort:
				c.Close()
			}
		}()
		laddr := l.Addr().(*net.TCPAddr)
		c2, err := net.DialTCP("tcp", nil, laddr)
		if err != nil {
			close(abort)
			l.Close()
			return nil, nil, err
		}
		lc := <-ch
		if lc.err != nil {
			c2.Close()
			l.Close()
			return nil, nil, err
		}
		return newTransport(lc.c), newTransport(c2), nil
	}

	t.Run("ServerToClient", func(t *testing.T) {
		testTransport(t, makePipe)
	})

	t.Run("ClientToServer", func(t *testing.T) {
		testTransport(t, func() (t1, t2 Transport, err error) {
			t2, t1, err = makePipe()
			return
		})
	})
}
