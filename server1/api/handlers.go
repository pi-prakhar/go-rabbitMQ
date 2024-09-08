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

func (app *App) handleTest(w http.ResponseWriter, r *http.Request) {
	res := "Hello"
	w.Header().Add("content-Type", "text/html")
	w.Write([]byte(res))
}

func (app *App) sendMessageToQueue(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(app.HandlerTimeout))
	defer cancel()

	var res utils.Response
	var message models.Message

	err := json.NewDecoder(r.Body).Decode(&message)
	if err != nil {
		res = &utils.ErrorResponse{
			Message:    "Failed to decode request body",
			StatusCode: http.StatusBadRequest,
			Error:      err.Error(),
		}
		res.Write(w)
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
	}

	log.Printf("Successfully pushed message to queue : %s", app.QueueTest.Name)

	res = &utils.SuccessResponse[string]{
		Message:    fmt.Sprintf("Succesfull Pushed Message to queue %s", app.QueueTest.Name),
		StatusCode: http.StatusOK,
		Data:       message.Message,
	}
	res.Write(w)
}
