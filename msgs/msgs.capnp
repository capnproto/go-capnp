struct Pointer {
}

struct Message {
	Cookie @0 : Pointer;
	Method @1 : UInt16;
	ReturnCookie @2 : Pointer;
	ReturnMethod @3 : UInt16;
	Arguments @4 : Pointer;
}
