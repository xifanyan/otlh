package importer

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/rs/zerolog/log"
	otlh "github.com/xifanyan/otlh/pkg"
	"github.com/xuri/excelize/v2"
)

var MatterTemplateHeader = []string{"Matter Name", "Matter Number", "Case Number", "PO Number", "Caption", "Region", "Business Unit", "Notes", "Inherit Email Config", "Email From", "Email Reply-To", "Name On Outgoing Emails", "Contacts"}

type MatterImporter struct {
	excel               string
	lineNnumberOfHeader int
	entries             []otlh.ImportMatterBody
	client              *otlh.Client
}

type MatterImporterBuilder struct {
	*MatterImporter
}

func NewMatterImporterBuilder() *MatterImporterBuilder {
	return &MatterImporterBuilder{
		MatterImporter: &MatterImporter{},
	}
}

func (b *MatterImporterBuilder) WithClient(client *otlh.Client) *MatterImporterBuilder {
	b.client = client
	return b
}

func (b *MatterImporterBuilder) WithExcel(excel string) *MatterImporterBuilder {
	b.excel = excel
	return b
}

func (b *MatterImporterBuilder) Build() *MatterImporter {
	return b.MatterImporter
}

func (e *MatterImporter) verifyHeader(row []string, header []string) error {
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

func (imptr *MatterImporter) LoadMatterData() error {
	var rows [][]string

	log.Info().Msg("Importing matters from " + imptr.excel)
	f, err := excelize.OpenFile(imptr.excel)
	if err != nil {
		return err
	}
	defer f.Close()

	// get the first sheet
	firstSheet := f.WorkBook.Sheets.Sheet[0].Name
	if rows, err = f.GetRows(firstSheet); err != nil {
		return err
	}

	for l, row := range rows {
		if imptr.verifyHeader(row, MatterTemplateHeader) == nil {
			imptr.lineNnumberOfHeader = l + 1
			log.Debug().Msgf("found header at line #%d", imptr.lineNnumberOfHeader)
			continue
		}

		// after header is found, load all non-empty rows
		if imptr.lineNnumberOfHeader > 0 && len(row) > 0 {
			// make sure the length of the row is equal to the length of the header
			data := make([]string, len(MatterTemplateHeader))
			copy(data, row)

			// IMPORTANT: trim all the values to avoid corner cases during especially when data is part of queryParam
			for i := range data {
				data[i] = strings.TrimSpace(data[i])
			}

			matter, err := imptr.client.FindMatterByName(data[0])
			if err != nil {
				return err
			}

			log.Debug().Msgf("Find Matter: %+v", matter)

			// parse contacts RFC 5322
			contacts, err := parseContacts(data[12])
			if err != nil {
				log.Error().Msgf("Error parsing contacts: %s, skipping ...", err)
				continue
			}

			entry := otlh.ImportMatterBody{
				ID:                       matter.ID,
				Name:                     data[0],
				Number:                   data[1],
				CaseNumber:               data[2],
				PoNumber:                 data[3],
				Caption:                  data[4],
				Region:                   data[5],
				BusinessUnit:             data[6],
				Notes:                    data[7],
				InheritEmailConfig:       data[8] == "TRUE",
				EmailFrom:                data[9],
				EmailReplyTo:             data[10],
				NameOnOutgoingEmails:     data[11],
				MatterContactsAttributes: contacts,
			}

			imptr.entries = append(imptr.entries, entry)
		}
	}

	return nil
}

func (imptr *MatterImporter) Import() error {
	log.Debug().Msg("[Start]: Importing Matters from Excel")
	err := imptr.LoadMatterData()
	if err != nil {
		return err
	}

	if err = imptr.checkDataIntegrity(); err != nil {
		return err
	}

	for _, entry := range imptr.entries {
		log.Debug().Msgf("Matter Input: %+v", entry)

		matter, err := imptr.client.ImportMatter(entry)
		if err != nil {
			return err
		}
		log.Debug().Msgf("Matter Output: %+v", matter)
	}

	return nil
}

func (imptr *MatterImporter) checkDataIntegrity() error {
	var uniqueNames map[string]struct{} = make(map[string]struct{})

	for _, entry := range imptr.entries {
		if _, ok := uniqueNames[entry.Name]; ok {
			return fmt.Errorf("duplicate name: %s", entry.Name)
		}
		uniqueNames[entry.Name] = struct{}{}
	}

	return nil
}

// parseContacts parses a string of contacts in the format "Name <email>" and returns a slice of Contact structs.
//
// Args:
//   input (string): A comma-separated string of contacts, where each contact is in the format "Name <email>".
//
// Returns:
//   ([]otlh.Contact, error): Returns a slice of Contact structs with parsed name and email fields.
//   If the input string is not in the correct format, an error is returned.

func parseContacts(input string) ([]otlh.Contact, error) {
	var contacts []otlh.Contact

	// Regular expression to match the RFC 5322 format
	re := regexp.MustCompile(`([^<]+)<([^>]+)>`)

	// Split the input string by commas
	parts := strings.Split(input, ",")

	for _, part := range parts {
		// Trim any leading or trailing whitespace
		part = strings.TrimSpace(part)

		// Match the pattern
		matches := re.FindStringSubmatch(part)
		if len(matches) != 3 {
			return nil, fmt.Errorf("invalid format: %s", part)
		}

		// Create a new Contact struct and append it to the slice
		contacts = append(contacts, otlh.Contact{
			Name:  strings.TrimSpace(matches[1]),
			Email: strings.TrimSpace(matches[2]),
		})
	}

	return contacts, nil
}
