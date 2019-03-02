package client

import "sync/atomic"

//AtomicBool stores boolean values atomically
type AtomicBool uint32

//StoreBool atomically stores val into addr.
func StoreBool(addr AtomicBool, b bool) {
	val := uint32(0)
	if b {
		val = uint32(1)
	}

	i := uint32(addr)
	atomic.StoreUint32(&i, val)
}

//LoadBool atomically loads addr.
func LoadBool(addr AtomicBool) bool {
	i := uint32(addr)
	return atomic.LoadUint32(&i) > 0
}
