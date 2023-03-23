package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/sirupsen/logrus"
	"github.com/wtkeqrf0/TgBot/ent"
	"github.com/wtkeqrf0/TgBot/internal/controller/dto"
	"net/http"
)

const (
	weather = "https://api.openweathermap.org/data/2.5/weather?lat=%v&lon=%v&appid=%s&units=metric"
	direct  = "http://api.openweathermap.org/geo/1.0/direct?q=%s&limit=1&appid=%s"
)

func (h Handler) starts(ctx context.Context, b *bot.Bot, update *models.Update) {

	user, err := h.getUserInfo(update.Message.Chat.ID)
	if err != nil {
		logrus.WithError(err).WithField("chat_id", update.Message.Chat.ID).Error("error occurred during user search")
		_, err = b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "An error has occurred. Try running the command again",
		})

		if err != nil {
			logrus.WithError(err).WithField("chat_id", update.Message.Chat.ID).Error("error occurred while sending a message")
		}
		return
	}

	_, err = b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,

		ReplyMarkup: models.ReplyKeyboardMarkup{
			ResizeKeyboard: true,
			Keyboard: [][]models.KeyboardButton{{
				models.KeyboardButton{Text: "❔Help"},
				models.KeyboardButton{Text: "\U0001FAAAGet my stats"},
			},
				{models.KeyboardButton{Text: "🗺Get weather form my location", RequestLocation: true}},
			},
		},

		Text: fmt.Sprintf("You are the %v user of GBot!\n"+
			"Write down the name of the city whose weather you want to know."+
			" You can also send the location!", user.ID),
	})
	if err != nil {
		logrus.WithError(err).WithField("chat_id", update.Message.Chat.ID).Error("error occurred while sending a message")
		return
	}
	logrus.WithField("chat_id", update.Message.Chat.ID).Info(update.Message.Text)
}

func (h Handler) myInfo(ctx context.Context, b *bot.Bot, update *models.Update) {
	user, err := h.getUserInfo(update.Message.Chat.ID)
	if err != nil {
		logrus.WithError(err).WithField("chat_id", update.Message.Chat.ID).Error("error occurred during user search")

		_, err = b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "An error has occurred. Try running the command again",
		})

		if err != nil {
			logrus.WithError(err).WithField("chat_id", update.Message.Chat.ID).Error("error occurred while sending a message")
		}
		return
	}

	res := fmt.Sprintf("You are the %v user of GBot!\n\n"+
		"At the moment you have requested weather information %v times\n\n"+
		"The last time you requested the weather was %v\n\n"+
		"For the first time you found out the weather with my help %v",
		user.ID, user.NumberOfRequests,
		user.UpdateTime.Format("2006.01.02 15:02:05"),
		user.CreateTime.Format("2006.01.02 15:02:05"))

	_, err = b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   res,
	})

	if err != nil {
		logrus.WithError(err).WithField("chat_id", update.Message.Chat.ID).Error("error occurred while sending a message")
		return
	}
	logrus.WithField("chat_id", update.Message.Chat.ID).Info(update.Message.Text)
}

func (h Handler) GetWeather(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.Message.Text == "" && update.Message.Location == nil {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Unsupported message type",
		})
		logrus.WithField("chat_id", update.Message.Chat.ID).Info("unsupported message type")
		return
	}

	user, err := h.getUserInfo(update.Message.Chat.ID)
	if err != nil {
		logrus.WithError(err).WithField("chat_id", update.Message.Chat.ID).Error("error occurred during user search")

		_, err = b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "An error has occurred. Try running the command again",
		})

		if err != nil {
			logrus.WithError(err).WithField("chat_id", update.Message.Chat.ID).Error("error occurred while sending a message")
		}
		return
	}
	res := new(dto.WeatherResp)
	if update.Message.Location != nil {

		if err = getJson(fmt.Sprintf(weather, update.Message.Location.Latitude,
			update.Message.Location.Longitude, cfg.Bot.WeatherKey), res); err != nil || len(res.Weather) < 1 {

			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.Message.Chat.ID,
				Text:   fmt.Sprintf("There is no such place!"),
			})

			logrus.WithField("chat_id", update.Message.Chat.ID).Infof("Invalid location: %v", update.Message.Location)
			return
		}
	} else {
		resp := new([]dto.CityParameters)

		if err = getJson(fmt.Sprintf(direct, update.Message.Text, cfg.Bot.WeatherKey), resp); err != nil {

			logrus.WithError(err).WithField("chat_id", update.Message.Chat.ID).Error("error was occurred")
			return
		}
		if len(*resp) < 1 {

			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.Message.Chat.ID,
				Text:   fmt.Sprintf("There is no such city!"),
			})

			logrus.WithField("chat_id", update.Message.Chat.ID).Infof("Invalid city: %s", update.Message.Text)
			return
		}

		if err = getJson(fmt.Sprintf(weather, (*resp)[0].Lat, (*resp)[0].Lon, cfg.Bot.WeatherKey), res); err != nil {
			logrus.WithError(err).WithField("chat_id", update.Message.Chat.ID).Error("error was occurred")
			return
		}
	}

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text: fmt.Sprintf("The weather is now in %s: %s. Temperature: %v°C",
			res.Name, res.Weather[0].Main, res.Main.Temp),
	})

	user.NumberOfRequests++
	if err = h.users.UpdateUser(user); err != nil {

		logrus.WithError(err).WithField("chat_id", update.Message.Chat.ID).Error("error occurred during save data")
		return
	}

	logrus.WithField("chat_id", update.Message.Chat.ID).Info(update.Message.Text)
}

// getUserInfo returns user from db. Creates user, if not exist.
func (h Handler) getUserInfo(chatId int64) (*ent.User, error) {
	user, err := h.users.FindUserByChatId(chatId)
	if err != nil {
		return h.users.CreateUser(chatId)
	}
	return user, err
}

// getJson writes JSON to the target
func getJson(url string, target interface{}) error {
	r, err := http.Get(url)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	return json.NewDecoder(r.Body).Decode(target)
}
