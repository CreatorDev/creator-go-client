

# hateoas
`import "github.com/CreatorKit/go-deviceserver-client/hateoas"`

* [Overview](#pkg-overview)
* [Index](#pkg-index)

## <a name="pkg-overview">Overview</a>
Generic-ish HATEOAS client for golang




## <a name="pkg-index">Index</a>
* [Variables](#pkg-variables)
* [type Client](#Client)
  * [func Create(config *Client) *Client](#Create)
  * [func (c *Client) Delete(url string, navigateLinks Navigate, headers Headers, body io.Reader, response interface{}) (*http.Response, error)](#Client.Delete)
  * [func (c *Client) Do(method string, url string, navigateLinks Navigate, headers Headers, body io.Reader, result interface{}) (*http.Response, error)](#Client.Do)
  * [func (c *Client) Get(url string, navigateLinks Navigate, headers Headers, body io.Reader, response interface{}) (*http.Response, error)](#Client.Get)
  * [func (c *Client) Post(url string, navigateLinks Navigate, headers Headers, body io.Reader, response interface{}) (*http.Response, error)](#Client.Post)
  * [func (c *Client) PostForm(url string, navigateLinks Navigate, headers Headers, data url.Values, response interface{}) (*http.Response, error)](#Client.PostForm)
* [type HTTPDoer](#HTTPDoer)
* [type Headers](#Headers)
* [type Link](#Link)
* [type Links](#Links)
  * [func (l Links) Get(rel string) (*Link, error)](#Links.Get)
  * [func (l Links) Self() string](#Links.Self)
* [type Navigate](#Navigate)
* [type SimpleEndpoint](#SimpleEndpoint)


#### <a name="pkg-files">Package files</a>
[hateoas.go](https://github.com/CreatorKit/go-deviceserver-client/blob/master/hateoas/hateoas.go) 



## <a name="pkg-variables">Variables</a>
``` go
var (
    ErrorLinkNotFound = errors.New("Link not found")
    ErrorHttpStatus   = errors.New("HTTP status error")
    ErrorBadConfig    = errors.New("bad config")
)
```



## <a name="Client">type</a> [Client](https://github.com/CreatorKit/go-deviceserver-client/blob/master/hateoas/hateoas.go#L69)
``` go
type Client struct {
    EntryURL       string
    DefaultHeaders Headers
    Http           HTTPDoer
}
```
Client is the main object for executing requests







### <a name="Create">func</a> [Create](https://github.com/CreatorKit/go-deviceserver-client/blob/master/hateoas/hateoas.go#L76)
``` go
func Create(config *Client) *Client
```
Create will populate some defaults into a provided Client structure





### <a name="Client.Delete">func</a> (\*Client) [Delete](https://github.com/CreatorKit/go-deviceserver-client/blob/master/hateoas/hateoas.go#L176)
``` go
func (c *Client) Delete(url string, navigateLinks Navigate, headers Headers, body io.Reader, response interface{}) (*http.Response, error)
```
Delete is a small wrapper around Do




### <a name="Client.Do">func</a> (\*Client) [Do](https://github.com/CreatorKit/go-deviceserver-client/blob/master/hateoas/hateoas.go#L95)
``` go
func (c *Client) Do(method string, url string, navigateLinks Navigate, headers Headers, body io.Reader, result interface{}) (*http.Response, error)
```
Do will start at the provided URL (or default to `Client.EntryURL`) and traverse the links specified by `navigateLinks` (with GET)
before finally issuing `method` to the resultant URL




### <a name="Client.Get">func</a> (\*Client) [Get](https://github.com/CreatorKit/go-deviceserver-client/blob/master/hateoas/hateoas.go#L157)
``` go
func (c *Client) Get(url string, navigateLinks Navigate, headers Headers, body io.Reader, response interface{}) (*http.Response, error)
```
Get is a small wrapper around Do




### <a name="Client.Post">func</a> (\*Client) [Post](https://github.com/CreatorKit/go-deviceserver-client/blob/master/hateoas/hateoas.go#L162)
``` go
func (c *Client) Post(url string, navigateLinks Navigate, headers Headers, body io.Reader, response interface{}) (*http.Response, error)
```
Post is a small wrapper around Do




### <a name="Client.PostForm">func</a> (\*Client) [PostForm](https://github.com/CreatorKit/go-deviceserver-client/blob/master/hateoas/hateoas.go#L167)
``` go
func (c *Client) PostForm(url string, navigateLinks Navigate, headers Headers, data url.Values, response interface{}) (*http.Response, error)
```
PostForm is a small wrapper around Do




## <a name="HTTPDoer">type</a> [HTTPDoer](https://github.com/CreatorKit/go-deviceserver-client/blob/master/hateoas/hateoas.go#L64)
``` go
type HTTPDoer interface {
    Do(*http.Request) (*http.Response, error)
}
```
HTTPDoer would normally be an instance of *http.Client but
by making it an interface we allow wrapping to permit
logging or caching etc










## <a name="Headers">type</a> [Headers](https://github.com/CreatorKit/go-deviceserver-client/blob/master/hateoas/hateoas.go#L54)
``` go
type Headers map[string]string
```
Headers is the HTTP headers to be added to a request










## <a name="Link">type</a> [Link](https://github.com/CreatorKit/go-deviceserver-client/blob/master/hateoas/hateoas.go#L21)
``` go
type Link struct {
    Rel  string `json:"rel"`
    Href string `json:"href"`
    Type string `json:"type"`
}
```
Link is the main HATEOAS link object










## <a name="Links">type</a> [Links](https://github.com/CreatorKit/go-deviceserver-client/blob/master/hateoas/hateoas.go#L28)
``` go
type Links []Link
```
Links is just an array of Link, but with some helper methods










### <a name="Links.Get">func</a> (Links) [Get](https://github.com/CreatorKit/go-deviceserver-client/blob/master/hateoas/hateoas.go#L31)
``` go
func (l Links) Get(rel string) (*Link, error)
```
Get returns the link matching the specified `rel` name




### <a name="Links.Self">func</a> (Links) [Self](https://github.com/CreatorKit/go-deviceserver-client/blob/master/hateoas/hateoas.go#L42)
``` go
func (l Links) Self() string
```
Self is a small wrapper around Get, mostly useful when you don't need to worry if the link isn't actually there
(e.g. when printing debug stuff to a CLI output)




## <a name="Navigate">type</a> [Navigate](https://github.com/CreatorKit/go-deviceserver-client/blob/master/hateoas/hateoas.go#L51)
``` go
type Navigate []string
```
Navigate specifies a list of links to be traversed










## <a name="SimpleEndpoint">type</a> [SimpleEndpoint](https://github.com/CreatorKit/go-deviceserver-client/blob/master/hateoas/hateoas.go#L57)
``` go
type SimpleEndpoint struct {
    Links Links `json:"Links"`
}
```
SimpleEndpoint is used when navigating links on endpoints














- - -
Generated by [godoc12md](http://godoc.org/github.com/davecheney/godoc2md)
