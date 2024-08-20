package otlh

import "fmt"

type Silenthold struct {
	ID              int    `json:"id,omitempty"`
	MatterID        int    `json:"matter_id,omitempty"`
	Name            string `json:"name,omitempty"`
	HoldDescription string `json:"hold_description,omitempty"`
	Draft           bool   `json:"draft,omitempty"`
	CreatedAt       string `json:"created_at,omitempty"`
	UpdatedAt       string `json:"updated_at,omitempty"`
	Status          string `json:"status,omitempty"`
	CreatedBy       struct {
		Name string `json:"name,omitempty"`
		Type string `json:"type,omitempty"`
	} `json:"created_by,omitempty"`
	ApprovalDetails struct {
		Enabled      bool `json:"enabled,omitempty"`
		Status       any  `json:"status,omitempty"`
		Requester    any  `json:"requester,omitempty"`
		RequestedAt  any  `json:"requested_at,omitempty"`
		Notes        any  `json:"notes,omitempty"`
		Comments     any  `json:"comments,omitempty"`
		LastApprover any  `json:"last_approver,omitempty"`
		RespondedAt  any  `json:"responded_at,omitempty"`
	} `json:"approval_details,omitempty"`
	Links struct {
		Self struct {
			Href string `json:"href,omitempty"`
		} `json:"self,omitempty"`
		Custodians struct {
			Href string `json:"href,omitempty"`
		} `json:"custodians,omitempty"`
		AdvisoryNotice struct {
			Href string `json:"href,omitempty"`
		} `json:"advisory_notice,omitempty"`
		ReleaseNotice struct {
			Href string `json:"href,omitempty"`
		} `json:"release_notice,omitempty"`
		AdvisoryCopies struct {
			Href string `json:"href,omitempty"`
		} `json:"advisory_copies,omitempty"`
		History struct {
			Href string `json:"href,omitempty"`
		} `json:"history,omitempty"`
		Stats struct {
			Href string `json:"href,omitempty"`
		} `json:"stats,omitempty"`
		Matter struct {
			Href string `json:"href,omitempty"`
		} `json:"matter,omitempty"`
	} `json:"_links,omitempty"`
}

type Silentholds []Silenthold

type SilentholdsResponse struct {
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
		Silentholds Silentholds `json:"silent_holds"`
	} `json:"_embedded"`
}

type SilentholdRequestBuilder struct {
	*SilentholdRequest
}

type SilentholdRequest struct {
	id int
	Request
}

func (b *SilentholdRequestBuilder) WithID(id int) *SilentholdRequestBuilder {
	b.id = id
	return b
}

func (b *SilentholdRequestBuilder) Build() (*SilentholdRequest, error) {
	return b.SilentholdRequest, nil
}

func (req *SilentholdRequest) Endpoint() string {
	if req.id > 0 {
		return fmt.Sprintf("/t/%s/api/%s/silent_holds/%d", req.tenant, APIVERSION, req.id)
	}
	return fmt.Sprintf("/t/%s/api/%s/silent_holds", req.tenant, APIVERSION)
}
