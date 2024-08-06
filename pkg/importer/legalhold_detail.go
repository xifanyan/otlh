package importer

import (
	"fmt"
	"time"

	"github.com/xuri/excelize/v2"
)

const MAX_LENGTH_OF_HOLDNAME int = 100

var HoldDetailsHeader = []string{"Matter id", "Hold Name", "Hold notice subject", "Hold notice body", "Hold notice title", "Hold notice attachment names"}
var CustodianDetailsHeader = []string{"Name", "Email", "sent_at", "acknowledged_at", "released_at"}

type HoldDetail struct {
	MatterName                string
	MatterID                  string `json:"Matter id"`
	HoldName                  string `json:"Hold Name"`
	HoldNoticeSubject         string `json:"Hold notice subject"`
	HoldNoticeTitle           string `json:"Hold notice title"`
	HoldNoticeBody            string `json:"Hold notice body"`
	HoldNoticeAttachmentNames string `json:"Hold notice attachment names"`
}

type CustodianDetail struct {
	Name          string `json:"name"`
	Email         string `json:"email"`
	SentAt        string `json:"sent_at"`
	AcknowlegedAt string `json:"acknowledged_at"`
	ReleasedAt    string `json:"released_at"`
}

type LegalholdDetail struct {
	FolderName       string
	HoldDetail       HoldDetail
	CustodianDetails []CustodianDetail
}

type LegalholdDetails []LegalholdDetail

func convertDateTimeFormat(tz string, input string) (string, error) {
	if input == "" {
		return "", fmt.Errorf("empty datetime")
	}

	if tz == "UTC" {
		return input, nil
	}

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

	f.NewSheet("hold_details")
	f.SetSheetRow("hold_details", "A1", &HoldDetailsHeader)
	f.SetSheetRow("hold_details", "A2",
		&[]interface{}{
			lhd.HoldDetail.MatterID,
			lhd.HoldDetail.HoldName,
			lhd.HoldDetail.HoldNoticeSubject,
			lhd.HoldDetail.HoldNoticeBody,
			lhd.HoldDetail.HoldNoticeTitle,
			lhd.HoldDetail.HoldNoticeAttachmentNames,
		},
	)

	f.NewSheet("custodian_details")
	f.SetSheetRow("custodian_details", "A1", &CustodianDetailsHeader)
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

		f.SetSheetRow("custodian_details", fmt.Sprintf("A%d", i+2), &row)
	}

	// delete detaful sheet "Sheet1"
	f.DeleteSheet("Sheet1")

	if err = f.SaveAs(fmt.Sprintf("%s/legal_hold_details.xlsx", dir)); err != nil {
		return err
	}

	return nil
}
