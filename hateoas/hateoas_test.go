package hateoas

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	ts *httptest.Server
)

func init() {
	ts = CreateTestServer()
}

type Bob struct {
	Name  string `json:"Name"`
	Links Links  `json:"Links"`
}

type Where struct {
	Where string `json:"Where"`
}

func MakeAbsolute(base, rel string) string {
	u, err := url.Parse(ts.URL)
	if err != nil {
		panic(err)
	}
	u.Path = rel
	return u.String()
}

func handleEntry(w http.ResponseWriter, r *http.Request) {
	ep := SimpleEndpoint{
		Links: Links{
			Link{
				Rel:  "bob",
				Href: MakeAbsolute(ts.URL, "/bob"),
			},
		},
	}
	buf, err := json.Marshal(ep)
	if err != nil {
		http.Error(w, "oops", http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, string(buf))
}

func handleBob(w http.ResponseWriter, r *http.Request) {
	ep := Bob{
		Name: "bob",
		Links: Links{
			Link{
				Rel:  "where",
				Href: MakeAbsolute(ts.URL, "/where"),
			},
		},
	}
	buf, err := json.Marshal(ep)
	if err != nil {
		http.Error(w, "oops", http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, string(buf))
}

func handleWhere(w http.ResponseWriter, r *http.Request) {
	ep := Where{
		Where: "over the rainbow",
	}
	buf, err := json.Marshal(ep)
	if err != nil {
		http.Error(w, "oops", http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, string(buf))
}

func CreateTestServer() *httptest.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/", handleEntry)
	mux.HandleFunc("/bob", handleBob)
	mux.HandleFunc("/where", handleWhere)
	ts := httptest.NewServer(mux)
	return ts
}

func TestEntry(t *testing.T) {
	assert := assert.New(t)

	ts := CreateTestServer()
	defer ts.Close()

	client := Create(&Client{
		EntryURL: ts.URL,
	})

	var entry SimpleEndpoint
	_, err := client.Get("", Navigate{}, nil, nil, &entry)
	assert.Nil(err)

	bob, err := entry.Links.Get("bob")
	assert.Nil(err)
	assert.Equal(MakeAbsolute(ts.URL, "/bob"), bob.Href)
}

func TestBob(t *testing.T) {
	assert := assert.New(t)

	client := Create(&Client{
		EntryURL: ts.URL,
	})

	var bob Bob
	resp, err := client.Get("", Navigate{"bob"}, nil, nil, &bob)
	assert.Nil(err)
	assert.Equal("bob", bob.Name)

	ep := resp.Request.URL.String()
	resp, err = client.Get(ep, nil, nil, nil, &bob)
	assert.Nil(err)
	assert.Equal("bob", bob.Name)

	var where Where
	resp, err = client.Get(ep, Navigate{"where"}, nil, nil, &where)
	assert.Nil(err)
	assert.Equal("over the rainbow", where.Where)

	resp, err = client.Get("", Navigate{"bob", "where"}, nil, nil, &where)
	assert.Nil(err)
	assert.Equal("over the rainbow", where.Where)

	resp, err = client.Get("", Navigate{"bob", "nowhere"}, nil, nil, &where)
	assert.Equal(ErrorLinkNotFound, err)
}

type httpLogger struct {
	http.Client
	logged []string
}

func (h *httpLogger) Do(req *http.Request) (*http.Response, error) {
	resp, err := h.Client.Do(req)
	if resp != nil {
		h.logged = append(h.logged, fmt.Sprintf("%s %s %d", req.Method, req.URL.String(), resp.StatusCode))
	} else {
		h.logged = append(h.logged, fmt.Sprintf("%s %s (%s)", req.Method, req.URL.String(), err.Error()))
	}
	return resp, err
}

func TestLogger(t *testing.T) {
	assert := assert.New(t)

	logger := &httpLogger{Client: http.Client{}}
	client := Create(&Client{
		EntryURL: ts.URL,
		Http:     logger,
	})

	var bob Bob
	_, err := client.Get("", Navigate{"bob"}, nil, nil, &bob)
	assert.Nil(err)
	assert.Equal("bob", bob.Name)
	assert.Equal(2, len(logger.logged))

	for _, s := range logger.logged {
		log.Println(s)
	}
}

func TestSkipTLS(t *testing.T) {
	assert := assert.New(t)

	client := Create(&Client{
		EntryURL: ts.URL,
		Http: &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
			},
		},
	})

	var bob Bob
	_, err := client.Get("", Navigate{"bob"}, nil, nil, &bob)
	assert.Nil(err)
	assert.Equal("bob", bob.Name)
}
