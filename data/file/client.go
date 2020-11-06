package file

import (
	"io/ioutil"
)

type Client struct {
	path string
}

func NewClient(path string) *Client {
	client := Client{
		path: path,
	}
	return &client
}

func (c *Client) Read() ([]byte, error) {
	data, err := ioutil.ReadFile(c.path)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (c *Client) Write(data []byte) error {
	return ioutil.WriteFile(c.path, data, 0644)
}
