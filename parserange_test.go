package httputils_test

import (
	. "git.koofr.lan/go-httputils.git"
	"git.koofr.lan/go-ioutils.git"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ParseRange", func() {
	It("should not parse malformed range", func() {
		_, err := ParseRange("malformed", 200)
		Expect(err).To(HaveOccurred())
	})

	It("should only parse bytes", func() {
		_, err := ParseRange("items=0-5", 1000)
		Expect(err).To(HaveOccurred())
	})

	It("should not parse range if start is greater than end", func() {
		_, err := ParseRange("bytes=500-20", 200)
		Expect(err).To(HaveOccurred())
	})

	It("should parse range if start and end exist", func() {
		spans, err := ParseRange("bytes=0-499", 200)
		Expect(err).NotTo(HaveOccurred())
		Expect(spans).To(Equal([]ioutils.FileSpan{{0, 199}}))

		spans, err = ParseRange("bytes=0-499", 1000)
		Expect(err).NotTo(HaveOccurred())
		Expect(spans).To(Equal([]ioutils.FileSpan{{0, 499}}))

		spans, err = ParseRange("bytes=40-80", 1000)
		Expect(err).NotTo(HaveOccurred())
		Expect(spans).To(Equal([]ioutils.FileSpan{{40, 80}}))
	})

	It("should parse range if start does not exist", func() {
		spans, err := ParseRange("bytes=-500", 1000)
		Expect(err).NotTo(HaveOccurred())
		Expect(spans).To(Equal([]ioutils.FileSpan{{500, 999}}))

		spans, err = ParseRange("bytes=-400", 1000)
		Expect(err).NotTo(HaveOccurred())
		Expect(spans).To(Equal([]ioutils.FileSpan{{600, 999}}))
	})

	It("should parse range if end does not exist", func() {
		spans, err := ParseRange("bytes=500-", 1000)
		Expect(err).NotTo(HaveOccurred())
		Expect(spans).To(Equal([]ioutils.FileSpan{{500, 999}}))

		spans, err = ParseRange("bytes=400-", 1000)
		Expect(err).NotTo(HaveOccurred())
		Expect(spans).To(Equal([]ioutils.FileSpan{{400, 999}}))
	})

	It("should parse range if both start and end equal 0", func() {
		spans, err := ParseRange("bytes=0-0", 1000)
		Expect(err).NotTo(HaveOccurred())
		Expect(spans).To(Equal([]ioutils.FileSpan{{0, 0}}))
	})

	It("should parse multiple ranges", func() {
		spans, err := ParseRange("bytes=40-80,-1", 1000)
		Expect(err).NotTo(HaveOccurred())
		Expect(spans).To(Equal([]ioutils.FileSpan{{40, 80}, {999, 999}}))
	})

	It("should not parse multiple ranges if any range is incorrect", func() {
		_, err := ParseRange("bytes=40-80,500-20", 200)
		Expect(err).To(HaveOccurred())
	})

})
