package hateoas

import "fmt"

type Link struct {
	Rel  string `json:"rel"`
	Href string `json:"href"`
	Type string `json:"type"`
}

type Links []Link

func (l Links) GetLink(rel string) (*Link, error) {
	for _, link := range l {
		if link.Rel == rel {
			return &link, nil
		}
	}
	return nil, fmt.Errorf("unable to find rel '%s'", rel)
}
