package http

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

type Client struct {
	method       string
	Url          string
	requestBody  []byte
	header       http.Header
	byteResponse []byte
	error        error
	statusCode   int
}

func NewHttpClient() *Client {
	return &Client{}
}

func (c *Client) SetMethod(method string) *Client {
	c.method = method
	return c
}

func (c *Client) SetUrl(url string) *Client {
	c.Url = url
	return c
}

func (c *Client) SetHeader(header http.Header) *Client {
	c.header = header
	return c
}

func (c *Client) SetRequestBody(body any) *Client {
	jsonBytes, err := json.Marshal(body)
	if err != nil {
		c.error = err
		return c
	}

	c.requestBody = jsonBytes
	return c
}

func (c *Client) SetJsonHeader() *Client {
	if c.header == nil {
		c.header = make(http.Header)
	}
	c.header.Set("Content-Type", "application/json")
	return c
}

func (c *Client) Do() *Client {
	if c.error != nil {
		return c
	}

	req, _ := http.NewRequest(c.method, c.Url, bytes.NewReader(c.requestBody))
	req.Header = c.header

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		c.error = err
		return c
	}
	defer resp.Body.Close()

	c.statusCode = resp.StatusCode
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		c.error = err
		return c
	}

	c.byteResponse = b
	return c
}

func (c *Client) UnmarshalResponse(v any) *Client {
	if c.error != nil {
		return c
	}

	err := json.Unmarshal(c.byteResponse, &v)
	if err != nil {
		c.error = err
		return c
	}
	return c
}

func (c *Client) Error() error {
	return c.error
}

func (c *Client) Status() int {
	return c.statusCode
}
