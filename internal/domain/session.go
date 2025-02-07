package domain

import tgBotAPI "github.com/go-telegram-bot-api/telegram-bot-api/v5"

type session interface {
	SetFlagValue(string, bool)
}
type Session struct {
	ChatID     int64
	State      string
	StateFlags struct {
		isMenuSended      bool
		isOrderApproved   bool
		isCustomValidated bool
	}
	OrderBag []*Order
}

type Order struct {
	ID      string
	Product string
	Count   int
	Price   float64
}

func NewSession(chatID int64) *Session {
	return &Session{ChatID: chatID}
}

func (s *Session) MarkMenuSended() {
	s.StateFlags.isMenuSended = true
}

func (s *Session) AddToOrderBag(ID string, Name string, Price float64) {
	s.OrderBag = append(s.OrderBag, &Order{
		ID:      ID,
		Product: Name,
		Count:   1,
		Price:   Price,
	})
}

func (s *Session) GetOrderBag() []*Order {
	return s.OrderBag
}

func (s *Session) ClearOrderBag() {
	s.OrderBag = nil
}

func (s *Session) GetOrderBagPrice() float64 {
	totalPrice := 0.0
	for _, order := range s.OrderBag {
		totalPrice += order.Price * float64(order.Count)
	}
	return totalPrice
}

func (s *Session) DropProductFromCart(productID string) {
	newBag := make([]*Order, 0, len(s.OrderBag))
	for _, order := range s.OrderBag {
		if order.ID != productID {
			newBag = append(newBag, order)
		}
	}
	s.OrderBag = newBag
}

func (s *Session) GetOrderBagKeyboard(prefix string) tgBotAPI.InlineKeyboardMarkup {
	var buttons [][]tgBotAPI.InlineKeyboardButton
	for _, order := range s.OrderBag {
		buttons = append(buttons, []tgBotAPI.InlineKeyboardButton{
			tgBotAPI.NewInlineKeyboardButtonData(order.Product, prefix+order.ID),
		})
	}
	return tgBotAPI.NewInlineKeyboardMarkup(buttons...)
}
