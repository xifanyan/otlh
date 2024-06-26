package otlh

import "fmt"

/*
Custodian represents a custodian record in the system. It contains identifying
information about the custodian as well as metadata about their status.
*/
type Custodian struct {
	ID                      int    `json:"id" csv:"id"`
	Name                    string `json:"name" csv:"name"`
	Email                   string `json:"email" csv:"email"`
	Synced                  bool   `json:"synced,omitempty" csv:"synced"`
	Phone                   string `json:"phone,omitempty" csv:"phone"`
	Notes                   string `json:"notes,omitempty" csv:"notes"`
	Title                   string `json:"title,omitempty" csv:"title"`
	EmployeeID              string `json:"employee_id,omitempty" csv:"employee_id"`
	EmployeeType            string `json:"employee_type,omitempty" csv:"employee_type"`
	EmployeeStatus          string `json:"employee_status,omitempty" csv:"employee_status"`
	EmployeeStatusChangedAt string `json:"employee_status_changed_at,omitempty" csv:"employee_status_changed_at"`
	Department              string `json:"department,omitempty" csv:"depoartment"`
	Location                string `json:"location,omitempty" csv:"location"`
	SupervisorEmail         string `json:"supervisor_email,omitempty" csv:"supervisor_email"`
	SupervisorName          string `json:"supervisor_name,omitempty" csv:"supervisor_name"`
	Function                string `json:"function,omitempty" csv:"function"`
	Business                string `json:"business,omitempty" csv:"business"`
	Country                 string `json:"country,omitempty" csv:"country"`
	DelegateEmail           string `json:"delegate_email,omitempty" csv:"delegate_email"`
	DelegateName            string `json:"delegate_name,omitempty" csv:"delegate_name"`
	CreatedAt               string `json:"created_at,omitempty" csv:"created_at"`
	UpdatedAt               string `json:"updated_at,omitempty" csv:"updated_at"`
	Links                   struct {
		Self struct {
			Href string `json:"href,omitempty" csv:"-"`
		} `json:"self,omitempty" csv:"-"`
		LegalHolds struct {
			Href string `json:"href,omitempty" csv:"-"`
		} `json:"legal_holds,omitempty" csv:"-"`
		Matters struct {
			Href string `json:"href,omitempty" csv:"-"`
		} `json:"matters,omitempty"`
		Stats struct {
			Href string `json:"href,omitempty" csv:"-"`
		} `json:"stats,omitempty" csv:"-"`
	} `json:"_links,omitempty" csv:"-"`
}

type Custodians []Custodian

type CustodiansResponse struct {
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
		Custodians Custodians `json:"custodians"`
	} `json:"_embedded"`
}

type CustodianRequestBuilder struct {
	*CustodianRequest
}

type CustodianRequest struct {
	id                int
	matterID          int
	legalHoldID       int
	silentHoldID      int
	custodianGroupsID int
	Request
}

func (b *CustodianRequestBuilder) WithID(id int) *CustodianRequestBuilder {
	b.id = id
	return b
}

func (b *CustodianRequestBuilder) WithMatterID(id int) *CustodianRequestBuilder {
	b.matterID = id
	return b
}

func (b *CustodianRequestBuilder) WithLegalHoldID(id int) *CustodianRequestBuilder {
	b.legalHoldID = id
	return b
}

func (b *CustodianRequestBuilder) WithCustodianGroupsID(id int) *CustodianRequestBuilder {
	b.custodianGroupsID = id
	return b
}

func (b *CustodianRequestBuilder) WithSilentHoldID(id int) *CustodianRequestBuilder {
	b.silentHoldID = id
	return b
}

func (b *CustodianRequestBuilder) Build() (*CustodianRequest, error) {
	return b.CustodianRequest, nil
}

func (req *CustodianRequest) Endpoint() string {
	if req.id == 0 {
		if req.matterID > 0 {
			// retrieves custodians under a matter: /t/{tenant}/api/{version}/matters/{id}/custodians
			return fmt.Sprintf("/t/%s/api/%s/matters/%d/custodians", req.tenant, APIVERSION, req.matterID)
		}

		if req.legalHoldID > 0 {
			// retrieves custodians under a legal hold: /t/{tenant}/api/{version}/legal_holds/{id}/custodians
			return fmt.Sprintf("/t/%s/api/%s/legal_holds/%d/custodians", req.tenant, APIVERSION, req.legalHoldID)
		}

		if req.custodianGroupsID > 0 {
			// retrieves custodians under a custodian group: /t/{tenant}/api/{version}/custodian_groups/{id}/custodians
			return fmt.Sprintf("/t/%s/api/%s/custodian_groups/%d/custodians", req.tenant, APIVERSION, req.custodianGroupsID)
		}

		if req.silentHoldID > 0 {
			// retrieves custodians under a silent hold: /t/{tenant}/api/{version}/silent_holds/{id}/custodians
			return fmt.Sprintf("/t/%s/api/%s/silent_holds/%d/custodians", req.tenant, APIVERSION, req.silentHoldID)
		}

		// list all custodians: /t/{tenant}/api/{version}/custodians
		return fmt.Sprintf("/t/%s/api/%s/custodians", req.tenant, APIVERSION)
	}

	// retrivees a single custodian: /t/{tenant}/api/{version}/custodians/{id}
	return fmt.Sprintf("/t/%s/api/%s/custodians/%d", req.tenant, APIVERSION, req.id)
}
