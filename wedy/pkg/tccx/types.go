package tccx

type Consumer interface {
	Start() error
}

type TccTransactionEvent struct {
	TxId      int64
	partition []string
	offset    []string
	timeStamp []int64
	topic     string
	retry     int64
	Order     OrderTX
}
type OrderTX struct {
	OrderId     int64    `json:"orderId"`
	UserId      int64    `json:"userId"`
	Mount       int64    `json:"mount"`
	Quantity    int64    `json:"quantity"`
	PromoteCode []string `json:"promoteCode"`
}
