// Generic-ish HATEOAS client for golang
package hateoas

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/pkg/errors"
)

var (
	ErrorLinkNotFound = "Link not found"
	ErrorHttpStatus   = "HTTP status error"
	ErrorBadConfig    = "bad config"
)

// Link is the main HATEOAS link object
type Link struct {
	Rel  string `json:"rel"`
	Href string `json:"href"`
	Type string `json:"type"`
}

// Links is just an array of Link, but with some helper methods
type Links []Link

// Get returns the link matching the specified `rel` name
func (l Links) Get(rel string) (*Link, error) {
	for _, link := range l {
		if link.Rel == rel {
			return &link, nil
		}
	}
	return nil, errors.New(ErrorLinkNotFound)
}

// Self is a small wrapper around Get, mostly useful when you don't need to worry if the link isn't actually there
// (e.g. when printing debug stuff to a CLI output)
func (l Links) Self() string {
	link, err := l.Get("self")
	if err != nil {
		return `(unable to find "self" link)`
	}
	return link.Href
}

func (l Links) String() string {
	s := "["
	for i, _ := range l {
		s += fmt.Sprintf("\n  %+v", l[i])
	}
	s += "]"
	return s
}

// Navigate specifies a list of links to be traversed
type Navigate []string

// Headers is the HTTP headers to be added to a request
type Headers map[string]string

// SimpleEndpoint is used when navigating links on endpoints
type SimpleEndpoint struct {
	Links Links `json:"Links"`
}

// HTTPDoer would normally be an instance of *http.Client but
// by making it an interface we allow wrapping to permit
// logging or caching etc
type HTTPDoer interface {
	Do(*http.Request) (*http.Response, error)
}

// Client is the main object for executing requests
type Client struct {
	EntryURL       string
	DefaultHeaders Headers
	Http           HTTPDoer
}

// Create will populate some defaults into a provided Client structure
func Create(config *Client) *Client {
	client := config
	if client == nil {
		client = &Client{}
	}
	if client.Http == nil {
		client.Http = &http.Client{}
	}
	if client.DefaultHeaders == nil {
		client.DefaultHeaders = Headers{
			"Accept":       "application/json",
			"Content-Type": "application/json",
		}
	}
	return client
}

// Do will start at the provided URL (or default to `Client.EntryURL`) and traverse the links specified by `navigateLinks` (with GET)
// before finally issuing `method` to the resultant URL
func (c *Client) Do(method string, url string, navigateLinks Navigate, headers Headers, body io.Reader, result interface{}) (*http.Response, error) {
	if url == "" {
		url = c.EntryURL
	}
	for _, nav := range navigateLinks {
		var ep SimpleEndpoint
		resp, err := c.Get(url, nil, nil, nil, &ep)
		if err != nil {
			return resp, err
		}
		if resp.StatusCode >= http.StatusBadRequest {
			return resp, errors.New(ErrorHttpStatus)
		}
		link, err := ep.Links.Get(nav)
		if err != nil {
			return resp, errors.New(ErrorLinkNotFound)
		}
		url = link.Href
	}

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	for n, v := range c.DefaultHeaders {
		req.Header.Set(n, v)
	}
	for n, v := range headers {
		if v == "" {
			req.Header.Del(n)
		} else {
			req.Header.Set(n, v)
		}
	}

	resp, err := c.Http.Do(req)
	if err != nil {
		return resp, err
	}

	respbody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return resp, err
	}
	resp.Body.Close()

	if resp.StatusCode >= http.StatusBadRequest {
		return resp, errors.New(ErrorHttpStatus)
	}

	if result != nil {
		err = json.Unmarshal(respbody, result)
		if err != nil {
			return resp, err
		}
	}

	return resp, nil
}

// Get is a small wrapper around Do
func (c *Client) Get(url string, navigateLinks Navigate, headers Headers, body io.Reader, response interface{}) (*http.Response, error) {
	return c.Do("GET", url, navigateLinks, headers, body, response)
}

// Post is a small wrapper around Do
func (c *Client) Post(url string, navigateLinks Navigate, headers Headers, body io.Reader, response interface{}) (*http.Response, error) {
	return c.Do("POST", url, navigateLinks, headers, body, response)
}

// PostForm is a small wrapper around Do
func (c *Client) PostForm(url string, navigateLinks Navigate, headers Headers, data url.Values, response interface{}) (*http.Response, error) {
	if headers == nil {
		headers = Headers{}
	}
	headers["Content-Type"] = "application/x-www-form-urlencoded"
	return c.Do("POST", url, navigateLinks, headers, strings.NewReader(data.Encode()), response)
}

// Delete is a small wrapper around Do
func (c *Client) Delete(url string, navigateLinks Navigate, headers Headers, body io.Reader, response interface{}) (*http.Response, error) {
	return c.Do("DELETE", url, navigateLinks, headers, body, response)
}
