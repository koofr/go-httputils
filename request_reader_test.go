package httputils_test

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	. "github.com/koofr/go-httputils"
	"github.com/koofr/go-ioutils"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("RequestReader", func() {
	Describe("MultipartRequestReader", func() {
		It("should read multipart file", func() {
			body := `--------------------------c8898eaa2e25254d
Content-Disposition: form-data; name="file"; filename="foo"
Content-Type: application/octet-stream

bar
--------------------------c8898eaa2e25254d--`

			req, err := http.NewRequest("POST", "/", bytes.NewReader([]byte(body)))
			Expect(err).NotTo(HaveOccurred())

			req.Header.Set("Content-Type", "multipart/form-data; boundary=------------------------c8898eaa2e25254d")

			r, name, err := MultipartRequestReader(req)
			Expect(err).NotTo(HaveOccurred())
			Expect(name).To(Equal("foo"))

			data, err := ioutil.ReadAll(r)
			Expect(err).NotTo(HaveOccurred())

			Expect(data).To(Equal([]byte("bar")))
		})

		It("should read multipart file with a full path", func() {
			body := `--------------------------c8898eaa2e25254d
Content-Disposition: form-data; name="file"; filename="path/to/foo"
Content-Type: application/octet-stream

bar
--------------------------c8898eaa2e25254d--`

			req, err := http.NewRequest("POST", "/", bytes.NewReader([]byte(body)))
			Expect(err).NotTo(HaveOccurred())

			req.Header.Set("Content-Type", "multipart/form-data; boundary=------------------------c8898eaa2e25254d")

			r, name, err := MultipartRequestReader(req)
			Expect(err).NotTo(HaveOccurred())
			Expect(name).To(Equal("path/to/foo"))

			data, err := ioutil.ReadAll(r)
			Expect(err).NotTo(HaveOccurred())

			Expect(data).To(Equal([]byte("bar")))
		})

		It("should not read multipart file if content-type is invalid", func() {
			body := `--------------------------c8898eaa2e25254d
Content-Disposition: form-data; name="file"; filename=""
Content-Type: application/octet-stream

bar
--------------------------c8898eaa2e25254d--`

			req, err := http.NewRequest("POST", "/", bytes.NewReader([]byte(body)))
			Expect(err).NotTo(HaveOccurred())

			req.Header.Set("Content-Type", "invalid")

			_, _, err = MultipartRequestReader(req)
			Expect(err.Error()).To(Equal("MultipartRequestReader multipart reader error: request Content-Type isn't multipart/form-data"))
		})

		It("should not read multipart file if boundary is missing", func() {
			body := `--------------------------c8898eaa2e25254d
Content-Disposition: form-data; name="file"; filename=""
Content-Type: application/octet-stream

bar
--------------------------c8898eaa2e25254d--`

			req, err := http.NewRequest("POST", "/", bytes.NewReader([]byte(body)))
			Expect(err).NotTo(HaveOccurred())

			req.Header.Set("Content-Type", "multipart/form-data")

			_, _, err = MultipartRequestReader(req)
			Expect(err.Error()).To(Equal("MultipartRequestReader multipart reader error: no multipart boundary param in Content-Type"))
		})

		It("should not read multipart file if filename is missing", func() {
			body := `--------------------------c8898eaa2e25254d
Content-Disposition: form-data; name="file"; filename=""
Content-Type: application/octet-stream

bar
--------------------------c8898eaa2e25254d--`

			req, err := http.NewRequest("POST", "/", bytes.NewReader([]byte(body)))
			Expect(err).NotTo(HaveOccurred())

			req.Header.Set("Content-Type", "multipart/form-data; boundary=------------------------c8898eaa2e25254d")

			_, _, err = MultipartRequestReader(req)
			Expect(err.Error()).To(Equal("MultipartRequestReader part is not a file"))
		})

		It("should not read multipart file if file is not the first field", func() {
			body := `--------------------------c8898eaa2e25254d
Content-Disposition: form-data; name="another"

value
--------------------------c8898eaa2e25254d
Content-Disposition: form-data; name="file"; filename="foo"
Content-Type: application/octet-stream

bar
--------------------------c8898eaa2e25254d--`

			req, err := http.NewRequest("POST", "/", bytes.NewReader([]byte(body)))
			Expect(err).NotTo(HaveOccurred())

			req.Header.Set("Content-Type", "multipart/form-data; boundary=------------------------c8898eaa2e25254d")

			_, _, err = MultipartRequestReader(req)
			Expect(err.Error()).To(Equal("MultipartRequestReader part is not a file"))
		})

		It("should get an error if body does not end with boundary", func() {
			body := `--------------------------c8898eaa2e25254d
Content-Disposition: form-data; name="file"; filename="foo"
Content-Type: application/octet-stream

ba`

			req, err := http.NewRequest("POST", "/", bytes.NewReader([]byte(body)))
			Expect(err).NotTo(HaveOccurred())

			req.Header.Set("Content-Type", "multipart/form-data; boundary=------------------------c8898eaa2e25254d")

			r, _, err := MultipartRequestReader(req)
			Expect(err).NotTo(HaveOccurred())

			_, err = ioutil.ReadAll(r)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("unexpected EOF"))
		})

		It("should get an error if file ends with EOF", func() {
			body := `--------------------------c8898eaa2e25254d
Content-Disposition: form-data; name="file"; filename="foo"
Content-Type: application/octet-stream

ba`

			reader := io.MultiReader(bytes.NewReader([]byte(body)), ioutils.NewErrorReader(fmt.Errorf("myerr")))

			req, err := http.NewRequest("POST", "/", reader)
			Expect(err).NotTo(HaveOccurred())

			req.Header.Set("Content-Type", "multipart/form-data; boundary=------------------------c8898eaa2e25254d")

			r, _, err := MultipartRequestReader(req)
			Expect(err).NotTo(HaveOccurred())

			_, err = ioutil.ReadAll(r)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("myerr"))
		})

		It("should get an error if body ends with EOF", func() {
			body := `--------------------------c8898eaa2e25254d`

			reader := io.MultiReader(bytes.NewReader([]byte(body)), ioutils.NewErrorReader(io.EOF))

			req, err := http.NewRequest("POST", "/", reader)
			Expect(err).NotTo(HaveOccurred())

			req.Header.Set("Content-Type", "multipart/form-data; boundary=------------------------c8898eaa2e25254d")

			_, _, err = MultipartRequestReader(req)
			Expect(err.Error()).To(Equal("MultipartRequestReader NextPart error: multipart: NextPart: EOF"))
		})

		It("should get an error if body is empty", func() {
			reader := ioutils.NewErrorReader(io.EOF)

			req, err := http.NewRequest("POST", "/", reader)
			Expect(err).NotTo(HaveOccurred())

			req.Header.Set("Content-Type", "multipart/form-data; boundary=------------------------c8898eaa2e25254d")

			_, _, err = MultipartRequestReader(req)
			Expect(err.Error()).To(Equal("MultipartRequestReader NextPart error: multipart: NextPart: EOF"))
		})
	})

})
