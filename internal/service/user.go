package service

import (
	"context"
	"github.com/wtkeqrf0/TgBot/ent"
)

type UserPostgres interface {
	CreateUser(ctx context.Context, chatId int64) (*ent.User, error)
	FindUserByChatId(ctx context.Context, chatId int64) (*ent.User, error)
	UpdateUser(ctx context.Context, user *ent.User) error
}

type UserService struct {
	postgres UserPostgres
}

func NewUserService(postgres UserPostgres) *UserService {
	return &UserService{postgres: postgres}
}

func (u UserService) CreateUser(chatId int64) (*ent.User, error) {
	return u.postgres.CreateUser(context.Background(), chatId)
}

func (u UserService) FindUserByChatId(chatId int64) (*ent.User, error) {
	return u.postgres.FindUserByChatId(context.Background(), chatId)
}

func (u UserService) UpdateUser(user *ent.User) error {
	return u.postgres.UpdateUser(context.Background(), user)
}
