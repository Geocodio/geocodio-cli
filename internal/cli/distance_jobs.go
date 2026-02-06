package cli

import (
	"context"
	"fmt"
	"os"

	"github.com/geocodio/geocodio-cli/internal/api"
	"github.com/geocodio/geocodio-cli/internal/ui"
	"github.com/urfave/cli/v3"
)

func distanceJobsCmd() *cli.Command {
	return &cli.Command{
		Name:  "distance-jobs",
		Usage: "Manage async distance calculation jobs",
		Commands: []*cli.Command{
			distanceJobsCreateCmd(),
			distanceJobsListCmd(),
			distanceJobsStatusCmd(),
			distanceJobsDownloadCmd(),
			distanceJobsDeleteCmd(),
		},
	}
}

func distanceJobsCreateCmd() *cli.Command {
	return &cli.Command{
		Name:  "create",
		Usage: "Create a new distance job",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "name",
				Aliases:  []string{"n"},
				Usage:    "Job name (required)",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "origins",
				Aliases:  []string{"o"},
				Usage:    "File containing origins (one per line)",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "destinations",
				Aliases:  []string{"d"},
				Usage:    "File containing destinations (one per line)",
				Required: true,
			},
			&cli.StringFlag{
				Name:    "mode",
				Aliases: []string{"m"},
				Usage:   "Routing mode: driving or straightline",
				Value:   "driving",
			},
			&cli.StringFlag{
				Name:    "units",
				Aliases: []string{"u"},
				Usage:   "Distance units: miles or km",
				Value:   "miles",
			},
			&cli.BoolFlag{
				Name:    "watch",
				Aliases: []string{"w"},
				Usage:   "Watch job progress until completion",
			},
		},
		Action: distanceJobsCreateAction,
	}
}

func distanceJobsCreateAction(ctx context.Context, cmd *cli.Command) error {
	app, err := newApp(cmd)
	if err != nil {
		return err
	}

	origins, err := readLines(cmd.String("origins"))
	if err != nil {
		return fmt.Errorf("reading origins file: %w", err)
	}

	destinations, err := readLines(cmd.String("destinations"))
	if err != nil {
		return fmt.Errorf("reading destinations file: %w", err)
	}

	req := &api.DistanceJobCreateRequest{
		Name:         cmd.String("name"),
		Origins:      origins,
		Destinations: destinations,
		Mode:         cmd.String("mode"),
		Units:        cmd.String("units"),
	}

	resp, err := app.client.CreateDistanceJob(ctx, req)
	if err != nil {
		return err
	}

	if err := app.formatter.FormatDistanceJob(resp); err != nil {
		return err
	}

	status := ""
	identifier := ""
	if resp.Data != nil {
		status = resp.Data.Status
		identifier = resp.Data.Identifier
	}
	if cmd.Bool("watch") && status != "COMPLETED" && status != "FAILED" && identifier != "" {
		return watchDistanceJob(ctx, app, identifier)
	}

	return nil
}

func distanceJobsListCmd() *cli.Command {
	return &cli.Command{
		Name:   "list",
		Usage:  "List all distance jobs",
		Action: distanceJobsListAction,
	}
}

func distanceJobsListAction(ctx context.Context, cmd *cli.Command) error {
	app, err := newApp(cmd)
	if err != nil {
		return err
	}

	resp, err := app.client.ListDistanceJobs(ctx)
	if err != nil {
		return err
	}

	return app.formatter.FormatDistanceJobList(resp)
}

func distanceJobsStatusCmd() *cli.Command {
	return &cli.Command{
		Name:      "status",
		Usage:     "Get status of a distance job",
		ArgsUsage: "<job-id>",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "watch",
				Aliases: []string{"w"},
				Usage:   "Watch job progress until completion",
			},
		},
		Action: distanceJobsStatusAction,
	}
}

func distanceJobsStatusAction(ctx context.Context, cmd *cli.Command) error {
	app, err := newApp(cmd)
	if err != nil {
		return err
	}

	if cmd.NArg() < 1 {
		return fmt.Errorf("job identifier required")
	}

	identifier := cmd.Args().First()

	resp, err := app.client.GetDistanceJob(ctx, identifier)
	if err != nil {
		return err
	}

	if err := app.formatter.FormatDistanceJob(resp); err != nil {
		return err
	}

	status := ""
	if resp.Data != nil {
		status = resp.Data.Status
	}
	if cmd.Bool("watch") && status != "COMPLETED" && status != "FAILED" {
		return watchDistanceJob(ctx, app, identifier)
	}

	return nil
}

func distanceJobsDownloadCmd() *cli.Command {
	return &cli.Command{
		Name:      "download",
		Usage:     "Download results of a completed distance job",
		ArgsUsage: "<job-id>",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "output",
				Aliases: []string{"o"},
				Usage:   "Output file path",
			},
		},
		Action: distanceJobsDownloadAction,
	}
}

func distanceJobsDownloadAction(ctx context.Context, cmd *cli.Command) error {
	app, err := newApp(cmd)
	if err != nil {
		return err
	}

	if cmd.NArg() < 1 {
		return fmt.Errorf("job identifier required")
	}

	identifier := cmd.Args().First()

	data, err := ui.WithSpinner(app.stderr, "Downloading results...", func() ([]byte, error) {
		return app.client.DownloadDistanceJob(ctx, identifier)
	})
	if err != nil {
		return err
	}

	outputPath := cmd.String("output")
	if outputPath != "" {
		if err := os.WriteFile(outputPath, data, 0600); err != nil {
			return fmt.Errorf("writing output file: %w", err)
		}
		return app.formatter.FormatMessage(fmt.Sprintf("Downloaded to %s", outputPath))
	}

	_, err = app.stdout.Write(data)
	return err
}

func distanceJobsDeleteCmd() *cli.Command {
	return &cli.Command{
		Name:      "delete",
		Usage:     "Delete a distance job",
		ArgsUsage: "<job-id>",
		Action:    distanceJobsDeleteAction,
	}
}

func distanceJobsDeleteAction(ctx context.Context, cmd *cli.Command) error {
	app, err := newApp(cmd)
	if err != nil {
		return err
	}

	if cmd.NArg() < 1 {
		return fmt.Errorf("job identifier required")
	}

	identifier := cmd.Args().First()

	if err := app.client.DeleteDistanceJob(ctx, identifier); err != nil {
		return err
	}

	return app.formatter.FormatMessage(fmt.Sprintf("Deleted job %s", identifier))
}

func watchDistanceJob(ctx context.Context, app *App, identifier string) error {
	display := ui.NewWatchDisplay(app.stderr)

	fmt.Fprintln(app.stderr, "\nWatching job progress...")

	resp, err := app.client.PollDistanceJob(ctx, identifier, func(job *api.DistanceJobResponse) {
		if job.Data == nil {
			return
		}
		display.Update(ui.WatchUpdate{
			Progress: float64(job.Data.Progress),
			Status:   job.Data.Status,
		})
	})

	display.Done()

	if err != nil {
		return err
	}

	return app.formatter.FormatDistanceJob(resp)
}
