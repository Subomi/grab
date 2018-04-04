package grab

import (
	"io"
	"strconv"
	"strings"
	"sync"
)

func NewRangeRequest(br [2]string, filename string) *rangeRequest {
	return &rangeRequest{
		bytesRange: br,
		filename:   filename,
	}
}

/*
	This Object encapsulates a full download request.
*/

// Download Files -> Compile into 1 -> Delete Sub-files
// Url -> Follow Redirects -> Create Download Request
func Staging(c *Config) (*DownloadReq, error) {
	head, err := getHead(c.URL)

	if err != nil {
		return nil, err
	}

	dReq := new(DownloadReq)
	dReq.file = new(fileObject)
	dReq.wg = new(sync.WaitGroup)

	dReq.conf = c
	dReq.URL = c.URL
	dReq.file.filename = c.Filename

	header := head.Header

	// Parse Headers
	fType, err := extractFileType(header)

	if err != nil {
		// Doesn't provide file type
		// TODO(subomi) Find other means to decode content & proceed.
		// This applies to other header fields retrieval mechanisms
		dReq.rangeDownload = false
	} else {
		dReq.rangeDownload = true
	}

	dReq.file.fileExt = fType
	cSize, err := extractContentSize(header)

	if err != nil {
		return nil, err
	}
	intSize, err := strconv.Atoi(cSize)

	if err != nil {
		return nil, err
	}
	dReq.file.size = int64(intSize)

	if dReq.rangeDownload {
		dReq.bytesPerRoutine = dReq.file.size / c.Routines
		dReq.rr, err = createRR(dReq)

		if err != nil {
			return nil, err
		}
	}

	return dReq, nil
}

func createRR(dr *DownloadReq) ([]*rangeRequest, error) {
	var lower, upper int64
	var rTrips []*rangeRequest

	for i := 0; i < int(dr.conf.Routines); i++ {

		if i == 0 {
			lower, upper = 0, dr.bytesPerRoutine-int64(1)
		} else if i+1 == int(dr.conf.Routines) {
			lower = upper + 1
			if rem := dr.file.size - (lower + dr.bytesPerRoutine); rem != 0 {
				upper = lower + dr.bytesPerRoutine + rem
			} else {
				upper = lower + dr.bytesPerRoutine
			}
		} else {
			lower = upper + 1
			upper = lower + (dr.bytesPerRoutine - 1)
		}

		//log.Println(i, "=>", lower, upper)
		bytesRange := [2]string{strconv.Itoa(int(lower)), strconv.Itoa(int(upper))}
		filename := dr.file.filename + "-" + strconv.Itoa(i)

		rr := NewRangeRequest(bytesRange, filename)
		rr.dr = dr
		rTrips = append(rTrips, rr)
	}

	return rTrips, nil
}

func checkCloser(c io.Closer) error {
	err := c.Close()
	if err != nil {
		return err
	}
	return nil
}

func makeRange(lower, upper string) []string {
	return []string{"bytes=" + lower + "-" + upper}
}

func addPath(c *Config, name string) string {
	return c.Path + name
}

func makeString(str []string) string {
	return strings.Join(str, "")
}
