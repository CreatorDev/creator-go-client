package hateoas

import (
	"encoding/json"
	"fmt"
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

	client, _ := Create(&Client{
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

	client, _ := Create(&Client{
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
