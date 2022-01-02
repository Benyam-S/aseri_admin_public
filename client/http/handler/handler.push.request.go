package handler

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/Benyam-S/aseri/entity"
	"github.com/gorilla/mux"
)

// BotHandlerErrorResponse is a type that defines a bot handler error response message
type BotHandlerErrorResponse struct {
	Error string `json:"error"`
}

// HandleRePushJob is a handler that handles the request for re-pushing job notification to channel
func (handler *StaffHandler) HandleRePushJob(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	_, ok := ctx.Value(entity.Key("staff_member")).(*entity.Staff)
	if !ok {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	id := mux.Vars(r)["id"]

	job, err := handler.jbService.FindJob(id)
	if err != nil {
		output, _ := json.MarshalIndent(map[string]string{"error": err.Error()}, "", "\t")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(output)
		return
	}

	// Adding data to push queue
	request := new(entity.ChannelRequest)
	request.Type = entity.PushToChannel
	request.Value = job.ID

	handler.pq.AddToQueue(request)
	handler.pushChan <- entity.StartPush
}

// HandlePushRequest is a method that handles push notification sending process to bot handler
func (handler *StaffHandler) HandlePushRequest() {

	domainName := os.Getenv("bot_domain_address")
	port := os.Getenv("bot_client_server_port")

	for range handler.pushChan {

		for {

			i := 0
			requestCount := 0

			if len(handler.pq.GetQueue()) == 0 {
				break
			}

			for {

				if i >= len(handler.pq.GetQueue()) {
					break
				}

				request := handler.pq.GetQueue()[i]

				if requestCount > 15 {
					time.Sleep(time.Second * 30)
					requestCount = 0
				}

				if request.Type == entity.PushForApproval {
					url := "https://" + domainName + ":" + port + "/approval/result/" + request.Value

					newPushRequest := new(entity.ChannelRequest)
					newPushRequest.Type = entity.PushToChannel
					newPushRequest.Value = request.Value

					i = handler.SendToBotHandler(url, i, newPushRequest)
				}

				if request.Type == entity.PushToChannel {
					url := "https://" + domainName + ":" + port + "/push/notification/channel/" + request.Value

					newPushRequest := new(entity.ChannelRequest)
					newPushRequest.Type = entity.PushToSubscribers
					newPushRequest.Value = request.Value

					i = handler.SendToBotHandler(url, i, newPushRequest)
				}

				if request.Type == entity.PushToSubscribers {
					url := "https://" + domainName + ":" + port + "/push/notification/subscriber/" + request.Value
					i = handler.SendToBotHandler(url, i, nil)
				}

				requestCount++
			}
		}

	}
}

// SendToBotHandler is method that sends a push notification to bot handler
func (handler *StaffHandler) SendToBotHandler(url string, index int, newRequest *entity.ChannelRequest) int {

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	response, err := client.Get(url)
	if err != nil {
		handler.logger.LogFileError(err.Error(), entity.ServerLogFile)
		handler.pq.RemoveFromQueueWithIndex(index)
		return index
	}

	if response.StatusCode == 200 {
		if newRequest != nil {
			handler.pq.AddToQueue(newRequest)
		}
		handler.pq.RemoveFromQueueWithIndex(index)
	} else if response.StatusCode == 400 {

		handlerResponse, _ := ioutil.ReadAll(response.Body)
		defer response.Body.Close()

		errorResponse := new(BotHandlerErrorResponse)
		json.Unmarshal(handlerResponse, errorResponse)

		if errorResponse.Error == "retry" {
			index++
			return index
		}

		handler.pq.RemoveFromQueueWithIndex(index)
	} else {
		handler.pq.RemoveFromQueueWithIndex(index)
		handler.logger.Log(
			fmt.Sprintf("Request to %s returned with %d status code.", url, response.StatusCode),
			"w", entity.ServerLogFile)
	}

	return index
}
