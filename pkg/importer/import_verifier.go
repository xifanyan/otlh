package importer

import (
	"fmt"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
	otlh "github.com/xifanyan/otlh/pkg"
)

func (e *ExcelImporter) verifyHeader(row []string, header []string) error {
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

func (e *ExcelImporter) VerifyCustodians() error {
	var verr *ValidationError = newValidationError(ErrorCustodianNotFound)

	for email, name := range e.collections.UniqueCustodians {
		err := e.FindCustodianByNameAndEmail(name, email)
		if err != nil {
			verr.add(fmt.Errorf("custodian: %s email: %s not found", name, email))
		}
	}

	if verr.hasErrors() {
		return verr
	}

	return nil
}

func (e *ExcelImporter) verifySameMatterUnderSameFolder() error {
	var verr *ValidationError = newValidationError(ErrorSameMatterUnderDifferentFolders)

	uniq := make(map[string]string)

	for i, entry := range e.entries {
		if _, ok := uniq[entry.MatterName]; !ok {
			uniq[entry.MatterName] = entry.FolderName
			continue
		}

		if uniq[entry.MatterName] != entry.FolderName {
			verr.add(fmt.Errorf("line #%d: matter %s is not under the same folder", e.lineNnumberOfHeader+i+1, entry.MatterName))
		}
	}

	if verr.hasErrors() {
		return verr
	}

	return nil
}

// VerifySameCustodianEmailUnderSameCustodianName verifies that each custodian
// email is under the same custodian name. Returns error if not.
func (e *ExcelImporter) verifySameCustodianEmailUnderSameCustodianName() error {
	var verr *ValidationError = newValidationError(ErrorSameCustodianEmailUnderDifferentCustodianName)

	// create a map to track custodian emails and their corresponding names
	uniq := make(map[string]string)

	for i, entry := range e.entries {
		// if custodian email is not in the map, add it
		if _, ok := uniq[entry.CustodianEmail]; !ok {
			uniq[entry.CustodianEmail] = entry.CustodianName
			continue
		}

		// if custodian email is already in the map, make sure it is under the
		// same custodian name
		if uniq[entry.CustodianEmail] != entry.CustodianName {
			verr.add(fmt.Errorf("line #%d: custodian email %s is not under the same name \"%s\" vs \"%s\"",
				e.lineNnumberOfHeader+i+1,
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

func (e *ExcelImporter) verifyLastIssued() error {
	var verr *ValidationError = newValidationError(ErrorRequiredLastIssued)

	for i, entry := range e.entries {
		if entry.LastIssued == "" {
			verr.add(fmt.Errorf("line #%d: matter [%s] - hold [%s] does not have LastIssued field", e.lineNnumberOfHeader+i+1, entry.MatterName, entry.HoldName))
		}

		if _, err := time.Parse(INPUT_TIME_FORMAT, entry.LastIssued); err != nil {
			verr.add(fmt.Errorf("line #%d: matter [%s] - hold [%s] - LastIssued field is invalid: %s", e.lineNnumberOfHeader+i+1, entry.MatterName, entry.HoldName, err))
		}

		d, _ := convertDateTimeFormat(e.timezone, entry.LastIssued)
		log.Debug().Msgf("- line #%d: [%s] %s -> [UTC] %s", e.lineNnumberOfHeader+i+1, e.timezone, entry.LastIssued, d)
	}

	if verr.hasErrors() {
		return verr
	}
	return nil
}

func (e *ExcelImporter) verifyResponseDate() error {
	var verr *ValidationError = newValidationError(ErrorRequiredLastIssued)

	for i, entry := range e.entries {
		if entry.ResponseDate != "" {
			if _, err := time.Parse(INPUT_TIME_FORMAT, entry.ResponseDate); err != nil {
				verr.add(fmt.Errorf("line #%d: matter [%s] - hold [%s] - ResponseDate field is invalid: %s", e.lineNnumberOfHeader+i+1, entry.MatterName, entry.HoldName, err))
			}
			d, _ := convertDateTimeFormat(e.timezone, entry.ResponseDate)
			log.Debug().Msgf("- line #%d: [%s] %s -> [UTC] %s", e.lineNnumberOfHeader+i+1, e.timezone, entry.ResponseDate, d)
		}
	}

	if verr.hasErrors() {
		return verr
	}
	return nil
}

func (e *ExcelImporter) verifyReleasedAt() error {
	var verr *ValidationError = newValidationError(ErrorRequiredLastIssued)

	for i, entry := range e.entries {
		if entry.ReleasedAt != "" {
			if _, err := time.Parse(INPUT_TIME_FORMAT, entry.ReleasedAt); err != nil {
				verr.add(fmt.Errorf("line #%d: matter [%s] - hold [%s] - ReleasedAt field is invalid: %s", e.lineNnumberOfHeader+i+1, entry.MatterName, entry.HoldName, err))
			}
		}
	}

	if verr.hasErrors() {
		return verr
	}
	return nil
}

func (e *ExcelImporter) verifyEmailAddress() error {
	var verr *ValidationError = newValidationError(ErrorInvalidEmailAddress)
	for i, entry := range e.entries {
		if !otlh.IsValidEmailAddress(entry.CustodianEmail) {
			verr.add(fmt.Errorf("line #%d: matter [%s] - hold [%s] - custodian email [%s] is invalid", e.lineNnumberOfHeader+i+1, entry.MatterName, entry.HoldName, entry.CustodianEmail))
		}
	}

	if verr.hasErrors() {
		return verr
	}
	return nil
}

func (e *ExcelImporter) verifyHoldName() error {
	var verr *ValidationError = newValidationError(ErrorHoldNameTooLong)
	for i, entry := range e.entries {
		if len(entry.HoldName) > MAX_LENGTH_OF_HOLDNAME {
			verr.add(fmt.Errorf("line #%d: matter [%s] - hold [%s] has exceeded number of characters allowed", e.lineNnumberOfHeader+i+1, entry.MatterName, entry.HoldName))
		}
	}

	if verr.hasErrors() {
		return verr
	}
	return nil
}

func (e *ExcelImporter) baselineDataIntegrityCheck() error {
	var err error

	// verify the same matter under same folder
	log.Debug().Msg("Verify same matter under same folder ...")
	if err = e.verifySameMatterUnderSameFolder(); err != nil {
		return err
	}

	log.Debug().Msg("Verify hold name ...")
	if err = e.verifyHoldName(); err != nil {
		return err
	}

	log.Debug().Msg("Verify custodian email ...")
	if err = e.verifyEmailAddress(); err != nil {
		return err
	}

	log.Debug().Msg("Verify same custodian under the same email ...")
	if err = e.verifySameCustodianEmailUnderSameCustodianName(); err != nil {
		return err
	}

	log.Debug().Msg("Verify Last Issued field ...")
	if err = e.verifyLastIssued(); err != nil {
		return err
	}

	log.Debug().Msg("Verify Response Date field ...")
	if err = e.verifyResponseDate(); err != nil {
		return err
	}

	log.Debug().Msg("Verify Released At field ...")
	if err = e.verifyReleasedAt(); err != nil {
		return err
	}

	log.Debug().Msg("Verify custodians ...")
	if err = e.VerifyCustodians(); err != nil {
		return err
	}

	return nil
}
