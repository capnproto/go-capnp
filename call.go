package capn

import (
	"time"
)

type Call struct {
	Message  Message
	Async    bool
	Deadline time.Time
	Reply    Pointer
	Send     func(*Call) error
}

/*
struct Message {
	Object @0 : Pointer;
	Method @1 : UInt16;
	Arguments @2 : Pointer;
	ReplyObject @3 : Pointer;
	ReplyMethod @4 : UInt16;
}
*/

type Message struct {
	P Pointer
}

const (
	Message_Object      = 0
	Message_Method      = 0
	Message_Arguments   = 1
	Message_ReplyObject = 2
	Message_ReplyMethod = 16
)

func (p Message) Object() Pointer      { return p.P.ReadPtr(0) }
func (p Message) Method() uint16       { return p.P.ReadStruct16(0) }
func (p Message) Arguments() Pointer   { return p.P.ReadPtr(1) }
func (p Message) ReplyObject() Pointer { return p.P.ReadPtr(2) }
func (p Message) ReplyMethod() uint16  { return p.P.ReadStruct16(16) }

func (p Message) SetObject(obj Pointer) error        { return p.P.WritePtr(0, obj) }
func (p Message) SetMethod(method uint16) error      { return p.P.WriteStruct16(0, method) }
func (p Message) SetArguments(args Pointer) error    { return p.P.WritePtr(1, args) }
func (p Message) SetReplyObject(obj Pointer) error   { return p.P.WritePtr(2, obj) }
func (p Message) SetReplyMethod(method uint16) error { return p.P.WriteStruct16(16, method) }
