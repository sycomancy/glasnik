package infra

import (
	"fmt"
	"io"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

type tor struct {
	MaxTimeout         time.Duration   // Max Timeout
	MaxIdleConnections int             // Max Idle Connections
	Transport          *http.Transport // Custom transport
	once               sync.Once       // Sync so we only set up 1 client
	netClient          *http.Client    // http client
}

// Create an instance of a Tor client
func (t *tor) new() *http.Client {
	// ensure that we only create one
	t.once.Do(func() {
		// Transport configuration
		t.Transport = &http.Transport{
			Proxy:        torProxy(),           // We can use Tor or Smart Proxy - rotating IP addresses - if nil no proxy is used
			MaxIdleConns: t.MaxIdleConnections, // max idle connections
			// TODO: Change to DialContext because Dial is deprecated
			Dial: (&net.Dialer{
				Timeout: t.MaxTimeout, // max dialer timeout
			}).Dial,
			TLSHandshakeTimeout: t.MaxTimeout, // transport layer security max timeout
		}

		// HTTPClient
		t.netClient = &http.Client{
			Timeout:   t.MaxTimeout, // round tripper timeout
			Transport: t.Transport,  // how our HTTP requests are made
		}
	})

	return t.netClient
}

func (t *tor) newIp() {
	// Get a new proxy
	t.Transport.Proxy = torProxy()
}

type IncognitoClient struct {
	tor             *tor
	client          *http.Client
	backoffSchedule BackoffSchedule
}

// A backoff schedule for when and how often to retry failed HTTP
// requests. The first element is the time to wait after the
// first failure, the second the time to wait after the second
// failure, etc. After reaching the last element, retries stop
// and the request is considered failed.
type BackoffSchedule = []time.Duration

// NewIP swaps a client's transport with a new one
func NewIncognitoClient(backoffSchedule BackoffSchedule) *IncognitoClient {
	t := &IncognitoClient{}

	if backoffSchedule == nil {
		t.backoffSchedule = []time.Duration{1 * time.Second, 3 * time.Second, 10 * time.Second}
	} else {
		t.backoffSchedule = backoffSchedule
	}

	t.tor = &tor{
		MaxTimeout:         20 * time.Second,
		MaxIdleConnections: 10,
	}

	t.client = t.tor.new()

	return t
}

func (t *IncognitoClient) GetURLData(url string, headers map[string]string) (status string, body io.ReadCloser, err error) {
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		// TODO(sycomancy): handle this
		fmt.Print("Unable to create request", err)
		return "unknown", nil, err
	}

	for header, value := range headers {
		req.Header.Add(header, value)
	}

	res, err := t.client.Do(req)

	if err != nil {
		if strings.Contains(err.Error(), "connection refused") {
			panic("Unable to connect to tor. Please check if tor client is running.")
		}
		return res.Status, nil, err
	}

	body = res.Body

	return res.Status, body, nil
}

func (t *IncognitoClient) GetURLDataWithRetries(url string, headers map[string]string) (string, io.ReadCloser, error) {
	var body io.ReadCloser
	var err error
	var status string

	for _, backoff := range t.backoffSchedule {
		status, body, err = t.GetURLData(url, headers)

		// TODO(sycomancy): Refactor this, to specific
		if err == nil && status != "400 Bad Request" {
			break
		}

		logrus.WithFields(logrus.Fields{
			"error":    err,
			"status":   status,
			"retry in": backoff,
		}).Warning("http-client handling error from with retry")

		t.tor.newIp()
		time.Sleep(backoff)
	}

	// All retries failed
	if err != nil {
		body.Close()
		return status, nil, err
	}

	return status, body, nil
}

// TorProxy initializezs and returns a TOR SOCKS proxy function for use in a Transport
func torProxy() func(*http.Request) (*url.URL, error) {
	// A source of uniformly-distributed pseudo-random
	rand.Seed(time.Now().UnixNano())

	// Pseudo-random int value
	num := rand.Intn(0x7fffffff-10000) + 10000

	// Tor listens for SOCKS connections on localhost port 9050
	// TODO(sycomancy: read this from env variable
	base := "socks5://%d:x@127.0.0.1:9050"
	// Proxy url with random credentials
	rawUrl := fmt.Sprintf(base, num)
	// Parse proxy url
	proxyUrl, err := url.Parse(rawUrl)
	if err != nil {
		fmt.Println("invalid url to parse when creating proxy transport. err: ", err)
		return nil
	}

	proxy := http.ProxyURL(proxyUrl)
	return proxy
}
