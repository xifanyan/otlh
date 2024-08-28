package importer

import (
	"fmt"
	"os"
	"strings"

	"github.com/rs/zerolog/log"
	otlh "github.com/xifanyan/otlh/pkg"
	"github.com/xuri/excelize/v2"
)

var silentholdTemplateHeader = []string{"Folder Name", "Matter Name", "Hold Name", "Advisory Notice Subject", "Advisory Notice Title", "Custodian Name", "Custodian Email", "Last Issued", "Release Date", "Advisory Notice Body"}

type SilentholdExcelImporter struct {
	ExcelImporter
}

func (imptr *SilentholdExcelImporter) collect(entry HoldEntry) {
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

	if _, ok := imptr.collections.HoldDetailMap[entry.MatterName+entry.HoldName]; !ok {
		imptr.collections.HoldDetailMap[entry.MatterName+entry.HoldName] = SilentholdInfo{
			MatterName:            entry.MatterName,
			HoldName:              entry.HoldName,
			AdvisoryNoticeSubject: entry.Subject,
			AdvisoryNoticeTitle:   entry.Title,
			AdvisoryNoticeBody:    entry.Body,
		}
	}

	imptr.collections.HoldToCustodiansMap[entry.MatterName+entry.HoldName] = append(imptr.collections.HoldToCustodiansMap[entry.MatterName+entry.HoldName], CustodianDetail{
		Name:       entry.CustodianName,
		Email:      entry.CustodianEmail,
		SentAt:     entry.LastIssued,
		ReleasedAt: entry.ReleasedAt,
	})
}

func (imptr *SilentholdExcelImporter) LoadSilentholdData() error {
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
			if imptr.verifyHeader(row, silentholdTemplateHeader) == nil {
				log.Debug().Msgf("found header at line #%d", l+1)
				imptr.lineNnumberOfHeader = l + 1
				continue
			}
		}

		// after header is found, load all non-empty rows
		if imptr.lineNnumberOfHeader > 0 && len(row) > 0 {
			// make sure the length of the row is equal to the length of the header
			data := make([]string, len(silentholdTemplateHeader))
			copy(data, row)

			// IMPORTANT: trim all the values to avoid corner cases during especially when data is part of queryParam
			for i := range data {
				data[i] = strings.TrimSpace(data[i])
			}

			entry := HoldEntry{
				FolderName:     data[0],
				MatterName:     data[1],
				HoldName:       data[2],
				Subject:        data[3],
				Title:          data[4],
				CustodianName:  data[5],
				CustodianEmail: data[6],
				LastIssued:     data[7],
				ReleasedAt:     data[8],
				Body:           data[9],
			}

			matterMatch := imptr.matterName == "" || imptr.matterName == entry.MatterName
			holdMatch := imptr.holdName == "" || imptr.holdName == entry.HoldName
			if (matterMatch && holdMatch) || (matterMatch && imptr.holdName == "") || (holdMatch && imptr.matterName == "") {
				imptr.entries = append(imptr.entries, entry)
				imptr.collect(entry)
			}
		}
	}

	log.Debug().Msgf("unique folders: %d, unique matters: %d, unique custodians: %d, matter to folder map size: %d, hold to custodians map size: %d",
		len(imptr.collections.UnqiueFolderNames),
		len(imptr.collections.UniqueMatterNames),
		len(imptr.collections.UniqueCustodians),
		len(imptr.collections.MatterToFolderMap),
		len(imptr.collections.HoldDetailMap),
	)

	return nil
}

func (imptr *SilentholdExcelImporter) transformToSilentholdDetails() (SilentholdDetails, error) {
	var silentholdDetails SilentholdDetails = []SilentholdDetail{}

	for _, hold := range imptr.collections.HoldDetailMap {
		holdInfo := hold.(SilentholdInfo)
		silentholdDetail := SilentholdDetail{
			FolderName:       imptr.collections.MatterToFolderMap[holdInfo.MatterName],
			SilentholdInfo:   holdInfo,
			CustodianDetails: imptr.collections.HoldToCustodiansMap[holdInfo.MatterName+holdInfo.HoldName],
		}
		silentholdDetails = append(silentholdDetails, silentholdDetail)
	}

	return silentholdDetails, nil
}

func (imptr *SilentholdExcelImporter) Import() error {

	log.Debug().Msg("[Start]: Importing Legalholds from Excel")

	silentholdDetails, err := imptr.transformToSilentholdDetails()
	if err != nil {
		return err
	}

	for _, silentholdDetail := range silentholdDetails {
		var tmpDir string
		var matterID int

		log.Debug().Msgf("[Processing]: Matter %s, Hold: %s, # of custodians: %d",
			silentholdDetail.SilentholdInfo.MatterName,
			silentholdDetail.SilentholdInfo.HoldName,
			len(silentholdDetail.CustodianDetails),
		)

		if matterID, err = imptr.getMatterID(silentholdDetail.SilentholdInfo.MatterName); err != nil {
			log.Error().Msgf("[%s - %s] not able to get matter id", silentholdDetail.SilentholdInfo.MatterName, silentholdDetail.SilentholdInfo.HoldName)
			continue
		}
		silentholdDetail.SilentholdInfo.MatterID = fmt.Sprintf("%d", matterID)

		_, err := imptr.client.FindSilenthold(silentholdDetail.SilentholdInfo.HoldName, matterID)
		if err == nil {
			log.Error().Msgf("[%s - %s] already exists", silentholdDetail.SilentholdInfo.MatterName, silentholdDetail.SilentholdInfo.HoldName)
			continue
		}

		if tmpDir, err = os.MkdirTemp("", "silenthold_"); err != nil {
			log.Error().Msgf("not able to create temp dir [%s - %s]", silentholdDetail.SilentholdInfo.MatterName, silentholdDetail.SilentholdInfo.HoldName)
			continue
		}
		log.Debug().Msgf("temp dir: %s", tmpDir)

		if err = silentholdDetail.saveToExcel(tmpDir, imptr.timezone); err != nil {
			log.Error().Msgf("not able to save to excel file [%s - %s]", silentholdDetail.SilentholdInfo.MatterName, silentholdDetail.SilentholdInfo.HoldName)
			continue
		}

		files := []string{fmt.Sprintf("%s/%s", tmpDir, "silent_hold_details.xlsx")}

		log.Debug().Msgf("Creating Zip file: %s", tmpDir+"/silent_hold_details.zip")
		if err = otlh.CreateZipArchive(tmpDir+"/silent_hold_details.zip", files); err != nil {
			log.Error().Msgf("not able to create zip file [%s - %s]: %s", silentholdDetail.SilentholdInfo.MatterName, silentholdDetail.SilentholdInfo.HoldName, err)
			continue
		}

		log.Debug().Msgf("Importing silenthold - [%s - %s]", silentholdDetail.SilentholdInfo.MatterName, silentholdDetail.SilentholdInfo.HoldName)
		_, err = imptr.client.ImportSilenthold(tmpDir + "/silent_hold_details.zip")
		if err != nil {
			log.Error().Msgf("silenthold import failed %s - [%s - %s]", err, silentholdDetail.SilentholdInfo.MatterName, silentholdDetail.SilentholdInfo.HoldName)
			continue
		}

	}

	log.Debug().Msg("[End]: Finished Importing Silentholds")

	return nil
}
