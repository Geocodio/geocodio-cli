package remove

import (
	"fmt"
	"github.com/geocodio/geocodio-cli/api"
	"github.com/geocodio/geocodio-cli/output"
	"github.com/urfave/cli/v2"
	"io"
	"net/http"
	"strings"
)

func RegisterCommand() *cli.Command {
	var command *cli.Command
	command = new(cli.Command)
	command.Name = "remove"
	command.Aliases = []string{"rm"}
	command.Usage = "Delete an existing geocoding job"
	command.Action = remove
	command.ArgsUsage = "id"

	return command
}

type DeleteResponse struct {
	Success bool   `json:"success,omitempty"`
	Message string `json:"message,omitempty"`
}

func remove(c *cli.Context) error {
	spreadsheetJobId, err := api.ValidateSpreadsheetJobId(c)
	if err != nil {
		return err
	}

	body, _, err := api.Request(http.MethodDelete, fmt.Sprintf("lists/%d", spreadsheetJobId), c)
	if err != nil {
		return output.ErrorAndExit(err)
	}

	response := DeleteResponse{}
	if err = api.ParseJson(body, &response); err != nil {
		return err
	}

	if err := outputOutcome(c.App.Writer, response); err != nil {
		return nil
	}

	return nil
}

func outputOutcome(w io.Writer, response DeleteResponse) error {
	message := response.Message

	if response.Success {
		if len(message) <= 0 {
			message = "Spreadsheet job was successfully deleted"
		}

		output.Success(w, message)

		return nil
	} else {
		if strings.Contains(message, "Resource not found") {
			message = "No spreadsheet job with that id found. Make sure that the id is correct and that you are using the expected API key"
		} else if len(message) <= 0 {
			message = "Spreadsheet job could not be deleted"
		}

		return output.ErrorStringAndExit(message)
	}
}
