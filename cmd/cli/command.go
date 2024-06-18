package main

import (
	"encoding/json"
	"os"

	otlh "github.com/xifanyan/otlh/pkg"
	importer "github.com/xifanyan/otlh/pkg/importer"
	"github.com/xifanyan/otlh/pkg/verifier"

	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
)

type ClientConfig struct {
	Domain    string `json:"domain"`
	Port      int    `json:"port"`
	Tenant    string `json:"tenant"`
	AuthToken string `json:"authToken"`
}

var (
	CreateCmd = &cli.Command{
		Name: "create",
		Subcommands: []*cli.Command{
			CreateFolderCmd,
			CreateMatterCmd,
		},
	}

	GetCmd = &cli.Command{
		Name: "get",
		Subcommands: []*cli.Command{
			GetCustodiansCmd,
			GetFoldersCmd,
			GetMattersCmd,
			GetLegalholdsCmd,
		},
	}

	ImportCmd = &cli.Command{
		Name: "import",
		Subcommands: []*cli.Command{
			ImportLegalholdsCmd,
		},
	}

	ImportLegalholdsCmd = &cli.Command{
		Name:     "legalholds",
		Category: "import",
		Action:   execute,
		Flags: []cli.Flag{
			AtttachmentDirectory,
			Excel,
			Timezone,
			HoldName,
			MatterName,
			CheckInputOnly,
		},
	}

	VerifyCmd = &cli.Command{
		Name: "verify",
		Subcommands: []*cli.Command{
			VerifyCustodiansCmd,
		},
	}

	VerifyCustodiansCmd = &cli.Command{
		Name:     "custodians",
		Category: "verify",
		Action:   execute,
		Flags: []cli.Flag{
			CSV,
		},
	}

	GetCustodiansCmd = &cli.Command{
		Name:     "custodians",
		Category: "get",
		Action:   execute,
		Flags:    DefaultListOptions,
	}

	GetFoldersCmd = &cli.Command{
		Name:     "folders",
		Category: "get",
		Action:   execute,
		Flags:    DefaultListOptions,
	}

	GetMattersCmd = &cli.Command{
		Name:     "matters",
		Category: "get",
		Action:   execute,
		Flags:    DefaultListOptions,
	}

	GetLegalholdsCmd = &cli.Command{
		Name:     "legalholds",
		Category: "get",
		Action:   execute,
		Flags:    DefaultListOptions,
	}

	CreateFolderCmd = &cli.Command{
		Name:     "folder",
		Category: "create",
		Action:   execute,
		Flags: []cli.Flag{
			Name,
		},
	}

	CreateMatterCmd = &cli.Command{
		Name:     "matter",
		Category: "create",
		Action:   execute,
		Flags: []cli.Flag{
			Name,
			FolderID,
		},
	}

	Commands = []*cli.Command{
		CreateCmd,
		GetCmd,
		ImportCmd,
		VerifyCmd,
	}
)

func loadConfig(file string) (*ClientConfig, error) {
	data, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	var config ClientConfig
	err = json.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

func execute(ctx *cli.Context) error {
	switch ctx.Command.Category {
	case "create":
		switch ctx.Command.Name {
		case "folder":
			return createFolder(ctx)
		case "matter":
			return createMatter(ctx)
		}
	case "import":
		if ctx.Command.Name == "legalholds" {
			return importLegalholds(ctx)
		}
	case "get":
		switch ctx.Command.Name {
		case "custodians":
			return getCustodians(ctx)
		case "folders":
			return getFolders(ctx)
		case "matters":
			return getMatters(ctx)
		case "legalholds":
			return getLegalholds(ctx)
		}
	case "verify":
		switch ctx.Command.Name {
		case "custodians":
			return verifyCustodians(ctx)
		}
	}
	return nil
}

func NewClient(ctx *cli.Context) *otlh.Client {

	if _, err := os.Stat(ctx.String("config")); os.IsNotExist(err) {
		log.Debug().Msgf("config file not found: %s, use command args", ctx.String("config"))
		return otlh.NewClientBuilder().
			WithDomain(ctx.String("domain")).
			WithPort(ctx.Int("port")).
			WithTenant(ctx.String("tenant")).
			WithAuthToken(ctx.String("authToken")).
			Build()
	}

	log.Debug().Msgf("use config file: %s", ctx.String("config"))
	cfg, err := loadConfig(ctx.String("config"))
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load config")
	}

	return otlh.NewClientBuilder().
		WithDomain(cfg.Domain).
		WithPort(cfg.Port).
		WithTenant(cfg.Tenant).
		WithAuthToken(cfg.AuthToken).
		Build()

}

