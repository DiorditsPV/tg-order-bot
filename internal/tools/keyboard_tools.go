package tools

import (
	"orderbot/internal/domain"

	tgBotAPI "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func AppendDefaultButtons(keyboard *tgBotAPI.InlineKeyboardMarkup) tgBotAPI.InlineKeyboardMarkup {
	keyboard.InlineKeyboard = append(keyboard.InlineKeyboard,
		[]tgBotAPI.InlineKeyboardButton{
			tgBotAPI.NewInlineKeyboardButtonData("« Назад", domain.CallbackCategoryMenu),
			tgBotAPI.NewInlineKeyboardButtonData("Оформить заказ", domain.CallbackMakeOrder),
		},
	)
	return *keyboard
}
