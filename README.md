

# deviceserver
`import "github.com/CreatorKit/go-deviceserver-client"`

* [Overview](#pkg-overview)
* [Index](#pkg-index)
* [Subdirectories](#pkg-subdirectories)

## <a name="pkg-overview">Overview</a>
Golang client for the creatordev.io deviceserver REST API.




## <a name="pkg-index">Index</a>
* [Variables](#pkg-variables)
* [func ParseVerify(serialized []byte, signingKey interface{}) ([]byte, error)](#ParseVerify)
* [func TokenFromPSK(psk string, orgID int) (token string, err error)](#TokenFromPSK)
* [type AccessKey](#AccessKey)
* [type AccessKeys](#AccessKeys)
* [type Client](#Client)
  * [func Create(hclient *h.Client) (*Client, error)](#Create)
  * [func (d *Client) Authenticate(credentials *AccessKey) error](#Client.Authenticate)
  * [func (d *Client) Close()](#Client.Close)
  * [func (d *Client) CreateAccessKey(name string) (*AccessKey, error)](#Client.CreateAccessKey)
  * [func (d *Client) Delete(endpoint string) error](#Client.Delete)
  * [func (d *Client) DeleteAccessKey(key *AccessKey) error](#Client.DeleteAccessKey)
  * [func (d *Client) DeleteSelf(links *h.Links) error](#Client.DeleteSelf)
  * [func (d *Client) GetAccessKeys(previous *AccessKeys) (*AccessKeys, error)](#Client.GetAccessKeys)
  * [func (d *Client) HATEOAS() *h.Client](#Client.HATEOAS)
  * [func (d *Client) RefreshAuth(refreshToken string) error](#Client.RefreshAuth)
  * [func (d *Client) SetBearerToken(token string)](#Client.SetBearerToken)
  * [func (d *Client) Subscribe(endpoint string, req *SubscriptionRequest, resp *SubscriptionResponse) error](#Client.Subscribe)
  * [func (d *Client) Unsubscribe(subscription *SubscriptionResponse) error](#Client.Unsubscribe)
* [type EntryPoint](#EntryPoint)
* [type Error](#Error)
* [type JwtSigner](#JwtSigner)
  * [func (s *JwtSigner) Init(alg jose.SignatureAlgorithm, signingKey interface{}) error](#JwtSigner.Init)
  * [func (s *JwtSigner) MarshallSignSerialize(in interface{}) (string, error)](#JwtSigner.MarshallSignSerialize)
* [type OAuthToken](#OAuthToken)
* [type OrgClaim](#OrgClaim)
* [type PageInfo](#PageInfo)
* [type SubscriptionRequest](#SubscriptionRequest)
* [type SubscriptionResponse](#SubscriptionResponse)


#### <a name="pkg-files">Package files</a>
[deviceserver.go](https://github.com/CreatorKit/go-deviceserver-client/blob/master/deviceserver.go) [jwt.go](https://github.com/CreatorKit/go-deviceserver-client/blob/master/jwt.go) [struct.go](https://github.com/CreatorKit/go-deviceserver-client/blob/master/struct.go) 



## <a name="pkg-variables">Variables</a>
``` go
var (
    // ErrorInvalidKeyName can be sent in response to CreateAccessKey
    ErrorInvalidKeyName = errors.New("Invalid key name")
)
```


## <a name="ParseVerify">func</a> [ParseVerify](https://github.com/CreatorKit/go-deviceserver-client/blob/master/jwt.go#L62)
``` go
func ParseVerify(serialized []byte, signingKey interface{}) ([]byte, error)
```
ParseVerify performs signature validation and returns byte string



## <a name="TokenFromPSK">func</a> [TokenFromPSK](https://github.com/CreatorKit/go-deviceserver-client/blob/master/jwt.go#L16)
``` go
func TokenFromPSK(psk string, orgID int) (token string, err error)
```
TokenFromPSK generates an JWT with signed OrgClaim




## <a name="AccessKey">type</a> [AccessKey](https://github.com/CreatorKit/go-deviceserver-client/blob/master/struct.go#L11)
``` go
type AccessKey struct {
    Links  hateoas.Links `json:"Links"`
    Name   string        `json:"Name,omitempty"`
    Key    string        `json:"Key,omitempty"`
    Secret string        `json:"Secret,omitempty"`
}
```









## <a name="AccessKeys">type</a> [AccessKeys](https://github.com/CreatorKit/go-deviceserver-client/blob/master/struct.go#L23)
``` go
type AccessKeys struct {
    PageInfo PageInfo      `json:"PageInfo"`
    Items    []AccessKey   `json:"Items"`
    Links    hateoas.Links `json:"Links"`
}
```









## <a name="Client">type</a> [Client](https://github.com/CreatorKit/go-deviceserver-client/blob/master/deviceserver.go#L21)
``` go
type Client struct {
    // contains filtered or unexported fields
}
```
Client is the main object for interacting with the deviceserver







### <a name="Create">func</a> [Create](https://github.com/CreatorKit/go-deviceserver-client/blob/master/deviceserver.go#L30)
``` go
func Create(hclient *h.Client) (*Client, error)
```
Create constructs a deviceserver client from a provided hateoas client.
If you want logging/caching etc, you should set those options during
hateoas client initialisation





### <a name="Client.Authenticate">func</a> (\*Client) [Authenticate](https://github.com/CreatorKit/go-deviceserver-client/blob/master/deviceserver.go#L109)
``` go
func (d *Client) Authenticate(credentials *AccessKey) error
```
Authenticate uses the provided key/secret to obtain an access_token/refresh_token




### <a name="Client.Close">func</a> (\*Client) [Close](https://github.com/CreatorKit/go-deviceserver-client/blob/master/deviceserver.go#L44)
``` go
func (d *Client) Close()
```
Close will clean things up as required




### <a name="Client.CreateAccessKey">func</a> (\*Client) [CreateAccessKey](https://github.com/CreatorKit/go-deviceserver-client/blob/master/deviceserver.go#L60)
``` go
func (d *Client) CreateAccessKey(name string) (*AccessKey, error)
```
CreateAccessKey does what it says on the tin. The client
should already be authenticated somehow, by calling either
Authenticate/RefreshAuth/SetBearerToken




### <a name="Client.Delete">func</a> (\*Client) [Delete](https://github.com/CreatorKit/go-deviceserver-client/blob/master/deviceserver.go#L171)
``` go
func (d *Client) Delete(endpoint string) error
```
Delete performs DELETE on the specified resource




### <a name="Client.DeleteAccessKey">func</a> (\*Client) [DeleteAccessKey](https://github.com/CreatorKit/go-deviceserver-client/blob/master/deviceserver.go#L78)
``` go
func (d *Client) DeleteAccessKey(key *AccessKey) error
```
DeleteAccessKey does what it says on the tin




### <a name="Client.DeleteSelf">func</a> (\*Client) [DeleteSelf](https://github.com/CreatorKit/go-deviceserver-client/blob/master/deviceserver.go#L177)
``` go
func (d *Client) DeleteSelf(links *h.Links) error
```
DeleteSelf will find the "self" link and DELETE that




### <a name="Client.GetAccessKeys">func</a> (\*Client) [GetAccessKeys](https://github.com/CreatorKit/go-deviceserver-client/blob/master/deviceserver.go#L83)
``` go
func (d *Client) GetAccessKeys(previous *AccessKeys) (*AccessKeys, error)
```
GetAccessKeys returns the list of accesskeys in this organisation




### <a name="Client.HATEOAS">func</a> (\*Client) [HATEOAS](https://github.com/CreatorKit/go-deviceserver-client/blob/master/deviceserver.go#L187)
``` go
func (d *Client) HATEOAS() *h.Client
```
HATEOAS exposes the underlying hateoas client so that you
can use that where necessary. Shouldn't be needed often.




### <a name="Client.RefreshAuth">func</a> (\*Client) [RefreshAuth](https://github.com/CreatorKit/go-deviceserver-client/blob/master/deviceserver.go#L129)
``` go
func (d *Client) RefreshAuth(refreshToken string) error
```
RefreshAuth uses the provided refresh_token obtain an access_token/refresh_token




### <a name="Client.SetBearerToken">func</a> (\*Client) [SetBearerToken](https://github.com/CreatorKit/go-deviceserver-client/blob/master/deviceserver.go#L49)
``` go
func (d *Client) SetBearerToken(token string)
```
SetBearerToken sets the Authorization header on the underlying hateoas client




### <a name="Client.Subscribe">func</a> (\*Client) [Subscribe](https://github.com/CreatorKit/go-deviceserver-client/blob/master/deviceserver.go#L151)
``` go
func (d *Client) Subscribe(endpoint string, req *SubscriptionRequest, resp *SubscriptionResponse) error
```
Subscribe sets up webhook subscriptions, i.e. COAP observations.
The `endpoint` can be
- "" (=entrypoint) to subscribe to ClientConnected/ClientDisconnected events
- a specific resource "self" URL to subscribe to observations on that resource




### <a name="Client.Unsubscribe">func</a> (\*Client) [Unsubscribe](https://github.com/CreatorKit/go-deviceserver-client/blob/master/deviceserver.go#L166)
``` go
func (d *Client) Unsubscribe(subscription *SubscriptionResponse) error
```



## <a name="EntryPoint">type</a> [EntryPoint](https://github.com/CreatorKit/go-deviceserver-client/blob/master/struct.go#L11)
``` go
type EntryPoint struct {
    Links hateoas.Links `json:"Links"`
}
```









## <a name="Error">type</a> [Error](https://github.com/CreatorKit/go-deviceserver-client/blob/master/struct.go#L61)
``` go
type Error struct {
    ErrorCode    string `json:"ErrorCode"`
    ErrorMessage string `json:"ErrorMessage"`
    ErrorDetails string `json:"ErrorDetails"`
}
```









## <a name="JwtSigner">type</a> [JwtSigner](https://github.com/CreatorKit/go-deviceserver-client/blob/master/jwt.go#L11)
``` go
type JwtSigner struct {
    // contains filtered or unexported fields
}
```
JwtSigner is the main object for simplified JWT operations










### <a name="JwtSigner.Init">func</a> (\*JwtSigner) [Init](https://github.com/CreatorKit/go-deviceserver-client/blob/master/jwt.go#L39)
``` go
func (s *JwtSigner) Init(alg jose.SignatureAlgorithm, signingKey interface{}) error
```
Init creates JOSE signer




### <a name="JwtSigner.MarshallSignSerialize">func</a> (\*JwtSigner) [MarshallSignSerialize](https://github.com/CreatorKit/go-deviceserver-client/blob/master/jwt.go#L46)
``` go
func (s *JwtSigner) MarshallSignSerialize(in interface{}) (string, error)
```
MarshallSignSerialize returns a compacted serialised JWT from a claims structure




## <a name="OAuthToken">type</a> [OAuthToken](https://github.com/CreatorKit/go-deviceserver-client/blob/master/struct.go#L54)
``` go
type OAuthToken struct {
    AccessToken  string `json:"access_token"`
    TokenType    string `json:"token_type"`
    ExpiresIn    int    `json:"expires_in"`
    RefreshToken string `json:"refresh_token"`
}
```









## <a name="OrgClaim">type</a> [OrgClaim](https://github.com/CreatorKit/go-deviceserver-client/blob/master/struct.go#L49)
``` go
type OrgClaim struct {
    OrgID int   `json:"OrgID"`
    Exp   int74 `json:"exp"`
}
```









## <a name="PageInfo">type</a> [PageInfo](https://github.com/CreatorKit/go-deviceserver-client/blob/master/struct.go#L16)
``` go
type PageInfo struct {
    TotalCount int           `json:"TotalCount"`
    ItemsCount int           `json:"ItemsCount"`
    StartIndex int           `json:"StartIndex"`
    Links      hateoas.Links `json:"Links,omitempty"`
}
```









## <a name="SubscriptionRequest">type</a> [SubscriptionRequest](https://github.com/CreatorKit/go-deviceserver-client/blob/master/struct.go#L29)
``` go
type SubscriptionRequest struct {
    SubscriptionType string `json:"SubscriptionType"`
    URL              string `json:"Url"`

    AcceptContentType string `json:"AcceptContentType,omitempty"`
    Property          string `json:"Property,omitempty"`
    Attributes        *struct {
        Pmin        string `json:"Pmin,omitempty"`
        Pmax        string `json:"Pmax,omitempty"`
        Step        string `json:"Step,omitempty"`
        LessThan    string `json:"LessThan,omitempty"`
        GreaterThan string `json:"GreaterThan,omitempty"`
    } `json:"Attributes,omitempty"`
}
```









## <a name="SubscriptionResponse">type</a> [SubscriptionResponse](https://github.com/CreatorKit/go-deviceserver-client/blob/master/struct.go#L44)
``` go
type SubscriptionResponse struct {
    ID    string        `json:"ID"`
    Links hateoas.Links `json:"Links"`
}
```













- - -
Generated by [godoc12md](http://godoc.org/github.com/davecheney/godoc2md)
