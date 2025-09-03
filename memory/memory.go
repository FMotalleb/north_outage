package memory

import "time"

type Memory[T any] interface {
	Pop(string) (T, bool)
	Put(string, T, time.Duration)
}
