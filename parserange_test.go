package httputils_test

import (
	. "github.com/koofr/go-httputils"
	"github.com/koofr/go-ioutils"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ParseRange", func() {
	It("should not parse malformed range", func() {
		_, _, err := ParseRange("malformed", 200)
		Expect(err).To(HaveOccurred())
	})

	It("should only parse bytes", func() {
		_, _, err := ParseRange("items=0-5", 1000)
		Expect(err).To(HaveOccurred())
	})

	It("should not parse range if start is greater than end", func() {
		_, _, err := ParseRange("bytes=500-20", 200)
		Expect(err).To(HaveOccurred())

		_, _, err = ParseRange("bytes=199-", 200)
		Expect(err).NotTo(HaveOccurred())
		_, _, err = ParseRange("bytes=200-", 200)
		Expect(err).To(HaveOccurred())
	})

	It("should parse range if start and end exist", func() {
		spans, hasEnd, err := ParseRange("bytes=0-499", 200)
		Expect(err).NotTo(HaveOccurred())
		Expect(spans).To(Equal([]ioutils.FileSpan{{0, 199}}))
		Expect(hasEnd).To(BeTrue())

		spans, hasEnd, err = ParseRange("bytes=0-499", 1000)
		Expect(err).NotTo(HaveOccurred())
		Expect(spans).To(Equal([]ioutils.FileSpan{{0, 499}}))
		Expect(hasEnd).To(BeTrue())

		spans, hasEnd, err = ParseRange("bytes=40-80", 1000)
		Expect(err).NotTo(HaveOccurred())
		Expect(spans).To(Equal([]ioutils.FileSpan{{40, 80}}))
		Expect(hasEnd).To(BeTrue())
	})

	It("should parse range if start does not exist", func() {
		spans, hasEnd, err := ParseRange("bytes=-500", 1000)
		Expect(err).NotTo(HaveOccurred())
		Expect(spans).To(Equal([]ioutils.FileSpan{{500, 999}}))
		Expect(hasEnd).To(BeTrue())

		spans, hasEnd, err = ParseRange("bytes=-400", 1000)
		Expect(err).NotTo(HaveOccurred())
		Expect(spans).To(Equal([]ioutils.FileSpan{{600, 999}}))
		Expect(hasEnd).To(BeTrue())
	})

	It("should parse range if end does not exist", func() {
		spans, hasEnd, err := ParseRange("bytes=500-", 1000)
		Expect(err).NotTo(HaveOccurred())
		Expect(spans).To(Equal([]ioutils.FileSpan{{500, 999}}))
		Expect(hasEnd).To(BeFalse())

		spans, hasEnd, err = ParseRange("bytes=400-", 1000)
		Expect(err).NotTo(HaveOccurred())
		Expect(spans).To(Equal([]ioutils.FileSpan{{400, 999}}))
		Expect(hasEnd).To(BeFalse())
	})

	It("should parse range if both start and end equal 0", func() {
		spans, hasEnd, err := ParseRange("bytes=0-0", 1000)
		Expect(err).NotTo(HaveOccurred())
		Expect(spans).To(Equal([]ioutils.FileSpan{{0, 0}}))
		Expect(hasEnd).To(BeTrue())
	})

	It("should parse multiple ranges", func() {
		spans, _, err := ParseRange("bytes=40-80,-1", 1000)
		Expect(err).NotTo(HaveOccurred())
		Expect(spans).To(Equal([]ioutils.FileSpan{{40, 80}, {999, 999}}))
	})

	It("should not parse multiple ranges if any range is incorrect", func() {
		_, _, err := ParseRange("bytes=40-80,500-20", 200)
		Expect(err).To(HaveOccurred())
	})

})
