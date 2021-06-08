package httputils_test

import (
	"net/http"

	. "github.com/koofr/go-httputils"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Forcedownload", func() {
	It("should set headers to force download", func() {
		h := make(http.Header)

		ForceDownload(`fi le".txt`, h)

		Expect(h).To(Equal(http.Header{
			"Content-Type":        {"application/force-download"},
			"Content-Disposition": {`attachment; filename="fi le\".txt"; filename*=UTF-8''fi%20le%22.txt`},
		}))
	})

	It("should escape unicode characters", func() {
		h := make(http.Header)

		ForceDownload("čšž,ČŠŽ.txt", h)

		Expect(h).To(Equal(http.Header{
			"Content-Type":        {"application/force-download"},
			"Content-Disposition": {`attachment; filename="???,???.txt"; filename*=UTF-8''%C4%8D%C5%A1%C5%BE%2C%C4%8C%C5%A0%C5%BD.txt`},
		}))
	})

	It("should escape semicolons", func() {
		h := make(http.Header)

		ForceDownload("foo; bar; baz.txt", h)

		Expect(h).To(Equal(http.Header{
			"Content-Type":        {"application/force-download"},
			"Content-Disposition": {`attachment; filename="foo; bar; baz.txt"; filename*=UTF-8''foo%3B%20bar%3B%20baz.txt`},
		}))
	})
})
