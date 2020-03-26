package pool

import (
	"errors"
	"fmt"
	"log"
)

type factory func() interface{}

type Pool struct {
	storage chan interface{}
}

func New(capacity int, f factory) (*Pool, error) {
	if capacity <= 0 {
		return nil, fmt.Errorf("capacity must be greater than 0, actual value: %v", capacity)
	}

	np := Pool{
		storage: make(chan interface{}, capacity),
	}

	for i := 0; i < capacity; i++ {
		inter := f()
		np.storage <- inter
	}

	return &np, nil
}

func (cp *Pool) Get() (interface{}, error) {
	channelInstances := cp.storage
	if channelInstances == nil {
		return nil, errors.New("get interface from pool, storage is not initialized")
	}

	inst := <-channelInstances
	if inst == nil {
		return nil, errors.New("get interface from pool, instance is nil")
	}

	return inst, nil
}

func (cp *Pool) Put(instance interface{}) {
	select {
	case cp.storage <- instance:
	default:
		log.Fatal("put resource back into the pool")
	}
}
