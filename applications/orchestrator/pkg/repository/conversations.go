package repository

import (
	"gorm.io/gorm"
	"yumiko_kawaii.com/yine/applications/orchestrator/pkg/models"
)

type IConversations interface {
	IRepository[models.Conversation]
}

type conversations struct {
	IRepository[models.Conversation]
	db *gorm.DB
}

func NewConversations(db *gorm.DB) IConversations {
	return &conversations{
		db:          db,
		IRepository: New[models.Conversation](db),
	}
}
