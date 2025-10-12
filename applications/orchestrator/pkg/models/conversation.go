package models

import "time"

type Conversation struct {
	Id        int
	CreatedAt time.Time
	UpdatedAt time.Time
}
