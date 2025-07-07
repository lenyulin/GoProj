package domain

type Order struct {
	UserId        int64
	OrderId       int64
	ActivityId    int64
	ProductId     int64
	Quantity      int64
	Price         float64
	PaymentMethod string
	PromoCode     []string
}
