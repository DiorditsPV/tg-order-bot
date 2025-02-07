package domain

import tgBotAPI "github.com/go-telegram-bot-api/telegram-bot-api/v5"

type Product struct {
	ID       string
	Name     string
	Category string
	Price    float64
}

// Получение клавиатуры с продуктами для категории
func GetProductsKeyboard(categoryID string) *tgBotAPI.InlineKeyboardMarkup {
	var buttons [][]tgBotAPI.InlineKeyboardButton

	// Ищем категорию
	var category Category
	for _, c := range Categories {
		if c.ID == categoryID {
			category = c
			break
		}
	}

	// Добавляем кнопки продуктов по 2 в ряд
	for i := 0; i < len(category.Products); i += 2 {
		row := make([]tgBotAPI.InlineKeyboardButton, 0)
		row = append(row, tgBotAPI.NewInlineKeyboardButtonData(
			category.Products[i].Name,
			ButtonPrefixProduct+category.Products[i].ID,
		))
		if i+1 < len(category.Products) {
			row = append(row, tgBotAPI.NewInlineKeyboardButtonData(
				category.Products[i+1].Name,
				ButtonPrefixProduct+category.Products[i+1].ID,
			))
		}
		buttons = append(buttons, row)
	}

	// Добавляем навигационные кнопки
	buttons = append(buttons, []tgBotAPI.InlineKeyboardButton{
		tgBotAPI.NewInlineKeyboardButtonData("« Назад", CallbackCategoryMenu),
		tgBotAPI.NewInlineKeyboardButtonData("Удалить из корзины", CallbackOrderDrop),
		tgBotAPI.NewInlineKeyboardButtonData("Оформить заказ", CallbackMakeOrder),
	})

	return &tgBotAPI.InlineKeyboardMarkup{InlineKeyboard: buttons}
}
