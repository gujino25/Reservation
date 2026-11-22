package httptransport

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func pathParam(r *http.Request, key string) string {
	if value := r.PathValue(key); value != "" {
		return value
	}
	return chi.URLParam(r, key)
}
