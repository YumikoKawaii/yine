package models

import "time"

type Message struct {
	Id             int
	Sender         string
	ConversationId int64
	Content        string
	Type           string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
