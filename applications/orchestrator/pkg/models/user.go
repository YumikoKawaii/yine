package models

import "time"

type User struct {
	Id             int
	Identification string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
