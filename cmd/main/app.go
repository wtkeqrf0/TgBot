package main

import (
	"context"
	"github.com/go-telegram/bot"
	"github.com/sirupsen/logrus"
	"github.com/wtkeqrf0/TgBot/ent"
	"github.com/wtkeqrf0/TgBot/internal/controller"
	"github.com/wtkeqrf0/TgBot/internal/middleware/exceptions"
	"github.com/wtkeqrf0/TgBot/internal/repo"
	"github.com/wtkeqrf0/TgBot/internal/service"
	"github.com/wtkeqrf0/TgBot/pkg/client/postgresql"
	"github.com/wtkeqrf0/TgBot/pkg/conf"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func init() {
	logrus.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:02:05",
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyLevel: "status",
			logrus.FieldKeyFunc:  "caller",
			logrus.FieldKeyMsg:   "message",
		},
	})

	logrus.SetReportCaller(true)
}

func main() {
	cfg := conf.GetConfig()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	pClient := postgresql.Open(cfg.Postgres.Username, cfg.Postgres.Password,
		cfg.Postgres.Host, cfg.Postgres.Port, cfg.Postgres.DBName)

	h := controller.NewHandler(service.NewUserService(repo.NewUserStorage(pClient.User)))

	opts := []bot.Option{
		bot.WithMiddlewares(exceptions.PanicRecovery),
		bot.WithErrorsHandler(exceptions.ErrorHandler),
		bot.WithDefaultHandler(h.GetWeather),
	}

	b, err := bot.New("6049622717:AAEps2leDafnL566U03ijOwhDS4hYIesVTg", opts...)
	if err != nil {
		logrus.Fatalf("error was occured while creating bot connection: %v", err)
	}

	h.InitRoutes(b)
	Run(ctx, b, pClient)
}

// Run the Server with graceful shutdown
func Run(ctx context.Context, b *bot.Bot, pClient *ent.Client) {

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go b.Start(ctx)
	logrus.Info("Server Started")

	<-quit

	logrus.Info("Server Shutting Down ...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := pClient.Close(); err != nil {
		logrus.Fatalf("PostgreSQL Connection Shutdown Failed: %s", err)
	}

	if ok, err := b.Close(ctx); err != nil && !ok {
		logrus.Fatalf("Bot Shutdown Failed: %v. This error returns for the first 10 minutes after the bot is launched", err)
	}

	logrus.Info("Server Exited Properly")
}
