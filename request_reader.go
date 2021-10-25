package httputils

import (
	"fmt"
	"io"
	"mime"
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

	// In Go 1.17 p.FileName() calls filepath.Base(filename) and we want the
	// original value with path
	_, dispositionParams, _ := mime.ParseMediaType(p.Header.Get("Content-Disposition"))
	filename := dispositionParams["filename"]

	if filename == "" {
		return nil, "", fmt.Errorf("MultipartRequestReader part is not a file")
	}

	return p, filename, nil
}
