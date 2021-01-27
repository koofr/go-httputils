package httputils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

func ResponseJSON(w http.ResponseWriter, r *http.Request, code int, v interface{}) error {
	jsonBytes, err := json.Marshal(v)
	if err != nil {
		return fmt.Errorf("response json marshal error: %w", err)
	}

	w.Header().Set("Content-Length", strconv.Itoa(len(jsonBytes)))
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	w.WriteHeader(code)

	_, err = w.Write(jsonBytes)
	return err
}
