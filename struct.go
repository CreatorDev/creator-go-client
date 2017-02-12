package deviceserver

import (
	"fmt"
	"net/url"
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

type ObjectDefinitionProperty struct {
	PropertyDefinitionID string `json:"PropertyDefinitionID"`
	PropertyID           string `json:"PropertyID"`
	Name                 string `json:"Name"`
	Description          string `json:"Description"`
	DataType             string `json:"DataType"`
	Units                string `json:"Units"`
	IsCollection         bool   `json:"IsCollection"`
	IsMandatory          bool   `json:"IsMandatory"`
	Access               string `json:"Access"`
	SerialisationName    string `json:"SerialisationName"`
}

type ObjectDefinitionProperties []ObjectDefinitionProperty

func (p ObjectDefinitionProperties) Get(nameOrID string) *ObjectDefinitionProperty {
	for _, pp := range p {
		if pp.PropertyID == nameOrID || pp.SerialisationName == nameOrID {
			return &pp
		}
	}
	return nil
}
func (p ObjectDefinitionProperties) String() string {
	s := "[\n"
	for _, pp := range p {
		s += fmt.Sprintf("%+v\n", pp)
	}
	s += "]\n"
	return s
}

type ObjectDefinition struct {
	ObjectDefinitionID string                     `json:"ObjectDefinitionID"`
	ObjectID           string                     `json:"ObjectID"`
	Name               string                     `json:"Name"`
	MIMEType           string                     `json:"MIMEType"`
	Description        string                     `json:"Description"`
	SerialisationName  string                     `json:"SerialisationName"`
	Singleton          bool                       `json:"Singleton"`
	Properties         ObjectDefinitionProperties `json:"Properties"`
	Links              hateoas.Links              `json:"Links"`
}

type ObjectDefinitionRegistry struct {
	href map[string]*ObjectDefinition
	id   map[int]*ObjectDefinition
}

func CreateObjectDefinitionRegistry() *ObjectDefinitionRegistry {
	r := ObjectDefinitionRegistry{}
	r.href = make(map[string]*ObjectDefinition)
	r.id = make(map[int]*ObjectDefinition)
	return &r
}

func (r *ObjectDefinitionRegistry) Set(href string, def *ObjectDefinition) {
	u, err := url.Parse(href)
	if err != nil {
		return
	}
	if def != nil {
		id, err := strconv.Atoi(def.ObjectID)
		if err != nil {
			return
		}

		r.href[u.Path] = def
		r.id[id] = def
	}
}

func (r *ObjectDefinitionRegistry) GetByHref(href string) *ObjectDefinition {
	u, err := url.Parse(href)
	if err != nil {
		return nil
	}
	def, exists := r.href[u.Path]
	if !exists {
		return nil
	}
	return def
}

func (r *ObjectDefinitionRegistry) GetByID(id int) *ObjectDefinition {
	def, exists := r.id[id]
	if !exists {
		return nil
	}
	return def
}

type WebhookItem struct {
	SubscriptionType string                 `json:"SubscriptionType"`
	TimeTriggered    string                 `json:"TimeTriggered"`
	Value            map[string]interface{} `json:"Value"`
	Links            hateoas.Links          `json:"Links"`
}

type Webhook struct {
	Items []WebhookItem `json:"Items"`
}
