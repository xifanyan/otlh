package importer

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
	otlh "github.com/xifanyan/otlh/pkg"
)

func (imptr *LegalholdExcelImporter) verifyHeader(row []string, header []string) error {
	if len(row) != len(header) {
		return fmt.Errorf("header length mismatch")
	}

	for i, col := range row {
		x := strings.ToLower(strings.TrimSpace(col))
		y := strings.ToLower(strings.TrimSpace(header[i]))
		if x != y {
			return fmt.Errorf("invalid header: %s", col)
		}
	}
	return nil
}

func (imptr *LegalholdExcelImporter) verifySameMatterUnderSameFolder() error {
	var verr *ValidationError = newValidationError(ErrorSameMatterUnderDifferentFolders)

	uniq := make(map[string]string)

	for i, entry := range imptr.entries {
		if _, ok := uniq[entry.MatterName]; !ok {
			uniq[entry.MatterName] = entry.FolderName
			continue
		}

		if uniq[entry.MatterName] != entry.FolderName {
			verr.add(fmt.Errorf("line #%d: matter %s is not under the same folder", imptr.lineNnumberOfHeader+i+1, entry.MatterName))
		}
	}

	if verr.hasErrors() {
		return verr
	}

	return nil
}

// VerifySameCustodianEmailUnderSameCustodianName verifies that each custodian
// email is under the same custodian name. Returns error if not.
func (imptr *LegalholdExcelImporter) verifySameCustodianEmailUnderSameCustodianName() error {
	var verr *ValidationError = newValidationError(ErrorSameCustodianEmailUnderDifferentCustodianName)

	// create a map to track custodian emails and their corresponding names
	uniq := make(map[string]string)

	for i, entry := range imptr.entries {
		// if custodian email is not in the map, add it
		if _, ok := uniq[entry.CustodianEmail]; !ok {
			uniq[entry.CustodianEmail] = entry.CustodianName
			continue
		}

		// if custodian email is already in the map, make sure it is under the
		// same custodian name
		if uniq[entry.CustodianEmail] != entry.CustodianName {
			verr.add(fmt.Errorf("line #%d: custodian email %s is not under the same name \"%s\" vs \"%s\"",
				imptr.lineNnumberOfHeader+i+1,
				entry.CustodianEmail,
				uniq[entry.CustodianEmail],
				entry.CustodianName,
			))
		}
	}

	if verr.hasErrors() {
		return verr
	}

	return nil
}

func (imptr *LegalholdExcelImporter) verifyAttachmentFiles() error {
	var verr *ValidationError = newValidationError(ErrorAttachmentFileNotFound)

	for attachment := range imptr.collections.UniqueAttachmentNames {
		path := filepath.Join(imptr.attachmentDirectory, attachment)
		log.Debug().Msgf("verifying attachment file %s", path)
		_, err := os.Stat(path)
		if os.IsNotExist(err) {
			verr.add(fmt.Errorf("attachment file %s does not exist", path))
		}
	}

	if verr.hasErrors() {
		return verr
	}
	return nil
}

func (imptr *LegalholdExcelImporter) verifyLastIssued() error {
	var verr *ValidationError = newValidationError(ErrorRequiredLastIssued)

	for i, entry := range imptr.entries {
		if entry.LastIssued == "" {
			verr.add(fmt.Errorf("line #%d: matter [%s] - hold [%s] does not have LastIssued field", imptr.lineNnumberOfHeader+i+1, entry.MatterName, entry.HoldName))
		}

		if _, err := time.Parse(INPUT_TIME_FORMAT, entry.LastIssued); err != nil {
			verr.add(fmt.Errorf("line #%d: matter [%s] - hold [%s] - LastIssued field is invalid: %s", imptr.lineNnumberOfHeader+i+1, entry.MatterName, entry.HoldName, err))
		}
	}

	if verr.hasErrors() {
		return verr
	}
	return nil
}

func (imptr *LegalholdExcelImporter) verifyResponseDate() error {
	var verr *ValidationError = newValidationError(ErrorRequiredLastIssued)

	for i, entry := range imptr.entries {
		if entry.ResponseDate != "" {
			if _, err := time.Parse(INPUT_TIME_FORMAT, entry.ResponseDate); err != nil {
				verr.add(fmt.Errorf("line #%d: matter [%s] - hold [%s] - ResponseDate field is invalid: %s", imptr.lineNnumberOfHeader+i+1, entry.MatterName, entry.HoldName, err))
			}
		}
	}

	if verr.hasErrors() {
		return verr
	}
	return nil
}

func (imptr *LegalholdExcelImporter) verifyReleasedAt() error {
	var verr *ValidationError = newValidationError(ErrorRequiredLastIssued)

	for i, entry := range imptr.entries {
		if entry.ReleasedAt != "" {
			if _, err := time.Parse(INPUT_TIME_FORMAT, entry.ReleasedAt); err != nil {
				verr.add(fmt.Errorf("line #%d: matter [%s] - hold [%s] - ReleasedAt field is invalid: %s", imptr.lineNnumberOfHeader+i+1, entry.MatterName, entry.HoldName, err))
			}
		}
	}

	if verr.hasErrors() {
		return verr
	}
	return nil
}

func (imptr *LegalholdExcelImporter) verifyEmailAddress() error {
	var verr *ValidationError = newValidationError(ErrorInvalidEmailAddress)
	for i, entry := range imptr.entries {
		if !otlh.IsValidEmailAddress(entry.CustodianEmail) {
			verr.add(fmt.Errorf("line #%d: matter [%s] - hold [%s] - custodian email [%s] is invalid", imptr.lineNnumberOfHeader+i+1, entry.MatterName, entry.HoldName, entry.CustodianEmail))
		}
	}

	if verr.hasErrors() {
		return verr
	}
	return nil
}

func (imptr *LegalholdExcelImporter) PerformDataIntegrityCheck() error {
	var err error
	// verify the same matter under same folder
	log.Debug().Msg("Verify same matter under same folder ...")
	if err = imptr.verifySameMatterUnderSameFolder(); err != nil {
		return err
	}

	log.Debug().Msg("Verify custodian email ...")
	if err = imptr.verifyEmailAddress(); err != nil {
		return err
	}

	log.Debug().Msg("Verify same custodian under the same email ...")
	if err = imptr.verifySameCustodianEmailUnderSameCustodianName(); err != nil {
		return err
	}

	log.Debug().Msg("Verify Last Issued field ...")
	if err = imptr.verifyLastIssued(); err != nil {
		return err
	}

	log.Debug().Msg("Verify Response Date field ...")
	if err = imptr.verifyResponseDate(); err != nil {
		return err
	}

	log.Debug().Msg("Verify Released At field ...")
	if err = imptr.verifyReleasedAt(); err != nil {
		return err
	}

	log.Debug().Msg("Verify attachment files ...")
	if err = imptr.verifyAttachmentFiles(); err != nil {
		return err
	}

	log.Debug().Msg("[PASS] Data integrity check")

	return err
}
