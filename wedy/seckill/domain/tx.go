package domain

type OrderTX struct {
	OrderId     int64    `json:"orderId"`
	ActivityId  int64    `json:"activityId"`
	ProductId   int64    `json:"productId"`
	UserId      int64    `json:"userId"`
	Price       float64  `json:"price"`
	Quantity    int64    `json:"quantity"`
	PromoteCode []string `json:"promoteCode"`
}
