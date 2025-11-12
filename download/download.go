package download

import (
	"fmt"
	"github.com/geocodio/geocodio-cli/api"
	"github.com/geocodio/geocodio-cli/output"
	"github.com/urfave/cli/v2"
	"net/http"
)

type downloadResponse struct {
	Message           string   `json:"message"`
	Success           bool `json:"success"`
}

func RegisterCommand() *cli.Command {
	var command *cli.Command
	command = new(cli.Command)
	command.Name = "download"
	command.Usage = "Download output data for a specific geocoding job"
	command.Action = download
	command.ArgsUsage = "id"

	return command
}

func download(c *cli.Context) error {
	spreadsheetJobId, err := api.ValidateSpreadsheetJobId(c)
	if err != nil {
		return err
	}

	body, isJson, err := api.Request(http.MethodGet, fmt.Sprintf("lists/%d/download", spreadsheetJobId), c)
	if err != nil {
		return output.ErrorAndExit(err)
	}

	if isJson {
		job := downloadResponse{}
		if err = api.ParseJson(body, &job); err != nil {
			return err
		}

		return output.ErrorStringAndExit(job.Message)
	} else {
		fmt.Fprintf(c.App.Writer, "%s", string(body))
	}

	return nil
}
