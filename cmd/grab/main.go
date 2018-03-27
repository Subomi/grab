package main

import (
	"flag"
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

	grab.Staging(c)
	grab.Download(c)
}

func parseArgs(c *grab.Config) {
	flag.StringVar(&c.Url, "url", "https://", "Url to download from ")
	flag.StringVar(&c.Filename, "filename", "grabbed", "Filename used to save file")
	flag.UintVar(&c.Routines, "routines", 4, "Number of goroutines to use to download")

	flag.Parse()
}
