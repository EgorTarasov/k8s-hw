package httptransport

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/EgorTarasov/shopx/internal/domain"
	api "github.com/EgorTarasov/shopx/internal/generated/api"
)

type errorMapping struct {
	status  int
	code    string
	message string
}

func mapDomainError(err error) (errorMapping, bool) {
	switch {
	case errors.Is(err, domain.ErrEmailTaken):
		return errorMapping{http.StatusConflict, "email_taken", "email already taken"}, true
	case errors.Is(err, domain.ErrBadCredentials):
		return errorMapping{http.StatusUnauthorized, "bad_credentials", "invalid email or password"}, true
	case errors.Is(err, domain.ErrSessionNotFound):
		return errorMapping{http.StatusUnauthorized, "session_not_found", "session is missing or expired"}, true
	case errors.Is(err, domain.ErrOrderNotFound):
		return errorMapping{http.StatusNotFound, "order_not_found", "order not found"}, true
	case errors.Is(err, domain.ErrUserNotFound):
		return errorMapping{http.StatusNotFound, "user_not_found", "user not found"}, true
	case errors.Is(err, domain.ErrInvalidInput):
		return errorMapping{http.StatusBadRequest, "invalid_input", err.Error()}, true
	case errors.Is(err, domain.ErrInvalidStatus):
		return errorMapping{http.StatusBadRequest, "invalid_status", err.Error()}, true
	case errors.Is(err, domain.ErrForbidden):
		return errorMapping{http.StatusForbidden, "forbidden", "forbidden"}, true
	}
	return errorMapping{}, false
}

func errorWriter(logger *slog.Logger) func(w http.ResponseWriter, r *http.Request, err error) {
	return func(w http.ResponseWriter, r *http.Request, err error) {
		mapping, known := mapDomainError(err)
		if !known {
			logger.Error("http handler internal error",
				"method", r.Method,
				"path", r.URL.Path,
				"err", err,
			)
			mapping = errorMapping{http.StatusInternalServerError, "internal_error", "internal server error"}
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(mapping.status)
		_ = json.NewEncoder(w).Encode(api.Error{Code: mapping.code, Message: mapping.message})
	}
}
