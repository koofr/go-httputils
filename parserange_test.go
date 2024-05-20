package httputils_test

import (
	"github.com/koofr/go-ioutils"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	. "github.com/koofr/go-httputils"
)

var _ = Describe("ParseRange", func() {
	It("should not parse malformed range", func() {
		_, _, err := ParseRange("malformed", 200)
		Expect(err).To(Equal(ErrInvalidRange))
	})

	It("should only parse bytes", func() {
		_, _, err := ParseRange("items=0-5", 1000)
		Expect(err).To(Equal(ErrInvalidRange))
	})

	It("should not parse range if start is greater than end", func() {
		_, _, err := ParseRange("bytes=500-20", 200)
		Expect(err).To(Equal(ErrInvalidRange))

		_, _, err = ParseRange("bytes=199-", 200)
		Expect(err).NotTo(HaveOccurred())
		_, _, err = ParseRange("bytes=200-", 200)
		Expect(err).To(Equal(ErrInvalidRange))
	})

	It("should parse range if start and end exist", func() {
		spans, hasEnd, err := ParseRange("bytes=0-499", 200)
		Expect(err).NotTo(HaveOccurred())
		Expect(spans).To(Equal([]ioutils.FileSpan{{Start: 0, End: 199}}))
		Expect(hasEnd).To(BeTrue())

		spans, hasEnd, err = ParseRange("bytes=0-499", 1000)
		Expect(err).NotTo(HaveOccurred())
		Expect(spans).To(Equal([]ioutils.FileSpan{{Start: 0, End: 499}}))
		Expect(hasEnd).To(BeTrue())

		spans, hasEnd, err = ParseRange("bytes=40-80", 1000)
		Expect(err).NotTo(HaveOccurred())
		Expect(spans).To(Equal([]ioutils.FileSpan{{Start: 40, End: 80}}))
		Expect(hasEnd).To(BeTrue())
	})

	It("should parse range if start does not exist", func() {
		spans, hasEnd, err := ParseRange("bytes=-500", 1000)
		Expect(err).NotTo(HaveOccurred())
		Expect(spans).To(Equal([]ioutils.FileSpan{{Start: 500, End: 999}}))
		Expect(hasEnd).To(BeTrue())

		spans, hasEnd, err = ParseRange("bytes=-400", 1000)
		Expect(err).NotTo(HaveOccurred())
		Expect(spans).To(Equal([]ioutils.FileSpan{{Start: 600, End: 999}}))
		Expect(hasEnd).To(BeTrue())
	})

	It("should parse range if end does not exist", func() {
		spans, hasEnd, err := ParseRange("bytes=500-", 1000)
		Expect(err).NotTo(HaveOccurred())
		Expect(spans).To(Equal([]ioutils.FileSpan{{Start: 500, End: 999}}))
		Expect(hasEnd).To(BeFalse())

		spans, hasEnd, err = ParseRange("bytes=400-", 1000)
		Expect(err).NotTo(HaveOccurred())
		Expect(spans).To(Equal([]ioutils.FileSpan{{Start: 400, End: 999}}))
		Expect(hasEnd).To(BeFalse())
	})

	It("should parse range if both start and end equal 0", func() {
		spans, hasEnd, err := ParseRange("bytes=0-0", 1000)
		Expect(err).NotTo(HaveOccurred())
		Expect(spans).To(Equal([]ioutils.FileSpan{{Start: 0, End: 0}}))
		Expect(hasEnd).To(BeTrue())
	})

	It("should parse multiple ranges", func() {
		spans, _, err := ParseRange("bytes=40-80,-1", 1000)
		Expect(err).NotTo(HaveOccurred())
		Expect(spans).To(Equal([]ioutils.FileSpan{{Start: 40, End: 80}, {Start: 999, End: 999}}))
	})

	It("should not parse multiple ranges if any range is incorrect", func() {
		_, _, err := ParseRange("bytes=40-80,500-20", 200)
		Expect(err).To(Equal(ErrInvalidRange))
	})

	DescribeTable(
		"Negative ranges",
		func(rng string, shouldSucceed bool, expectedStart int, expectedEnd int, expectedHasEnd bool) {
			spans, hasEnd, err := ParseRange(rng, 200)
			if shouldSucceed {
				Expect(err).NotTo(HaveOccurred())
				Expect(hasEnd).To(Equal(expectedHasEnd))
				if expectedStart == -1 {
					Expect(spans).To(BeEmpty())
				} else {
					Expect(spans).To(HaveLen(1))
					Expect(spans[0].Start).To(Equal(int64(expectedStart)))
					Expect(spans[0].End).To(Equal(int64(expectedEnd)))
				}
			} else {
				Expect(err).To(Equal(ErrInvalidRange))
			}
		},
		Entry("negative end", "bytes=--6", false, 0, 0, false),
		Entry("negative zero end", "bytes=--0", false, 0, 0, false),
		Entry("double negative zero end", "bytes=---0", false, 0, 0, false),
		Entry("end only", "bytes=-6", true, 194, 199, true),
		Entry("start only", "bytes=6-", true, 6, 199, false),
		Entry("invalid end", "bytes=-6-", false, 0, 0, false),
		Entry("zero end", "bytes=-0", true, 200, 199, true),
		Entry("empty range", "bytes=", true, -1, -1, true),
	)
})
