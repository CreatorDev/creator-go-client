package deviceserver

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.flowcloud.systems/creator-ops/logger"
)

func TestCreateKey(t *testing.T) {
	ds, err := Create(&Config{
		BaseUrl: "https://deviceserver-mv.flowcloud.systems/",
		PSK:     os.Getenv("DEVICESERVER_PSK"),
		Log:     &logger.LogLogger{},
	})
	assert.Nil(t, err)
	assert.NotNil(t, ds)
	defer ds.Close()

	token, _ := ds.TokenPSK()
	ds.SetBearerToken(token)

	k, err := ds.CreateAccessKey("bob")
	assert.Nil(t, err)
	assert.NotZero(t, k)

	k2, err := ds.CreateAccessKey("bob")
	assert.Nil(t, err)
	assert.NotZero(t, k2)

	err = ds.DeleteAccessKey(k)
	assert.Nil(t, err)

	err = ds.DeleteAccessKey(k2)
	assert.Nil(t, err)

}

func TestSubscriptions(t *testing.T) {
	ds, err := Create(&Config{
		BaseUrl: "https://deviceserver-mv.flowcloud.systems/",
		PSK:     os.Getenv("DEVICESERVER_PSK"),
		Log:     &logger.LogLogger{},
	})
	assert.Nil(t, err)
	assert.NotNil(t, ds)
	defer ds.Close()

	token, _ := ds.TokenPSK()
	ds.SetBearerToken(token)

	k, err := ds.CreateAccessKey("bob")
	assert.Nil(t, err)
	assert.NotZero(t, k)

	err = ds.Authenticate(k)
	assert.Nil(t, err)

	sub := SubscriptionRequest{
		SubscriptionType: "ClientConnected",
		URL:              "http://127.0.0.1/mywebhook",
	}
	var resp SubscriptionResponse

	err = ds.Subscribe("", &sub, &resp)
	assert.Nil(t, err)
	assert.NotEqual(t, "", resp.ID)

	err = ds.Unsubscribe(&resp)
	assert.Nil(t, err)

	err = ds.DeleteAccessKey(k)
	assert.Nil(t, err)
}
