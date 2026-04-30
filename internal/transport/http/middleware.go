package httptransport

import (
	"context"
	"net/http"

	"github.com/EgorTarasov/shopx/internal/domain"
)

const SessionHeader = "X-Session-Token"

type ctxKey string

const tokenCtxKey ctxKey = "session-token"

func WithToken(ctx context.Context, token domain.SessionToken) context.Context {
	return context.WithValue(ctx, tokenCtxKey, token)
}

func TokenFromContext(ctx context.Context) (domain.SessionToken, bool) {
	t, ok := ctx.Value(tokenCtxKey).(domain.SessionToken)
	return t, ok && t != ""
}

func ExtractToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if raw := r.Header.Get(SessionHeader); raw != "" {
			r = r.WithContext(WithToken(r.Context(), domain.SessionToken(raw)))
		}
		next.ServeHTTP(w, r)
	})
}

func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, "+SessionHeader)
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}
