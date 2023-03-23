package exceptions

import (
	"context"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/sirupsen/logrus"
)

// ServerError err
var (
	ServerError = newError("Server exception was occurred")
)

// ErrorHandler used for error handling. Handles only MyError type errors
func ErrorHandler(err error) {
	if my, ok := err.(MyError); ok {
		logrus.WithError(my.Err).Error(my.Msg)
	} else {
		logrus.WithError(err).Error("UNEXPECTED ERROR")
	}
}

// PanicRecovery recovers any panic
func PanicRecovery(next bot.HandlerFunc) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		defer func() {
			if r := recover(); r != nil {
				logrus.Errorf("panic was recovered: %v", r)
			}
		}()
		next(ctx, b, update)
	}
}
