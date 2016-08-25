package deviceserver

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/patrickmn/go-cache"
	"github.com/square/go-jose"
	h "gitlab.flowcloud.systems/creator-ops/go-deviceserver-client/hateoas"
	l "gitlab.flowcloud.systems/creator-ops/logger"
)

type Config struct {
	BaseUrl       string
	PSK           string
	SkipTLSVerify bool
	Log           l.Logger
}

type Client struct {
	signer       JwtSigner
	client       *http.Client
	cache        *cache.Cache
	authGetJson  h.Headers
	authPostJson h.Headers
	log          l.Logger
	hclient      *h.Client
	token        OAuthToken
	tokenExpires time.Time
}

func Create(config *Config) (*Client, error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: config.SkipTLSVerify,
		},
	}

	hclient, err := h.Create(&h.Client{
		EntryURL: config.BaseUrl,

		Http: &http.Client{
			Transport: tr,
		},
	})
	if err != nil {
		return nil, err
	}

	ds := Client{
		hclient: hclient,
		cache:   cache.New(120*time.Second, 30*time.Second),
		log:     config.Log,
	}

	err = ds.signer.Init(jose.HS256, []byte(config.PSK))
	if err != nil {
		return nil, err
	}

	return &ds, nil
}

func (d *Client) Close() {

}

func (d *Client) SetBearerToken(token string) {
	if token != "" {
		d.log.Info("token=%s", token)
		d.hclient.DefaultHeaders["Authorization"] = "Bearer " + token
	} else {
		delete(d.hclient.DefaultHeaders, "Authorization")
	}
}

func (d *Client) TokenPSK() (token string, err error) {

	// the lifetime should be shorter, but think I'm hitting some timezone issues at the moment
	orgClaim := OrgClaim{
		OrgID: 0,
		Exp:   time.Now().Add(60 * time.Minute).Unix(),
	}

	serialized, err := d.signer.MarshallSignSerialize(orgClaim)
	if err != nil {
		return "", err
	}
	//fmt.Println(serialized)

	return serialized, nil
}

func (d *Client) CreateAccessKey(name string) (*AccessKey, error) {
	var key AccessKey
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
		d.SetBearerToken(token.AccessToken)
	}
	return err
}

func (d *Client) Subscribe(endpoint string, req *SubscriptionRequest, resp *SubscriptionResponse) error {
	buf, err := json.Marshal(req)
	if err != nil {
		d.log.Error("Marshal(): %s", err.Error())
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

func (d *Client) Delete(endpoint string) error {
	_, err := d.hclient.Delete(endpoint, nil, nil, nil, nil)
	return err
}

func (d *Client) DeleteSelf(links *h.Links) error {
	self, err := links.Get("self")
	if err != nil {
		return nil
	}
	return d.Delete(self.Href)
}
