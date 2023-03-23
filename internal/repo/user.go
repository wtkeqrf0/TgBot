package repo

import (
	"context"
	"github.com/wtkeqrf0/TgBot/ent"
	"github.com/wtkeqrf0/TgBot/ent/user"
)

type UserStorage struct {
	userClient *ent.UserClient
}

func NewUserStorage(userClient *ent.UserClient) *UserStorage {
	return &UserStorage{userClient: userClient}
}

func (r *UserStorage) CreateUser(ctx context.Context, chatId int64) (*ent.User, error) {
	return r.userClient.Create().SetChatID(chatId).SetNumberOfRequests(1).Save(ctx)
}

func (r *UserStorage) FindUserByChatId(ctx context.Context, chatId int64) (*ent.User, error) {
	return r.userClient.Query().Where(user.ChatID(chatId)).Only(ctx)
}

func (r *UserStorage) UpdateUser(ctx context.Context, user *ent.User) error {
	return r.userClient.UpdateOne(user).SetNumberOfRequests(user.NumberOfRequests).Exec(ctx)
}
