package models

import "time"

type UserConversation struct {
	Id                 int       `gorm:"column:id;primaryKey;autoIncrement"`
	UserIdentification string    `gorm:"column:user_identification;type:varchar(255);not null;index"`
	ConversationId     int       `gorm:"column:conversation_id;not null;index"`
	Role               string    `gorm:"column:role;type:varchar(50);not null"`
	CreatedAt          time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt          time.Time `gorm:"column:updated_at;autoUpdateTime"`

	User *User `gorm:"foreignKey:UserIdentification;references:Identification"`
}
