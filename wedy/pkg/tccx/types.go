package tccx

import "context"

type TCC interface {
	Try(ctx context.Context) error
	Confirm(ctx context.Context) error
	Cancel(ctx context.Context) error
}

const (
	AddTransaction = 1 + iota
	TransactionInProgress
	TransactionCompleted
	TransactionFailed
)
