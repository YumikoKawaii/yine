package models

import "time"

type UserConversation struct {
	Id                 int
	UserIdentification string
	ConversationId     int
	Role               string
	CreatedAt          time.Time
	UpdatedAt          time.Time

	User *User
}
