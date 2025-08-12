package payment_method

type PaymentMethod struct {
	ID          int64  `json:"id"`          // ID - айди сущности
	UserID      int64  `json:"userId"`      // UserID - ID пользователя, которому принадлежит метод
	Name        string `json:"name"`        // Name - название метода (альфа, сбп Тиньк, ибан и тд.)
	Description string `json:"description"` // Description - описание метода
	Recipient   string `json:"recipient"`   // Recipient - реквизит, будет копироваться на фронте
}
