package importer

import (
	"github.com/rs/zerolog/log"
)

func (imptr *SilentholdExcelImporter) PerformDataIntegrityCheck() error {
	var err error

	if err = imptr.baselineDataIntegrityCheck(); err != nil {
		return err
	}

	log.Debug().Msg("[PASS] Data integrity check")

	return nil
}
