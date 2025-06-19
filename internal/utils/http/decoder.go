package http

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
	gojson "github.com/goccy/go-json"
)

var (
	validate *validator.Validate
)

func init() {
	validate = validator.New()
}

func DecodeReq[T any](r *http.Request) (*T, error) {
	var v T

	decoder := gojson.NewDecoder(r.Body)
	if err := decoder.Decode(&v); err != nil {
		return nil, err
	}

	err := validate.Struct(&v)
	if err != nil {
		return nil, err
	}

	return &v, nil
}

func MustMarshal(v any) []byte {
	out, _ := json.Marshal(v)
	return out
}
