package models

type Message struct {
	Message string `json:"message" validate:"required"`
}

type MessageWithPrirority struct {
	Message  string `json:"message" validate:"required"`
	Priority string `json:"priority"`
}
