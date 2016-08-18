package deviceserver

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/patrickmn/go-cache"
	"github.com/square/go-jose"
	l "gitlab.flowcloud.systems/creator-ops/logger"
)

type Config struct {
	BaseUrl       string
	PSK           string
	SkipTLSVerify bool
	Log           l.Logger
}

type Client struct {
	baseUrl      string
	signer       JwtSigner
	client       *http.Client
	cache        *cache.Cache
	authGetJson  map[string]string
	authPostJson map[string]string
	log          l.Logger
}

func Create(config *Config) (*Client, error) {
	ds := Client{
		baseUrl: config.BaseUrl,
		cache:   cache.New(120*time.Second, 30*time.Second),
		log:     config.Log,
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: config.SkipTLSVerify},
	}

	ds.client = &http.Client{
		Transport: tr,
	}

	ds.authGetJson = map[string]string{
		"Accept": "application/json",
	}

	ds.authPostJson = map[string]string{
		"Accept":       "application/json",
		"Content-Type": "application/json",
	}

	err := ds.signer.Init(jose.HS256, []byte(config.PSK))
	if err != nil {
		return nil, err
	}

	return &ds, nil
}

func (d *Client) Close() {

}

func (d *Client) SetBearerToken(token string) {
	if token != "" {
		d.authGetJson["Authorization"] = "Bearer " + token
		d.authPostJson["Authorization"] = "Bearer " + token
	} else {
		delete(d.authGetJson, "Authorization")
		delete(d.authPostJson, "Authorization")
	}
}

func (d *Client) AuthorizePSK(req *http.Request) error {

	// the lifetime should be shorter, but think I'm hitting some timezone issues at the moment
	orgClaim := OrgClaim{
		OrgID: 0,
		Exp:   time.Now().Add(60 * time.Minute).Unix(),
	}

	serialized, err := d.signer.MarshallSignSerialize(orgClaim)
	if err != nil {
		return err
	}
	//fmt.Println(serialized)

	req.Header.Set("Authorization", "Bearer "+serialized)
	return nil
}

func (d *Client) Get(url string, headers map[string]string, result interface{}) error {
	// cached, ok := d.cache.Get(url)
	// if ok {
	// 	err := interfacetools.CopyOut(cached, result)
	// 	fmt.Println(cached)
	// 	fmt.Println(result)
	// 	return err
	// }

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		d.log.Request(req, "NewRequest(): %s", err.Error())
		return err
	}

	for n, v := range headers {
		req.Header.Set(n, v)
	}

	if _, ok := headers["Authorization"]; !ok {
		err = d.AuthorizePSK(req)
		if err != nil {
			d.log.Request(req, "AuthorizePSK(): %s", err.Error())
			return err
		}
	}

	resp, err := d.client.Do(req)
	if err != nil {
		d.log.Request(req, "Do(): %s", err.Error())
		return err
	}
	if resp.StatusCode > 400 {
		d.log.Request(req, "status: %s(%d)", resp.Status, resp.StatusCode)
		return fmt.Errorf("http status: %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		d.log.Request(req, "ReadAll(): %s", err.Error())
		return err
	}
	resp.Body.Close()

	err = json.Unmarshal(body, result)
	if err != nil {
		d.log.Request(req, "Unmarshal(): %s", err.Error())
		return err
	}

	d.log.Request(req, "Get(): OK")
	d.cache.Set(url, result, cache.DefaultExpiration)
	return nil
}

func (d *Client) Post(url string, headers map[string]string, postbody io.Reader, result interface{}) error {

	req, err := http.NewRequest("POST", url, postbody)
	if err != nil {
		d.log.Request(req, "NewRequest(): %s", err.Error())
		return err
	}

	for n, v := range headers {
		req.Header.Set(n, v)
	}

	if _, ok := headers["Authorization"]; !ok {
		err = d.AuthorizePSK(req)
		if err != nil {
			d.log.Request(req, "AuthorizePSK(): %s", err.Error())
			return err
		}
	}

	resp, err := d.client.Do(req)
	if err != nil {
		d.log.Request(req, "Do(): %s", err.Error())
		return err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		d.log.Request(req, "ReadAll(): %s", err.Error())
		return err
	}
	resp.Body.Close()

	if resp.StatusCode > 400 {
		var e Error
		err = json.Unmarshal(body, &e)
		if err != nil {
			d.log.Request(req, "status: %s(%d)", resp.Status, resp.StatusCode)
		} else {
			d.log.Request(req, "[%s][%s][%s]", e.ErrorCode, e.ErrorMessage, e.ErrorDetails)
		}
		return fmt.Errorf("http status: %d", resp.StatusCode)
	}

	err = json.Unmarshal(body, result)
	if err != nil {
		d.log.Request(req, "Unmarshal(): %s", err.Error())
		return err
	}

	d.log.Request(req, "Post(): OK")
	return nil
}

func (d *Client) Delete(url string, headers map[string]string) error {
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		d.log.Request(req, "NewRequest(): %s", err.Error())
		return err
	}

	for n, v := range headers {
		req.Header.Set(n, v)
	}

	if _, ok := headers["Authorization"]; !ok {
		err = d.AuthorizePSK(req)
		if err != nil {
			d.log.Request(req, "AuthorizePSK(): %s", err.Error())
			return err
		}
	}

	resp, err := d.client.Do(req)
	if err != nil {
		d.log.Request(req, "Do(): %s", err.Error())
		return err
	}
	if resp.StatusCode >= 400 {
		d.log.Request(req, "status: %s(%d)", resp.Status, resp.StatusCode)
		return fmt.Errorf("http status: %d", resp.StatusCode)
	}

	d.log.Request(req, "Delete(): %d", resp.StatusCode)
	return nil
}

func (d *Client) CreateAccessKey(name string) (*AccessKey, error) {
	var entry EntryPoint
	err := d.Get(d.baseUrl, d.authGetJson, &entry)
	if err != nil {
		d.log.Error("Get(): %s", err.Error())
		return nil, err
	}

	accesskeys, err := entry.Links.GetLink("accesskeys")
	if err != nil {
		d.log.Error("GetLink(): %s", err.Error())
		return nil, err
	}

	var key AccessKey
	err = d.Post(accesskeys.Href,
		d.authPostJson,
		bytes.NewBuffer([]byte(fmt.Sprintf(`{"Name":"%s"}`, name))),
		&key)

	return &key, nil
}

func (d *Client) Subscribe(endpoint string, req *SubscriptionRequest, resp *SubscriptionResponse) error {
	var entry EntryPoint
	err := d.Get(endpoint, d.authGetJson, &entry)
	if err != nil {
		d.log.Error("Get(): %s", err.Error())
		return err
	}

	d.log.Info("subscribe: %s", endpoint)
	for _, l := range entry.Links {
		d.log.Info("link: Rel %s Href %s Type %s", l.Rel, l.Href, l.Type)
	}

	subscriptions, err := entry.Links.GetLink("subscriptions")
	if err != nil {
		d.log.Error("GetLink(): %s", err.Error())
		return err
	}

	buf, err := json.Marshal(req)
	if err != nil {
		d.log.Error("Marshal(): %s", err.Error())
		return err
	}

	err = d.Post(subscriptions.Href,
		d.authPostJson,
		bytes.NewBuffer(buf),
		resp)

	return err
}

func (d *Client) Unsubscribe(endpoint string) error {
	return d.Delete(endpoint, d.authPostJson)
}
