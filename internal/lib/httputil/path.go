package httputil

import (
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
)

// PathInt64 returns the value of the path parameter as an int64.
func PathInt64(w http.ResponseWriter, r *http.Request, log *slog.Logger, name string) (int64, bool) {
	value := chi.URLParam(r, name)
	if value == "" {
		WriteRequestError(w, r, log, "invalid path parameter", fmt.Errorf("%s path parameter is required", name), http.StatusBadRequest)
		return 0, false
	}

	id, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		WriteRequestError(w, r, log, "invalid path parameter", err, http.StatusBadRequest)
		return 0, false
	}

	return id, true
}