func importLegalholds(ctx *cli.Context) error {
	var err error

	tz := otlh.GetTimezoneLocation(ctx.String("timezone"))
	log.Debug().Msgf("timezone: %s", tz)

	imp := importer.NewLegalholdExcelImporter().
		WithClient(NewClient(ctx)).
		WithExcel(ctx.String("excel")).
		WithTimezone(tz).
		WithMatterName(ctx.String("matterName")).
		WithHoldName(ctx.String("holdName")).
		WithAttachmentDirectory(ctx.String("attachmentDirectory"))

	err = imp.LoadHoldData()
	if err != nil {
		return err
	}

	err = imp.PerformDataIntegrityCheck()
	if err != nil {
		return err
	}

	if ctx.Bool("checkInputOnly") {
		return nil
	}

	err = imp.Import()
	return err
}

func listOptions(ctx *cli.Context) *otlh.ListOptions {
	return otlh.NewListOptions().
		WithPageNumber(ctx.Int("pageNumber")).
		WithPageSize(ctx.Int("pageSize")).
		WithSort(ctx.String("sort")).
		WithFilterName(ctx.String("filterName")).
		WithFilterTerm(ctx.String("filterTerm"))
}

func getCustodians(ctx *cli.Context) error {
	var err error
	var v any

	client := NewClient(ctx)

	if ctx.Bool("all") {
		client.PrintAllCustodians()
		return nil
	}

	if ctx.Int("id") > 0 {
		v, err = client.GetCustodian(ctx.Int("id"))
	} else {
		v, err = client.GetCustodians(listOptions(ctx))
	}

	if err != nil {
		return err
	}

	printer := otlh.NewPrinter().JSON().Build()
	printer.Print(v)
	return nil
}

func getFolders(ctx *cli.Context) error {
	var err error
	var v any

	client := NewClient(ctx)

	if ctx.Int("id") > 0 {
		v, err = client.GetFolder(ctx.Int("id"))
	} else {
		v, err = client.GetFolders(listOptions(ctx))
	}

	if err != nil {
		return err
	}

	printer := otlh.NewPrinter().JSON().Build()
	printer.Print(v)
	return nil
}

func getMatters(ctx *cli.Context) error {
	var err error
	var v any

	client := NewClient(ctx)

	if ctx.Int("id") > 0 {
		v, err = client.GetMatter(ctx.Int("id"))
	} else {
		v, err = client.GetMatters(listOptions(ctx))
	}

	if err != nil {
		return err
	}

	printer := otlh.NewPrinter().JSON().Build()
	printer.Print(v)
	return nil
}

func getLegalholds(ctx *cli.Context) error {
	var err error
	var v any

	client := NewClient(ctx)

	if ctx.Int("id") > 0 {
		v, err = client.GetLegalhold(ctx.Int("id"))
	} else {
		v, err = client.GetLegalholds(listOptions(ctx))
	}

	if err != nil {
		return err
	}

	printer := otlh.NewPrinter().JSON().Build()
	printer.Print(v)
	return nil
}

func createFolder(ctx *cli.Context) error {
	var err error

	client := NewClient(ctx)
	_, err = client.FindOrCreateFolder(ctx.String("name"))

	return err
}

func createMatter(ctx *cli.Context) error {
	var err error

	client := NewClient(ctx)
	_, err = client.FindOrCreateMatter(ctx.String("name"), ctx.Int("folderID"))

	return err
}

func verifyCustodians(ctx *cli.Context) error {
	var err error

	vf := verifier.NewCustodianVerifier().
		WithClient(NewClient(ctx))

		/*
			custodiansFromOTLH, err := vf.LoadAllCustodiansFromOTLH()
			if err != nil {
				return err
			}
		*/

	// custodiansFromCSV, err := vf.LoadCustodiansFromCSV(ctx.String("csv"))
	_, err = vf.LoadCustodiansFromCSV(ctx.String("csv"))
	if err != nil {
		return err
	}

	/*
		for custodian := range custodiansFromCSV {
			if _, ok := custodiansFromOTLH[custodian]; ok {
				continue
			} else {
				log.Error().Msgf("Custodian %s not found in OTLH", custodian)
			}
		}
	*/

	return nil
}
