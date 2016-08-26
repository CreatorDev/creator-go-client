package deviceserver

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"time"

	h "gitlab.flowcloud.systems/creator-ops/go-deviceserver-client/hateoas"
)

var (
	ErrorInvalidKeyName = errors.New("Invalid key name")
)

type Client struct {
	hclient      *h.Client
	token        OAuthToken
	tokenExpires time.Time
}

func Create(hclient *h.Client) (*Client, error) {
	if hclient == nil ||
		hclient.EntryURL == "" {
		return nil, h.ErrorBadConfig
	}

	d := Client{
		hclient: hclient,
	}

	return &d, nil
}

func (d *Client) Close() {

}

func (d *Client) SetBearerToken(token string) {
	if token != "" {
		d.hclient.DefaultHeaders["Authorization"] = "Bearer " + token
	} else {
		delete(d.hclient.DefaultHeaders, "Authorization")
	}
}

func (d *Client) CreateAccessKey(name string) (*AccessKey, error) {
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

func (d *Client) DeleteAccessKey(key *AccessKey) error {
	return d.DeleteSelf(&key.Links)
}

func (d *Client) GetAccessKeys() (*AccessKeys, error) {
	var keys AccessKeys
	_, err := d.hclient.Get("",
		h.Navigate{"accesskeys"},
		nil,
		nil,
		&keys,
	)
	return &keys, err
}

func (d *Client) Authenticate(credentials *AccessKey) error {
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

func (d *Client) RefreshAuth(refresh_token string) error {
	var token OAuthToken
	_, err := d.hclient.PostForm("",
		h.Navigate{"authenticate"},
		nil,
		url.Values{
			"grant_type":    []string{"refresh_token"},
			"refresh_token": []string{refresh_token},
		},
		&token)
	if err == nil {
		d.token = token
		d.tokenExpires = time.Now().Add(time.Duration(token.ExpiresIn) * time.Second)
		d.SetBearerToken(token.AccessToken)
	}
	return err
}

// Subscribe sets up webhook subscriptions, i.e. COAP observations.
// The `endpoint` can be
// - "" (=entrypoint) to subscribe to ClientConnected/ClientDisconnected events
// - a specific resource "self" URL to subscribe to observations on that resource
func (d *Client) Subscribe(endpoint string, req *SubscriptionRequest, resp *SubscriptionResponse) error {
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

func (d *Client) Unsubscribe(subscription *SubscriptionResponse) error {
	return d.DeleteSelf(&subscription.Links)
}

// Delete performs DELETE on the specified resource
func (d *Client) Delete(endpoint string) error {
	_, err := d.hclient.Delete(endpoint, nil, nil, nil, nil)
	return err
}

// DeleteSelf will find the "self" link and DELETE that
func (d *Client) DeleteSelf(links *h.Links) error {
	self, err := links.Get("self")
	if err != nil {
		return nil
	}
	return d.Delete(self.Href)
}

func (d *Client) HATEOAS() *h.Client {
	return d.hclient
}
