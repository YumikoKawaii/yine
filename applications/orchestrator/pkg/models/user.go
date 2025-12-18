package models

import "time"

type User struct {
	Id             int       `gorm:"column:id;primaryKey;autoIncrement"`
	Identification string    `gorm:"column:identification;type:varchar(255);unique;not null"`
	CreatedAt      time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt      time.Time `gorm:"column:updated_at;autoUpdateTime"`
}
