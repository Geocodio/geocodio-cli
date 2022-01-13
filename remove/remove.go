package remove

import (
	"encoding/json"
	"fmt"
	"github.com/geocodio/geocodio-cli/api"
	"github.com/urfave/cli/v2"
	"net/http"
	"strconv"
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
	Success       bool     `json:"success,omitempty"`
	Message       string     `json:"message,omitempty"`
}

func remove(c *cli.Context) error {
	hostname := c.String("hostname")
	apiKey := c.String("apikey")

	error := api.Validate(hostname, apiKey)
	if error != nil {
		return error
	}

	spreadsheetJobId, err := strconv.Atoi(c.Args().First())
	if err != nil || spreadsheetJobId <= 0 {
		return cli.Exit("Invalid spreadsheet job id specified", 1)
	}

	body := api.Request(http.MethodDelete, fmt.Sprintf("lists/%d", spreadsheetJobId), hostname, apiKey)

	response := DeleteResponse{}
	jsonErr := json.Unmarshal(body, &response)
	if jsonErr != nil {
		return cli.Exit("Could not parse JSON from the Geocodio API", 1)
	}

	message := response.Message

	if response.Success {
		if len(message) <= 0 {
			message = "Spreadsheet job was successfully deleted"
		}
		return cli.Exit(message, 0)
	} else {
		if strings.Contains(message, "Resource not found") {
			message = "No spreadsheet job with that id found. Make sure that the id is correct and that you are using the expected API key"
		} else if len(message) <= 0 {
			message = "Spreadsheet job could not be deleted"
		}
		return cli.Exit(fmt.Sprintf("Error: %s", message), 1)
	}

	return nil
}
