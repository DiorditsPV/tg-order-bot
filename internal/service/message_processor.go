package service

import (
	"fmt"

	tgBotAPI "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (s *Service) ProcessMessage(update *tgBotAPI.Update) []MessageResponse {
	chatID := update.Message.Chat.ID
	var responses []MessageResponse
	responses = append(responses, MessageResponse{
		ChatID: chatID,
		Text:   "Извините, я не понимаю эту команду. Используйте кнопки меню или /help для списка команд",
	})

	fmt.Println(update.Message.From.UserName, chatID, update.Message.Text)
	return responses
}
