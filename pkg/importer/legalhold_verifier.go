package importer

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/rs/zerolog/log"
)

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

func (imptr *LegalholdExcelImporter) PerformDataIntegrityCheck() error {
	var err error

	if err = imptr.baselineDataIntegrityCheck(); err != nil {
		return err
	}

	log.Debug().Msg("Verify attachment files ...")
	if err = imptr.verifyAttachmentFiles(); err != nil {
		return err
	}

	log.Debug().Msg("[PASS] Data integrity check")

	return err
}
