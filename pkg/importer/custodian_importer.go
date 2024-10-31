package importer

import (
	"encoding/json"
	"os"
	"strings"

	"github.com/gocarina/gocsv"
	"github.com/rs/zerolog/log"
	otlh "github.com/xifanyan/otlh/pkg"
)

type CustodianImporter struct {
	input      string
	batchSize  int
	custodians []otlh.CustodianInputData
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

func (b *CustodianImporterBuilder) WithInput(name string) *CustodianImporterBuilder {
	b.input = name
	return b
}

func (b *CustodianImporterBuilder) WithBatchSize(size int) *CustodianImporterBuilder {
	b.batchSize = size
	return b
}

func parseCustodianJson(input string) ([]otlh.CustodianInputData, error) {
	data, err := os.ReadFile(input)
	if err != nil {
		return nil, err
	}

	var custodians []otlh.CustodianInputData

	err = json.Unmarshal(data, &custodians)
	if err != nil {
		return nil, err
	}
	return custodians, nil
}

func parseCustodianCsv(input string) ([]otlh.CustodianInputData, error) {
	data, err := os.ReadFile(input)
	if err != nil {
		return nil, err
	}

	var custodians []otlh.CustodianInputData

	err = gocsv.UnmarshalBytes(data, &custodians)
	if err != nil {
		return nil, err
	}
	return custodians, nil
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

	log.Debug().Msgf("custodians %+v:", b.custodians)
	return b.CustodianImporter, nil
}

func (c *CustodianImporter) Import() error {
	return c.client.ImportCustodians(c.custodians, c.batchSize)
}
