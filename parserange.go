package httputils

import (
	"errors"
	"net/textproto"
	"strconv"
	"strings"

	"github.com/koofr/go-ioutils"
)

var ErrInvalidRange = errors.New("invalid range")

func ParseRange(s string, size int64) (spans []ioutils.FileSpan, hasEnd bool, err error) {
	if s == "" {
		return nil, false, nil // header not present
	}

	const b = "bytes="
	if !strings.HasPrefix(s, b) {
		return nil, false, ErrInvalidRange
	}

	hasEnd = true

	for _, ra := range strings.Split(s[len(b):], ",") {
		ra = textproto.TrimString(ra)
		if ra == "" {
			continue
		}
		i := strings.Index(ra, "-")
		if i < 0 {
			return nil, false, ErrInvalidRange
		}
		start, end := textproto.TrimString(ra[:i]), textproto.TrimString(ra[i+1:])
		var s ioutils.FileSpan
		if start == "" {
			// If no start is specified, end specifies the
			// range start relative to the end of the file,
			// and we are dealing with <suffix-length>
			// which has to be a non-negative integer as per
			// RFC 7233 Section 2.1 "Byte-Ranges".
			if end == "" || end[0] == '-' {
				return nil, false, ErrInvalidRange
			}
			i, err := strconv.ParseInt(end, 10, 64)
			if i < 0 || err != nil {
				return nil, false, ErrInvalidRange
			}
			if i > size {
				i = size
			}
			s.Start = size - i
			s.End = size - 1
		} else {
			i, err := strconv.ParseInt(start, 10, 64)
			if err != nil || i >= size || i < 0 {
				return nil, false, ErrInvalidRange
			}
			s.Start = i
			if end == "" {
				// If no end is specified, range extends to end of the file.
				s.End = size - 1
				hasEnd = false
			} else {
				i, err := strconv.ParseInt(end, 10, 64)
				if err != nil || s.Start > i {
					return nil, false, ErrInvalidRange
				}
				if i >= size {
					i = size - 1
				}
				s.End = i
			}
		}
		spans = append(spans, s)
	}

	return spans, hasEnd, nil
}
