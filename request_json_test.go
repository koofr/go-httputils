package httputils_test

import (
	"bytes"
	"errors"
	"io"
	"net/http/httptest"
	"strings"

	"github.com/koofr/go-ioutils"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/koofr/go-httputils"
)

var _ = Describe("RequestJSON", func() {
	Describe("RequestJSONError", func() {
		It("should close the request body", func() {
			closed := false

			r := httptest.NewRequest("GET", "/", ioutils.NewPassCloseReader(bytes.NewReader(nil), func() error {
				closed = true
				return nil
			}))

			_, _ = RequestJSONError(r, nil, 1)

			Expect(closed).To(BeTrue())
		})

		It("should handle empty Content-Type", func() {
			r := httptest.NewRequest("GET", "/", nil)
			r.Header.Set("Content-Type", "")
			_, err := RequestJSONError(r, nil, 1)
			Expect(err).To(HaveOccurred())
			var errInvalidContentType *ErrInvalidContentType
			Expect(errors.As(err, &errInvalidContentType)).To(BeTrue())
			Expect(err.Error()).To(Equal("request json error: invalid content type: mime: no media type"))
		})

		It("should handle invalid Content-Type", func() {
			r := httptest.NewRequest("GET", "/", nil)
			r.Header.Set("Content-Type", "foo;;")
			_, err := RequestJSONError(r, nil, 1)
			Expect(err).To(HaveOccurred())
			var errInvalidContentType *ErrInvalidContentType
			Expect(errors.As(err, &errInvalidContentType)).To(BeTrue())
			Expect(err.Error()).To(Equal("request json error: invalid content type: mime: invalid media parameter"))
		})

		It("should handle incorrect Content-Type", func() {
			r := httptest.NewRequest("GET", "/", nil)
			r.Header.Set("Content-Type", "text/plain")
			_, err := RequestJSONError(r, nil, 1)
			Expect(err).To(HaveOccurred())
			var errInvalidContentType *ErrInvalidContentType
			Expect(errors.As(err, &errInvalidContentType)).To(BeTrue())
			Expect(err.Error()).To(Equal("request json error: invalid content type: expected Content-Type to be application/json but got: text/plain"))
		})

		It("should handle broken body", func() {
			r := httptest.NewRequest("GET", "/", io.MultiReader(bytes.NewReader(make([]byte, 1*1024*1024)), ioutils.NewErrorReader(io.ErrUnexpectedEOF)))
			r.Header.Set("Content-Type", "application/json")
			_, err := RequestJSONError(r, nil, 2*1024*1024)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("request json error: unexpected EOF"))
		})

		It("should handle too large body", func() {
			r := httptest.NewRequest("GET", "/", bytes.NewReader(make([]byte, 1*1024*1024)))
			r.Header.Set("Content-Type", "application/json")
			_, err := RequestJSONError(r, nil, 1*1024)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("request json error: request body too large"))
		})

		It("should handle invalid json", func() {
			r := httptest.NewRequest("GET", "/", strings.NewReader("invalid"))
			r.Header.Set("Content-Type", "application/json")
			_, err := RequestJSONError(r, nil, 1*1024)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("request json error: invalid JSON: invalid character 'i' looking for beginning of value"))
		})

		It("should parse a valid json", func() {
			v := map[string]interface{}{}
			r := httptest.NewRequest("GET", "/", strings.NewReader(`{"foo": "bar"}`))
			r.Header.Set("Content-Type", "application/json")
			_, err := RequestJSONError(r, &v, 1*1024)
			Expect(err).NotTo(HaveOccurred())
			Expect(v).To(Equal(map[string]interface{}{
				"foo": "bar",
			}))
		})
	})
})
