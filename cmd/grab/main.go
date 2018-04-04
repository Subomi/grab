package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/play/grab"
)

const (
	DownloadPath = "/home/subomi/Documents/grabber/"
)

func main() {
	c := new(grab.Config)
	parseArgs(c)

	// TODO: Customize this value
	c.Path = DownloadPath

	err := os.Mkdir(c.Path, 0777)

	if os.IsExist(err) {
		// Its fine.
	}

	dr, err := grab.Staging(c)

	if err != nil {
		log.Fatal("An error occurred", err)
	}

	grab.Download(dr)

	if ok := grab.Compile(dr.Downloads); !ok {
		log.Fatal("Compilation Error Occurred!")
	}
	
	fmt.Println("Download Completed")
}

func parseArgs(c *grab.Config) {
	flag.StringVar(&c.URL, "url", "https://", "Url to download from ")
	flag.StringVar(&c.Filename, "filename", "grabbed", "Filename used to save file")
	flag.Int64Var(&c.Routines, "routines", 4, "Number of goroutines to use to download")

	flag.Parse()
}
