package deviceserver

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEntry(t *testing.T) {
	ds, err := Create(&DSConfig{
		BaseUrl: "https://deviceserver-mv.flowcloud.systems/",
		PSK:     os.Getenv("DEVICESERVER_PSK"),
	})
	assert.Nil(t, err)
	assert.NotNil(t, ds)

	var entry EntryPoint
	err = ds.Get(ds.baseUrl, ds.authGetJson, &entry)
	assert.Nil(t, err)
	assert.NotZero(t, entry)

	auth, err := entry.Links.GetLink("authenticate")
	assert.Nil(t, err)
	assert.NotZero(t, auth.Href)
	assert.NotZero(t, auth.Rel)

	bob, err := entry.Links.GetLink("bob")
	assert.NotNil(t, err)
	assert.Nil(t, bob)

	accesskeys, err := entry.Links.GetLink("accesskeys")
	assert.Nil(t, err)
	assert.NotZero(t, accesskeys.Href)
	assert.NotZero(t, accesskeys.Rel)
}

func TestCreateKey(t *testing.T) {
	ds, err := Create(&DSConfig{
		BaseUrl: "https://deviceserver-mv.flowcloud.systems/",
		PSK:     os.Getenv("DEVICESERVER_PSK"),
	})
	assert.Nil(t, err)
	assert.NotNil(t, ds)

	k, err := ds.CreateAccessKey("bob")
	assert.Nil(t, err)
	assert.NotZero(t, k)

	k2, err := ds.CreateAccessKey("bob")
	assert.Nil(t, err)
	assert.NotZero(t, k2)
}
