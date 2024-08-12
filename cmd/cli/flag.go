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

	SkipInputCheck = &cli.BoolFlag{
		Name:    "skipInputCheck",
		Aliases: []string{"sc"},
		Usage:   "skip input check",
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
