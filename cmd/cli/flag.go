package main

import "github.com/urfave/cli/v2"

var (
	Excel = &cli.StringFlag{
		Name:    "excel",
		Aliases: []string{"e"},
		Usage:   "excel file used for legalhold import",
	}

	CSV = &cli.StringFlag{
		Name:  "csv",
		Usage: "path for csv input",
	}

	Name = &cli.StringFlag{
		Name:  "name",
		Usage: "name",
	}

	All = &cli.BoolFlag{
		Name:  "all",
		Usage: "all",
	}

	ID = &cli.IntFlag{
		Name:  "id",
		Usage: "id",
	}

	FolderID = &cli.IntFlag{
		Name:  "folderID",
		Usage: "folderID",
	}

	GroupID = &cli.IntFlag{
		Name:  "groupID",
		Usage: "group id",
	}

	CustodianID = &cli.IntFlag{
		Name:  "custodianID",
		Usage: "custodian id",
	}

	LegalHoldID = &cli.IntFlag{
		Name:  "legalHoldID",
		Usage: "legalhold id",
	}

	SilentHoldID = &cli.IntFlag{
		Name:  "silentHoldID",
		Usage: "silenthold id",
	}

	CustodianGroupID = &cli.IntFlag{
		Name:  "custodianGroupID",
		Usage: "custodian group id",
	}

	MatterID = &cli.IntFlag{
		Name:  "matterID",
		Usage: "matter id",
	}

	MatterName = &cli.StringFlag{
		Name:    "matterName",
		Aliases: []string{"mn"},
		Usage:   "matter name",
	}

	HoldName = &cli.StringFlag{
		Name:    "holdName",
		Aliases: []string{"hn"},
		Usage:   "hold name",
	}

	PageNumber = &cli.IntFlag{
		Name:    "pageNumber",
		Aliases: []string{"pn"},
		Usage:   "page number",
	}

	PageSize = &cli.IntFlag{
		Name:    "pageSize",
		Aliases: []string{"ps"},
		Usage:   "page size",
		Value:   50,
	}

	Sort = &cli.StringFlag{
		Name:    "sort",
		Aliases: []string{"s"},
		Usage:   "sort",
	}

	FilterTerm = &cli.StringFlag{
		Name:    "filterTerm",
		Aliases: []string{"t"},
		Usage:   "filter[term]",
	}

	FilterName = &cli.StringFlag{
		Name:    "filterName",
		Aliases: []string{"n"},
		Usage:   "filter[name]",
	}

	CheckInputOnly = &cli.BoolFlag{
		Name:    "checkInputOnly",
		Aliases: []string{"ci"},
		Usage:   "check input only",
	}

	AtttachmentDirectory = &cli.StringFlag{
		Name:    "attachmentDirectory",
		Aliases: []string{"ad"},
		Usage:   "attachment directory",
		Value:   ".",
	}

	Timezone = &cli.StringFlag{
		Name:    "timezone",
		Aliases: []string{"tz"},
		Usage:   "timezone for dates used in input file, supproted timezones: PST|EST|MST|CST",
		Value:   "UTC",
	}

	Zipfile = &cli.StringFlag{
		Name:    "zipfile",
		Aliases: []string{"z"},
		Usage:   "zip package for importing holds e.g., legal_hold_details.zip",
		Value:   "",
	}
)

var DefaultListOptions = []cli.Flag{
	All,
	ID,
	PageNumber,
	PageSize,
	Sort,
	FilterTerm,
	FilterName,
}
