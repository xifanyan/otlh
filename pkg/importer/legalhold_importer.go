package importer

import (
	"fmt"
	"os"
	"strings"

	"github.com/rs/zerolog/log"
	otlh "github.com/xifanyan/otlh/pkg"
	"github.com/xuri/excelize/v2"
)

var legalholdTemplateHeader = []string{"Folder Name", "Matter Name", "Hold Name", "Hold Notice Subject", "Hold Notice Title", "Custodian Name", "Custodian Email", "Last Issued", "Response Date", "Release Date", "Legal Hold Text", "Attachment Names"}

const INPUT_TIME_FORMAT = "1/2/06 3:04 PM"
const OUTPUT_TIME_FORMAT = "01/02/2006 03:04 PM"

type LegalholdEntry struct {
	FolderName        string
	MatterName        string
	HoldName          string
	HoldNoticeSubject string
	HoldNoticeTitle   string
	CustodianName     string
	CustodianEmail    string
	LastIssued        string
	ResponseDate      string
	ReleasedAt        string
	LegalholdText     string
	AttachmentNames   string
}

type LegalholdEntries []LegalholdEntry

type Collections struct {
	UnqiueFolderNames     map[string]int
	UniqueMatterNames     map[string]int
	UniqueCustodians      map[string]string
	UniqueAttachmentNames map[string]struct{}
	MatterToFolderMap     map[string]string
	HoldDetailMap         map[string]LegalholdInfo
	HoldToCustodiansMap   map[string][]CustodianDetail
}

type LegalholdExcelImporter struct {
	excel               string
	matterName          string
	holdName            string
	attachmentDirectory string
	entries             LegalholdEntries
	lineNnumberOfHeader int
	collections         Collections
	client              *otlh.Client
	timezone            string
}

func NewLegalholdExcelImporter() *LegalholdExcelImporter {
	return &LegalholdExcelImporter{
		entries: LegalholdEntries{},
		collections: Collections{
			UnqiueFolderNames:     make(map[string]int),
			UniqueMatterNames:     make(map[string]int),
			UniqueCustodians:      make(map[string]string),
			UniqueAttachmentNames: make(map[string]struct{}),
			MatterToFolderMap:     make(map[string]string),
			HoldDetailMap:         make(map[string]LegalholdInfo),
			HoldToCustodiansMap:   make(map[string][]CustodianDetail),
		},
	}
}

func (imptr *LegalholdExcelImporter) WithClient(client *otlh.Client) *LegalholdExcelImporter {
	imptr.client = client
	return imptr
}

func (imptr *LegalholdExcelImporter) WithExcel(excel string) *LegalholdExcelImporter {
	imptr.excel = excel
	return imptr
}

func (imptr *LegalholdExcelImporter) WithTimezone(tz string) *LegalholdExcelImporter {
	imptr.timezone = tz
	return imptr
}

func (imptr *LegalholdExcelImporter) WithAttachmentDirectory(dir string) *LegalholdExcelImporter {
	imptr.attachmentDirectory = dir
	return imptr
}

func (imptr *LegalholdExcelImporter) WithMatterName(name string) *LegalholdExcelImporter {
	imptr.matterName = name
	return imptr
}

func (imptr *LegalholdExcelImporter) WithHoldName(name string) *LegalholdExcelImporter {
	imptr.holdName = name
	return imptr
}

