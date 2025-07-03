package domain

type Order struct {
	UserId        int64
	OrderId       int64
	ActivityId    string
	ProductId     string
	PaymentMethod string
	PromoCode     string
}
