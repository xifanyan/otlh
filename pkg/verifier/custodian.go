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
	var custodianMap map[string]struct{} = make(map[string]struct{})

	req, _ := otlh.NewRequest().WithTenant(cv.client.Tenant()).Get().Custodian().Build()
	opts := otlh.NewListOptions().WithPageSize(100)

	custodians, err := cv.client.GetAllCustodians(req, opts)
	if err != nil {
		return nil, err
	}

	for _, custodian := range custodians {
		custodianMap[custodian.Email] = struct{}{}
	}

	log.Debug().Msgf("Loaded %d custodians", len(custodianMap))
	return custodianMap, nil
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