func (imptr *LegalholdExcelImporter) collect(entry LegalholdEntry) {
	if _, ok := imptr.collections.UnqiueFolderNames[entry.FolderName]; !ok {
		imptr.collections.UnqiueFolderNames[entry.FolderName] = 0
	}

	if _, ok := imptr.collections.UniqueMatterNames[entry.MatterName]; !ok {
		imptr.collections.UniqueMatterNames[entry.MatterName] = 0
	}

	if _, ok := imptr.collections.UniqueCustodians[entry.CustodianEmail]; !ok {
		imptr.collections.UniqueCustodians[entry.CustodianEmail] = entry.CustodianName
	}

	if _, ok := imptr.collections.MatterToFolderMap[entry.MatterName]; !ok {
		imptr.collections.MatterToFolderMap[entry.MatterName] = entry.FolderName
	}

	attachmentNames := strings.TrimSpace(entry.AttachmentNames)
	if attachmentNames != "" {
		for _, name := range strings.Split(attachmentNames, ",") {
			name = strings.TrimSpace(name)
			if _, ok := imptr.collections.UniqueAttachmentNames[name]; !ok {
				imptr.collections.UniqueAttachmentNames[name] = struct{}{}
			}
		}
	}

	if _, ok := imptr.collections.HoldDetailMap[entry.MatterName+entry.HoldName]; !ok {
		imptr.collections.HoldDetailMap[entry.MatterName+entry.HoldName] = LegalholdInfo{
			MatterName:                entry.MatterName,
			HoldName:                  entry.HoldName,
			HoldNoticeSubject:         entry.HoldNoticeSubject,
			HoldNoticeTitle:           entry.HoldNoticeTitle,
			HoldNoticeBody:            entry.LegalholdText,
			HoldNoticeAttachmentNames: entry.AttachmentNames,
		}
	}

	imptr.collections.HoldToCustodiansMap[entry.MatterName+entry.HoldName] = append(imptr.collections.HoldToCustodiansMap[entry.MatterName+entry.HoldName], CustodianDetail{
		Name:          entry.CustodianName,
		Email:         entry.CustodianEmail,
		SentAt:        entry.LastIssued,
		AcknowlegedAt: entry.ResponseDate,
		ReleasedAt:    entry.ReleasedAt,
	})
}

func (imptr *LegalholdExcelImporter) LoadLegalholdData() error {
	var rows [][]string

	log.Debug().Msgf("open excel file %s", imptr.excel)
	f, err := excelize.OpenFile(imptr.excel)
	if err != nil {
		return err
	}
	defer f.Close()

	// get the first sheet
	firstSheet := f.WorkBook.Sheets.Sheet[0].Name
	if rows, err = f.GetRows(firstSheet); err != nil {
		return err
	}

	for l, row := range rows {
		if imptr.lineNnumberOfHeader == 0 && len(row) > 0 && row[0] == "Folder Name" {
			if imptr.verifyHeader(row, legalholdTemplateHeader) == nil {
				log.Debug().Msgf("found header at line #%d", l+1)
				imptr.lineNnumberOfHeader = l + 1
				continue
			}
		}

		// after header is found, load all non-empty rows
		if imptr.lineNnumberOfHeader > 0 && len(row) > 0 {
			// make sure the length of the row is equal to the length of the header
			data := make([]string, len(legalholdTemplateHeader))
			copy(data, row)

			// IMPORTANT: trim all the values to avoid corner cases during especially when data is part of queryParam
			for i := range data {
				data[i] = strings.TrimSpace(data[i])
			}

			entry := LegalholdEntry{
				FolderName:        data[0],
				MatterName:        data[1],
				HoldName:          data[2],
				HoldNoticeSubject: data[3],
				HoldNoticeTitle:   data[4],
				CustodianName:     data[5],
				CustodianEmail:    data[6],
				LastIssued:        data[7],
				ResponseDate:      data[8],
				ReleasedAt:        data[9],
				LegalholdText:     data[10],
				AttachmentNames:   data[11],
			}

			matterMatch := imptr.matterName == "" || imptr.matterName == entry.MatterName
			holdMatch := imptr.holdName == "" || imptr.holdName == entry.HoldName
			if (matterMatch && holdMatch) || (matterMatch && imptr.holdName == "") || (holdMatch && imptr.matterName == "") {
				imptr.entries = append(imptr.entries, entry)
				imptr.collect(entry)
			}
		}
	}

	log.Debug().Msgf("unique folders: %d, unique matters: %d, unique custodians: %d, unique attachments: %d, matter to folder map size: %d, hold to custodians map size: %d",
		len(imptr.collections.UnqiueFolderNames),
		len(imptr.collections.UniqueMatterNames),
		len(imptr.collections.UniqueCustodians),
		len(imptr.collections.UniqueAttachmentNames),
		len(imptr.collections.MatterToFolderMap),
		len(imptr.collections.HoldDetailMap),
	)

	return nil
}

func (imptr *LegalholdExcelImporter) transformToLegalholdDetails() (LegalholdDetails, error) {
	var legalholdDetails LegalholdDetails = []LegalholdDetail{}

	for _, holdDetail := range imptr.collections.HoldDetailMap {
		legalholdDetail := LegalholdDetail{
			FolderName:       imptr.collections.MatterToFolderMap[holdDetail.MatterName],
			LegalholdInfo:    holdDetail,
			CustodianDetails: imptr.collections.HoldToCustodiansMap[holdDetail.MatterName+holdDetail.HoldName],
		}
		legalholdDetails = append(legalholdDetails, legalholdDetail)
	}

	return legalholdDetails, nil
}

