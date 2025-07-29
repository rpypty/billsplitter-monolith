package payment_method

type PaymentMethod struct {
	ID        int    // ID - айди сущности
	Name      string // Name - название метода (альфа, сбп Тиньк, ибан и тд.)
	Recipient string // Recipient - реквизит, будет копироваться на фронте
}
