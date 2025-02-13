package queue

type Key interface {
	uint | uint8 | uint16 | uint64 | string
}

type IQueue[K Key, T any] interface {
	Process(K, func(T, ack func(K)) error) error
	Add(K, T) error
}

type RedisQueue struct {
}
