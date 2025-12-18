package models

import "time"

type Message struct {
	Id             int       `gorm:"column:id;primaryKey;autoIncrement"`
	Sender         string    `gorm:"column:sender;type:varchar(255);not null"`
	ConversationId int64     `gorm:"column:conversation_id;not null;index"`
	Content        string    `gorm:"column:content;type:text;not null"`
	Type           string    `gorm:"column:type;type:varchar(50);not null"`
	CreatedAt      time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt      time.Time `gorm:"column:updated_at;autoUpdateTime"`
}
