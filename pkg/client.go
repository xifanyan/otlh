package otlh

import (
	"crypto/tls"
	"encoding/json"
	"fmt"

	"github.com/go-resty/resty/v2"
	"github.com/rs/zerolog/log"
)

const (
	DEFAULT_DOMAIN = "api.otlegalhold.com"
	DEFAULT_PORT   = 443
	APIVERSION     = "v3"
)

/*
ClientBuilder is a builder struct for creating an instance of Client.
It allows for a more fluent way of setting fields on the Client.
*/
type ClientBuilder struct {
	*Client
}

/*
Client represents a client for the OpenText Legal Hold API. It contains
fields that define how to connect to the API as well as a resty client
to perform actual requests.
*/
type Client struct {
	// skipVerify skips verifying the server's TLS certificate. This should
	// only be used for testing or for APIs that do not have a valid TLS
	// certificate.
	skipVerify bool

	// domain is the domain to connect to.
	domain string

	// port is the port to connect to.
	port int

	// tenant name
	tenant string

	// AuthToken is the auth token to use when connecting to the API.
	authToken string

	// RestyClient is the resty client used to perform requests.
	RestyClient *resty.Client
}

func NewClientBuilder() *ClientBuilder {
	return &ClientBuilder{&Client{
		skipVerify: false,
		domain:     "localhost",
		port:       DEFAULT_PORT,
		authToken:  "",
	}}
}

func (b *ClientBuilder) WithDomain(domain string) *ClientBuilder {
	b.domain = domain
	return b
}

func (b *ClientBuilder) WithPort(port int) *ClientBuilder {
	b.port = port
	return b
}

func (b *ClientBuilder) WithTenant(tenant string) *ClientBuilder {
	b.tenant = tenant
	return b
}

func (b *ClientBuilder) WithAuthToken(authToken string) *ClientBuilder {
	b.authToken = authToken
	return b
}

func (b *ClientBuilder) SkipVerify() *ClientBuilder {
	b.skipVerify = true
	return b
}

/*
Build constructs a new OpenText Legal Hold API client from a ClientBuilder.

Parameters:
- b: ClientBuilder instance

Returns:
- *Client: Constructed client instance

This function builds a resty client with the base URL and headers set and then
sets the RestyClient field on the ClientBuilder. It returns the Client field
of the ClientBuilder.
*/
func (b *ClientBuilder) Build() *Client {
	r := resty.New().
		SetBaseURL(fmt.Sprintf("https://%s:%d", b.domain, b.port)).
		SetHeaders(map[string]string{
			"accept":       "application/json",
			"X-AUTH-TOKEN": b.authToken,
			"Content-Type": "application/json",
		})

	// r.SetDebug(true)

	if b.skipVerify {
		r.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	}

	b.RestyClient = r

	return b.Client
}

func handleOptions(r *resty.Request, opts ...Options) (bool, error) {
	var isMultipart bool
	for _, opt := range opts {
		switch opt.optionType() {
		case FILE:
			isMultipart = true
			r.SetFiles(opt.options())
		case QUERYPARAM:
			r.SetQueryParams(opt.options())
		case BODY:
			r.SetBody(opt.options()["body"])
		}
	}
	return isMultipart, nil
}

/*
Send sends a request to the server and returns the response body.

Parameters:
- req: The request to send
- opts: Optional parameters for the request

Returns:
- []byte: The response body
- error: Any error that occurred during the request
*/
func (c *Client) Send(req Requestor, opts ...Options) ([]byte, error) {
	var resp *resty.Response
	var err error

	r := c.RestyClient.R()

	isMultipart, _ := handleOptions(r, opts...)
	if isMultipart {
		r.SetHeader("Content-Type", "multipart/form-data")
	}

	// log.Debug().Msgf("endpoint: %s", req.Endpoint())

	switch req.Method() {
	case GET:
		resp, err = r.Get(req.Endpoint())
	case POST:
		resp, err = r.Post(req.Endpoint())
	case PATCH:
		resp, err = r.Patch(req.Endpoint())
	default:
		return nil, fmt.Errorf("unsupported method")
	}

	if err != nil {
		return nil, err
	}

	// Return an error if the status code is not 200.
	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode())
	}

	return resp.Body(), nil
}

func (c *Client) GetCustodian(id int) (Custodian, error) {
	var v []byte
	var err error

	custodian := Custodian{}

	req, _ := NewRequest().WithTenant(c.tenant).Custodian().WithID(id).Build()
	if v, err = c.Send(req); err != nil {
		return custodian, err
	}
	if err = json.Unmarshal(v, &custodian); err != nil {
		return custodian, err
	}

	return custodian, nil
}

