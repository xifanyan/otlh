package importer

import (
	"fmt"

	"github.com/xuri/excelize/v2"
)

var SilentholdDetailsHeader = []string{"Matter id", "Hold Name", "Advisory notice subject", "Advisory notice body", "Advisory notice title"}
var SilentholdCustodianDetailsHeader = []string{"Name", "Email", "sent_at", "released_at"}

type SilentholdInfo struct {
	MatterName            string
	MatterID              string `json:"Matter id"`
	HoldName              string `json:"Hold Name"`
	AdvisoryNoticeSubject string `json:"Advisory notice subject"`
	AdvisoryNoticeTitle   string `json:"Advisory notice title"`
	AdvisoryNoticeBody    string `json:"Advisory notice body"`
}

type SilentholdDetail struct {
	FolderName       string
	SilentholdInfo   SilentholdInfo
	CustodianDetails []CustodianDetail
}

type SilentholdDetails []SilentholdDetail

func (shd SilentholdDetail) saveToExcel(dir string, tz string) error {
	var err error

	f := excelize.NewFile()
	defer f.Close()

	if _, err = f.NewSheet(SHEET_NAME_HOLD_DETAILS); err != nil {
		return err
	}

	_ = f.SetSheetRow(SHEET_NAME_HOLD_DETAILS, "A1", &SilentholdDetailsHeader)
	_ = f.SetSheetRow(SHEET_NAME_HOLD_DETAILS, "A2",
		&[]interface{}{
			shd.SilentholdInfo.MatterID,
			shd.SilentholdInfo.HoldName,
			shd.SilentholdInfo.AdvisoryNoticeSubject,
			shd.SilentholdInfo.AdvisoryNoticeBody,
			shd.SilentholdInfo.AdvisoryNoticeTitle,
		},
	)

	if _, err = f.NewSheet(SHEET_NAME_CUSTODIANS_DETAILS); err != nil {
		return err
	}

	_ = f.SetSheetRow(SHEET_NAME_CUSTODIANS_DETAILS, "A1", &SilentholdCustodianDetailsHeader)
	for i, custodianDetail := range shd.CustodianDetails {
		row := []any{
			custodianDetail.Name,
			custodianDetail.Email,
			convertDateTimeOrDefault(tz, custodianDetail.SentAt),
			convertDateTimeOrDefault(tz, custodianDetail.ReleasedAt),
		}

		_ = f.SetSheetRow(SHEET_NAME_CUSTODIANS_DETAILS, fmt.Sprintf("A%d", i+2), &row)
	}

	// delete detaful sheet "Sheet1"
	if err = f.DeleteSheet("Sheet1"); err != nil {
		return err
	}

	if err = f.SaveAs(fmt.Sprintf("%s/silent_hold_details.xlsx", dir)); err != nil {
		return err
	}

	return nil
}
