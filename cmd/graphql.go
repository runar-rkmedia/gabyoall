package cmd

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/runar-rkmedia/gabyoall/logger"
	"github.com/runar-rkmedia/gabyoall/queries"
	"golang.org/x/net/http2"
)

type Endpoint struct {
	l       logger.AppLogger
	Url     string
	Headers http.Header
	client  HttpClient
}

func NewMockedEndpoint(l logger.AppLogger, url string) Endpoint {
	return NewEndpointWithClient(l, url, MockHttpClient{})
}
func NewEndpoint(l logger.AppLogger, url string) Endpoint {
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
		MaxIdleConnsPerHost: 1000,
		DisableCompression:  true,
		DisableKeepAlives:   true,
		// Proxy:               http.ProxyURL(b.ProxyAddr),
	}
	useHTTP2 := false
	if useHTTP2 {
		http2.ConfigureTransport(transport)
	} else {
		transport.TLSNextProto = make(map[string]func(string, *tls.Conn) http.RoundTripper)
	}
	client := &http.Client{
		Transport: transport,
		// intentionally long timeout
		Timeout: time.Minute * 30,
	}
	return NewEndpointWithClient(l, url, client)
}
func NewEndpointWithClient(l logger.AppLogger, url string, client HttpClient) Endpoint {
	if url == "" {
		l.Fatal().Str("url", url).Msg("Got empty url")
		return Endpoint{}
	}
	if !strings.Contains(url, "://") {
		url = "https://" + url
	}
	if client == nil {
		l.Fatal().Msg("client is required")
	}

	return Endpoint{
		l:       l,
		Url:     url,
		Headers: http.Header{},
		client:  client,
	}
}

type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
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

func (g *Endpoint) RunQuery(query queries.GraphQLQuery, okStatusCodes []int) (*http.Response, RequestStat, error) {
	stat := NewStat()
	l := logger.AppLogger{Logger: g.l.With().Str("operationName", query.OperationName).Str("endpoint", g.Url).Str("requestId", stat.RequestID).Logger()}
	b, err := json.Marshal(query)
	if err != nil {
		l.Fatal().Err(err).Msg("Failed to marshal query")
		return nil, stat.End(ServerTestError, err), err
	}

	r, err := http.NewRequest(http.MethodPost, g.Url, bytes.NewReader(b))
	if err != nil {
		l.Error().Err(err).Msg("Failed to create request")
		return nil, stat.End(ServerTestError, err), err
	}
	r.Body.Close()
	return g.DoRequest(l, r, stat, okStatusCodes)
}
func (g *Endpoint) DoRequest(l logger.AppLogger, r *http.Request, stat RequestStat, okStatusCodes []int) (*http.Response, RequestStat, error) {
	debug := g.l.HasDebug()
	if debug {
		l.Debug().Msg("Creating request")
	}
	// Setting our own request-id is nice for easy tracing
	for k, v := range g.Headers {
		r.Header.Set(k, v[0])
	}
	r.Header.Set("Connection", "close")
	r.Close = true
	if r.Header.Get("X-Request-Id") == "" {
		r.Header.Set("X-Request-Id", stat.RequestID)
	}
	if r.Header.Get("Content-Type") == "" {
		r.Header.Set("Content-Type", "application/json")
	}
	if debug {
		l.Debug().Interface("headers", r.Header).Msg("Doing request")
	}
	res, err := g.client.Do(r)
	if err != nil {
		l.ErrErr(err).Msg("Failed to run request")
		return nil, stat.End(Unknwon+"Request", err), err
	}
	stat.StatusCode = res.StatusCode
	contentType := res.Header.Get("Content-Type")
	l = logger.AppLogger{Logger: l.With().Str("contentType", contentType).Int("statusCode", res.StatusCode).Logger()}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		l.ErrErr(err).Msg("failed to ready body")
		err = fmt.Errorf("failed to read body")
		return nil, stat.End(Unknwon+"Body", err), err
	}
	res.Body.Close()
	if l.HasTrace() {
		l.Trace().Str("rawBody", string(body)).Msg("got raw body")
	}
	if strings.Contains(contentType, "json") {
		var gqlResponse GqlResponse
		var gqlResponseRaw map[string]interface{}

		err = json.Unmarshal(body, &gqlResponseRaw)
		if err != nil {
			l.ErrWarn(err).Msg("Failed to unmarshal body (raw)")
		} else {
			stat.Response = gqlResponseRaw
			if debug {
				l.Debug().Interface("json-response", gqlResponseRaw).Msg("got json body")
			}
		}
		err = json.Unmarshal(body, &gqlResponse)
		if err != nil {
			l.ErrWarn(err).Msg("Failed to unmarshal body")
		} else {
			if gqlResponse.Errors != nil && len(gqlResponse.Errors) > 0 {
				firstMessage := gqlResponse.Errors[0].Message
				l.Error().Str("firstMessage", firstMessage).Interface("json-response", gqlResponseRaw).Msg("got errors in request")
				return nil, stat.End(ErrorType(firstMessage), err), err
			}

		}

	}
	// In case the server changed the request-id
	if id := res.Header.Get("X-Request-Id"); id != "" && id != stat.RequestID {
		stat.RequestID = res.Header.Get("X-Request-Id")
		l = l.WithStringPairs("requestId", stat.RequestID)
	}
	l.Debug().Int("statusCode", res.StatusCode).Interface("headers", res.Header).Msg("Got response")

	statusOk := false
	for _, s := range okStatusCodes {
		if res.StatusCode == s {
			statusOk = true
			break
		}

	}
	if len(okStatusCodes) == 0 {
		statusOk = res.StatusCode < 299
	}
	if !statusOk {
		stat.RawResponse = string(body)
		err := fmt.Errorf("Got non-ok-statusCode: %d", res.StatusCode)
		l.Error().Msg("Statuscode is not 2xx")
		return res, stat.End(ErrorType(fmt.Sprintf("%s-%d", NonOK, res.StatusCode)), err), err
	}
	switch contentType {
	case "text/html":
		l.Warn().Msg("Looks like an html-page. Is the endpoint correct")
	}
	return res, stat.End("", err), err

}
