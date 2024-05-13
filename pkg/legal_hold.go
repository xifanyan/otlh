package otlh

import "fmt"

type Legalhold struct {
	ID              int    `json:"id,omitempty"`
	MatterID        int    `json:"matter_id,omitempty"`
	Name            string `json:"name,omitempty"`
	HoldDescription string `json:"hold_description,omitempty"`
	Draft           bool   `json:"draft,omitempty"`
	ResponseDueDate string `json:"response_due_date,omitempty"`
	CreatedAt       string `json:"created_at,omitempty"`
	UpdatedAt       string `json:"updated_at,omitempty"`
	Status          string `json:"status,omitempty"`
	CreatedBy       struct {
		Name string `json:"name,omitempty"`
		Type string `json:"type,omitempty"`
	} `json:"created_by,omitempty"`
	CanSendAcknowledgementReminders struct {
		Code   int    `json:"code,omitempty"`
		Reason string `json:"reason,omitempty"`
	} `json:"can_send_acknowledgement_reminders,omitempty"`
	CanSendHoldReminders struct {
		Code   int    `json:"code,omitempty"`
		Reason string `json:"reason,omitempty"`
	} `json:"can_send_hold_reminders,omitempty"`
	Links struct {
		Self struct {
			Href string `json:"href,omitempty"`
		} `json:"self,omitempty"`
		Custodians struct {
			Href string `json:"href,omitempty"`
		} `json:"custodians,omitempty"`
		HoldNotice struct {
			Href string `json:"href,omitempty"`
		} `json:"hold_notice,omitempty"`
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

type Legalholds []Legalhold

type LegalholdsResponse struct {
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
		Legalholds Legalholds `json:"legal_holds"`
	} `json:"_embedded"`
}

type LegalholdAction int

const (
	IMPORT LegalholdAction = iota + 1
	SEND_NOTICE
)

type LegalholdRequestBuilder struct {
	*LegalholdRequest
}

type LegalholdRequest struct {
	id          int
	custodianID int
	action      LegalholdAction
	Request
}

func (b *LegalholdRequestBuilder) WithID(id int) *LegalholdRequestBuilder {
	b.id = id
	return b
}

func (b *LegalholdRequestBuilder) WithCustodianID(custodianID int) *LegalholdRequestBuilder {
	b.custodianID = custodianID
	return b
}

func (b *LegalholdRequestBuilder) Import() *LegalholdRequestBuilder {
	b.action = IMPORT
	return b
}

func (b *LegalholdRequestBuilder) SendNotice() *LegalholdRequestBuilder {
	b.action = SEND_NOTICE
	return b
}

func (b *LegalholdRequestBuilder) Build() (*LegalholdRequest, error) {
	return b.LegalholdRequest, nil
}

func (req *LegalholdRequest) Endpoint() string {
	switch req.action {
	case IMPORT:
		return fmt.Sprintf("/t/%s/api/%s/legal_holds/import", req.tenant, APIVERSION)
	case SEND_NOTICE:
		return fmt.Sprintf("/t/%s/api/%s/legal_holds/send_notice", req.tenant, APIVERSION)
	}

	if req.id > 0 {
		return fmt.Sprintf("/t/%s/api/%s/legal_holds/%d", req.tenant, APIVERSION, req.id)
	}
	return fmt.Sprintf("/t/%s/api/%s/legal_holds", req.tenant, APIVERSION)
}
