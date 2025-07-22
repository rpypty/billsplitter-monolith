package http

import (
	"fmt"
	"net/http"

	"billsplitter-monolith/internal/cfg"
)

func RespondErrWithStatus(w http.ResponseWriter, status int, msg string) {
	RespondErrWithStatusf(w, status, msg)
}

func RespondErrWithStatusf(w http.ResponseWriter, status int, msg string, args ...interface{}) {
	// FIXME: bad
	if status == http.StatusInternalServerError && !cfg.IsDebug() {
		// в проде не показываем детали ошибки
		msg = "Internal Server Error"
	}
	res := ErrorResponse{
		ErrorMessage: fmt.Sprintf(msg, args...),
	}
	RespondJsonWithStatus(w, status, res)
}

func RespondJson(w http.ResponseWriter, data any) {
	RespondJsonWithStatus(w, http.StatusOK, data)
}

func RespondJsonWithStatus(w http.ResponseWriter, status int, data any) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	bytes := MustMarshal(data)
	_, _ = w.Write(bytes)
}
