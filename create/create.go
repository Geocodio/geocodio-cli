package create

import (
	"fmt"
	"github.com/dustin/go-humanize"
	"github.com/geocodio/geocodio-cli/api"
	"github.com/geocodio/geocodio-cli/output"
	"github.com/urfave/cli/v2"
	"io"
	"os"
	"strings"
)

func RegisterCommand() *cli.Command {
	var command *cli.Command
	command = new(cli.Command)
	command.Name = "create"
	command.Usage = "Geocode a new spreadsheet"
	command.Action = geocode
	command.UsageText = "filename format"
	command.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:  "direction",
			Value: "forward",
			Usage: "Direction can either be \"forward\" (address to coordinate) or \"reverse\" (coordinate to address)",
		},
		&cli.StringFlag{
			Name:  "fields",
			Usage: "A comma-separated list of fields to append to the geocoding job. Read more here https://www.geocod.io/docs/#fields",
		},
	}

	return command
}

func geocode(c *cli.Context) error {
	direction, format, file, err := validateInput(c)
	defer file.Close()

	if err != nil {
		return err
	}

	fields := c.String("fields")

	body, err := api.Upload(file, direction, format, fields, c)
	if err != nil {
		return cli.Exit(err, 1)
	}

	job := api.SpreadsheetJob{}
	if err = api.ParseJson(body, &job); err != nil {
		return err
	}

	if err := validateResponse(job); err != nil {
		return err
	}

	outputJob(c.App.Writer, job)

	return nil
}

func validateInput(c *cli.Context) (string, string, *os.File, error) {
	direction := c.String("direction")
	if direction != "forward" && direction != "reverse" {
		return "", "", nil, output.ErrorStringAndExit("Please specify a valid direction. Valid values are \"forward\" or \"reverse\"")
	}

	filename := c.Args().Get(0)
	format := c.Args().Get(1)

	if len(filename) <= 0 {
		return "", "", nil, output.ErrorStringAndExit("Please specify a file to geocode")
	}

	if len(format) <= 0 || !strings.Contains(format, "{{") {
		return "", "", nil, output.ErrorStringAndExit("Please specify a geocoding format. The format is used to tell the geocoder which spreadsheet columns to use for geocoding. Read more here https://www.geocod.io/docs/#format-syntax")
	}

	file, err := os.Open(filename)

	if err != nil {
		return "", "", nil, output.ErrorStringAndExit(fmt.Sprintf("Could not open %s", filename))
	}

	return direction, format, file, nil
}

func validateResponse(job api.SpreadsheetJob) error {
	if job.Message != "" {
		return output.ErrorStringAndExit(job.Message)
	}

	if job.Error != "" {
		return output.ErrorStringAndExit(job.Error)
	}

	return nil
}

func outputJob(w io.Writer, job api.SpreadsheetJob) {
	output.Success(w, "Spreadsheet job created")
	fmt.Fprint(w, "\n")

	fmt.Fprintf(w, "Id: %d\n", job.Id)
	fmt.Fprintf(w, "Filename: %s\n", job.File.Filename)
	fmt.Fprintf(w, "Headers: %s\n", strings.Join(job.File.Headers, " | "))
	fmt.Fprintf(w, "Rows: %s\n", humanize.Comma(int64(job.File.EstimatedRowsCount)))


	fmt.Fprint(w, "\n\tTo see job status\n")
	fmt.Fprintf(w, "\t$ %s status %d", os.Args[0], job.Id)
}