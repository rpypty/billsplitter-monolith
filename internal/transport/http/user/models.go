package user

type UpdateUserProfileReq struct {
	Username  string `json:"username" example:"johndoe"`
	FirstName string `json:"firstName" example:"John"`
	LastName  string `json:"lastName" example:"Doe"`
}

// CreatePaymentMethodReq - запрос на создание платежного метода
type CreatePaymentMethodReq struct {
	Name        string `json:"name" validate:"required" example:"СБП"`
	Description string `json:"description" example:"Система быстрых платежей"`
	Recipient   string `json:"recipient" validate:"required" example:"+79001234567"`
}

// UpdatePaymentMethodReq - запрос на обновление платежного метода
type UpdatePaymentMethodReq struct {
	Name        string `json:"name" validate:"required" example:"СБП"`
	Description string `json:"description" example:"Система быстрых платежей"`
	Recipient   string `json:"recipient" validate:"required" example:"+79001234567"`
}

// PaymentMethodResponse - ответ с данными платежного метода
type PaymentMethodResponse struct {
	ID          int64  `json:"id" example:"1"`
	UserID      int64  `json:"userId" example:"123"`
	Name        string `json:"name" example:"СБП"`
	Description string `json:"description" example:"Система быстрых платежей"`
	Recipient   string `json:"recipient" example:"+79001234567"`
}

// PaymentMethodsResponse - ответ со списком платежных методов
type PaymentMethodsResponse struct {
	PaymentMethods []PaymentMethodResponse `json:"paymentMethods"`
}
