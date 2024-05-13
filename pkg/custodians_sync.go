package otlh

import (
	"fmt"
)

type CustodiansSyncRequest struct {
	tenant   string
	matterID int
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
	return fmt.Sprintf("/t/%s/api/%s/custodians/sync", req.tenant, APIVERSION)
}