func (c *Client) GetCustodians(opts Options) (Custodians, error) {
	var v []byte
	var err error

	var resp CustodiansResponse

	req, _ := NewRequest().WithTenant(c.tenant).Get().Custodian().Build()
	if v, err = c.Send(req, opts); err != nil {
		return nil, err
	}
	if err = json.Unmarshal(v, &resp); err != nil {
		return nil, err
	}
	return resp.Embedded.Custodians, nil
}

func (c *Client) GetFolder(id int) (Folder, error) {
	var v []byte
	var err error

	folder := Folder{}

	req, _ := NewRequest().WithTenant(c.tenant).Folder().WithID(id).Build()
	if v, err = c.Send(req); err != nil {
		return folder, err
	}
	if err = json.Unmarshal(v, &folder); err != nil {
		return folder, err
	}

	return folder, nil
}

func (c *Client) GetFolders(opts Options) (Folders, error) {
	var err error

	var body []byte
	var resp FoldersResponse

	req, _ := NewRequest().WithTenant(c.tenant).Get().Folder().Build()
	if body, err = c.Send(req, opts); err != nil {
		return nil, err
	}

	if err = json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}

	return resp.Embedded.Folders, nil
}

func (c *Client) GetGroup(id int) (Group, error) {
	var v []byte
	var err error

	group := Group{}

	req, _ := NewRequest().WithTenant(c.tenant).Group().WithID(id).Build()
	if v, err = c.Send(req); err != nil {
		return group, err
	}
	if err = json.Unmarshal(v, &group); err != nil {
		return group, err
	}

	return group, nil
}

func (c *Client) GetGroups(opts Options) (Groups, error) {
	var err error

	var body []byte
	var resp GroupsResponse

	req, _ := NewRequest().WithTenant(c.tenant).Get().Group().Build()
	if body, err = c.Send(req, opts); err != nil {
		return nil, err
	}

	if err = json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}

	return resp.Embedded.Groups, nil
}

func (c *Client) GetMatter(id int) (Matter, error) {
	var v []byte
	var err error

	matter := Matter{}

	req, _ := NewRequest().WithTenant(c.tenant).Matter().WithID(id).Build()
	if v, err = c.Send(req); err != nil {
		return matter, err
	}
	if err = json.Unmarshal(v, &matter); err != nil {
		return matter, err
	}

	return matter, nil
}

func (c *Client) GetMatters(opts Options) (Matters, error) {
	var err error

	var body []byte
	var resp MattersResponse

	req, _ := NewRequest().WithTenant(c.tenant).Get().Matter().Build()

	if body, err = c.Send(req, opts); err != nil {
		return nil, err
	}

	if err = json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}

	return resp.Embedded.Matters, nil
}

func (c *Client) GetLegalhold(id int) (Legalhold, error) {
	var v []byte
	var err error

	legalhold := Legalhold{}

	req, _ := NewRequest().WithTenant(c.tenant).Legalhold().WithID(id).Build()
	if v, err = c.Send(req); err != nil {
		return legalhold, err
	}
	if err = json.Unmarshal(v, &legalhold); err != nil {
		return legalhold, err
	}

	return legalhold, nil
}

func (c *Client) GetLegalholds(opts Options) (Legalholds, error) {
	var err error

	var body []byte
	var resp LegalholdsResponse

	req, _ := NewRequest().WithTenant(c.tenant).Get().Legalhold().Build()

	if body, err = c.Send(req, opts); err != nil {
		return nil, err
	}

	if err = json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}

	return resp.Embedded.Legalholds, nil
}

/**
 * GetAllCustodians retrieves all custodians from the API, paginating through results as needed.
 *
 * Parameters:
 * - req: the request to send to the API
 * - opts: options for the request
 *
 * Returns:
 * - custodiansChan: a channel that will receive all custodians
 * - errsChan: a channel that will receive any errors that occur
 */
func (c *Client) GetAllCustodians(opts Options) (<-chan Custodian, <-chan error) {

	var resp CustodiansResponse

	custodiansChan := make(chan Custodian)
	errsChan := make(chan error)

	req, _ := NewRequest().WithTenant(c.tenant).Get().Custodian().Build()

	go func() {
		var v []byte
		var err error

		defer close(custodiansChan)
		defer close(errsChan)

		page := 1
		for {
			log.Debug().Msgf("GetAllCustodians: page %d", page)
			opts.(*ListOptions).WithPageNumber(page)
			if v, err = c.Send(req, opts); err != nil {
				errsChan <- err
				break
			}

			if err = json.Unmarshal(v, &resp); err != nil {
				errsChan <- err
				break
			}

			for _, custodian := range resp.Embedded.Custodians {
				custodiansChan <- custodian
			}

			if !resp.Page.HasMore {
				return
			}
			page++
		}
	}()

	return custodiansChan, errsChan
}

