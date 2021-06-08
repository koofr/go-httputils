package httputils

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

func ForceDownload(filename string, h http.Header) {
	escapedName := strings.Replace(filename, `"`, `\"`, -1)
	asciiName := asciiEscape(escapedName)

	u := &url.URL{
		Path: filename,
	}

	encodedName := strings.Replace(u.String(), ",", "%2C", -1)
	encodedName = strings.Replace(encodedName, ";", "%3B", -1)

	d := fmt.Sprintf(`attachment; filename="%s"; filename*=UTF-8''%s`, asciiName, encodedName)

	h.Set("Content-Type", "application/force-download")
	h.Set("Content-Disposition", d)
}

func asciiEscape(s string) string {
	ascii := make([]byte, 0, len(s))

	for _, runeValue := range s {
		if runeValue >= 0x20 && runeValue <= 0x7e {
			ascii = append(ascii, byte(runeValue))
		} else {
			ascii = append(ascii, '?')
		}
	}

	return string(ascii)
}
