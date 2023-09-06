package utilaio

import (
	"errors"
	"sync"
)

type Queue struct {
	Elements chan interface{}
	size     int
	mutex    sync.Mutex
}

func InitQueue(size int) *Queue {
	queue := &Queue{
		Elements: make(chan interface{}, size),
		size:     size,
	}
	return queue
}

func (queue *Queue) Enqueue(element interface{}) error {
	var err error
	queue.mutex.Lock()
	if len(queue.Elements) >= queue.size {
		err = errors.New("queue size full")
	} else {
		queue.Elements <- element
	}
	queue.mutex.Unlock()
	return err
}

func (queue *Queue) Dequeue() interface{} {
	var element interface{}
	queue.mutex.Lock()
	if len(queue.Elements) <= 0 {
		element = nil
	} else {
		element = <-queue.Elements
	}
	queue.mutex.Unlock()
	return element
}
