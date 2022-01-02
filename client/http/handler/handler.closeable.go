package handler

import (
	"time"

	"github.com/Benyam-S/aseri/entity"
)

// HandleCloseable is a method that closes opened jobs when they reach their due date
func (handler *StaffHandler) HandleCloseable() {

	for {
		time.Sleep(time.Hour * 24)

		now := time.Now()
		closedJobs := handler.jbService.CloseDueJobs(now, entity.JobStatusOpened)

		for _, closedJob := range closedJobs {

			// Adding data to push queue
			request := new(entity.ChannelRequest)
			request.Type = entity.PushToChannel
			request.Value = closedJob.ID

			handler.pq.AddToQueue(request)
			handler.pushChan <- entity.StartPush
		}
	}
}
