package verifier

import (
	"os"

	"github.com/gocarina/gocsv"
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

func (cv *CustodianVerifier) LoadAllCustodiansFromOTLH() (map[string]struct{}, error) {
	var custodians map[string]struct{} = make(map[string]struct{})

	opts := otlh.NewListOptions().WithPageSize(100)

	custodianCh, errCh := cv.client.GetAllCustodians(opts)
	for {
		select {
		case custodian, ok := <-custodianCh:
			if !ok {
				custodianCh = nil
			} else {
				custodians[custodian.Email] = struct{}{}
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

	log.Debug().Msgf("Loaded %d custodians", len(custodians))
	return custodians, nil
}

type CustodianInput struct {
	EmployeeID      string `csv:"EmployeeID"`
	Name            string `csv:"Name"`
	Emailaddress    string `csv:"Emailaddress"`
	EmployeeStatus  string `csv:"EmployeeStatus`
	EmployeeType    string `csv:"EmployeeType`
	Title           string `csv:"Title"`
	OfficePhone     string `csv:"OfficePhone`
	Department      string `csv:"Department`
	Location        string `csv:"Location`
	SupervisorName  string `csv:"SupervisorName`
	SupervisorEmail string `csv:"SupervisorEmail`
	Function        string `csv:"Function`
	Business        string `csv:"Business`
	Notes           string `csv:"Notes`
}

func (cv *CustodianVerifier) LoadCustodiansFromCSV(fileName string) (map[string]struct{}, error) {

	file, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Read the CSV file into a slice of Person structs
	var custodians []*CustodianInput
	if err := gocsv.UnmarshalFile(file, &custodians); err != nil {
		return nil, err
	}

	// Create a map to store the data
	dataMap := make(map[string]struct{})

	// Iterate over the slice and use the Name field as the key
	for _, custodian := range custodians {
		if _, ok := dataMap[custodian.Emailaddress]; ok {
			log.Error().Msgf("Duplicate custodian: %s", custodian.Emailaddress)
			continue
		}
		dataMap[custodian.Emailaddress] = struct{}{}
	}

	log.Debug().Msgf("Custodians loaded: %d", len(dataMap))
	return dataMap, nil
}
