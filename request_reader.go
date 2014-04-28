package httputils

import (
	"fmt"
	"io"
	"net/http"
)

func MultipartRequestReader(r *http.Request) (io.Reader, string, error) {
	reader, err := r.MultipartReader()

	if err != nil {
		err = fmt.Errorf("MultipartRequestReader multipart reader error: %s", err)
		return nil, "", err
	}

	p, err := reader.NextPart()

	if err != nil {
		err = fmt.Errorf("MultipartRequestReader NextPart error: %s", err)
		return nil, "", err
	}

	name := p.FormName()

	if name == "" {
		return nil, "", fmt.Errorf("MultipartRequestReader missing field name")
	}

	filename := p.FileName()

	if filename == "" {
		return nil, "", fmt.Errorf("MultipartRequestReader part is not a file")
	}

	return p, filename, nil
}
