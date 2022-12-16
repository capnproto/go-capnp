package rpc

import "capnproto.org/go/capnp/v3"

type releaseList []capnp.ReleaseFunc

func (rl releaseList) Release() {
	for _, r := range rl {
		r()
	}
	for i := range rl {
		rl[i] = func() {}
	}
}
