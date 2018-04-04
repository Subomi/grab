package grab

import (
	"io"
	"log"
	"net/http"
	"os"
)

/*
	Fan Out to download workers
*/
func Download(dr *DownloadReq) {
	rt, err := NewdRoundTrip()

	if err != nil {
		log.Fatal(err)
	}

	client := &http.Client{
		Transport: rt,
	}
	for _, rr := range dr.rr {
		dr.wg.Add(1)
		dr.Downloads = append(dr.Downloads, rr.filename)
		go downloadWorker(client, rr)
	}

	dr.wg.Wait()
}

func downloadWorker(client *http.Client, rr *rangeRequest) (err error) {
	// Download a file with range bytes

	defer rr.dr.wg.Done()

	req, err := http.NewRequest("GET", rr.dr.URL, nil)

	if err != nil {
		log.Fatal("An error occurred creating new request", err)
		return err
	}

	req.Header["Range"] = makeRange(rr.bytesRange[0], rr.bytesRange[1])
	resp, err := client.Do(req)

	if err != nil {
		log.Fatal("An error occurred in the client", err)
		return err
	}

	defer func() {
		err = checkCloser(resp.Body)
	}()
	tmpFile, err := os.Create(rr.filename)

	if err != nil {
		log.Fatal("An error occurred .. ", err)
		return err
	}

	defer func() {
		err = checkCloser(tmpFile)
	}()

	// Pipe resp to file
	_, err = io.Copy(tmpFile, resp.Body)

	if err != nil {
		log.Fatal("An error occurred ..", err)
		return err
	}

	//log.Printf("%d bytes has been saved to %s", written, rr.filename)

	return nil
}
