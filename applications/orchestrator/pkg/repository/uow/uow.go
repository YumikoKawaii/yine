package uow

import (
	"context"

	"gorm.io/gorm"
	"yumiko_kawaii.com/yine/applications/orchestrator/pkg/repository"
)

type IStore interface {
	Users() repository.IUsers
	Messages() repository.IMessages
	Conversations() repository.IConversations
	UserConversations() repository.IUserConversations
}
type store struct {
	users             repository.IUsers
	messages          repository.IMessages
	conversations     repository.IConversations
	userConversations repository.IUserConversations
}

func (s *store) Users() repository.IUsers {
	return s.users
}

func (s *store) Messages() repository.IMessages {
	return s.messages
}

func (s *store) Conversations() repository.IConversations {
	return s.conversations
}

func (s *store) UserConversations() repository.IUserConversations {
	return s.userConversations
}

type worker struct {
	db *gorm.DB
}

type Block func(store IStore) error

type IWorker interface {
	Do(context.Context, Block) error
}

func New(db *gorm.DB) IWorker {
	return &worker{db: db}
}

func (s *worker) Do(_ context.Context, block Block) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		newStore := &store{
			users:             repository.NewUsers(tx),
			messages:          repository.NewMessages(tx),
			conversations:     repository.NewConversations(tx),
			userConversations: repository.NewUserConversations(tx),
		}
		return block(newStore)
	})
}
