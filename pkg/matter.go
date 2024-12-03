package otlh

import "fmt"

type Matter struct {
	ID                   int    `json:"id,omitempty"`
	Name                 string `json:"name,omitempty"`
	Number               any    `json:"number,omitempty"`
	Notes                any    `json:"notes,omitempty"`
	Caption              any    `json:"caption,omitempty"`
	PoNumber             any    `json:"po_number,omitempty"`
	CaseNumber           any    `json:"case_number,omitempty"`
	InheritEmailConfig   bool   `json:"inherit_email_config,omitempty"`
	EmailFrom            string `json:"email_from,omitempty"`
	NameOnOutgoingEmails string `json:"name_on_outgoing_emails,omitempty"`
	EmailReplyTo         string `json:"email_reply_to,omitempty"`
	Region               any    `json:"region,omitempty"`
	BusinessUnit         any    `json:"business_unit,omitempty"`
	CreatedAt            string `json:"created_at,omitempty"`
	UpdatedAt            string `json:"updated_at,omitempty"`
	CanBeDeleted         bool   `json:"can_be_deleted,omitempty"`
	MatterContacts       []any  `json:"matter_contacts,omitempty"`
	CreatedBy            struct {
		Name  string `json:"name,omitempty"`
		Email string `json:"email,omitempty"`
		Type  string `json:"type,omitempty"`
	} `json:"created_by,omitempty"`
	Links struct {
		Self struct {
			Href string `json:"href,omitempty"`
		} `json:"self,omitempty"`
		LegalHolds struct {
			Href string `json:"href,omitempty"`
		} `json:"legal_holds,omitempty"`
		Custodians struct {
			Href string `json:"href,omitempty"`
		} `json:"custodians,omitempty"`
		Stats struct {
			Href string `json:"href,omitempty"`
		} `json:"stats,omitempty"`
		Folder struct {
			Href string `json:"href,omitempty"`
		} `json:"folder,omitempty"`
	} `json:"_links,omitempty"`
}

type Matters []Matter

type MatterRequest struct {
	id int
	Request
}

type MatterRequestBuilder struct {
	*MatterRequest
}

func (b *MatterRequestBuilder) WithID(id int) *MatterRequestBuilder {
	b.id = id
	return b
}

func (b *MatterRequestBuilder) Build() (*MatterRequest, error) {
	return b.MatterRequest, nil
}

func (req *MatterRequest) Endpoint() string {
	if req.id == 0 {
		return fmt.Sprintf("/t/%s/api/%s/matters", req.tenant, APIVERSION)
	}
	return fmt.Sprintf("/t/%s/api/%s/matters/%d", req.tenant, APIVERSION, req.id)
}

type MattersResponse struct {
	DefaultEntityListInfo
	Embedded struct {
		Matters []Matter `json:"matters"`
	} `json:"_embedded"`
}

type CreateMatterBody struct {
	Name                     string `json:"name"`
	Number                   string `json:"number,omitempty"`
	FolderID                 int    `json:"folder_id"`
	Notes                    string `json:"notes,omitempty"`
	Caption                  string `json:"caption,omitempty"`
	PoNumber                 string `json:"po_number,omitempty"`
	CaseNumber               string `json:"case_number,omitempty"`
	InheritEmailConfig       bool   `json:"inherit_email_config"`
	EmailFrom                string `json:"email_from,omitempty"`
	NameOnOutgoingEmails     string `json:"name_on_outgoing_emails,omitempty"`
	EmailReplyTo             string `json:"email_reply_to,omitempty"`
	Region                   string `json:"region,omitempty"`
	BusinessUnit             string `json:"business_unit,omitempty"`
	CustodianIds             []int  `json:"custodian_ids,omitempty"`
	MatterContactsAttributes []struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	} `json:"matter_contacts_attributes"`
}

func NewCreateMatterBody() *CreateMatterBody {
	return &CreateMatterBody{
		InheritEmailConfig: true,
		MatterContactsAttributes: []struct {
			Name  string `json:"name"`
			Email string `json:"email"`
		}{},
	}
}

func (b *CreateMatterBody) WithName(name string) *CreateMatterBody {
	b.Name = name
	return b
}

func (b *CreateMatterBody) WithFolderID(id int) *CreateMatterBody {
	b.FolderID = id
	return b
}

type MatterContactsAttribute struct {
	Name  string
	Email string
}

func (b *CreateMatterBody) WithMatterContactsAttributes(attr MatterContactsAttribute) *CreateMatterBody {
	b.MatterContactsAttributes = []struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}{
		{Name: attr.Name, Email: attr.Email},
	}

	return b
}