func (c *Client) PrintAllCustodians() {

	printer := NewPrinter().JSON().Build()

	opts := NewListOptions().WithPageSize(100)

	custodianCh, errCh := c.GetAllCustodians(opts)

	fmt.Println("[")

	isFirstElement := true
	for {
		select {
		case custodian, ok := <-custodianCh:

			if !ok {
				custodianCh = nil
			} else {
				if !isFirstElement {
					fmt.Println(",")
				}
				isFirstElement = false

				printer.Print(custodian)
			}
		case err, ok := <-errCh:
			if !ok {
				errCh = nil
			} else {
				log.Error().Err(err)
			}
		}

		if custodianCh == nil && errCh == nil {
			break
		}
	}

	fmt.Println("]")
}

/**
 * ImportLegalhold imports a legal hold from a ZIP file.
 *
 * Parameters:
 * - zipFile: the path to the ZIP file containing the legal hold details
 *
 * Returns:
 * - Legalhold: the imported legal hold
 * - error: any error that occurred during the import process
 *
 * This function first creates a new request to import a legal hold using the
 * NewRequest, WithTenant, Post, Legalhold, and Import methods. It then creates
 * a new file option using NewFileOptions and WithFile to specify the ZIP file
 * to be uploaded.
 *
 * The function then sends the request with the file option using the Send method.
 * If the request is successful, the response body is unmarshaled into a Legalhold
 * struct and returned. If any errors occur during the process, they are returned
 * as an error.
 */
func (c *Client) ImportLegalhold(zipFile string) (Legalhold, error) {
	var err error

	var respBody []byte
	var resp Legalhold = Legalhold{}

	req, _ := NewRequest().WithTenant(c.tenant).Post().Legalhold().Import().Build()
	opts := NewFileOptions().WithFile("legal_hold_details", zipFile)

	if respBody, err = c.Send(req, opts); err != nil {
		return resp, err
	}

	if err = json.Unmarshal(respBody, &resp); err != nil {
		return resp, err
	}

	return resp, nil
}

/**
 * FindFolderByName finds a folder by its exact name.
 *
 * Parameters:
 * - name: the name of the folder to find
 *
 * Returns:
 * - Folder: the found folder, or an empty Folder if not found
 * - error: any error that occurred during the search
 *
 * This function first retrieves a list of folders using the GetFolders method with a filter
 * for the provided name. It then iterates through the list of folders and returns the first
 * one that matches the exact name. If no folder is found, an error is returned.
 */
func (c *Client) FindFolderByName(name string) (Folder, error) {
	var err error
	var folders Folders = Folders{}

	log.Debug().Msgf("searching folder by name [%s]", name)

	/*
		filterName helps to limit the results to folders CONTAINING the provided name.
	*/
	opts := NewListOptions().WithFilterName(name)

	if folders, err = c.GetFolders(opts); err != nil {
		return Folder{}, err
	}

	// find folder name matches exactly and return it
	for _, folder := range folders {
		if folder.Name == name {
			return folder, nil
		}
	}

	return Folder{}, fmt.Errorf("folder [%s] not found", name)
}

/**
 * CreateFolder creates a new folder with the given name.
 *
 * Parameters:
 * - name: the name of the folder to create
 * - groupIDs: list of group ids to associate with the folder
 *
 * Returns:
 * - Folder: the newly created folder
 * - error: any error that occurred during the creation
 *
 * This function sends a request to create a new folder with the given name. If the
 * creation is successful, the new folder is returned. If an error occurs, the error
 * is returned.
 */
func (c *Client) CreateFolder(name string, groupIDs []int) (Folder, error) {
	var err error
	var respBody []byte

	var folder Folder = Folder{}

	var createFolder *CreateFolderBody = NewCreateFolderBody().WithName(name).WithGroupIDs(groupIDs)

	req, _ := NewRequest().WithTenant(c.tenant).Post().Folder().Build()

	body, _ := json.Marshal(createFolder)
	opts := NewBodyOptions().WithBody(string(body))

	if respBody, err = c.Send(req, opts); err != nil {
		return folder, err
	}

	if err = json.Unmarshal(respBody, &folder); err != nil {
		return folder, err
	}

	log.Debug().Msgf("created folder %s with id %d", name, folder.ID)

	return folder, nil
}

/**
 * FindOrCreateFolder attempts to find a folder with the given name, and if it doesn't exist, creates a new folder with that name.
 *
 * @param name - The name of the folder to find or create.
 * @returns The found or created folder, and any error that occurred.
 */
