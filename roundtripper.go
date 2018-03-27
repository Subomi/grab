package grab

/*
	Package that simply downloads.
	- Should Implement the Interface RoundTripper tho
*/

type dRoundTrip struct {
	bytesRange [2]string
}

func NewdRoundTrip(bR [2]string) *dRoundTrip {
	return &dRoundTrip{
		bytesRange: bR,
	}
}
