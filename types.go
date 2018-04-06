package grab

import (
	"sync"
)

type DownloadReq struct {
	conf            *Config
	wg              *sync.WaitGroup
	rangeDownload   bool
	bytesPerRoutine int64
	URL             string
	file            *fileObject
	rr              []*rangeRequest
	Downloads       []string
}

type fileObject struct {
	filename string
	fileExt  string
	size     int64
}

type rangeRequest struct {
	filename   string
	dr         *DownloadReq
	bytesRange [2]string
}
