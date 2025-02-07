package domain

const (
	// Префиксы для callback данных кнопок
	ButtonPrefixProduct  = "add_product_"     // Префикс для продуктов (добавление в корзину)
	ButtonPrefixCategory = "choose_category_" // Префикс для категорий (выбор продуктов)
	ButtonPrefixDrop     = "drop_"            // Префикс для отправки коллбеков с деталями заказа
)

const (
	// Полные значения callback данных
	CallbackStartOrder   = "start_order"    // Начало оформления заказа
	CallbackCategoryMenu = "category_menu"  // Возврат в главное меню
	CallbackOrderDrop    = "drop_from_cart" // Просмотр корзины
	CallbackMakeOrder    = "make_order"     // Оформление заказа
	CallbackOrdersList   = "old_orders"     // Просмотр прошлых заказов
)

const (
	// Полные значения callback данных
	CloseOrder   = "close_order"
	ApproveOrder = "approve_order"
)

var (
	Categories = []Category{
		{
			ID:   "coffee",
			Name: "☕️ Кофе",
			Products: []Product{
				{ID: "espresso", Name: "Эспрессо", Price: 120},
				{ID: "cappuccino", Name: "Капучино", Price: 180},
				{ID: "latte", Name: "Латте", Price: 200},
				{ID: "americano", Name: "Американо", Price: 150},
				{ID: "raf", Name: "Раф", Price: 250},
				{ID: "flatwhite", Name: "Флэт Уайт", Price: 220},
			},
		},
		{
			ID:   "food",
			Name: "🍽 Еда",
			Products: []Product{
				{ID: "croissant", Name: "Круассан", Price: 150},
				{ID: "sandwich", Name: "Сэндвич", Price: 250},
				{ID: "salad", Name: "Салат", Price: 350},
				{ID: "cake", Name: "Пирожное", Price: 200},
				{ID: "muffin", Name: "Маффин", Price: 180},
				{ID: "cookie", Name: "Печенье", Price: 100},
			},
		},
	}

	// Карта для быстрого поиска продукта по ID
	ProductsMap = make(map[string]Product)
)

func init() {
	// Заполняем карту продуктов
	for _, category := range Categories {
		for _, product := range category.Products {
			product.Category = category.ID
			ProductsMap[product.ID] = product
		}
	}
}
