package hateoas

import (
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

var (
	ErrorLinkNotFound = errors.New("Link not found")
	ErrorHttpStatus   = errors.New("HTTP status error")
	ErrorBadConfig    = errors.New("bad config")
)

type Link struct {
	Rel  string `json:"rel"`
	Href string `json:"href"`
	Type string `json:"type"`
}

type Links []Link

func (l Links) Get(rel string) (*Link, error) {
	for _, link := range l {
		if link.Rel == rel {
			return &link, nil
		}
	}
	return nil, ErrorLinkNotFound
}

type Navigate []string
type Headers map[string]string

type SimpleEndpoint struct {
	Links Links `json:"Links"`
}

// HTTPDoer would normally be an instance of *http.Client but
// by making it an interface we allow wrapping to permit
// logging or caching etc
type HTTPDoer interface {
	Do(*http.Request) (*http.Response, error)
}

type Client struct {
	EntryURL       string
	DefaultHeaders Headers
	Http           HTTPDoer
}

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
			return resp, ErrorHttpStatus
		}
		link, err := ep.Links.Get(nav)
		if err != nil {
			return resp, ErrorLinkNotFound
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
		return resp, ErrorHttpStatus
	}

	if result != nil {
		err = json.Unmarshal(respbody, result)
		if err != nil {
			return resp, err
		}
	}

	return resp, nil
}

func (c *Client) Get(url string, navigateLinks Navigate, headers Headers, body io.Reader, response interface{}) (*http.Response, error) {
	return c.Do("GET", url, navigateLinks, headers, body, response)
}

func (c *Client) Post(url string, navigateLinks Navigate, headers Headers, body io.Reader, response interface{}) (*http.Response, error) {
	return c.Do("POST", url, navigateLinks, headers, body, response)
}

func (c *Client) PostForm(url string, navigateLinks Navigate, headers Headers, data url.Values, response interface{}) (*http.Response, error) {
	if headers == nil {
		headers = Headers{}
	}
	headers["Content-Type"] = "application/x-www-form-urlencoded"
	return c.Do("POST", url, navigateLinks, headers, strings.NewReader(data.Encode()), response)
}

func (c *Client) Delete(url string, navigateLinks Navigate, headers Headers, body io.Reader, response interface{}) (*http.Response, error) {
	return c.Do("DELETE", url, navigateLinks, headers, body, response)
}
