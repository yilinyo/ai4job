package postgres

import (
	"context"

	"gorm.io/gorm"
)

type transactionContextKey struct{}

type TransactionRunner struct {
	db *gorm.DB
}

func NewTransactionRunner(db *gorm.DB) *TransactionRunner {
	return &TransactionRunner{db: db}
}

func (r *TransactionRunner) WithinTransaction(ctx context.Context, fn func(context.Context) error) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return fn(context.WithValue(ctx, transactionContextKey{}, tx))
	})
}

func dbFromContext(ctx context.Context, fallback *gorm.DB) *gorm.DB {
	if tx, ok := ctx.Value(transactionContextKey{}).(*gorm.DB); ok && tx != nil {
		return tx
	}
	return fallback
}
