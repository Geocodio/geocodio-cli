package api

import "github.com/urfave/cli/v2"

func Validate(hostname string, apiKey string) error {
	if len(hostname) <= 0 {
		return cli.Exit("Please specify a valid hostname", 1)
	}

	if len(apiKey) <= 0 {
		return cli.Exit("Please specify a valid apikey", 1)
	}

	return nil
}
