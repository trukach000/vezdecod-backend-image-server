package inmemorydb

import (
	"context"
	"errors"

	"github.com/tarantool/go-tarantool"
)

const (
	missingTarantoolInstanceError = "missing tarantool instance"
	TARANTOOL_INSTANCE_CTX_NAME   = "TarantoolInstance"
)

func GetTarantoolFromContext(ctx context.Context) (*tarantool.Connection, error) {
	dbEntry, ok := ctx.Value(TARANTOOL_INSTANCE_CTX_NAME).(*tarantool.Connection)
	if !ok {
		return nil, errors.New(missingTarantoolInstanceError)
	}
	return dbEntry, nil
}

func PutTarantoolToContext(entry *tarantool.Connection, ctx context.Context) context.Context {
	return context.WithValue(ctx, TARANTOOL_INSTANCE_CTX_NAME, entry)
}
