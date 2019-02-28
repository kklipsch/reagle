package main

import "sync/atomic"

type atomicBool uint32

func (a atomicBool) Set(b bool) {
	val := uint32(0)
	if b {
		val = uint32(1)
	}

	i := uint32(a)
	atomic.StoreUint32(&i, val)
}

func (a atomicBool) Get() bool {
	i := uint32(a)
	return atomic.LoadUint32(&i) > 0
}
