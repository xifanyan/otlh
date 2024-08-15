package otlh

import "fmt"

type Folder struct {
	ID                                                 int    `json:"id,omitempty"`
	Name                                               string `json:"name,omitempty"`
	Number                                             any    `json:"number,omitempty"`
	ContactName                                        any    `json:"contact_name,omitempty"`
	ContactEmail                                       any    `json:"contact_email,omitempty"`
	ContactPhone                                       any    `json:"contact_phone,omitempty"`
	Notes                                              any    `json:"notes,omitempty"`
	Address1                                           any    `json:"address_1,omitempty"`
	Address2                                           any    `json:"address_2,omitempty"`
	City                                               any    `json:"city,omitempty"`
	State                                              any    `json:"state,omitempty"`
	Zip                                                any    `json:"zip,omitempty"`
	CreatedAt                                          string `json:"created_at,omitempty"`
	UpdatedAt                                          string `json:"updated_at,omitempty"`
	NameOnOutgoingEmails                               string `json:"name_on_outgoing_emails,omitempty"`
	EmailReplyTo                                       string `json:"email_reply_to,omitempty"`
	CombineHoldReminders                               bool   `json:"combine_hold_reminders,omitempty"`
	IncludeActiveLegalHoldsInReleaseNotice             string `json:"include_active_legal_holds_in_release_notice,omitempty"`
	EmailFrom                                          string `json:"email_from,omitempty"`
	InheritEmailConfig                                 bool   `json:"inherit_email_config,omitempty"`
	IncludeActiveLegalHoldsInReleaseNoticeInheritValue string `json:"include_active_legal_holds_in_release_notice_inherit_value,omitempty"`
	CreatedBy                                          struct {
		Name  string `json:"name,omitempty"`
		Type  string `json:"type,omitempty"`
		Email string `json:"email,omitempty"`
	} `json:"created_by,omitempty"`
	CanBeDeleted bool `json:"can_be_deleted,omitempty"`
	Links        struct {
		Self struct {
			Href string `json:"href,omitempty"`
		} `json:"self,omitempty"`
		Matters struct {
			Href string `json:"href,omitempty"`
		} `json:"matters,omitempty"`
		Groups struct {
			Href string `json:"href,omitempty"`
		} `json:"groups,omitempty"`
		Site struct {
			Href string `json:"href,omitempty"`
		} `json:"site,omitempty"`
		Stats struct {
			Href string `json:"href,omitempty"`
		} `json:"stats,omitempty"`
	} `json:"_links,omitempty"`
}

type Folders []Folder

type FoldersResponse struct {
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
		Folders []Folder `json:"folders"`
	} `json:"_embedded"`
}

type FolderRequest struct {
	id int
	Request
}

type FolderRequestBuilder struct {
	*FolderRequest
}

func (b *FolderRequestBuilder) WithID(id int) *FolderRequestBuilder {
	b.id = id
	return b
}

func (b *FolderRequestBuilder) Build() (*FolderRequest, error) {
	return b.FolderRequest, nil
}

func (req *FolderRequest) Endpoint() string {
	if req.id == 0 {
		return fmt.Sprintf("/t/%s/api/%s/folders", req.tenant, APIVERSION)
	}
	return fmt.Sprintf("/t/%s/api/%s/folders/%d", req.tenant, APIVERSION, req.id)
}

// minimum requirement for creating a legal hold: name, inheritEmailConfig
type CreateFolderBody struct {
	Name                                   string `json:"name"`
	Number                                 string `json:"number,omitempty"`
	Address1                               string `json:"address_1,omitempty"`
	Address2                               string `json:"address_2,omitempty"`
	City                                   string `json:"city,omitempty"`
	State                                  string `json:"state,omitempty"`
	Zip                                    string `json:"zip,omitempty"`
	InheritEmailConfig                     bool   `json:"inherit_email_config"`
	EmailFrom                              string `json:"email_from,omitempty"`
	NameOnOutgoingEmails                   string `json:"name_on_outgoing_emails,omitempty"`
	EmailReplyTo                           string `json:"email_reply_to,omitempty"`
	IncludeActiveLegalHoldsInReleaseNotice string `json:"include_active_legal_holds_in_release_notice,omitempty"`
	ContactName                            string `json:"contact_name,omitempty"`
	ContactEmail                           string `json:"contact_email,omitempty"`
	ContactPhone                           string `json:"contact_phone,omitempty"`
	Notes                                  string `json:"notes,omitempty"`
	GroupIDs                               []int  `json:"group_ids,omitempty"`
}

// TODO: enable group name support
func NewCreateFolderBody() *CreateFolderBody {
	return &CreateFolderBody{
		InheritEmailConfig: true,
	}
}

func (b *CreateFolderBody) WithName(name string) *CreateFolderBody {
	b.Name = name
	return b
}

func (b *CreateFolderBody) WithGroupIDs(groupIDs []int) *CreateFolderBody {
	b.GroupIDs = groupIDs
	return b
}
