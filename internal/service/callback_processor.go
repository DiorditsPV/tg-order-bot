package service

import (
	"fmt"
	"orderbot/internal/domain"
	"orderbot/internal/tools"
	"strings"

	tgBotAPI "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (s *Service) ProcessCallback(update *tgBotAPI.Update) []MessageResponse {
	chatID := update.CallbackQuery.Message.Chat.ID
	lastMessageID := update.CallbackQuery.Message.MessageID
	callbackData := update.CallbackQuery.Data
	username := update.CallbackQuery.From.UserName

	logCtx := &tools.LogContext{
		ChatID:   chatID,
		Username: username,
	}

	tools.LogInfoWithContext(logCtx, "Получен callback: %s", callbackData)

	if !s.repo.CheckSession(chatID) {
		chatSession := domain.NewSession(chatID)
		s.repo.Save(chatSession)
		tools.LogInfoWithContext(logCtx, "Создана новая сессия")
	}
	chatSession := s.repo.Get(chatID)

	var responses []MessageResponse
	switch {
	case callbackData == domain.CloseOrder:
		chatSession.ClearOrderBag()
		tools.LogInfoWithContext(logCtx, "Заказ отменен")

	case callbackData == domain.ApproveOrder:
		tools.LogInfoWithContext(logCtx, "Заказ подтвержден")
		s.sendOrderToAdmin()

	case callbackData == domain.CallbackMakeOrder:
		tools.LogInfoWithContext(logCtx, "Оформление заказа")
		responses = s.sendOrder(chatID, chatSession)

	case callbackData == domain.CallbackStartOrder || callbackData == domain.CallbackCategoryMenu:
		tools.LogInfoWithContext(logCtx, "Переход к выбору категории")
		responses = s.sendCategoryMenu(chatID, callbackData)

	case strings.HasPrefix(callbackData, domain.ButtonPrefixCategory):
		categoryID := strings.TrimPrefix(callbackData, domain.ButtonPrefixCategory)
		tools.LogInfoWithContext(logCtx, "Выбрана категория: %s", categoryID)
		responses = s.chooseProducts(chatID, categoryID)

	case strings.HasPrefix(callbackData, domain.ButtonPrefixProduct):
		productID := strings.TrimPrefix(callbackData, domain.ButtonPrefixProduct)
		if product, ok := domain.ProductsMap[productID]; ok {
			tools.LogInfoWithContext(logCtx, "Добавлен товар в корзину: %s (%.2f руб.)", product.Name, product.Price)
		}
		responses = s.addProductToCart(chatID, productID, chatSession, lastMessageID)

	case callbackData == domain.CallbackOrderDrop:
		tools.LogInfoWithContext(logCtx, "Открыто меню удаления товаров")
		responses = []MessageResponse{*s.sendOrderDetails(chatID, chatSession, domain.ButtonPrefixDrop)}

	case strings.HasPrefix(callbackData, domain.ButtonPrefixDrop):
		productID := strings.TrimPrefix(callbackData, domain.ButtonPrefixDrop)
		if product, ok := domain.ProductsMap[productID]; ok {
			tools.LogInfoWithContext(logCtx, "Удален товар из корзины: %s", product.Name)
		}
		dropResponse := s.dropProductFromCart(chatID, productID, chatSession)
		responses = []MessageResponse{*dropResponse}
	}

	return responses
}

func (s *Service) dropProductFromCart(chatID int64, productID string, chatSession *domain.Session) *MessageResponse {
	if product, ok := domain.ProductsMap[productID]; ok {
		chatSession.DropProductFromCart(productID)
		s.repo.Save(chatSession)

		text := fmt.Sprintf("Товар %s удален из корзины", product.Name)
		payload := UpdatePayload{
			Text: text,
		}
		response := s.sendOrderDetails(chatID, chatSession, domain.ButtonPrefixDrop)
		response.UpdatePayload = &payload
		response.Text = text
		return response
	}

	return nil
}

func (s *Service) sendOrderToAdmin() []MessageResponse {
	return nil
}

