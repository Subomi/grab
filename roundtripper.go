package grab

import "net/http"

/*
	Package that simply downloads.
	- Should Implement the Interface RoundTripper tho
*/

type dRoundTrip struct {
	transport http.RoundTripper
	req       *http.Request
	dr        *DownloadReq
}

func NewdRoundTrip() (*dRoundTrip, error) {
	return &dRoundTrip{
		transport: http.DefaultTransport,
	}, nil
}

func (rt *dRoundTrip) RoundTrip(req *http.Request) (*http.Response, error) {
	return rt.transport.RoundTrip(req)
}
