package grab

import "strings"

/*
	Parsing headers logic stays here.
*/

func extractFileType(str []string) string {
	ext := strings.Join(str, "")
	extArray := strings.Split(ext, "/")

	return extArray[1]
}
