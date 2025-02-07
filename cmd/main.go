package main

import (
	"fmt"
	"orderbot/internal/handler"
	"orderbot/internal/repository"
	"orderbot/internal/service"
	"os"

	"orderbot/internal/tools"

	tgBotAPI "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	fmt.Println("Старт Go")

	stopFlag := make(chan bool)

	// Инициализация бота
	bot, err := initBot()
	if err != nil {
		panic(err)
	}

	// Инициализация логгера
	tools.InitLogger(bot)

	// Инициализация слоев приложения
	app := initApp(bot)

	// Запуск обработчика
	app.handler.Start()

	<-stopFlag
	fmt.Println("Стоп Go")
}

func initBot() (*tgBotAPI.BotAPI, error) {
	bot, err := tgBotAPI.NewBotAPI(os.Getenv("BOT_TOKEN"))
	if err != nil {
		fmt.Printf("Ошибка при создании бота %s", err)
		return nil, err
	}
	return bot, nil
}

type application struct {
	repository *repository.SessionRepository
	service    *service.Service
	handler    *handler.Handler
}

func initApp(bot *tgBotAPI.BotAPI) *application {
	repo := repository.NewRepository()
	srv := service.NewService(repo)
	hndlr := handler.NewHandler(srv, bot)

	return &application{
		repository: repo,
		service:    srv,
		handler:    hndlr,
	}
}
