package list

import (
	"encoding/json"
	"fmt"
	"github.com/dustin/go-humanize"
	"github.com/geocodio/geocodio-cli/api"
	"github.com/olekukonko/tablewriter"
	"github.com/urfave/cli/v2"
	"net/http"
	"os"
)

type ListsResult struct {
	CurrentPage  int              `json:"current_page"`
	FirstPageUrl string           `json:"first_page_url,omitempty"`
	From         int              `json:"from"`
	NextPageUrl  string           `json:"next_page_url,omitempty"`
	Path         string           `json:"path"`
	PerPage      int              `json:"per_page"`
	PrevPageUrl  string           `json:"prev_page_url,omitempty"`
	To           int              `json:"to"`
	Jobs         []api.SpreadsheetJob `json:"data"`
	Error        string           `json:"error,omitempty"`
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
	hostname := c.String("hostname")
	apiKey := c.String("apikey")

	if len(hostname) <= 0 {
		return cli.Exit("Please specify a valid hostname", 1)
	}

	if len(apiKey) <= 0 {
		return cli.Exit("Please specify a valid apikey", 1)
	}

	body := api.Request(http.MethodGet, "lists", hostname, apiKey)

	listResults := ListsResult{}
	jsonErr := json.Unmarshal(body, &listResults)
	if jsonErr != nil {
		return cli.Exit("Could not parse JSON from the Geocodio API", 1)
	}

	if listResults.Error != "" {
		return cli.Exit(fmt.Sprintf("Error: %s", listResults.Error), 1)
	}

	var data [][]string

	if len(listResults.Jobs) <= 0 {
		return cli.Exit("No previous spreadsheet jobs found", 1)
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
			job.ExpiresAt.Format("Mon Jan _2 15:04:05 2006"),
		}
		data = append(data, row)
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Id", "Filename", "Rows", "State", "Progress", "Message", "Time left", "Expires"})

	for _, v := range data {
		table.Append(v)
	}
	table.Render()

	return nil
}
