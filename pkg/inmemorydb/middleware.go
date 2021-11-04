package inmemorydb

import (
	"net/http"

	"github.com/tarantool/go-tarantool"
)

func NewTarantoolMiddleware(tarInst *tarantool.Connection) *TarantoolMiddleware {
	return &TarantoolMiddleware{
		tarInst: tarInst,
	}
}

type TarantoolMiddleware struct {
	tarInst *tarantool.Connection
}

func (cm *TarantoolMiddleware) Attach(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		ctx = PutTarantoolToContext(cm.tarInst, ctx)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}
