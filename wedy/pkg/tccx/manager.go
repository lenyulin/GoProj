package tccx

import (
	"GoProj/wedy/pkg/logger"
	"context"
	"github.com/redis/go-redis/v9"
)

const (
	TccTransactionWatchTopic = "seckill_tcc_transaction_watch"
	SubmitCancelRequest      = "SubmitCancelRequest"
	AddTransaction           = "AddTransaction"
	TransactionComplete      = "TransactionComplete"
	TransactionFailed        = "TransactionFailed"
)

type TransactionResult struct{}

type TccManager struct {
	l     logger.LoggerV1
	redis redis.Cmdable
}

func NewTccManager(l logger.LoggerV1, redis redis.Cmdable) *TccManager {
	return &TccManager{
		l:     l,
		redis: redis,
	}
}
func (m *TccManager) Process(ctx context.Context, tcc OrderTX) error {

}
