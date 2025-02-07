package service

import (
	"fmt"

	tgBotAPI "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (s *Service) ProcessMessage(update *tgBotAPI.Update) []MessageResponse {
	chatID := update.Message.Chat.ID
	var responses []MessageResponse

	switch update.Message.Text {
	case "Сделать заказ":
		responses = append(responses, MessageResponse{
			ChatID: chatID,
			Text:   "Пожалуйста, выберите блюда из меню",
		})
	case "Мои заказы":
		responses = append(responses, MessageResponse{
			ChatID: chatID,
			Text:   "У вас пока нет заказов",
		})
	default:
		responses = append(responses, MessageResponse{
			ChatID: chatID,
			Text:   "Извините, я не понимаю эту команду. Используйте меню или /help для списка команд",
		})
	}

	fmt.Println(update.Message.From.UserName, chatID, update.Message.Text)
	return responses
}
