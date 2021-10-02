package requests

import (
	"io"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/runar-rkmedia/gabyoall/logger"
)

func NewMockedEndpoint(l logger.AppLogger, url string) Endpoint {
	return NewEndpointWithClient(l, url, MockHttpClient{})
}

type MockHttpClient struct{}

func (m MockHttpClient) Do(req *http.Request) (*http.Response, error) {

	// TODO: replace with a mocked http-client-interface
	time.Sleep(time.Millisecond * time.Duration(rand.Int63n(80)+1))
	errorType := Unknwon
	n := rand.Intn(7)
	switch n {
	case 1:
		errorType = NonOK
	case 2:
		errorType = ServerTestError
	case 3:
		errorType = "RandomErr"
	case 4:
		errorType = "OtherErr"
	case 6:
		errorType = "MadeUpError"

	}
	res := http.Response{
		StatusCode: rand.Intn(500) + 100,
		Body:       io.NopCloser(strings.NewReader(string(errorType))),
	}

	return &res, nil
}
