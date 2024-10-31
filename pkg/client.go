package otlh

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/schollz/progressbar/v3"
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

	if zerolog.GlobalLevel() == zerolog.TraceLevel {
		r.SetDebug(true)
	}

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

func (c *Client) Tenant() string {
	return c.tenant
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

/*
Do sends a request and unmarshals the response into a given resource.

Parameters:
- req: The request to send
- resource: The object to unmarshal the response into
- opts: Optional parameters for the request

Returns:
- error: Any error that occurred during the request or unmarshalling
*/
func (c *Client) Do(req Requestor, resource any, opts ...Options) error {
	var v []byte
	var err error

	if v, err = c.Send(req, opts...); err != nil {
		return err
	}

	return json.Unmarshal(v, resource)
}

func (c *Client) GetCustodian(req Requestor) (Custodian, error) {
	var custodian Custodian
	return custodian, c.Do(req, &custodian)
}

func (c *Client) GetCustodians(req Requestor, opts ...Options) (Custodians, error) {
	var resp CustodiansResponse
	return resp.Embedded.Custodians, c.Do(req, &resp, opts...)
}

func (c *Client) GetAllCustodians(req Requestor, opts Options) (Custodians, error) {
	return getAllEntities(c, req, opts, unmarshalCustodians)
}

func (c *Client) GetCustodianGroup(req Requestor) (CustodianGroup, error) {
	var custodianGroup CustodianGroup
	return custodianGroup, c.Do(req, &custodianGroup)
}

func (c *Client) GetCustodianGroups(req Requestor, opts ...Options) (CustodianGroups, error) {
	var resp CustodianGroupsResponse
	return resp.Embedded.CustodianGroups, c.Do(req, &resp, opts...)
}

func (c *Client) GetAllCustodianGroups(req Requestor, opts Options) (CustodianGroups, error) {
	return getAllEntities(c, req, opts, unmarshalCustodianGroups)
}

func (c *Client) GetGroup(req Requestor) (Group, error) {
	var group Group
	return group, c.Do(req, &group)
}

func (c *Client) GetGroups(req Requestor, opts ...Options) (Groups, error) {
	var resp GroupsResponse
	return resp.Embedded.Groups, c.Do(req, &resp, opts...)
}

func (c *Client) GetAllGroups(req Requestor, opts Options) (Groups, error) {
	return getAllEntities(c, req, opts, unmarshalGroups)
}

func (c *Client) GetFolder(req Requestor) (Folder, error) {
	var folder Folder
	return folder, c.Do(req, &folder)
}

func (c *Client) GetFolders(req Requestor, opts ...Options) (Folders, error) {
	var resp FoldersResponse
	return resp.Embedded.Folders, c.Do(req, &resp, opts...)
}

func (c *Client) GetAllFolders(req Requestor, opts Options) (Folders, error) {
	return getAllEntities(c, req, opts, unmarshalFolders)
}

func (c *Client) GetMatter(req Requestor) (Matter, error) {
	var matter Matter
	return matter, c.Do(req, &matter)
}

func (c *Client) GetMatters(req Requestor, opts ...Options) (Matters, error) {
	var resp MattersResponse
	return resp.Embedded.Matters, c.Do(req, &resp, opts...)
}

func (c *Client) GetAllMatters(req Requestor, opts Options) (Matters, error) {
	return getAllEntities(c, req, opts, unmarshalMatters)
}

func (c *Client) GetLegalhold(req Requestor) (Legalhold, error) {
	var legalhold Legalhold
	return legalhold, c.Do(req, &legalhold)
}

func (c *Client) GetLegalholds(req Requestor, opts ...Options) (Legalholds, error) {
	var resp LegalholdsResponse
	err := c.Do(req, &resp, opts...)
	return resp.Embedded.Legalholds, err
}

func (c *Client) GetAllLegalholds(req Requestor, opts Options) (Legalholds, error) {
	return getAllEntities(c, req, opts, unmarshalLegalholds)
}

func (c *Client) GetSilenthold(req Requestor) (Silenthold, error) {
	var silenthold Silenthold
	return silenthold, c.Do(req, &silenthold)
}

func (c *Client) GetSilentholds(req Requestor, opts ...Options) (Silentholds, error) {
	var resp SilentholdsResponse
	err := c.Do(req, &resp, opts...)
	return resp.Embedded.Silentholds, err
}

func (c *Client) GetAllSilentholds(req Requestor, opts Options) (Silentholds, error) {
	return getAllEntities(c, req, opts, unmarshalSilentholds)
}

func (c *Client) ImportCustodians(custodians []CustodianInputData, batchSize int) error {
	bar := progressbar.Default(int64(len(custodians)))
	defer bar.Finish()

	req := NewRequest().WithTenant(c.tenant).Post().CustodiansSync().Build()

	for i := 0; i < len(custodians); i += batchSize {
		end := i + batchSize
		if end > len(custodians) {
			end = len(custodians)
		}

		log.Debug().Msgf("Batch: %d - %d", i, end)

		body := CustodianSyncBody{
			Custodians: custodians[i:end],
		}

		custodianBody, err := json.Marshal(body)
		if err != nil {
			return err
		}

		opts := NewBodyOptions().WithBody(string(custodianBody))
		var resp CustodiansSyncResponse
		if err := c.Do(req, &resp, opts); err != nil {
			return err
		}

		bar.Add(batchSize)
		time.Sleep(5 * time.Millisecond)
	}

	return nil
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

func (c *Client) ImportSilenthold(zipFile string) (Silenthold, error) {
	var err error

	var respBody []byte
	var resp Silenthold = Silenthold{}

	req, _ := NewRequest().WithTenant(c.tenant).Post().Silenthold().Import().Build()
	opts := NewFileOptions().WithFile("silent_hold_details", zipFile)

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
	req, _ := NewRequest().WithTenant(c.tenant).Get().Folder().Build()
	opts := NewListOptions().WithFilterName(name)

	if folders, err = c.GetFolders(req, opts); err != nil {
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
	var folder Folder = Folder{}

	req, _ := NewRequest().WithTenant(c.tenant).Post().Folder().Build()

	var createFolder *CreateFolderBody = NewCreateFolderBody().WithName(name).WithGroupIDs(groupIDs)
	body, _ := json.Marshal(createFolder)
	opts := NewBodyOptions().WithBody(string(body))

	if err = c.Do(req, &folder, opts); err != nil {
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
	req, _ := NewRequest().WithTenant(c.tenant).Get().Matter().Build()
	opts := NewListOptions().WithFilterName(name)

	if matters, err = c.GetMatters(req, opts); err != nil {
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
	var matter Matter = Matter{}

	req, _ := NewRequest().WithTenant(c.tenant).Post().Matter().Build()

	var createMatter *CreateMatterBody = NewCreateMatterBody().WithName(name).WithFolderID(folderID)
	body, _ := json.Marshal(createMatter)
	opts := NewBodyOptions().WithBody(string(body))

	if err = c.Do(req, &matter, opts); err != nil {
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

	req, _ := NewRequest().WithTenant(c.tenant).Get().Legalhold().Build()
	opts := NewListOptions().WithFilterName(name)

	if legalholds, err = c.GetLegalholds(req, opts); err != nil {
		return Legalhold{}, err
	}

	for _, legalhold := range legalholds {
		if legalhold.Name == name && legalhold.MatterID == matterID {
			return legalhold, nil
		}
	}

	return Legalhold{}, fmt.Errorf("legalhold [%s] not found", name)
}

func (c *Client) FindSilenthold(name string, matterID int) (Silenthold, error) {
	var err error
	var silentholds Silentholds = Silentholds{}

	log.Debug().Msgf("searching legalhold by name [%s] and matterID [%d]", name, matterID)

	req, _ := NewRequest().WithTenant(c.tenant).Get().Silenthold().Build()
	opts := NewListOptions().WithFilterName(name)

	if silentholds, err = c.GetSilentholds(req, opts); err != nil {
		return Silenthold{}, err
	}

	for _, silenthold := range silentholds {
		if silenthold.Name == name && silenthold.MatterID == matterID {
			return silenthold, nil
		}
	}

	return Silenthold{}, fmt.Errorf("silent [%s] not found", name)
}

// FindGroupByName searches for a group by name.
//
// Parameters:
// - name: the name of the group to search for.
//
// Returns:
// - Group: the group found, or an empty Group if not found.
// - error: any error that occurred during the search.
func (c *Client) FindGroupByName(name string) (Group, error) {
	var err error
	var groups Groups = Groups{}

	log.Debug().Msgf("searching groups by name [%s]", name)

	req, _ := NewRequest().WithTenant(c.tenant).Get().Group().Build()
	opts := NewListOptions().WithFilterName(name)

	if groups, err = c.GetGroups(req, opts); err != nil {
		return Group{}, err
	}

	for _, group := range groups {
		if group.Name == name {
			return group, nil
		}
	}

	return Group{}, fmt.Errorf("group [%s] not found", name)
}
