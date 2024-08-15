package otlh

import "fmt"

type Group struct {
	ID          int    `json:"id,omitempty"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	Type        string `json:"type,omitempty"`
	CreatedAt   string `json:"created_at,omitempty"`
	UpdatedAt   string `json:"updated_at,omitempty"`
	Links       struct {
		Self struct {
			Href string `json:"href,omitempty"`
		} `json:"self,omitempty"`
		Users struct {
			Href string `json:"href,omitempty"`
		} `json:"users,omitempty"`
		Folders struct {
			Href string `json:"href,omitempty"`
		} `json:"folders,omitempty"`
		Site struct {
			Href string `json:"href,omitempty"`
		} `json:"site,omitempty"`
	} `json:"_links,omitempty"`
}

type Groups []Group

type GroupsResponse struct {
	Links struct {
		Self struct {
			Href string `json:"href"`
		} `json:"self"`
	} `json:"_links"`
	Page struct {
		HasMore    bool `json:"has-more"`
		TotalCount int  `json:"total-count"`
	} `json:"page"`
	Embedded struct {
		Groups []Group `json:"groups"`
	} `json:"_embedded"`
}

type GroupRequest struct {
	id int
	Request
}

type GroupRequestBuilder struct {
	*GroupRequest
}

func (b *GroupRequestBuilder) WithID(id int) *GroupRequestBuilder {
	b.id = id
	return b
}

func (b *GroupRequestBuilder) Build() (*GroupRequest, error) {
	return b.GroupRequest, nil
}

func (req *GroupRequest) Endpoint() string {
	if req.id == 0 {
		return fmt.Sprintf("/t/%s/api/%s/groups", req.tenant, APIVERSION)
	}
	return fmt.Sprintf("/t/%s/api/%s/group/%d", req.tenant, APIVERSION, req.id)
}
