package deviceserver

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.flowcloud.systems/creator-ops/go-deviceserver-client/hateoas"
)

var (
	deviceserverURL = os.Getenv("DEVICESERVER_URL")
	deviceserverPSK = os.Getenv("DEVICESERVER_PSK")
)

func init() {
	if deviceserverURL == "" {
		fmt.Println("Please provide env var DEVICESERVER_URL")
		os.Exit(1)
	}
	if deviceserverPSK == "" {
		fmt.Println("Please provide env var DEVICESERVER_PSK")
		os.Exit(1)
	}
}

type httpLogger struct {
	http.Client
	logged []string
	dump   bool
}

func (h *httpLogger) Do(req *http.Request) (*http.Response, error) {
	resp, err := h.Client.Do(req)
	var s string
	if resp != nil {
		s = fmt.Sprintf("%s %s %d", req.Method, req.URL.String(), resp.StatusCode)
	} else {
		s = fmt.Sprintf("%s %s (%s)", req.Method, req.URL.String(), err.Error())
	}
	h.logged = append(h.logged, s)
	if h.dump {
		log.Println(s)
	}
	return resp, err
}

// TestAuth handles everything to do with keys and tokens etc
func TestAuth(t *testing.T) {
	logger := &httpLogger{dump: true}
	d, err := Create(hateoas.Create(&hateoas.Client{
		EntryURL: deviceserverURL,
		Http:     logger,
	}))
	assert.Nil(t, err)
	assert.NotNil(t, d)
	defer d.Close()

	// set token from admin PSK in order to create the first key
	token, _ := TokenFromPSK(deviceserverPSK, 0)
	d.SetBearerToken(token)

	k, err := d.CreateAccessKey("bob")
	assert.Nil(t, err)
	assert.NotZero(t, k)

	k2, err := d.CreateAccessKey("bob")
	assert.Nil(t, err)
	assert.NotZero(t, k2)

	// clear admin token
	d.SetBearerToken("")

	var entry hateoas.SimpleEndpoint
	_, err = d.HATEOAS().Get("", nil, nil, nil, &entry)
	assert.Nil(t, err)
	_, err = entry.Links.Get("accesskeys")
	assert.NotNil(t, err) // when not authenticated, should not be able to get accesskeys

	err = d.Authenticate(k)
	assert.Nil(t, err)

	keys, err := d.GetAccessKeys(nil)
	assert.Nil(t, err)
	justKeys := []string{}
	for _, k := range keys.Items {
		justKeys = append(justKeys, k.Key)
	}
	assert.Contains(t, justKeys, k.Key)

	// var entry hateoas.SimpleEndpoint
	_, err = d.HATEOAS().Get("", nil, nil, nil, &entry)
	assert.Nil(t, err)
	accesskeys, err := entry.Links.Get("accesskeys")
	assert.Nil(t, err) // when authenticated, should be able to get accesskeys
	assert.NotNil(t, accesskeys)
	assert.NotEqual(t, "", accesskeys.Href)
	assert.NotEqual(t, "", accesskeys.Rel)

	err = d.RefreshAuth(d.token.RefreshToken)
	assert.Nil(t, err)

	err = d.DeleteAccessKey(k)
	assert.Nil(t, err)

	err = d.Authenticate(k2)
	err = d.DeleteAccessKey(k2)
	assert.Nil(t, err)
}

func TestSubscriptions(t *testing.T) {
	logger := &httpLogger{dump: true}
	d, err := Create(hateoas.Create(&hateoas.Client{
		EntryURL: deviceserverURL,
		Http:     logger,
	}))
	assert.Nil(t, err)
	assert.NotNil(t, d)
	defer d.Close()

	token, _ := TokenFromPSK(deviceserverPSK, 0)
	d.SetBearerToken(token)

	k, err := d.CreateAccessKey("bob")
	assert.Nil(t, err)
	assert.NotZero(t, k)

	err = d.Authenticate(k)
	assert.Nil(t, err)

	sub := SubscriptionRequest{
		SubscriptionType: "ClientConnected",
		URL:              "http://127.0.0.1/mywebhook",
	}
	var resp SubscriptionResponse

	err = d.Subscribe("", &sub, &resp)
	assert.Nil(t, err)
	assert.NotEqual(t, "", resp.ID)

	err = d.Unsubscribe(&resp)
	assert.Nil(t, err)

	err = d.DeleteAccessKey(k)
	assert.Nil(t, err)
}
