package create

import (
	"encoding/json"
	"fmt"
	"github.com/dustin/go-humanize"
	"github.com/geocodio/geocodio-cli/api"
	"github.com/urfave/cli/v2"
	"os"
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
			Name:    "direction",
			Value: "forward",
			Usage:   "Direction can either be \"forward\" (address to coordinate) or \"reverse\" (coordinate to address)",
		},
	}

	return command
}

func geocode(c *cli.Context) error {
	hostname := c.String("hostname")
	apiKey := c.String("apikey")

	error := api.Validate(hostname, apiKey)
	if error != nil {
		cli.ShowAppHelp(c)
		return error
	}

	direction := c.String("direction")
	if direction != "forward" && direction != "reverse" {
		cli.ShowAppHelp(c)
		return cli.Exit("Please specify a valid direction. Valid values are \"forward\" or \"reverse\"", 1)
	}

	filename := c.Args().Get(0)
	format := c.Args().Get(1)

	if len(filename) <= 0 {
		cli.ShowAppHelp(c)
		return cli.Exit("Please specify a filename to geocode", 1)
	}

	if len(format) <= 0 {
		cli.ShowAppHelp(c)
		return cli.Exit("Please specify a geocoding format", 1)
	}

	file, err := os.Open(filename)

	if err != nil {
		cli.ShowAppHelp(c)
		return cli.Exit(fmt.Sprintf("Could not open %s", filename), 1)
	}
	defer file.Close()

	body := api.Upload(file, direction, format, hostname, apiKey)

	fmt.Println(string(body))

	job := api.SpreadsheetJob{}
	jsonErr := json.Unmarshal(body, &job)
	if jsonErr != nil {
		return cli.Exit("Could not parse JSON from the Geocodio API", 1)
	}

	if job.Message != "" {
		return cli.Exit(job.Message, 1)
	}

	outputJob(job)

	return nil
}


func outputJob(job api.SpreadsheetJob) {
	fmt.Println("Spreadsheet job created")
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