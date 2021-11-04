package database

import (
	"context"
	"errors"
)

const (
	missingDbInstanceError = "missing db instance"
	DB_INSTANCE_CTX_NAME   = "DbInstance"
)

func GetDatabaseFromContext(ctx context.Context) (Db, error) {
	dbEntry, ok := ctx.Value(DB_INSTANCE_CTX_NAME).(Db)
	if !ok {
		return nil, errors.New(missingDbInstanceError)
	}
	return dbEntry, nil
}

func PutDatabaseToContext(entry Db, ctx context.Context) context.Context {
	return context.WithValue(ctx, DB_INSTANCE_CTX_NAME, entry)
}
