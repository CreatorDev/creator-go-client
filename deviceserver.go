// Golang client for the creatordev.io deviceserver REST API.
package deviceserver

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"time"

	h "github.com/CreatorKit/go-deviceserver-client/hateoas"
)

var (
	// ErrorInvalidKeyName can be sent in response to CreateAccessKey
	ErrorInvalidKeyName = errors.New("Invalid key name")
)

// Client is the main object for interacting with the deviceserver
type RESTClient struct {
	hclient      *h.Client
	token        OAuthToken
	tokenExpires time.Time
}

// Create constructs a deviceserver client from a provided hateoas client.
// If you want logging/caching etc, you should set those options during
// hateoas client initialisation
func Create(hclient *h.Client) (*RESTClient, error) {
	if hclient == nil ||
		hclient.EntryURL == "" {
		return nil, h.ErrorBadConfig
	}

	d := RESTClient{
		hclient: hclient,
	}

	return &d, nil
}

// Close will clean things up as required
func (d *RESTClient) Close() {

}

// SetBearerToken sets the Authorization header on the underlying hateoas client
func (d *RESTClient) SetBearerToken(token string) {
	if token != "" {
		d.hclient.DefaultHeaders["Authorization"] = "Bearer " + token
	} else {
		delete(d.hclient.DefaultHeaders, "Authorization")
	}
}

// CreateAccessKey does what it says on the tin. The client
// should already be authenticated somehow, by calling either
// Authenticate/RefreshAuth/SetBearerToken
func (d *RESTClient) CreateAccessKey(name string) (*AccessKey, error) {
	var key AccessKey

	// key names are not required, but make life much easier
	if name == "" {
		return nil, ErrorInvalidKeyName
	}

	_, err := d.hclient.Post("",
		h.Navigate{"accesskeys"},
		nil,
		bytes.NewBuffer([]byte(fmt.Sprintf(`{"Name":"%s"}`, name))),
		&key)

	return &key, err
}

// DeleteAccessKey does what it says on the tin
func (d *RESTClient) DeleteAccessKey(key *AccessKey) error {
	return d.DeleteSelf(&key.Links)
}

// GetAccessKeys returns the list of accesskeys in this organisation
func (d *RESTClient) GetAccessKeys(previous *AccessKeys) (*AccessKeys, error) {
	if previous == nil {
		var keys AccessKeys
		_, err := d.hclient.Get("",
			h.Navigate{"accesskeys"},
			nil,
			nil,
			&keys)
		return &keys, err
	}

	next, err := previous.PageInfo.Links.Get("next")
	if err == h.ErrorLinkNotFound {
		return nil, nil
	}

	var keys AccessKeys
	_, err = d.hclient.Get(next.Href,
		nil,
		nil,
		nil,
		&keys)
	return &keys, err
}

// Authenticate uses the provided key/secret to obtain an access_token/refresh_token
func (d *RESTClient) Authenticate(credentials *AccessKey) error {
	var token OAuthToken
	_, err := d.hclient.PostForm("",
		h.Navigate{"authenticate"},
		nil,
		url.Values{
			"grant_type": []string{"password"},
			"username":   []string{credentials.Key},
			"password":   []string{credentials.Secret},
		},
		&token)
	if err == nil {
		d.token = token
		d.tokenExpires = time.Now().Add(time.Duration(token.ExpiresIn) * time.Second)
		d.SetBearerToken(token.AccessToken)
	}
	return err
}

// RefreshAuth uses the provided refresh_token obtain an access_token/refresh_token
func (d *RESTClient) RefreshAuth(refreshToken string) error {
	var token OAuthToken
	_, err := d.hclient.PostForm("",
		h.Navigate{"authenticate"},
		nil,
		url.Values{
			"grant_type":    []string{"refresh_token"},
			"refresh_token": []string{refreshToken},
		},
		&token)
	if err == nil {
		d.token = token
		d.tokenExpires = time.Now().Add(time.Duration(token.ExpiresIn) * time.Second)
		d.SetBearerToken(token.AccessToken)
	}
	return err
}

func (d *RESTClient) GetClients(previous *Clients) (*Clients, error) {
	if previous == nil {
		var clients Clients
		_, err := d.hclient.Get("",
			h.Navigate{"clients"},
			nil,
			nil,
			&clients)

		return &clients, err
	}

	next, err := previous.PageInfo.Links.Get("next")
	if err == h.ErrorLinkNotFound {
		return nil, nil
	}

	var clients Clients
	_, err = d.hclient.Get(next.Href,
		nil,
		nil,
		nil,
		&clients)
	return &clients, err
}

func (d *RESTClient) GetSubscriptions(endpoint string, previous *Subscriptions) (*Subscriptions, error) {
	if endpoint != "" && previous != nil {
		return nil, errors.New("cannot get subscriptions for endpoint and previous")
	}

	if previous == nil {
		var subs Subscriptions
		_, err := d.hclient.Get(endpoint,
			h.Navigate{"subscriptions"},
			nil,
			nil,
			&subs)

		return &subs, err
	}

	next, err := previous.PageInfo.Links.Get("next")
	if err == h.ErrorLinkNotFound {
		return nil, nil
	}

	var subs Subscriptions
	_, err = d.hclient.Get(next.Href,
		nil,
		nil,
		nil,
		&subs)
	return &subs, err
}

// Subscribe sets up webhook subscriptions, i.e. COAP observations.
// The `endpoint` can be
// - "" (=entrypoint) to subscribe to ClientConnected/ClientDisconnected events
// - a specific resource "self" URL to subscribe to observations on that resource
func (d *RESTClient) Subscribe(endpoint string, req *SubscriptionRequest, resp *SubscriptionResponse) error {
	buf, err := json.Marshal(req)
	if err != nil {
		return err
	}

	_, err = d.hclient.Post(endpoint,
		h.Navigate{"subscriptions"},
		h.Headers{"Content-Type": "application/vnd.oma.lwm2m.subscription+json"},
		bytes.NewBuffer(buf),
		resp)

	return err
}

func (d *RESTClient) Unsubscribe(subscription *SubscriptionResponse) error {
	return d.DeleteSelf(&subscription.Links)
}

// Delete performs DELETE on the specified resource
func (d *RESTClient) Delete(endpoint string) error {
	_, err := d.hclient.Delete(endpoint, nil, nil, nil, nil)
	return err
}

// DeleteSelf will find the "self" link and DELETE that
func (d *RESTClient) DeleteSelf(links *h.Links) error {
	self, err := links.Get("self")
	if err != nil {
		return nil
	}
	return d.Delete(self.Href)
}

// HATEOAS exposes the underlying hateoas client so that you
// can use that where necessary. Shouldn't be needed often.
func (d *RESTClient) HATEOAS() *h.Client {
	return d.hclient
}
