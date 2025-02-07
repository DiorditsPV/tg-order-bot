package service

import (
	"orderbot/internal/domain"

	tgBotAPI "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const menuPhotoPath = "menu/menu.jpg"

func (s *Service) sendMenu(chatID int64, messageText string) []MessageResponse {

	keyboard := tgBotAPI.NewReplyKeyboard(
		tgBotAPI.NewKeyboardButtonRow(
			tgBotAPI.NewKeyboardButton("Сделать заказ"),
			tgBotAPI.NewKeyboardButton("Мои заказы"),
		),
	)

	return []MessageResponse{
		{
			ChatID:    chatID,
			PhotoPath: menuPhotoPath,
		},
		{
			ChatID:      chatID,
			Text:        messageText,
			ReplyMarkup: &keyboard,
		},
	}
}

func (s *Service) sendInlineMenu(chatID int64, messageText string) []MessageResponse {
	keyboard := tgBotAPI.NewInlineKeyboardMarkup(
		tgBotAPI.NewInlineKeyboardRow(
			tgBotAPI.NewInlineKeyboardButtonData("Сделать заказ", domain.CallbackStartOrder),
			tgBotAPI.NewInlineKeyboardButtonData("Мои заказы", domain.CallbackOrdersList),
		),
	)

	return []MessageResponse{
		{
			ChatID:            chatID,
			PhotoPath:         menuPhotoPath,
			Text:              messageText,
			InlineReplyMarkup: &keyboard,
		},
	}
}

func (s *Service) ProcessCommand(update *tgBotAPI.Update) []MessageResponse {
	chatID := update.Message.Chat.ID
	command := update.Message.Command()
	var responses []MessageResponse

	switch command {
	case "start":
		if !s.repo.CheckSession(chatID) {
			chatSession := domain.NewSession(chatID)
			s.repo.Save(chatSession)
		}

		responses = s.sendInlineMenu(chatID, "Выберите действие:")

		session := s.repo.Get(chatID)
		session.MarkMenuSended()
		s.repo.Save(session)

	case "resendMenu":
		responses = s.sendInlineMenu(chatID, "Вот меню и доступные действия:")

	case "help":
		responses = append(responses, MessageResponse{
			ChatID: chatID,
			Text:   "Доступные команды:\n/start - Начать заказ\n/help - Помощь\n/resend_menu - Показать меню снова",
		})
	}

	return responses
}
