package gorm

import (
	"context"

	"gorm.io/gorm"
)

type txCtxKey struct{}

// ContextWithTx binds a *gorm.DB (typically a transaction) to ctx so repository methods use it.
func ContextWithTx(ctx context.Context, db *gorm.DB) context.Context {
	if db == nil {
		return ctx
	}
	return context.WithValue(ctx, txCtxKey{}, db)
}

func dbWithContext(ctx context.Context, defaultDB *gorm.DB) *gorm.DB {
	if v, ok := ctx.Value(txCtxKey{}).(*gorm.DB); ok && v != nil {
		return v.WithContext(ctx)
	}
	return defaultDB.WithContext(ctx)
}
