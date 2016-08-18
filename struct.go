package deviceserver

import "gitlab.flowcloud.systems/creator-ops/go-deviceserver-client/hateoas"

type EntryPoint struct {
	Links hateoas.Links `json:"Links"`
}

type AccessKey struct {
	Links  hateoas.Links `json:"Links"`
	Name   string        `json:"Name,omitempty"`
	Key    string        `json:"Key,omitempty"`
	Secret string        `json:"Secret,omitempty"`
}

type SubscriptionRequest struct {
	SubscriptionType string `json:"SubscriptionType"`
	URL              string `json:"Url"`

	AcceptContentType string `json:"AcceptContentType,omitempty"`
	Property          string `json:"Property,omitempty"`
	Attributes        struct {
		Pmin        string `json:"Pmin,omitempty"`
		Pmax        string `json:"Pmax,omitempty"`
		Step        string `json:"Step,omitempty"`
		LessThan    string `json:"LessThan,omitempty"`
		GreaterThan string `json:"GreaterThan,omitempty"`
	} `json:"Attributes,omitempty"`
}

type SubscriptionResponse struct {
	ID    string        `json:"ID"`
	Links hateoas.Links `json:"Links"`
}

type OrgClaim struct {
	OrgID int   `json:"OrgID"`
	Exp   int64 `json:"exp"`
}

type Error struct {
	ErrorCode    string `json:"ErrorCode"`
	ErrorMessage string `json:"ErrorMessage"`
	ErrorDetails string `json:"ErrorDetails"`
}
