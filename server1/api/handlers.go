package api

import (
	"context"
	"encoding/json"
	"fmt"
	"go-rabbitmq-server1/internal/models"
	"go-rabbitmq-server1/utils"
	"log"
	"net/http"
	"time"
)

const (
	MESSAGE_EXCHANGE   = "messages"
	BROADCAST_EXCHANGE = "broadcast"
)

func (app *App) handleTest(w http.ResponseWriter, r *http.Request) {
	res := "Hello"
	w.Header().Add("content-Type", "text/html")
	w.Write([]byte(res))
}

func (app *App) sendMessageToQueue(w http.ResponseWriter, r *http.Request) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(app.HandlerTimeout))
	defer cancel()

	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var res utils.Response
	var message models.Message

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&message)
	if err != nil {
		res = &utils.ErrorResponse{
			Message:    "Failed to decode request body",
			StatusCode: http.StatusBadRequest,
			Error:      err.Error(),
		}
		res.Write(w)
		return
	}

	//validate
	err = utils.ValidateStruct(message)
	if err != nil {
		res = &utils.ErrorResponse{
			Message:    "Required fields are missing",
			StatusCode: http.StatusBadRequest,
			Error:      err.Error(),
		}
		res.Write(w)
		return
	}

	//push to queue
	if err = app.Producer.Push(ctx,
		app.QueueTest.Name,
		message.Message,
		"",
		false,
		false,
	); err != nil {
		res = &utils.ErrorResponse{
			Message:    fmt.Sprintf("Failed to push message to queue : %s", app.QueueTest.Name),
			StatusCode: http.StatusInternalServerError,
			Error:      err.Error(),
		}
		res.Write(w)
		return
	}

	res = &utils.SuccessResponse[string]{
		Message:    fmt.Sprintf("Succesfull Pushed Message to queue %s", app.QueueTest.Name),
		StatusCode: http.StatusOK,
		Data:       message.Message,
	}
	res.Write(w)

	log.Printf("Successfully pushed message to queue : %s", app.QueueTest.Name)
}

func (app *App) sendMessageToExchange(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(app.HandlerTimeout))
	defer cancel()

	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var res utils.Response
	var message models.MessageWithPrirority
	var err error

	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&message)
	if err != nil {
		res = &utils.ErrorResponse{
			Message:    "Failed to decode request body",
			StatusCode: http.StatusBadRequest,
			Error:      err.Error(),
		}
		res.Write(w)
		return
	}

	//validate
	err = utils.ValidateStruct(message)
	if err != nil {
		res = &utils.ErrorResponse{
			Message:    "Required fields are missing",
			StatusCode: http.StatusBadRequest,
			Error:      err.Error(),
		}
		res.Write(w)
		return
	}

	//push to exchange
	err = app.Producer.PushToExchange(ctx,
		message.Priority,
		message.Message,
		MESSAGE_EXCHANGE,
		false,
		false,
	)
	if err != nil {
		res = &utils.ErrorResponse{
			Message:    fmt.Sprintf("Failed to push message to exchange : %s with priority : %s", MESSAGE_EXCHANGE, message.Priority),
			StatusCode: http.StatusInternalServerError,
			Error:      err.Error(),
		}
		res.Write(w)
		return
	}

	res = &utils.SuccessResponse[string]{
		Message:    fmt.Sprintf("Succesfull Pushed Message to exchange %s with priority : %s", MESSAGE_EXCHANGE, message.Priority),
		StatusCode: http.StatusOK,
		Data:       message.Message,
	}
	res.Write(w)

	log.Printf("Successfully pushed message to exchange : %s with priority : %s", MESSAGE_EXCHANGE, message.Priority)
}
