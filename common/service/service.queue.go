package service

import (
	"github.com/Benyam-S/aseri/common"
	"github.com/Benyam-S/aseri/entity"
)

// PushQueue is a struct that holds all the requests that needs to be pushed
type PushQueue struct {
	Queue []*entity.ChannelRequest
}

// NewPushQueue is a function that returns a new push queue
func NewPushQueue() common.IPushQueue {
	return &PushQueue{Queue: make([]*entity.ChannelRequest, 0)}
}

// GetQueue is a method that returns all the elements inside the push queue
func (pq *PushQueue) GetQueue() []*entity.ChannelRequest {
	return pq.Queue
}

// AddToQueue is a method that adds new request to the push queue
func (pq *PushQueue) AddToQueue(request *entity.ChannelRequest) {
	pq.Queue = append(pq.Queue, request)
}

// RemoveFromQueueWithIndex is a method that removes an request from the push queue using an index
func (pq *PushQueue) RemoveFromQueueWithIndex(index int) {
	if len(pq.Queue) > index {
		pq.Queue = append(pq.Queue[:index], pq.Queue[index+1:]...)
	}
}
