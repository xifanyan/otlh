package otlh

import "fmt"

type CustodianGroup struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	CreatedBy struct {
		Name string `json:"name"`
		Type string `json:"type"`
	} `json:"created_by"`
	UpdatedBy struct {
		Name string `json:"name"`
		Type string `json:"type"`
	} `json:"updated_by"`
	CustodiansCount int    `json:"custodians_count"`
	CreatedAt       string `json:"created_at"`
	UpdatedAt       string `json:"updated_at"`
	Links           struct {
		Self struct {
			Href string `json:"href"`
		} `json:"self"`
		Custodians struct {
			Href string `json:"href"`
		} `json:"custodians"`
		Site struct {
			Href string `json:"href"`
		} `json:"site"`
	} `json:"_links"`
}

type CustodianGroupsResponse struct {
	DefaultEntityListInfo
	Embedded struct {
		CustodianGroups []CustodianGroup `json:"custodian_groups"`
	} `json:"_embedded"`
}

type CustodianGroups []CustodianGroup

type CustodianGroupRequestBuilder struct {
	*CustodianGroupRequest
}

type CustodianGroupRequest struct {
	id          int
	custodianID int
	Request
}

func (b *CustodianGroupRequestBuilder) WithID(id int) *CustodianGroupRequestBuilder {
	b.id = id
	return b
}

func (b *CustodianGroupRequestBuilder) WithCustodianID(id int) *CustodianGroupRequestBuilder {
	b.custodianID = id
	return b
}

func (b *CustodianGroupRequestBuilder) Build() (*CustodianGroupRequest, error) {
	return b.CustodianGroupRequest, nil
}

func (req *CustodianGroupRequest) Endpoint() string {
	if req.id == 0 {
		if req.custodianID > 0 {
			// retrieves custodians under a custodian group: /t/{tenant}/api/{version}/custodian_groups/{id}/custodians
			return fmt.Sprintf("/t/%s/api/%s/custodians/%d/custodian_groups", req.tenant, APIVERSION, req.custodianID)
		}
		return fmt.Sprintf("/t/%s/api/%s/custodian_groups", req.tenant, APIVERSION)
	}
	return fmt.Sprintf("/t/%s/api/%s/custodian_groups/%d", req.tenant, APIVERSION, req.id)
}
