package memory

import (
	"sync"
	"time"
)

type Runtime[T any] struct {
	mem  map[string]T
	lock sync.Locker
}

func NewRuntime[T any]() Memory[T] {
	mem := new(Runtime[T])
	mem.lock = new(sync.Mutex)
	mem.mem = make(map[string]T, 0)

	return mem
}

func (r *Runtime[T]) Pop(key string) (T, bool) {
	r.lock.Lock()
	defer r.lock.Unlock()
	v, ok := r.mem[key]
	if ok {
		delete(r.mem, key)
	}
	return v, ok
}

func (r *Runtime[T]) Put(key string, data T, ttl time.Duration) {
	r.lock.Lock()
	r.mem[key] = data
	r.lock.Unlock()

	if ttl > 0 {
		go func() {
			<-time.After(ttl)
			r.Pop(key)
		}()
	}
}
