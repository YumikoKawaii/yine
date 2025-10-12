package repository

import (
	"gorm.io/gorm"
	"yumiko_kawaii.com/yine/applications/orchestrator/pkg/models"
)

type IMessages interface {
	IRepository[models.Message]
}

type messages struct {
	IRepository[models.Message]
	db *gorm.DB
}

func NewMessages(db *gorm.DB) IMessages {
	return &messages{
		db:          db,
		IRepository: New[models.Message](db),
	}
}
