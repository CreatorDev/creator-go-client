package deviceserver

import (
	"strconv"

	"github.com/CreatorKit/go-deviceserver-client/hateoas"
)

type EntryPoint struct {
	Links hateoas.Links `json:"Links"`
}

type AccessKey struct {
	Links  hateoas.Links `json:"Links"`
	Name   string        `json:"Name,omitempty"`
	Key    string        `json:"Key,omitempty"`
	Secret string        `json:"Secret,omitempty"`
}

type PageInfo struct {
	TotalCount int           `json:"TotalCount"`
	ItemsCount int           `json:"ItemsCount"`
	StartIndex int           `json:"StartIndex"`
	Links      hateoas.Links `json:"Links,omitempty"`
}

type AccessKeys struct {
	PageInfo PageInfo      `json:"PageInfo"`
	Items    []AccessKey   `json:"Items"`
	Links    hateoas.Links `json:"Links"`
}

type SubscriptionAttributes struct {
	Pmin        string `json:"Pmin,omitempty"`
	Pmax        string `json:"Pmax,omitempty"`
	Step        string `json:"Step,omitempty"`
	LessThan    string `json:"LessThan,omitempty"`
	GreaterThan string `json:"GreaterThan,omitempty"`
}

type SubscriptionRequest struct {
	SubscriptionType string `json:"SubscriptionType"`
	URL              string `json:"Url"`

	AcceptContentType string                  `json:"AcceptContentType,omitempty"`
	Property          string                  `json:"Property,omitempty"`
	Attributes        *SubscriptionAttributes `json:"Attributes,omitempty"`
	Links             hateoas.Links           `json:"Links"`
}

type SubscriptionResponse struct {
	ID    string        `json:"ID"`
	Links hateoas.Links `json:"Links"`
}

type Subscriptions struct {
	PageInfo PageInfo              `json:"PageInfo"`
	Items    []SubscriptionRequest `json:"Items"`
	Links    hateoas.Links         `json:"Links"`
}

type OrgClaim struct {
	OrgID int   `json:"OrgID"`
	Exp   int64 `json:"exp"`
}

type OAuthToken struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
}

type Error struct {
	ErrorCode    string `json:"ErrorCode"`
	ErrorMessage string `json:"ErrorMessage"`
	ErrorDetails string `json:"ErrorDetails"`
}

type Client struct {
	Name  string
	Links hateoas.Links `json:"Links"`
}

type Clients struct {
	PageInfo PageInfo      `json:"PageInfo"`
	Items    []Client      `json:"Items"`
	Links    hateoas.Links `json:"Links"`
}

type ObjectType struct {
	ObjectTypeID string        `json:"ObjectTypeID"`
	Links        hateoas.Links `json:"Links"`
}

type ObjectTypes struct {
	PageInfo PageInfo     `json:"PageInfo"`
	Items    []ObjectType `json:"Items"`
}

type ObjectInstance map[string]interface{}

func (i ObjectInstance) InstanceID() int {
	id, _ := strconv.Atoi(i["InstanceID"].(string))
	return id
}

func (i ObjectInstance) Links() *hateoas.Links {
	result := hateoas.Links{}
	links := i["Links"].([]interface{})
	for _, l := range links {
		ll := hateoas.Link{}
		lmap := l.(map[string]interface{})
		href := lmap["href"]
		rel := lmap["rel"]
		_type := lmap["type"]
		if href != nil {
			ll.Href = href.(string)
		}
		if rel != nil {
			ll.Rel = rel.(string)
		}
		if _type != nil {
			ll.Type = _type.(string)
		}
		result = append(result, ll)
	}
	return &result
}

type ObjectInstances struct {
	// see https://github.com/CreatorDev/DeviceServer/issues/25
	// PageInfo PageInfo         `json:"PageInfo"`
	Items []ObjectInstance `json:"Items"`
	Links hateoas.Links    `json:"Links"`
}
