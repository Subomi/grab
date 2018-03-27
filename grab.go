package grab

import (
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

type fileObject struct {
	filename string
	fileExt string
	size    int64
}

/*
	This Object encapsulates a full download request.
*/
type DownloadReq struct {
	conf *Config
	wg *sync.WaitGroup
	rangeDownload   bool
	bytesPerRoutine int64
	Url             string
	file            *fileObject
	rTrips []*dRoundTrip
}

// Download Files -> Compile into 1 -> Delete Sub-files

// Url -> Follow Redirects ->
func Staging(c *Config) (*DownloadReq, error) {
	head, err := getHead(c.Url)

	if err != nil {
		return nil, err
	}

	dReq := new(DownloadReq)
	dReq.conf = c

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
		dReq.rTrips = createRange(dReq)
	}

	return dReq, nil
}

func createRange(dr *DownloadReq) []*dRoundTrip {
	var lower, upper int64
	var rTrips []*dRoundTrip

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
		dr.wg.Add(1)

		rT := &dRoundTrip{
			bytesRange: [2]string{strconv.Itoa(int(lower)), strconv.Itoa(int(upper))},
		}

		files := append(rTrips, rT)


		downloadReq = DownloadReq{
			client:   hC,
			Url:      c.Url,
			Range:    [2]string{strconv.Itoa(int(lower)), strconv.Itoa(int(upper))},
			Filename: addPath(c, c.Filename+"-"+strconv.Itoa(i)),
			file:     obj,
		}

		go downloadWorker(downloadReq)

		files = append(dr.files, downloadReq.Filename)
	}
}

/*
	- Send HEAD Request & Get ContentSize & Check if Server supports range requests
*/

func Download(c *Config) {
	var files []string
	obj, contentSize, err := getContentSize(c.Url)

	if err != nil {
		log.Fatal("Couldn't get content size.", err)
	}

	if contentSize == -1 {
		// Size cannot be determined - Unsure how to handle this.
		log.Fatal("Server does not support content size header in 2018? ")
	}

	// shared wg
	var wg sync.WaitGroup

	var downloadReq DownloadReq
	var lower, upper uint64

	hC := &http.Client{
		Timeout: 1000 * time.Second,
	}

	bytesPerRoutine := uint64(contentSize) / uint64(c.Routines)

	log.Println("Total File Size is", contentSize)
	log.Printf("Bytes per %d Routines is %d", c.Routines, bytesPerRoutine)

	log.Println("Waiting for goroutines to finish downloading ..")
	wg.Wait() // Wait for all downloads to finish

	// Compilation stage
	log.Println("Compiling .. ")

	for _, file := range files {
		log.Println(file)
	}

	// Block thread for compilation
	done := compile(files)
	log.Println("Done!! .. Boo ya!!", done)
}

func getContentSize(url string) (*Object, int64, error) {
	var obj Object
	resp, err := http.Head(url)

	if err != nil {
		return nil, 0, err
	}

	// TODO: Move this to its own space
	if makeString(resp.Header["Accept-Ranges"]) != "bytes" {
		return nil, 0, errors.New("Server doesn't support range requests")
	}

	// Properties of downloaded object.
	log.Println(resp.Header["Content-Type"], resp.Header["Accept-Ranges"])
	obj.fileExt = extractFileType(resp.Header["Content-Type"])
	obj.size = resp.ContentLength

	log.Println("==>", obj.fileExt)

	defer resp.Body.Close()
	return &obj, resp.ContentLength, nil
}

func downloadWorker(dr DownloadReq) (err error) {
	// Download a file with range bytes

	defer dr.wg.Done()

	req, err := http.NewRequest("GET", dr.Url, nil)

	if err != nil {
		return err
	}

	req.Header["Range"] = makeRange(dr.Range[0], dr.Range[1])

	resp, err := dr.client.Do(req)

	if err != nil {
		return err
	}

	defer func() {
		err = checkCloser(resp.Body)
	}()

	tmpFile, err := os.Create(dr.Filename)

	if err != nil {
		log.Println(err)
		return err
	}

	defer func() {
		err = checkCloser(tmpFile)
	}()

	// Pipe resp to file
	written, err := io.Copy(tmpFile, resp.Body)

	if err != nil {
		return err
	}

	log.Printf("%d bytes has been saved to %s", written, dr.Filename)

	return nil
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

func compile(files []string) bool {
	var done bool

	file, err := os.OpenFile(files[0], os.O_APPEND|os.O_WRONLY, 0755)

	if err != nil {
		log.Fatal(err)
	}

	for s := 1; s < len(files); s++ {
		part, err := os.OpenFile(files[s], os.O_RDONLY, 0755)

		if err != nil {
			log.Fatal(err)
		}

		_, err = io.Copy(file, part)

		if err != nil {
			log.Fatalf(err.Error())
		}

		if err = part.Close(); err != nil {
			log.Fatal(err.Error())
		}

		if err = os.Remove(files[s]); err != nil {
			log.Fatal(err.Error())
		}
	}

	done = true
	return done
}
