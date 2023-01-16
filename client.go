package gt

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"golang.org/x/net/http2"
	"gopkg.in/yaml.v3"
)

// Client wrap http.Client, used for Do request
type Client struct {
	ht     *http.Client
	Method string
	URL    string
	Body   io.Reader
	Header http.Header
	Resp   *http.Response
	Err    error
	Option
}

const (
	Timeout10 time.Duration = 10
	Timeout20 time.Duration = 20
)

const (
	contentType     = "Content-Type"
	xmlContentType  = "application/xml"
	jsonContentType = "application/json"
)

// Option http request option
type Option struct {
	Timeout time.Duration
}

var (
	defaultHeader = http.Header{
		contentType: []string{jsonContentType},
	}
)

// NewClient return http default client
func NewClient() *Client {
	return &Client{
		ht:     http.DefaultClient,
		Header: http.Header{},
	}
}

// NewDefaultClient return http default client and set Json header
func NewDefaultClient() *Client {
	return &Client{
		ht:     http.DefaultClient,
		Header: defaultHeader,
		Option: Option{
			Timeout: time.Second * Timeout10,
		},
	}
}

func (c *Client) Get() *Client {
	c.Method = "GET"
	return c
}

func (c *Client) Delete() *Client {
	c.Method = "DELETE"
	return c
}

func (c *Client) PUT() *Client {
	c.Method = "PUT"
	return c
}

func (c *Client) UPDATE() *Client {
	c.Method = "UPDATE"
	return c
}

func (c *Client) HEADER() *Client {
	c.Method = "HEADER"
	return c
}

func (c *Client) OPTIONS() *Client {
	c.Method = "OPTIONS"
	return c
}

func (c *Client) SetURL(URL string) *Client {
	c.URL = URL
	return c
}

func (c *Client) SetTimeout(duration time.Duration) *Client {
	c.Option.Timeout = duration
	return c
}

func (c *Client) SetLog(title string) *Client {
	c.ht.Transport = NewLogTrace(title)
	return c
}

func transformResponse(resp *http.Response) error {
	var body []byte
	if resp.Body != nil {
		data, err := io.ReadAll(resp.Body)
		switch err.(type) {
		case nil:
			body = data
		case http2.StreamError:
			return fmt.Errorf("stream error when reading response body, may be caused by closed connection. Please retry. Original error: %w", err)
		default:
			return fmt.Errorf("unexpected error when reading response body. Please retry. Original error: %w", err)
		}
	}
	if len(body) == 0 {
		return errors.New(resp.Status)
	}
	return errors.New(string(body))
}

func (c *Client) Do() *Client {
	req, err := http.NewRequest(c.Method, c.URL, c.Body)
	if err != nil {
		c.Err = err
		return c
	}

	c.ht.Timeout = c.Option.Timeout

	c.Resp, err = c.ht.Do(req)
	if err != nil {
		c.Err = err
		return c
	}

	switch {
	case c.Resp.StatusCode == http.StatusSwitchingProtocols:
		// no-op, we've been upgraded
	case c.Resp.StatusCode >= http.StatusInternalServerError:
		c.Err = errors.New("服务器内部错误")
		return c
	case c.Resp.StatusCode < http.StatusOK || c.Resp.StatusCode > http.StatusPartialContent:
		c.Err = transformResponse(c.Resp)
		return c
	}

	return c
}

func (c *Client) SetBody(body io.Reader) *Client {
	c.Body = body
	return c
}

func (c *Client) SetHeader(key string, values ...string) *Client {
	for _, v := range values {
		c.Header.Add(key, v)
	}
	return c
}

const (
	JSON DecodeFormat = "json"
	YAML DecodeFormat = "yaml"
)

var (
	DecoderTypeNotSupport = errors.New("decoder type not support")
)

type DecodeFormat string

func (c *Client) InTo(object interface{}, format DecodeFormat) error {
	if c.Err != nil {
		return c.Err
	}
	if object == nil {
		return nil
	}
	switch format {
	case JSON:
		decode := json.NewDecoder(c.Resp.Body)
		return decode.Decode(object)
	case YAML:
		decode := yaml.NewDecoder(c.Resp.Body)
		return decode.Decode(object)
	}
	return DecoderTypeNotSupport
}
