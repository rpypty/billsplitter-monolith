package valueobject

// UserID - айди пользователя
type UserID int64

// CurrencyCode - код валюты, список в valueobject/enums.go
type CurrencyCode string

// Amount - сумма денег
type Amount int64

type SplitType string

const (
	SplitTypeShares      = "shares"
	SplitTypePercentages = "percentage"
	SplitTypeCustom      = "custom"
	SplitTypeEven        = "even"
)
