package verifier

import (
	"github.com/rs/zerolog/log"
	otlh "github.com/xifanyan/otlh/pkg"
)

type CustodianVerifier struct {
	client *otlh.Client
}

func NewCustodianVerifier() *CustodianVerifier {
	return &CustodianVerifier{}
}

func (cv *CustodianVerifier) WithClient(client *otlh.Client) *CustodianVerifier {
	cv.client = client
	return cv
}

func (cv *CustodianVerifier) LoadAllCustodiansFromOTLH() error {
	var custodians map[string]otlh.Custodian = make(map[string]otlh.Custodian)

	opts := otlh.NewListOptions().WithPageSize(100)

	custodianCh, errCh := cv.client.GetAllCustodians(opts)
	for {
		select {
		case custodian, ok := <-custodianCh:
			if !ok {
				custodianCh = nil
			} else {
				custodians[custodian.Email] = custodian
			}
		case err, ok := <-errCh:
			if !ok {
				errCh = nil
			} else {
				log.Error().Err(err)
			}
		}

		if custodianCh == nil && errCh == nil {
			break
		}
	}

	return nil
}
