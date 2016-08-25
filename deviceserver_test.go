package deviceserver

import (
	"fmt"
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

func TestCreateKey(t *testing.T) {
	d, err := Create(hateoas.Create(&hateoas.Client{
		EntryURL: deviceserverURL,
	}))
	assert.Nil(t, err)
	assert.NotNil(t, d)
	defer d.Close()

	token, _ := TokenFromPSK(deviceserverPSK)
	d.SetBearerToken(token)

	k, err := d.CreateAccessKey("bob")
	assert.Nil(t, err)
	assert.NotZero(t, k)

	k2, err := d.CreateAccessKey("bob")
	assert.Nil(t, err)
	assert.NotZero(t, k2)

	err = d.DeleteAccessKey(k)
	assert.Nil(t, err)

	err = d.DeleteAccessKey(k2)
	assert.Nil(t, err)

}

func TestSubscriptions(t *testing.T) {
	d, err := Create(hateoas.Create(&hateoas.Client{
		EntryURL: deviceserverURL,
	}))
	assert.Nil(t, err)
	assert.NotNil(t, d)
	defer d.Close()

	token, _ := TokenFromPSK(deviceserverPSK)
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
