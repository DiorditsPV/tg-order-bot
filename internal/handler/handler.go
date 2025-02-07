package handler

import (
	"orderbot/internal/service"
	"orderbot/internal/tools"
	"os"

	tgBotAPI "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Handler struct {
	srv service.ServiceInterface
	bot *tgBotAPI.BotAPI
}

func NewHandler(srv service.ServiceInterface, bot *tgBotAPI.BotAPI) *Handler {
	return &Handler{srv: srv, bot: bot}
}

func handlePhoto(response service.MessageResponse) tgBotAPI.Chattable {
	photoBytes, err := os.ReadFile(response.PhotoPath)
	if err != nil {
		tools.LogError("Ошибка чтения фото: %v", err)
		return nil
	}
	photoMsg := tgBotAPI.NewPhoto(response.ChatID, tgBotAPI.FileBytes{
		Name:  "menu",
		Bytes: photoBytes,
	})
	if response.InlineReplyMarkup != nil {
		photoMsg.ReplyMarkup = response.InlineReplyMarkup
	}
	return photoMsg
}

func handleText(response service.MessageResponse) tgBotAPI.Chattable {
	textMsg := tgBotAPI.NewMessage(response.ChatID, response.Text)
	if response.InlineReplyMarkup != nil {
		textMsg.ReplyMarkup = *response.InlineReplyMarkup
	}
	return textMsg
}

func (h *Handler) Start() {
	u := tgBotAPI.NewUpdate(0)
	u.Timeout = 30

	tools.LogInfo("Бот запущен и ожидает сообщений")
	upd_chan := h.bot.GetUpdatesChan(u)
	for update := range upd_chan {
		responses := h.srv.Process(&update)
		for _, response := range responses {
			var msg tgBotAPI.Chattable

			if response.HasUpdatePayload() {
				if !response.HasPhoto() {
					if update.CallbackQuery.Message.Text != response.Text {
						msg = tgBotAPI.NewEditMessageTextAndMarkup(response.ChatID, update.CallbackQuery.Message.MessageID, response.Text, *response.InlineReplyMarkup)
					} else {
						tools.LogInfo("Сообщение не изменилось, пропускаем обновление")
					}
				} else {
					tools.LogWarn("Не реализовано редактирование фото")
				}
			} else {
				if response.HasPhoto() {
					msg = handlePhoto(response)
				} else {
					msg = handleText(response)
				}
			}

			if msg != nil {
				if _, err := h.bot.Send(msg); err != nil {
					tools.LogError("Ошибка отправки сообщения: %v", err)
				}
			}
		}
	}
}
