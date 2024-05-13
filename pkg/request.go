package otlh

type Method int

const (
	GET Method = iota
	POST
	PATCH
)

type Requestor interface {
	Method() Method
	Endpoint() string
}

type Request struct {
	method Method
	tenant string
}

func NewRequest() *Request {
	return &Request{}
}

func (req *Request) WithTenant(tenant string) *Request {
	req.tenant = tenant
	return req
}

func (req *Request) Method() Method {
	return req.method
}

func (req *Request) Get() *Request {
	req.method = GET
	return req
}

func (req *Request) Post() *Request {
	req.method = POST
	return req
}

func (req *Request) Custodian() *CustodianRequestBuilder {
	return &CustodianRequestBuilder{CustodianRequest: &CustodianRequest{Request: *req}}
}

func (req *Request) Legalhold() *LegalholdRequestBuilder {
	return &LegalholdRequestBuilder{LegalholdRequest: &LegalholdRequest{Request: *req}}
}

func (req *Request) Folder() *FolderRequestBuilder {
	return &FolderRequestBuilder{FolderRequest: &FolderRequest{Request: *req}}
}

func (req *Request) Matter() *MatterRequestBuilder {
	return &MatterRequestBuilder{MatterRequest: &MatterRequest{Request: *req}}
}