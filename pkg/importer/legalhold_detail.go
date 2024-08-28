package importer

import (
	"fmt"
	"time"

	"github.com/xuri/excelize/v2"
)

var LegalholdDetailsHeader = []string{"Matter id", "Hold Name", "Hold notice subject", "Hold notice body", "Hold notice title", "Hold notice attachment names"}
var LegalholdCustodianDetailsHeader = []string{"Name", "Email", "sent_at", "acknowledged_at", "released_at"}

type LegalholdInfo struct {
	MatterName                string
	MatterID                  string `json:"Matter id"`
	HoldName                  string `json:"Hold Name"`
	HoldNoticeSubject         string `json:"Hold notice subject"`
	HoldNoticeTitle           string `json:"Hold notice title"`
	HoldNoticeBody            string `json:"Hold notice body"`
	HoldNoticeAttachmentNames string `json:"Hold notice attachment names"`
}

type LegalholdDetail struct {
	FolderName       string
	LegalholdInfo    LegalholdInfo
	CustodianDetails []CustodianDetail
}

type LegalholdDetails []LegalholdDetail

func convertDateTimeFormat(tz string, input string) (string, error) {
	if input == "" {
		return "", fmt.Errorf("empty datetime")
	}

	//	if tz == "UTC" {
	//		return input, nil
	//	}

	loc, err := time.LoadLocation(tz)
	if err != nil {
		return "", err
	}

	inputTime, err := time.ParseInLocation(INPUT_TIME_FORMAT, input, loc)
	if err != nil {
		return "", err
	}

	outputTime := inputTime.In(time.UTC)
	output := outputTime.Format(OUTPUT_TIME_FORMAT)

	return output, nil
}

func (lhd LegalholdDetail) saveToExcel(dir string, tz string) error {
	var err error

	f := excelize.NewFile()
	defer f.Close()

	f.NewSheet(SHEET_NAME_HOLD_DETAILS)
	f.SetSheetRow(SHEET_NAME_HOLD_DETAILS, "A1", &LegalholdDetailsHeader)
	f.SetSheetRow(SHEET_NAME_HOLD_DETAILS, "A2",
		&[]interface{}{
			lhd.LegalholdInfo.MatterID,
			lhd.LegalholdInfo.HoldName,
			lhd.LegalholdInfo.HoldNoticeSubject,
			lhd.LegalholdInfo.HoldNoticeBody,
			lhd.LegalholdInfo.HoldNoticeTitle,
			lhd.LegalholdInfo.HoldNoticeAttachmentNames,
		},
	)

	f.NewSheet(SHEET_NAME_CUSTODIANS_DETAILS)
	f.SetSheetRow(SHEET_NAME_CUSTODIANS_DETAILS, "A1", &LegalholdCustodianDetailsHeader)
	for i, custodianDetail := range lhd.CustodianDetails {
		row := []interface{}{
			custodianDetail.Name,
			custodianDetail.Email,
		}

		if sentAt, err := convertDateTimeFormat(tz, custodianDetail.SentAt); err == nil {
			row = append(row, sentAt)
		}

		if acknowlegedAt, err := convertDateTimeFormat(tz, custodianDetail.AcknowlegedAt); err == nil {
			row = append(row, acknowlegedAt)
		}

		if releasedAt, err := convertDateTimeFormat(tz, custodianDetail.ReleasedAt); err == nil {
			row = append(row, releasedAt)
		}

		f.SetSheetRow(SHEET_NAME_CUSTODIANS_DETAILS, fmt.Sprintf("A%d", i+2), &row)
	}

	// delete detaful sheet "Sheet1"
	f.DeleteSheet("Sheet1")

	if err = f.SaveAs(fmt.Sprintf("%s/legal_hold_details.xlsx", dir)); err != nil {
		return err
	}

	return nil
}
