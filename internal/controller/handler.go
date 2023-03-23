package controller

import (
	"github.com/go-telegram/bot"
	"github.com/wtkeqrf0/TgBot/ent"
	"github.com/wtkeqrf0/TgBot/pkg/conf"
)

var cfg = conf.GetConfig()

// UserService interacts with the users table
type UserService interface {
	CreateUser(chatId int64) (*ent.User, error)
	FindUserByChatId(chatId int64) (*ent.User, error)
	UpdateUser(user *ent.User) error
}

type Handler struct {
	users UserService
}

func NewHandler(users UserService) *Handler {
	return &Handler{users: users}
}

func (h Handler) InitRoutes(b *bot.Bot) {
	b.RegisterHandler(bot.HandlerTypeMessageText, "/start", bot.MatchTypeExact, h.starts)
	b.RegisterHandler(bot.HandlerTypeMessageText, "❔Help", bot.MatchTypeExact, h.starts)
	b.RegisterHandler(bot.HandlerTypeMessageText, "\U0001FAAAGet my stats", bot.MatchTypeExact, h.myInfo)
}
