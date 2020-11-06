package http

import (
	"bytes"
	"io/ioutil"
	"net"
	"net/http"
	"time"
)

type Client struct {
	httpClient *http.Client
	url        string
}

func NewClient(url string) *Client {
	client := Client{
		httpClient: newHttpClient(),
		url:        url,
	}
	return &client
}

func newHttpClient() *http.Client {
	transport := &http.Transport{
		Dial: (&net.Dialer{
			Timeout: 3 * time.Second,
		}).Dial,
		TLSHandshakeTimeout: 5 * time.Second,
	}
	return &http.Client{
		Timeout:   time.Minute * 3,
		Transport: transport,
	}
}

func (c *Client) Read() ([]byte, error) {
	requestURL := c.url
	req, err := http.NewRequest(http.MethodGet, requestURL, bytes.NewBuffer([]byte{}))
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()

	req.URL.RawQuery = q.Encode()

	responseBody, err := c.request(req)
	if err != nil {
		return nil, err
	}

	return responseBody, nil
}

func (c *Client) request(req *http.Request) ([]byte, error) {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, err
	}

	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return responseBody, nil
}
