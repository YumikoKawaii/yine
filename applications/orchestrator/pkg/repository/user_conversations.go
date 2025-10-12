package repository

import (
	"gorm.io/gorm"
	"yumiko_kawaii.com/yine/applications/orchestrator/pkg/models"
)

type IUserConversations interface {
	IRepository[models.UserConversation]
}

type userConversations struct {
	IRepository[models.UserConversation]
	db *gorm.DB
}

func NewUserConversations(db *gorm.DB) IUserConversations {
	return &userConversations{
		db:          db,
		IRepository: New[models.UserConversation](db),
	}
}

type UserConversationFilter struct {
	ConversationId *int64

	PreloadOption *UserConversationPreloadOption
}

type UserConversationPreloadOption struct {
	Conversation *bool
	User         *bool
}

func (u UserConversationFilter) ApplyFilter(db *gorm.DB) *gorm.DB {
	if u.ConversationId != nil {
		db = db.Where("conversation_id = ?", *u.ConversationId)
	}

	return db
}