func (imptr *LegalholdExcelImporter) GetFolderID(name string) (int, error) {
	var err error
	var folder otlh.Folder

	folderID := imptr.collections.UnqiueFolderNames[name]
	if folderID > 0 {
		return folderID, nil
	}

	if folder, err = imptr.client.FindOrCreateFolder(name); err != nil {
		return 0, err
	}
	return folder.ID, nil
}

func (imptr *LegalholdExcelImporter) GetMatterID(name string) (int, error) {
	var err error
	var matter otlh.Matter
	var folderID int

	if folderID, err = imptr.GetFolderID(imptr.collections.MatterToFolderMap[name]); err != nil {
		return 0, err
	}

	if matter, err = imptr.client.FindOrCreateMatter(name, folderID); err != nil {
		return 0, err
	}

	return matter.ID, nil
}

func (imptr *LegalholdExcelImporter) Import() error {

	log.Debug().Msg("[Start]: Importing Legalholds from Excel")

	legalholdDetails, err := imptr.transformToLegalholdDetails()
	if err != nil {
		return err
	}

	for _, legalholdDetail := range legalholdDetails {
		var tmpDir string
		var matterID int

		log.Debug().Msgf("[Processing]: Matter %s, Hold: %s, # of custodians: %d",
			legalholdDetail.LegalholdInfo.MatterName,
			legalholdDetail.LegalholdInfo.HoldName,
			len(legalholdDetail.CustodianDetails),
		)

		if matterID, err = imptr.GetMatterID(legalholdDetail.LegalholdInfo.MatterName); err != nil {
			log.Error().Msgf("[%s - %s] not able to get matter id", legalholdDetail.LegalholdInfo.MatterName, legalholdDetail.LegalholdInfo.HoldName)
			continue
		}
		legalholdDetail.LegalholdInfo.MatterID = fmt.Sprintf("%d", matterID)

		_, err := imptr.client.FindLegalhold(legalholdDetail.LegalholdInfo.HoldName, matterID)
		if err == nil {
			log.Error().Msgf("[%s - %s] already exists", legalholdDetail.LegalholdInfo.MatterName, legalholdDetail.LegalholdInfo.HoldName)
			continue
		}

		if tmpDir, err = os.MkdirTemp("", "legalhold_"); err != nil {
			log.Error().Msgf("not able to create temp dir [%s - %s]", legalholdDetail.LegalholdInfo.MatterName, legalholdDetail.LegalholdInfo.HoldName)
			continue
		}
		log.Debug().Msgf("temp dir: %s", tmpDir)

		if err = legalholdDetail.saveToExcel(tmpDir, imptr.timezone); err != nil {
			log.Error().Msgf("not able to save to excel file [%s - %s]", legalholdDetail.LegalholdInfo.MatterName, legalholdDetail.LegalholdInfo.HoldName)
			continue
		}

		files := []string{fmt.Sprintf("%s/%s", tmpDir, "legal_hold_details.xlsx")}
		if legalholdDetail.LegalholdInfo.HoldNoticeAttachmentNames != "" {
			for _, attachment := range strings.Split(legalholdDetail.LegalholdInfo.HoldNoticeAttachmentNames, ",") {
				attachment = strings.TrimSpace(attachment)
				files = append(files, fmt.Sprintf("%s/%s", imptr.attachmentDirectory, attachment))
			}
		}

		log.Debug().Msgf("Creating Zip file: %s", tmpDir+"/legal_hold_details.zip")
		if err = otlh.CreateZipArchive(tmpDir+"/legal_hold_details.zip", files); err != nil {
			log.Error().Msgf("not able to create zip file [%s - %s]: %s", legalholdDetail.LegalholdInfo.MatterName, legalholdDetail.LegalholdInfo.HoldName, err)
			continue
		}

		log.Debug().Msgf("Importing legalhold - [%s - %s]", legalholdDetail.LegalholdInfo.MatterName, legalholdDetail.LegalholdInfo.HoldName)
		_, err = imptr.client.ImportLegalhold(tmpDir + "/legal_hold_details.zip")
		if err != nil {
			log.Error().Msgf("legalhold import failed %s - [%s - %s]", err, legalholdDetail.LegalholdInfo.MatterName, legalholdDetail.LegalholdInfo.HoldName)
			continue
		}

	}

	log.Debug().Msg("[End]: Finished Importing Legalholds")

	return nil
}
