package http

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func GetQueryParamInt(r *http.Request, key string) (int64, error) {
	idStr := chi.URLParam(r, key)

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return 0, err
	}

	return id, nil
}
