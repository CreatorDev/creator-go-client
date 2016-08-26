

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



## <a name="Client">type</a> [Client](https://github.com/CreatorKit/go-deviceserver-client/blob/master/hateoas/hateoas.go#L59)
``` go
type Client struct {
    EntryURL       string
    DefaultHeaders Headers
    Http           HTTPDoer
}
```






### <a name="Create">func</a> [Create](https://github.com/CreatorKit/go-deviceserver-client/blob/master/hateoas/hateoas.go#L65)
``` go
func Create(config *Client) *Client
```




### <a name="Client.Delete">func</a> (\*Client) [Delete](https://github.com/CreatorKit/go-deviceserver-client/blob/master/hateoas/hateoas.go#L159)
``` go
func (c *Client) Delete(url string, navigateLinks Navigate, headers Headers, body io.Reader, response interface{}) (*http.Response, error)
```



### <a name="Client.Do">func</a> (\*Client) [Do](https://github.com/CreatorKit/go-deviceserver-client/blob/master/hateoas/hateoas.go#L82)
``` go
func (c *Client) Do(method string, url string, navigateLinks Navigate, headers Headers, body io.Reader, result interface{}) (*http.Response, error)
```



### <a name="Client.Get">func</a> (\*Client) [Get](https://github.com/CreatorKit/go-deviceserver-client/blob/master/hateoas/hateoas.go#L143)
``` go
func (c *Client) Get(url string, navigateLinks Navigate, headers Headers, body io.Reader, response interface{}) (*http.Response, error)
```



### <a name="Client.Post">func</a> (\*Client) [Post](https://github.com/CreatorKit/go-deviceserver-client/blob/master/hateoas/hateoas.go#L147)
``` go
func (c *Client) Post(url string, navigateLinks Navigate, headers Headers, body io.Reader, response interface{}) (*http.Response, error)
```



### <a name="Client.PostForm">func</a> (\*Client) [PostForm](https://github.com/CreatorKit/go-deviceserver-client/blob/master/hateoas/hateoas.go#L151)
``` go
func (c *Client) PostForm(url string, navigateLinks Navigate, headers Headers, data url.Values, response interface{}) (*http.Response, error)
```



## <a name="HTTPDoer">type</a> [HTTPDoer](https://github.com/CreatorKit/go-deviceserver-client/blob/master/hateoas/hateoas.go#L55)
``` go
type HTTPDoer interface {
    Do(*http.Request) (*http.Response, error)
}
```
HTTPDoer would normally be an instance of *http.Client but
by making it an interface we allow wrapping to permit
logging or caching etc










## <a name="Headers">type</a> [Headers](https://github.com/CreatorKit/go-deviceserver-client/blob/master/hateoas/hateoas.go#L46)
``` go
type Headers map[string]string
```









## <a name="Link">type</a> [Link](https://github.com/CreatorKit/go-deviceserver-client/blob/master/hateoas/hateoas.go#L20)
``` go
type Link struct {
    Rel  string `json:"rel"`
    Href string `json:"href"`
    Type string `json:"type"`
}
```









## <a name="Links">type</a> [Links](https://github.com/CreatorKit/go-deviceserver-client/blob/master/hateoas/hateoas.go#L26)
``` go
type Links []Link
```









### <a name="Links.Get">func</a> (Links) [Get](https://github.com/CreatorKit/go-deviceserver-client/blob/master/hateoas/hateoas.go#L28)
``` go
func (l Links) Get(rel string) (*Link, error)
```



### <a name="Links.Self">func</a> (Links) [Self](https://github.com/CreatorKit/go-deviceserver-client/blob/master/hateoas/hateoas.go#L37)
``` go
func (l Links) Self() string
```



## <a name="Navigate">type</a> [Navigate](https://github.com/CreatorKit/go-deviceserver-client/blob/master/hateoas/hateoas.go#L45)
``` go
type Navigate []string
```









## <a name="SimpleEndpoint">type</a> [SimpleEndpoint](https://github.com/CreatorKit/go-deviceserver-client/blob/master/hateoas/hateoas.go#L48)
``` go
type SimpleEndpoint struct {
    Links Links `json:"Links"`
}
```













- - -
Generated by [godoc12md](http://godoc.org/github.com/davecheney/godoc2md)
