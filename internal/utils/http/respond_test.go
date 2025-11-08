package http

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
	_ "unsafe"
)

//go:linkname cfgDebug billsplitter-monolith/internal/cfg.debug
var cfgDebug bool

func TestRespondErrWithStatus(t *testing.T) {
	t.Run("delegates to formatted version", func(t *testing.T) {
		rec := httptest.NewRecorder()

		RespondErrWithStatus(rec, http.StatusBadRequest, "bad request")

		require.Equal(t, http.StatusBadRequest, rec.Code)
		body := parseErrorResponse(t, rec.Body.Bytes())
		require.Equal(t, "bad request", body.ErrorMessage)
	})
}

func TestRespondErrWithStatusf(t *testing.T) {
	t.Run("hides details for 500 when not in debug", func(t *testing.T) {
		setDebug(t, false)
		rec := httptest.NewRecorder()

		RespondErrWithStatusf(rec, http.StatusInternalServerError, "detailed %s", "msg")

		require.Equal(t, http.StatusInternalServerError, rec.Code)
		body := parseErrorResponse(t, rec.Body.Bytes())
		require.Contains(t, body.ErrorMessage, "Internal Server Error")
	})

	t.Run("keeps details for 500 when debug", func(t *testing.T) {
		setDebug(t, true)
		rec := httptest.NewRecorder()

		RespondErrWithStatusf(rec, http.StatusInternalServerError, "detailed %s", "msg")

		body := parseErrorResponse(t, rec.Body.Bytes())
		require.Equal(t, "detailed msg", body.ErrorMessage)
	})

	t.Run("uses provided message for other statuses", func(t *testing.T) {
		setDebug(t, false)
		rec := httptest.NewRecorder()

		RespondErrWithStatusf(rec, http.StatusBadRequest, "bad %s", "input")

		body := parseErrorResponse(t, rec.Body.Bytes())
		require.Equal(t, "bad input", body.ErrorMessage)
	})
}

func TestRespondJson(t *testing.T) {
	rec := httptest.NewRecorder()
	payload := map[string]string{"status": "ok"}

	RespondJson(rec, payload)

	require.Equal(t, http.StatusOK, rec.Code)
	require.Equal(t, "application/json", rec.Header().Get("Content-Type"))
	var decoded map[string]string
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &decoded))
	require.Equal(t, payload, decoded)
}

func TestRespondJsonWithStatus(t *testing.T) {
	rec := httptest.NewRecorder()
	payload := map[string]int{"status": 201}

	RespondJsonWithStatus(rec, http.StatusCreated, payload)

	require.Equal(t, http.StatusCreated, rec.Code)
	require.Equal(t, "application/json", rec.Header().Get("Content-Type"))
	var decoded map[string]int
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &decoded))
	require.Equal(t, payload, decoded)
}

func TestRespondOK(t *testing.T) {
	rec := httptest.NewRecorder()

	RespondOK(rec)

	require.Equal(t, http.StatusOK, rec.Code)
	var decoded ResponseOK
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &decoded))
	require.Equal(t, "OK", decoded.Message)
}

func parseErrorResponse(t *testing.T, body []byte) ErrorResponse {
	t.Helper()
	var resp ErrorResponse
	require.NoError(t, json.Unmarshal(body, &resp))
	return resp
}

func setDebug(t *testing.T, value bool) {
	t.Helper()
	prev := cfgDebug
	cfgDebug = value
	t.Cleanup(func() {
		cfgDebug = prev
	})
}