func (s *Service) sendOrder(chatID int64, chatSession *domain.Session) []MessageResponse {
	orderBag := chatSession.GetOrderBag()
	orderBagPrice := chatSession.GetOrderBagPrice()
	text := "Ваш заказ:\n"
	for _, order := range orderBag {
		text += fmt.Sprintf("%s - %d шт. (%.2f руб.)\n", order.Product, order.Count, order.Price)
	}
	text += fmt.Sprintf("Сумма заказа: %.2f руб.", orderBagPrice)

	return []MessageResponse{
		{
			ChatID: chatID,
			Text:   text,
			InlineReplyMarkup: &tgBotAPI.InlineKeyboardMarkup{
				InlineKeyboard: [][]tgBotAPI.InlineKeyboardButton{
					{
						tgBotAPI.NewInlineKeyboardButtonData("Подтвердить заказ", domain.ApproveOrder),
						tgBotAPI.NewInlineKeyboardButtonData("Отменить заказ", domain.CloseOrder),
					},
				},
			},
			UpdatePayload: &UpdatePayload{
				Text: text,
			},
		},
	}
}

func (s *Service) sendCategoryMenu(chatID int64, callbackData string) []MessageResponse {
	text := "Выберите категорию:"
	tools.LogInfo("Выбрана категория: %s", callbackData)
	payload := UpdatePayload{
		Text: text,
	}
	if callbackData == domain.CallbackCategoryMenu {
		return []MessageResponse{
			{
				ChatID:            chatID,
				Text:              text,
				InlineReplyMarkup: domain.GetCategoryKeyboard(),
				UpdatePayload:     &payload,
			},
		}
	} else if callbackData == domain.CallbackStartOrder {
		return []MessageResponse{
			{
				ChatID:            chatID,
				Text:              text,
				InlineReplyMarkup: domain.GetCategoryKeyboard(),
			},
		}
	}
	return []MessageResponse{}
}

func (s *Service) sendOrderDetails(chatID int64, chatSession *domain.Session, prefix string) *MessageResponse {
	keyboard := chatSession.GetOrderBagKeyboard(prefix)
	keyboard = tools.AppendDefaultButtons(&keyboard)
	price := chatSession.GetOrderBagPrice()
	text := "Ваш заказ пока пуст."
	if price > 0 {
		text = fmt.Sprintf("Стоимость заказа: %.2f руб.", price)
	}

	payload := UpdatePayload{
		Text: text,
	}
	return &MessageResponse{
		ChatID:            chatID,
		Text:              text,
		InlineReplyMarkup: &keyboard,
		UpdatePayload:     &payload,
	}
}

func (s *Service) chooseProducts(chatID int64, categoryID string) []MessageResponse {
	var categoryName string
	for _, c := range domain.Categories {
		if c.ID == categoryID {
			categoryName = c.Name
			break
		}
	}
	text := fmt.Sprintf("Выберите из категории %s:", categoryName)
	payload := UpdatePayload{
		Text: text,
	}
	return []MessageResponse{
		{
			ChatID:            chatID,
			Text:              text,
			InlineReplyMarkup: domain.GetProductsKeyboard(categoryID),
			UpdatePayload:     &payload,
		},
	}
}

func (s *Service) addProductToCart(chatID int64, productID string, chatSession *domain.Session, lastMessageID int) []MessageResponse {
	if product, ok := domain.ProductsMap[productID]; ok {
		chatSession.AddToOrderBag(product.ID, product.Name, product.Price)
		s.repo.Save(chatSession)

		payload := UpdatePayload{
			MessageID: lastMessageID,
			Text:      fmt.Sprintf("Добавлено в корзину: %s (%.2f руб.)\n\nВыберите из категории %s:", product.Name, product.Price, product.Category),
		}
		responses := []MessageResponse{
			{
				ChatID:            chatID,
				Text:              fmt.Sprintf("Добавлено в корзину: %s (%.2f руб.)\n\nВыберите из категории %s:", product.Name, product.Price, product.Category),
				InlineReplyMarkup: domain.GetProductsKeyboard(product.Category),
				UpdatePayload:     &payload,
			},
		}
		return responses
	}
	return nil
}
