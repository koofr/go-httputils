package httputils_test

import (
	. "git.koofr.lan/go-httputils.git"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"net/http"
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

		ForceDownload("čšžČŠŽ.txt", h)

		Expect(h).To(Equal(http.Header{
			"Content-Type":        {"application/force-download"},
			"Content-Disposition": {`attachment; filename="??????.txt"; filename*=UTF-8''%C4%8D%C5%A1%C5%BE%C4%8C%C5%A0%C5%BD.txt`},
		}))
	})
})
