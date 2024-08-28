package importer

import (
	"github.com/rs/zerolog/log"
	otlh "github.com/xifanyan/otlh/pkg"
)

const MAX_LENGTH_OF_HOLDNAME int = 100
const SHEET_NAME_HOLD_DETAILS string = "hold_details"
const SHEET_NAME_CUSTODIANS_DETAILS string = "custodian_details"

type HoldEntry struct {
	FolderName      string
	MatterName      string
	HoldName        string
	Subject         string
	Title           string
	CustodianName   string
	CustodianEmail  string
	LastIssued      string
	ResponseDate    string
	ReleasedAt      string
	Body            string
	AttachmentNames string
}

type HoldEntries []HoldEntry

type CustodianDetail struct {
	Name          string `json:"name"`
	Email         string `json:"email"`
	SentAt        string `json:"sent_at"`
	AcknowlegedAt string `json:"acknowledged_at"`
	ReleasedAt    string `json:"released_at"`
}

type Collections struct {
	UnqiueFolderNames     map[string]int
	UniqueMatterNames     map[string]int
	UniqueCustodians      map[string]string
	UniqueAttachmentNames map[string]struct{}
	MatterToFolderMap     map[string]string
	// HoldDetailMap         map[string]LegalholdInfo
	HoldDetailMap       map[string]interface{}
	HoldToCustodiansMap map[string][]CustodianDetail
}

type ExcelImporter struct {
	excel               string
	matterName          string
	holdName            string
	lineNnumberOfHeader int
	collections         Collections
	client              *otlh.Client
	timezone            string
	entries             HoldEntries
}

func NewExcelImporter() *ExcelImporter {
	return &ExcelImporter{
		collections: Collections{
			UnqiueFolderNames:     make(map[string]int),
			UniqueMatterNames:     make(map[string]int),
			UniqueCustodians:      make(map[string]string),
			UniqueAttachmentNames: make(map[string]struct{}),
			MatterToFolderMap:     make(map[string]string),
			HoldDetailMap:         make(map[string]any),
			HoldToCustodiansMap:   make(map[string][]CustodianDetail),
		},
	}
}

func (e *ExcelImporter) WithClient(client *otlh.Client) *ExcelImporter {
	e.client = client
	return e
}

func (e *ExcelImporter) WithExcel(excel string) *ExcelImporter {
	e.excel = excel
	return e
}

func (e *ExcelImporter) WithTimezone(tz string) *ExcelImporter {
	e.timezone = tz
	return e
}

func (e *ExcelImporter) WithMatterName(name string) *ExcelImporter {
	e.matterName = name
	return e
}

func (e *ExcelImporter) WithHoldName(name string) *ExcelImporter {
	e.holdName = name
	return e
}

func (e *ExcelImporter) Legalhold() *LegalholdExcelImporter {
	return &LegalholdExcelImporter{
		ExcelImporter: *e,
	}
}

func (e *ExcelImporter) Silenthold() *SilentholdExcelImporter {
	return &SilentholdExcelImporter{
		ExcelImporter: *e,
	}
}

// FindCustodianByNameAndEmail finds a custodian by name and email.
//
// Parameters:
// - name: the name of the custodian to find
// - email: the email of the custodian to find
// Return type:
// - error: any error that occurred during the search
func (e *ExcelImporter) FindCustodianByNameAndEmail(name, email string) error {
	var err error

	req, _ := otlh.NewRequest().WithTenant(e.client.Tenant()).Get().Custodian().Build()
	opts := otlh.NewListOptions().WithFilterName(name)
	custodians, err := e.client.GetCustodians(req, opts)
	if err != nil {
		return err
	}

	for _, custodian := range custodians {
		if (custodian.Name == name) && (custodian.Email == email) {
			log.Debug().Msgf("- [found] %s - %s", name, email)
			return nil
		}
	}

	return ErrorCustodianNotFound
}

func (e *ExcelImporter) getFolderID(name string) (int, error) {
	var err error
	var folder otlh.Folder

	folderID := e.collections.UnqiueFolderNames[name]
	if folderID > 0 {
		return folderID, nil
	}

	if folder, err = e.client.FindOrCreateFolder(name); err != nil {
		return 0, err
	}
	return folder.ID, nil
}

func (e *ExcelImporter) getMatterID(name string) (int, error) {
	var err error
	var matter otlh.Matter
	var folderID int

	if folderID, err = e.getFolderID(e.collections.MatterToFolderMap[name]); err != nil {
		return 0, err
	}

	if matter, err = e.client.FindOrCreateMatter(name, folderID); err != nil {
		return 0, err
	}

	return matter.ID, nil
}
