package pool

import (
	"errors"
)

const TimeLayout string = "2006/01/02 15:04:05"

var (
	newPoolError   = errors.New("invalid capasity settings. capasity must be greater than 0")
	getFromPoolErr = errors.New("Can't get interface from pool, because it's closed")
)

type Factory func() interface{}

type Pool struct {
	storage      chan interface{}
	fact_in_pool Factory
}

func NewPool(capacity int, factory Factory) (*Pool, error) {
	if capacity <= 0 {
		return nil, newPoolError
	}

	np := &Pool{
		storage:      make(chan interface{}, capacity),
		fact_in_pool: factory,
	}

	for i := 0; i < capacity; i++ {
		inter := factory()
		np.storage <- inter
	}

	return np, nil
}

func (cp *Pool) Get() (interface{}, error) {
	channelInstances := cp.storage
	if channelInstances == nil {
		return nil, getFromPoolErr
	}

	select {
	case inst := <-channelInstances:
		if inst == nil {
			return nil, getFromPoolErr
		}
		return inst, nil
	}
}

func (cp *Pool) Put(instance interface{}) {
	select {
	case cp.storage <- instance:
	}
}
