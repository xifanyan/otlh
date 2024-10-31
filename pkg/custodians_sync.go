package otlh

import (
	"fmt"
)

type CustodianSyncBody struct {
	Custodians []CustodianInputData `json:"custodians"`
}

type CustodianInputData struct {
	Name            string `json:"name" csv:"name"`
	Email           string `json:"email" csv:"email"`
	Phone           string `json:"phone,omitempty" csv:"phone"`
	Notes           string `json:"notes,omitempty" csv:"notes"`
	Title           string `json:"title,omitempty" csv:"title"`
	EmployeeID      string `json:"employee_id,omitempty" csv:"employee_id"`
	EmployeeType    string `json:"employee_type,omitempty" csv:"employee_type"`
	EmployeeStatus  string `json:"employee_status,omitempty" csv:"employee_status"`
	Department      string `json:"department,omitempty" csv:"depoartment"`
	Location        string `json:"location,omitempty" csv:"location"`
	SupervisorEmail string `json:"supervisor_email,omitempty" csv:"supervisor_email"`
	SupervisorName  string `json:"supervisor_name,omitempty" csv:"supervisor_name"`
	Function        string `json:"function,omitempty" csv:"function"`
	Business        string `json:"business,omitempty" csv:"business"`
	Country         string `json:"country,omitempty" csv:"country"`
	DelegateEmail   string `json:"delegate_email,omitempty" csv:"delegate_email"`
	DelegateName    string `json:"delegate_name,omitempty" csv:"delegate_name"`
}

type CustodiansSyncRequest struct {
	matterID int
	Request
}

type CustodiansSyncRequestBuilder struct {
	*CustodiansSyncRequest
}

func (b *CustodiansSyncRequestBuilder) WithTenant(tenant string) *CustodiansSyncRequestBuilder {
	b.CustodiansSyncRequest.tenant = tenant
	return b
}

func (b *CustodiansSyncRequestBuilder) WithMatterID(matterID int) *CustodiansSyncRequestBuilder {
	b.CustodiansSyncRequest.matterID = matterID
	return b
}

func (b *CustodiansSyncRequestBuilder) Build() *CustodiansSyncRequest {
	return b.CustodiansSyncRequest
}

func (req *CustodiansSyncRequest) Endpoint() string {
	if req.matterID > 0 {
		return fmt.Sprintf("/t/%s/api/%s/matters/%d/custodians/import", req.tenant, APIVERSION, req.matterID)
	}
	return fmt.Sprintf("/t/%s/api/%s/custodians/import", req.tenant, APIVERSION)
}

type CustodiansSyncResponse struct {
	Error  any `json:"error"`
	Errors any `json:"errors,omitempty"`
}
