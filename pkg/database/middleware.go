package database

import (
	"net/http"
)

func NewDatabaseMiddleware(db Db) *DatabaseMiddleware {
	return &DatabaseMiddleware{
		db: db,
	}
}

type DatabaseMiddleware struct {
	db Db
}

func (cm *DatabaseMiddleware) Attach(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		ctx = PutDatabaseToContext(cm.db, ctx)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}
