package tccx

import (
	"context"
	"time"
)

type TransactionResult struct{}

type TccManager struct {
	ttl    time.Duration
	result chan *TransactionResult
}

func NewTccManager(ttl time.Duration, queueSize int) *TccManager {
	return &TccManager{
		ttl,
		make(chan *TransactionResult, queueSize),
	}
}
func (m *TccManager) Process(ctx context.Context, tcc []TCC) (string, error) {
	//TODO implement me
	panic("implement me")
}
