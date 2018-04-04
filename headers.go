package grab

import (
	"net/http"
	"strings"
)

/*
	Parsing headers logic stays here.

	TODO(subomi): Remove hard-coded content sizes and be more
	aggressive in getting header fields.
*/

// Retrieve from a map or from a string array
func extractFileType(header interface{}) (string, error) {
	switch h := header.(type) {
	case []string:
		ext := strings.Join(h, "")
		extArray := strings.Split(ext, "/")
		return extArray[1], nil

	case http.Header:
		return extractFileType([]string{h.Get("Content-Type")})

	default:
		return "", errExtractingFileType
	}
}

func extractContentSize(header interface{}) (string, error) {
	switch h := header.(type) {
	case []string:
		return h[0], nil

	case http.Header:
		return extractContentSize([]string{h.Get("Content-Length")})

	default:
		return "", errExtractingContentSize
	}
}
