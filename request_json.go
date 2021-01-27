package httputils

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"mime"
	"net/http"

	"github.com/koofr/go-ioutils"
)

type ErrRequestJSON struct {
	Err error
}

func NewErrRequestJSON(err error) *ErrRequestJSON {
	return &ErrRequestJSON{Err: err}
}

func (e *ErrRequestJSON) Error() string {
	return "request json error: " + e.Err.Error()
}

func (e *ErrRequestJSON) Unwrap() error {
	return e.Err
}

type ErrInvalidContentType struct {
	Err error
}

func NewErrInvalidContentType(err error) *ErrInvalidContentType {
	return &ErrInvalidContentType{Err: err}
}

func (e *ErrInvalidContentType) Error() string {
	return "invalid content type: " + e.Err.Error()
}

func (e *ErrInvalidContentType) Unwrap() error {
	return e.Err
}

var ErrRequestBodyTooLarge = errors.New("request body too large")

type ErrInvalidJSON struct {
	Err   error
	Bytes []byte
}

func NewErrInvalidJSON(err error, bytes []byte) *ErrInvalidJSON {
	return &ErrInvalidJSON{
		Err:   err,
		Bytes: bytes,
	}
}

func (e *ErrInvalidJSON) Error() string {
	return "invalid JSON: " + e.Err.Error()
}

func (e *ErrInvalidJSON) Unwrap() error {
	return e.Err
}

func RequestJSONError(r *http.Request, v interface{}, maxRequestJSONSize int) (jsonBytes []byte, err error) {
	defer r.Body.Close()

	contentType := r.Header.Get("Content-Type")
	mediaType, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		return nil, NewErrRequestJSON(NewErrInvalidContentType(err))
	}
	if mediaType != "application/json" {
		return nil, NewErrRequestJSON(NewErrInvalidContentType(fmt.Errorf("expected Content-Type to be application/json but got: %s", mediaType)))
	}

	reader := ioutils.NewSizeLimitedReader(r.Body, int64(maxRequestJSONSize))

	jsonBytes, err = ioutil.ReadAll(reader)
	if err != nil {
		if errors.Is(err, ioutils.ErrMaxSizeExceeded) {
			return nil, NewErrRequestJSON(ErrRequestBodyTooLarge)
		}
		return nil, NewErrRequestJSON(err)
	}

	err = json.Unmarshal(jsonBytes, v)
	if err != nil {
		return nil, NewErrRequestJSON(NewErrInvalidJSON(err, jsonBytes))
	}

	return jsonBytes, nil
}
