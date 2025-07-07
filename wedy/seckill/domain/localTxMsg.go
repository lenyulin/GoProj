package domain

type LocalTxMessage struct {
	Status        int64
	RetryCount    int64
	CTime         int64
	UTime         int64
	NextRetryTime int64
	Tx            OrderTX
}
