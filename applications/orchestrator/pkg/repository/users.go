package repository

import (
	"gorm.io/gorm"
	"yumiko_kawaii.com/yine/applications/orchestrator/pkg/models"
)

type IUsers interface {
	IRepository[models.User]
}

type users struct {
	IRepository[models.User]
	db *gorm.DB
}

func NewUsers(db *gorm.DB) IUsers {
	return &users{
		db:          db,
		IRepository: New[models.User](db),
	}
}
