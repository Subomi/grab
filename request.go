package grab

import (
	"errors"
	"log"
	"net/http"
	"time"
)

func getHead(url string) (*http.Response, error) {
	client := &http.Client{
		Timeout:       1000 * time.Second,
		CheckRedirect: redirectPolicy,
	}

	return client.Head(url)
}

/*
	Our redirection policy code. - 10 Redirections maximum
*/
func redirectPolicy(req *http.Request, via []*http.Request) error {
	log.Println("Last request was ", via[len(via)-1].URL)
	if len(via) >= 10 {
		return errors.New("stopped after 10 redirects")
	}
	return nil
}
