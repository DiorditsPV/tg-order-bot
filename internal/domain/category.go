package domain

import tgBotAPI "github.com/go-telegram-bot-api/telegram-bot-api/v5"

type Category struct {
	ID       string
	Name     string
	Products []Product
}

// Получение клавиатуры для категории
func GetCategoryKeyboard() *tgBotAPI.InlineKeyboardMarkup {
	buttons := [][]tgBotAPI.InlineKeyboardButton{}

	for _, category := range Categories {

		buttons = append(buttons, []tgBotAPI.InlineKeyboardButton{
			tgBotAPI.NewInlineKeyboardButtonData(
				category.Name,
				ButtonPrefixCategory+category.ID,
			),
		})
	}

	return &tgBotAPI.InlineKeyboardMarkup{InlineKeyboard: buttons}
}
