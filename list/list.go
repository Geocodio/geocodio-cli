package list

import (
	"fmt"
	"github.com/dustin/go-humanize"
	"github.com/geocodio/geocodio-cli/api"
	"github.com/geocodio/geocodio-cli/output"
	"github.com/olekukonko/tablewriter"
	"github.com/urfave/cli/v2"
	"io"
	"net/http"
)

type ListsResult struct {
	CurrentPage  int                  `json:"current_page"`
	FirstPageUrl string               `json:"first_page_url,omitempty"`
	From         int                  `json:"from"`
	NextPageUrl  string               `json:"next_page_url,omitempty"`
	Path         string               `json:"path"`
	PerPage      int                  `json:"per_page"`
	PrevPageUrl  string               `json:"prev_page_url,omitempty"`
	To           int                  `json:"to"`
	Jobs         []api.SpreadsheetJob `json:"data"`
	Error        string               `json:"error,omitempty"`
}

func RegisterCommand() *cli.Command {
	var command *cli.Command
	command = new(cli.Command)
	command.Name = "list"
	command.Usage = "List existing geocoding jobs"
	command.Action = list

	return command
}

func list(c *cli.Context) error {
	body, _, err := api.Request(http.MethodGet, "lists", c)
	if err != nil {
		return output.ErrorAndExit(err)
	}

	listResults := ListsResult{}
	if err = api.ParseJson(body, &listResults); err != nil {
		return err
	}

	if err := outputResults(c.App.Writer, listResults); err != nil {
		return err
	}

	return nil
}

func outputResults(w io.Writer, listResults ListsResult) error {
	var rows [][]string

	if len(listResults.Jobs) <= 0 {
		return output.WarningAndExit("No previous spreadsheet jobs found")
	}

	for _, job := range listResults.Jobs {
		row := []string{
			fmt.Sprintf("%d", job.Id),
			job.File.Filename,
			humanize.Comma(int64(job.File.EstimatedRowsCount)),
			job.Status.State,
			fmt.Sprintf("%.0f%%", job.Status.Progress),
			job.Status.Message,
			job.Status.TimeLeftDescription,
			job.ExpiresAt.Format("Jan _2 15:04:05 2006"),
		}
		rows = append(rows, row)
	}

	table := tablewriter.NewWriter(w)
	table.Header("Id", "Filename", "Rows", "State", "Progress", "Message", "Time left", "Expires")

	for _, v := range rows {
		table.Append(v)
	}
	table.Render()

	return nil
}
