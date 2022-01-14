package status

import (
	"fmt"
	"github.com/dustin/go-humanize"
	"github.com/geocodio/geocodio-cli/api"
	"github.com/geocodio/geocodio-cli/output"
	"github.com/urfave/cli/v2"
	"net/http"
	"strconv"
)

func RegisterCommand() *cli.Command {
	var command *cli.Command
	command = new(cli.Command)
	command.Name = "status"
	command.Usage = "Query the status for a specific geocoding job"
	command.Action = status
	command.ArgsUsage = "id"

	return command
}

func status(c *cli.Context) error {
	spreadsheetJobId, err := validateInput(c)
	if err != nil {
		return err
	}

	body, err := api.Request(http.MethodGet, fmt.Sprintf("lists/%d", spreadsheetJobId), c)
	if err != nil {
		return output.ErrorAndExit(err)
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

func validateInput(c *cli.Context) (int, error) {
	spreadsheetJobId, err := strconv.Atoi(c.Args().First())
	if err != nil || spreadsheetJobId <= 0 {
		return 0, output.ErrorStringAndExit("Invalid spreadsheet job id specified")
	}
	return spreadsheetJobId, nil
}

func validateResponse(job api.SpreadsheetJob) error {
	if job.Id == 0 {
		return output.ErrorStringAndExit("No spreadsheet job with that id found. Make sure that the id is correct and that you are using the intended API key")
	}

	return nil
}

func outputJob(job api.SpreadsheetJob) {
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
