package domain

type OrderTX struct {
	OrderId     int64    `json:"orderId"`
	UserId      int64    `json:"userId"`
	Mount       int64    `json:"mount"`
	Quantity    int64    `json:"quantity"`
	PromoteCode []string `json:"promoteCode"`
}
