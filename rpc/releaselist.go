package rpc

import "capnproto.org/go/capnp/v3"

type releaseList []capnp.ReleaseFunc

func (rl releaseList) Release() {
	for i, r := range rl {
		if r != nil {
			r()
			rl[i] = nil
		}
	}
}

func (rl *releaseList) Add(r capnp.ReleaseFunc) {
	*rl = append(*rl, r)
}
