package status

import (
	"fmt"
	"github.com/dustin/go-humanize"
	"github.com/geocodio/geocodio-cli/api"
	"github.com/geocodio/geocodio-cli/output"
	"github.com/schollz/progressbar/v3"
	"github.com/urfave/cli/v2"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

func RegisterCommand() *cli.Command {
	var command *cli.Command
	command = new(cli.Command)
	command.Name = "status"
	command.Usage = "Query the status for a specific geocoding job"
	command.Action = status
	command.ArgsUsage = "id"
	command.Flags = []cli.Flag{
		&cli.BoolFlag{
			Name:    "follow",
			Aliases: []string{"f"},
			Usage:   "Whether to follow the job status until it is completed",
		},
	}

	return command
}

func status(c *cli.Context) error {
	spreadsheetJobId, err := api.ValidateSpreadsheetJobId(c)
	if err != nil {
		return err
	}

	var job api.SpreadsheetJob
	if err, job = fetchStatus(c, spreadsheetJobId); err != nil {
		return err
	}

	outputJob(c.App.Writer, job)

	if c.Bool("follow") {
		fmt.Printf("\n")
		bar := progressbar.NewOptions(100,
			progressbar.OptionEnableColorCodes(true),
			progressbar.OptionShowBytes(false),
			progressbar.OptionSetPredictTime(false),
			progressbar.OptionSetDescription(fmt.Sprintf("%s", job.File.Filename)),
			progressbar.OptionSetMaxDetailRow(1),
		)

		for {
			if err, job = fetchStatus(c, spreadsheetJobId); err != nil {
				return err
			}

			if job.Status.State != "ENQUEUED" && job.Status.State != "PROCESSING" {
				break
			}

			bar.Set(int(job.Status.Progress))
			bar.AddDetail(job.Status.TimeLeftDescription)

			time.Sleep(5 * time.Second)
		}
	}

	outputDownloadHelpText(c.App.Writer, job)

	return nil
}

func fetchStatus(c *cli.Context, spreadsheetJobId int) (error, api.SpreadsheetJob) {
	body, _, err := api.Request(http.MethodGet, fmt.Sprintf("lists/%d", spreadsheetJobId), c)
	if err != nil {
		return output.ErrorAndExit(err), api.SpreadsheetJob{}
	}

	job := api.SpreadsheetJob{}
	if err = api.ParseJson(body, &job); err != nil {
		return err, api.SpreadsheetJob{}
	}

	if err = validateResponse(job); err != nil {
		return err, api.SpreadsheetJob{}
	}

	return nil, job
}

func validateResponse(job api.SpreadsheetJob) error {
	if job.Id == 0 {
		return output.ErrorStringAndExit("No spreadsheet job with that id found. Make sure that the id is correct and that you are using the intended API key")
	}

	return nil
}

func outputJob(w io.Writer, job api.SpreadsheetJob) {
	fmt.Fprintf(w, "Id: %d\n", job.Id)
	fmt.Fprintf(w, "Filename: %s\n", job.File.Filename)
	fmt.Fprintf(w, "Fields: %s\n", strings.Join(job.Fields, ","))
	fmt.Fprintf(w, "Rows: %s\n", humanize.Comma(int64(job.File.EstimatedRowsCount)))
	fmt.Fprintf(w, "State: %s\n", job.Status.State)
	fmt.Fprintf(w, "Progress: %.0f%%\n", job.Status.Progress)
	fmt.Fprintf(w, "Message: %s\n", job.Status.Message)

	timeLeft := job.Status.TimeLeftDescription
	if len(timeLeft) <= 0 {
		timeLeft = "-"
	}
	fmt.Fprintf(w, "Time left: %s\n", timeLeft)

	fmt.Fprintf(w, "Expires: %s\n", job.ExpiresAt.Format("Jan _2 15:04:05 2006"))
}

func outputDownloadHelpText(w io.Writer, job api.SpreadsheetJob) {
	if job.Status.State == "COMPLETED" && len(job.DownloadUrl) > 0 {
		fmt.Fprint(w, "\n\tTo download geocoded file\n")
		fmt.Fprintf(w, "\t$ %s download %d > geocoded_%s", os.Args[0], job.Id, job.File.Filename)
	}
}
