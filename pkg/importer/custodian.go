package importer

import (
	"encoding/json"
	"os"
	"strings"

	"github.com/gocarina/gocsv"
	otlh "github.com/xifanyan/otlh/pkg"
)

type CustodianImporter struct {
	input      string
	custodians []otlh.Custodian
	client     *otlh.Client
}

type CustodianImporterBuilder struct {
	*CustodianImporter
}

func NewCustodianImporterBuilder() *CustodianImporterBuilder {
	return &CustodianImporterBuilder{
		CustodianImporter: &CustodianImporter{},
	}
}

func (b *CustodianImporterBuilder) WithClient(client *otlh.Client) *CustodianImporterBuilder {
	b.client = client
	return b
}

func (b *CustodianImporterBuilder) WithCSV(csv string) *CustodianImporterBuilder {
	b.input = csv
	return b
}

func (b *CustodianImporterBuilder) WithJSON(json string) *CustodianImporterBuilder {
	b.input = json
	return b
}

func parseCustodianJson(input string) ([]otlh.Custodian, error) {
	data, err := os.ReadFile(input)
	if err != nil {
		return nil, err
	}

	var custodians []otlh.Custodian

	err = json.Unmarshal(data, &custodians)
	if err != nil {
		return nil, err
	}
	return custodians, nil
}

func parseCustodianCsv(input string) ([]otlh.Custodian, error) {
	data, err := os.ReadFile(input)
	if err != nil {
		return nil, err
	}

	var custodians []otlh.Custodian

	err = gocsv.UnmarshalBytes(data, &custodians)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (b *CustodianImporterBuilder) Build() (*CustodianImporter, error) {
	var err error

	switch {
	case strings.HasSuffix(b.input, ".json"):
		b.custodians, err = parseCustodianJson(b.input)
	case strings.HasSuffix(b.input, ".csv"):
		b.custodians, err = parseCustodianCsv(b.input)
	}

	if err != nil {
		return nil, err
	}

	return b.CustodianImporter, nil
}
