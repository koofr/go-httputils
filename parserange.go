package httputils

import (
	"errors"
	"git.koofr.lan/go-ioutils.git"
	"strconv"
	"strings"
)

func ParseRange(s string, size int64) ([]ioutils.FileSpan, error) {
	if s == "" {
		return nil, nil // header not present
	}

	const b = "bytes="
	if !strings.HasPrefix(s, b) {
		return nil, errors.New("invalid range")
	}
	var spans []ioutils.FileSpan
	for _, ra := range strings.Split(s[len(b):], ",") {
		ra = strings.TrimSpace(ra)
		if ra == "" {
			continue
		}
		i := strings.Index(ra, "-")
		if i < 0 {
			return nil, errors.New("invalid range")
		}
		start, end := strings.TrimSpace(ra[:i]), strings.TrimSpace(ra[i+1:])
		var s ioutils.FileSpan
		if start == "" {
			// If no start is specified, end specifies the
			// range start relative to the end of the file.
			i, err := strconv.ParseInt(end, 10, 64)
			if err != nil {
				return nil, errors.New("invalid range")
			}
			if i > size {
				i = size
			}
			s.Start = size - i
			s.End = size - 1
		} else {
			i, err := strconv.ParseInt(start, 10, 64)
			if err != nil || i > size || i < 0 {
				return nil, errors.New("invalid range")
			}
			s.Start = i
			if end == "" {
				// If no end is specified, range extends to end of the file.
				s.End = size - 1
			} else {
				i, err := strconv.ParseInt(end, 10, 64)
				if err != nil || s.Start > i {
					return nil, errors.New("invalid range")
				}
				if i >= size {
					i = size - 1
				}
				s.End = i
			}
		}
		spans = append(spans, s)
	}
	return spans, nil
}
