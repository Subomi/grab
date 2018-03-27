package grab

import "errors"

var (
	errExtractingFileType    = errors.New("Error occurred extracting file types")
	errExtractingContentSize = errors.New("Error occurred extracting content size")
)
