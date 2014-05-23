package httputils

import (
	"github.com/koofr/go-httpclient"
	"io"
	"net/http"
)

func UploadFile(url string, reader io.Reader, name string, expectedStatus int, respValue interface{}) (res *http.Response, err error) {
	respEncoding := httpclient.Encoding("")

	if respValue != nil {
		respEncoding = httpclient.EncodingJSON
	}

	req := &httpclient.RequestData{
		Method:          "POST",
		FullURL:         url,
		ExpectedStatus:  []int{expectedStatus},
		RespEncoding:    respEncoding,
		RespValue:       respValue,
		IgnoreRedirects: true,
	}

	err = req.UploadFile("file", name, reader)

	if err != nil {
		return
	}

	res, err = httpclient.New().Request(req)

	return
}
