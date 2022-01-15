package api

import (
	"github.com/geocodio/geocodio-cli/output"
	"github.com/urfave/cli/v2"
	"strconv"
)

func ValidateSpreadsheetJobId(c *cli.Context) (int, error) {
	spreadsheetJobId, err := strconv.Atoi(c.Args().First())
	if err != nil || spreadsheetJobId <= 0 {
		return 0, output.ErrorStringAndExit("Invalid spreadsheet job id specified")
	}
	return spreadsheetJobId, nil
}
