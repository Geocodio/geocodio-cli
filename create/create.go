package create

import (
	"fmt"
	"github.com/dustin/go-humanize"
	"github.com/geocodio/geocodio-cli/api"
	"github.com/geocodio/geocodio-cli/output"
	"github.com/urfave/cli/v2"
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

	outputJob(job)

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
	defer file.Close()

	return direction, format, file, nil
}

func validateResponse(job api.SpreadsheetJob) error {
	if job.Message != "" {
		return output.ErrorStringAndExit(job.Message)
	}

	return nil
}

func outputJob(job api.SpreadsheetJob) {
	output.Success("Spreadsheet job created")
	fmt.Println("--------------------------")

	fmt.Println("Id:")
	fmt.Printf("\t%d\n\n", job.Id)

	fmt.Println("Filename:")
	fmt.Printf("\t%s\n\n", job.File.Filename)

	fmt.Println("Rows:")
	fmt.Printf("\t%s\n\n", humanize.Comma(int64(job.File.EstimatedRowsCount)))

	fmt.Println("State:")
	fmt.Printf("\t%s\n\n", job.Status.State)

	fmt.Println("Progress:")
	fmt.Printf("\t%.0f%%\n\n", job.Status.Progress)

	fmt.Println("Message:")
	fmt.Printf("\t%s\n\n", job.Status.Message)

	fmt.Println("Time left:")
	timeLeft := job.Status.TimeLeftDescription
	if len(timeLeft) <= 0 {
		timeLeft = "-"
	}
	fmt.Printf("\t%s\n\n", timeLeft)

	fmt.Println("Expires:")
	fmt.Printf("\t%s\n\n", job.ExpiresAt.Format("Mon Jan _2 15:04:05 2006"))
}