func (c *Client) FindOrCreateFolder(name string) (Folder, error) {
	if folder, err := c.FindFolderByName(name); err == nil {
		log.Debug().Msgf("found folder [%s] with id [%d]", name, folder.ID)
		return folder, nil
	}

	log.Debug().Msgf("folder [%s] not found, creating", name)

	group, err := c.FindGroupByName("All Admins")
	if err != nil {
		log.Debug().Msg("failed to find defalt admin group [All Admins]")
		return Folder{}, err
	}

	return c.CreateFolder(name, []int{group.ID})
}

/**
 * FindOrCreateMatter attempts to find a matter by the given name, and if not found, creates a new matter with the given name and folder ID.
 *
 * @param name - The name of the matter to find or create.
 * @param folderID - The ID of the folder to associate the new matter with.
 * @returns The found or created matter, and any error that occurred.
 */
func (c *Client) FindMatterByName(name string) (Matter, error) {
	var err error
	var matters Matters = Matters{}

	log.Debug().Msgf("searching matter by name [%s]", name)

	/*
		filterName helps to limit the results to folders CONTAINING the provided name.
	*/
	opts := NewListOptions().WithFilterName(name)

	if matters, err = c.GetMatters(opts); err != nil {
		return Matter{}, err
	}

	for _, matter := range matters {
		if matter.Name == name {
			return matter, nil
		}
	}

	return Matter{}, fmt.Errorf("matter [%s] not found", name)
}

/**
 * CreateMatter creates a new matter with the given name and folder ID.
 *
 * @param name - The name of the matter to create.
 * @param folderID - The ID of the folder to associate the new matter with.
 * @returns The created matter, and any error that occurred.
 */
func (c *Client) CreateMatter(name string, folderID int) (Matter, error) {
	var err error
	var respBody []byte

	var matter Matter = Matter{}

	var createMatter *CreateMatterBody = NewCreateMatterBody().WithName(name).WithFolderID(folderID)

	req, _ := NewRequest().WithTenant(c.tenant).Post().Matter().Build()

	body, _ := json.Marshal(createMatter)
	opts := NewBodyOptions().WithBody(string(body))

	if respBody, err = c.Send(req, opts); err != nil {
		return matter, err
	}

	if err = json.Unmarshal(respBody, &matter); err != nil {
		return matter, err
	}

	log.Debug().Msgf("created matter %s with id %d", name, matter.ID)

	return matter, nil
}

/**
 * FindOrCreateMatter finds a matter by name, or creates a new matter if it doesn't exist.
 *
 * @param name - The name of the matter to find or create.
 * @param folderID - The ID of the folder to associate the new matter with if it needs to be created.
 * @returns The found or created matter, and any error that occurred.
 */
func (c *Client) FindOrCreateMatter(name string, folderID int) (Matter, error) {
	if matter, err := c.FindMatterByName(name); err == nil {
		log.Debug().Msgf("found matter [%s] with id [%d]", name, matter.ID)
		return matter, nil
	}

	log.Debug().Msgf("matter [%s] not found, creating", name)

	return c.CreateMatter(name, folderID)
}

/**
 * FindLegalhold finds a legalhold by name and matter ID.
 *
 * @param name - The name of the legalhold to find.
 * @param matterID - The ID of the matter associated with the legalhold.
 * @returns The found legalhold, and any error that occurred.
 */
func (c *Client) FindLegalhold(name string, matterID int) (Legalhold, error) {
	var err error
	var legalholds Legalholds = Legalholds{}

	log.Debug().Msgf("searching legalhold by name [%s] and matterID [%d]", name, matterID)

	/*
		filterName helps to limit the results to folders CONTAINING the provided name.
	*/
	opts := NewListOptions().WithFilterName(name)

	if legalholds, err = c.GetLegalholds(opts); err != nil {
		return Legalhold{}, err
	}

	for _, legalhold := range legalholds {
		if legalhold.Name == name && legalhold.MatterID == matterID {
			return legalhold, nil
		}
	}

	return Legalhold{}, fmt.Errorf("legalhold [%s] not found", name)
}

func (c *Client) FindGroupByName(name string) (Group, error) {
	var err error
	var groups Groups = Groups{}

	log.Debug().Msgf("searching groups by name [%s]", name)

	opts := NewListOptions().WithFilterName(name)

	if groups, err = c.GetGroups(opts); err != nil {
		return Group{}, err
	}

	for _, group := range groups {
		if group.Name == name {
			return group, nil
		}
	}

	return Group{}, fmt.Errorf("group [%s] not found", name)
}
