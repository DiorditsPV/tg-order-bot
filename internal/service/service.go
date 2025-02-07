package service

import (
	"orderbot/internal/repository"

	tgBotAPI "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type UpdatePayload struct {
	MessageID int
	Text      string
}

type MessageResponse struct {
	ChatID            int64
	Text              string
	PhotoPath         string
	InlineReplyMarkup *tgBotAPI.InlineKeyboardMarkup
	ReplyMarkup       *tgBotAPI.ReplyKeyboardMarkup
	DeletePrevious    bool
	UpdatePayload     *UpdatePayload
}

func (m *MessageResponse) HasUpdatePayload() bool {
	return m.UpdatePayload != nil
}

func (m *MessageResponse) HasPhoto() bool {
	return m.PhotoPath != ""
}

// есть текст
func (m *MessageResponse) IsText() bool {
	return m.Text != ""
}

type ServiceInterface interface {
	Process(update *tgBotAPI.Update) []MessageResponse
	ProcessMessage(update *tgBotAPI.Update) []MessageResponse
	ProcessCommand(update *tgBotAPI.Update) []MessageResponse
	ProcessCallback(update *tgBotAPI.Update) []MessageResponse
}

type Service struct {
	repo repository.Repository
}

func NewService(repo *repository.SessionRepository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Process(update *tgBotAPI.Update) []MessageResponse {
	if update.Message != nil && update.Message.IsCommand() {
		return s.ProcessCommand(update)
	} else if update.CallbackQuery != nil {
		return s.ProcessCallback(update)
	} else if update.Message != nil {
		return s.ProcessMessage(update)
	}
	return nil
}
