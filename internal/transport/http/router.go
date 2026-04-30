package httptransport

import (
	"log/slog"
	"net/http"

	api "github.com/EgorTarasov/shopx/internal/generated/api"
)

func Router(h *Handler, logger *slog.Logger) http.Handler {
	mux := http.NewServeMux()
	writer := errorWriter(logger)
	strict := api.NewStrictHandlerWithOptions(h, nil, api.StrictHTTPServerOptions{
		RequestErrorHandlerFunc:  writer,
		ResponseErrorHandlerFunc: writer,
	})
	api.HandlerFromMux(strict, mux)
	return CORS(ExtractToken(mux))
}
