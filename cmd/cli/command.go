package main

import (
	"encoding/json"
	"fmt"
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

func checkTimezone(tz string) error {
	allowedTimezones := map[string]bool{
		"CST": true,
		"EST": true,
		"MST": true,
		"PST": true,
		"UTC": true,
	}
	if !allowedTimezones[tz] {
		return fmt.Errorf("timezone %s is not supported (UTC|CST|EST|MST|PST only)", tz)

	}
	return nil
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
			GetGroupsCmd,
			GetMattersCmd,
			GetLegalholdsCmd,
			GetSilentholdsCmd,
		},
	}

	ImportCmd = &cli.Command{
		Name: "import",
		Subcommands: []*cli.Command{
			ImportLegalholdsCmd,
			ImportSilentholdsCmd,
		},
	}

	ImportLegalholdsCmd = &cli.Command{
		Name:     "legalholds",
		Category: "import",
		Action:   execute,
		Flags: []cli.Flag{
			AtttachmentDirectory,
			Excel,
			Zipfile,
			Timezone,
			HoldName,
			MatterName,
			CheckInputOnly,
		},
		Before: func(c *cli.Context) error {
			return checkTimezone(c.String("timezone"))
		},
	}

	ImportSilentholdsCmd = &cli.Command{
		Name:     "silentholds",
		Category: "import",
		Action:   execute,
		Flags: []cli.Flag{
			Excel,
			Zipfile,
			Timezone,
			HoldName,
			MatterName,
			CheckInputOnly,
		},
		Before: func(c *cli.Context) error {
			return checkTimezone(c.String("timezone"))
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
		Flags: append(
			DefaultListOptions,
			GroupID,
		),
	}

	GetGroupsCmd = &cli.Command{
		Name:     "groups",
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

	GetSilentholdsCmd = &cli.Command{
		Name:     "silentholds",
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
		switch ctx.Command.Name {
		case "legalholds":
			return importLegalholds(ctx)
		case "silentholds":
			return importSilentholds(ctx)
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
		case "silentholds":
			return getSilentholds(ctx)
		case "groups":
			return getGroups(ctx)
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

	client := NewClient(ctx)

	zip := ctx.String("zipfile")
	if len(zip) > 0 {
		if _, err := os.Stat(zip); err == nil || os.IsExist(err) {
			if _, err = client.ImportLegalhold(zip); err != nil {
				log.Error().Msgf("failed to import Legalhold from zip file: %s, %s", zip, err)
				return err
			}
		}
		return err
	}

	tz := otlh.GetTimezoneLocation(ctx.String("timezone"))
	log.Debug().Msgf("timezone: %s", tz)

	imp := importer.NewExcelImporter().
		WithClient(client).
		WithExcel(ctx.String("excel")).
		WithTimezone(tz).
		WithMatterName(ctx.String("matterName")).
		WithHoldName(ctx.String("holdName")).
		Legalhold(). // convert to legalholdExcelImporter type
		WithAttachmentDirectory(ctx.String("attachmentDirectory"))

	err = imp.LoadLegalholdData()
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

func importSilentholds(ctx *cli.Context) error {
	var err error

	client := NewClient(ctx)

	zip := ctx.String("zipfile")
	if len(zip) > 0 {
		if _, err := os.Stat(zip); err == nil || os.IsExist(err) {
			if _, err = client.ImportSilenthold(zip); err != nil {
				log.Error().Msgf("failed to import Silenthold from zip file: %s, %s", zip, err)
				return err
			}
		}
		return err
	}

	tz := otlh.GetTimezoneLocation(ctx.String("timezone"))
	log.Debug().Msgf("timezone: %s", tz)

	imp := importer.NewExcelImporter().
		WithClient(NewClient(ctx)).
		WithExcel(ctx.String("excel")).
		WithTimezone(tz).
		WithMatterName(ctx.String("matterName")).
		WithHoldName(ctx.String("holdName")).
		Silenthold()

	err = imp.LoadSilentholdData()
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

	var req otlh.Requestor

	client := NewClient(ctx)

	b := otlh.NewRequest().WithTenant(client.Tenant()).Get().Custodian()
	opts := listOptions(ctx)

	if ctx.Int("id") > 0 {
		req, _ = b.WithID(ctx.Int("id")).Build()
		v, err = client.GetCustodian(req)
	} else {
		req, _ = b.Build()
		if ctx.Bool("all") {
			v, err = client.GetAllCustodians(req, opts)
		} else {
			v, err = client.GetCustodians(req, opts)
		}
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

	var req otlh.Requestor

	client := NewClient(ctx)

	b := otlh.NewRequest().WithTenant(client.Tenant()).Get().Folder()
	opts := listOptions(ctx)

	if ctx.Int("id") > 0 {
		req, _ = b.WithID(ctx.Int("id")).Build()
		v, err = client.GetFolder(req)
	} else {
		if ctx.Int("groupID") > 0 {
			b.WithGroupID(ctx.Int("groupID"))
		}

		req, _ = b.Build()
		if ctx.Bool("all") {
			v, err = client.GetAllFolders(req, opts)
		} else {
			v, err = client.GetFolders(req, opts)
		}
	}

	if err != nil {
		return err
	}

	printer := otlh.NewPrinter().JSON().Build()
	printer.Print(v)
	return nil
}

func getGroups(ctx *cli.Context) error {
	var err error
	var v any

	var req otlh.Requestor

	client := NewClient(ctx)
	opts := listOptions(ctx)

	b := otlh.NewRequest().WithTenant(client.Tenant()).Get().Group()
	if ctx.Int("id") > 0 {
		req, _ = b.WithID(ctx.Int("id")).Build()
		v, err = client.GetGroup(req)
	} else {
		req, _ = b.Build()
		if ctx.Bool("all") {
			v, err = client.GetAllGroups(req, opts)
		} else {
			v, err = client.GetGroups(req, opts)
		}
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

	var req otlh.Requestor

	client := NewClient(ctx)
	opts := listOptions(ctx)

	b := otlh.NewRequest().WithTenant(client.Tenant()).Get().Matter()
	if ctx.Int("id") > 0 {
		req, _ = b.WithID(ctx.Int("id")).Build()
		v, err = client.GetMatter(req)
	} else {
		req, _ = b.Build()
		if ctx.Bool("all") {
			v, err = client.GetAllMatters(req, opts)
		} else {
			v, err = client.GetMatters(req, opts)
		}
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

	var req otlh.Requestor

	client := NewClient(ctx)
	opts := listOptions(ctx)

	b := otlh.NewRequest().WithTenant(client.Tenant()).Get().Legalhold()
	if ctx.Int("id") > 0 {
		req, _ = b.WithID(ctx.Int("id")).Build()
		v, err = client.GetLegalhold(req)
	} else {
		req, _ = b.Build()
		if ctx.Bool("all") {
			v, err = client.GetAllLegalholds(req, opts)
		} else {
			v, err = client.GetLegalholds(req, opts)
		}
	}

	if err != nil {
		return err
	}

	printer := otlh.NewPrinter().JSON().Build()
	printer.Print(v)
	return nil
}

func getSilentholds(ctx *cli.Context) error {
	var err error
	var v any

	var req otlh.Requestor

	client := NewClient(ctx)
	opts := listOptions(ctx)

	b := otlh.NewRequest().WithTenant(client.Tenant()).Get().Silenthold()
	if ctx.Int("id") > 0 {
		req, _ = b.WithID(ctx.Int("id")).Build()
		v, err = client.GetSilenthold(req)
	} else {
		req, _ = b.Build()
		if ctx.Bool("all") {
			v, err = client.GetAllSilentholds(req, opts)
		} else {
			v, err = client.GetSilentholds(req, opts)
		}
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
