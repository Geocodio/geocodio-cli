package cli

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/geocodio/geocodio-cli/internal/api"
	"github.com/geocodio/geocodio-cli/internal/ui"
	"github.com/urfave/cli/v3"
)

func listsCmd() *cli.Command {
	return &cli.Command{
		Name:  "lists",
		Usage: "Manage spreadsheet geocoding jobs",
		Commands: []*cli.Command{
			listsUploadCmd(),
			listsListCmd(),
			listsStatusCmd(),
			listsDownloadCmd(),
			listsDeleteCmd(),
		},
	}
}

func listsUploadCmd() *cli.Command {
	return &cli.Command{
		Name:      "upload",
		Usage:     "Upload a spreadsheet for geocoding",
		ArgsUsage: "<file>",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "direction",
				Aliases:  []string{"d"},
				Usage:    "Geocoding direction: forward or reverse",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "format",
				Aliases:  []string{"f"},
				Usage:    "Column format template (e.g., {{A}} {{B}}, {{C}})",
				Required: true,
			},
			&cli.BoolFlag{
				Name:    "watch",
				Aliases: []string{"w"},
				Usage:   "Watch job progress until completion",
			},
			&cli.StringFlag{
				Name:  "callback",
				Usage: "Callback URL for completion notification",
			},
			&cli.StringFlag{
				Name:    "fields",
				Aliases: []string{"F"},
				Usage:   "Data append fields (comma-separated)",
			},
		},
		Action: listsUploadAction,
	}
}

func listsUploadAction(ctx context.Context, cmd *cli.Command) error {
	app, err := newApp(cmd)
	if err != nil {
		return err
	}

	if cmd.NArg() < 1 {
		return fmt.Errorf("file path required")
	}

	filePath := cmd.Args().First()
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("reading file: %w", err)
	}

	direction := cmd.String("direction")
	if direction != "forward" && direction != "reverse" {
		return fmt.Errorf("direction must be 'forward' or 'reverse'")
	}

	req := &api.ListUploadRequest{
		Filename:  filepath.Base(filePath),
		Data:      data,
		Direction: direction,
		Format:    cmd.String("format"),
		Callback:  cmd.String("callback"),
	}

	if fields := cmd.String("fields"); fields != "" {
		req.Fields = strings.Split(fields, ",")
	}

	resp, err := app.client.UploadList(ctx, req)
	if err != nil {
		return err
	}

	if cmd.Bool("watch") && (resp.Status == nil || (resp.Status.State != "COMPLETED" && resp.Status.State != "FAILED")) {
		fmt.Fprintf(app.stderr, "Uploaded list %d, watching progress...\n", resp.ID)
		return watchList(ctx, app, resp.ID)
	}

	return app.formatter.FormatList(resp)
}

func listsListCmd() *cli.Command {
	return &cli.Command{
		Name:   "list",
		Usage:  "List all uploaded spreadsheets",
		Action: listsListAction,
	}
}

func listsListAction(ctx context.Context, cmd *cli.Command) error {
	app, err := newApp(cmd)
	if err != nil {
		return err
	}

	resp, err := app.client.ListLists(ctx)
	if err != nil {
		return err
	}

	return app.formatter.FormatListList(resp)
}

func listsStatusCmd() *cli.Command {
	return &cli.Command{
		Name:      "status",
		Usage:     "Get status of a spreadsheet job",
		ArgsUsage: "<list-id>",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "watch",
				Aliases: []string{"w"},
				Usage:   "Watch job progress until completion",
			},
		},
		Action: listsStatusAction,
	}
}

func listsStatusAction(ctx context.Context, cmd *cli.Command) error {
	app, err := newApp(cmd)
	if err != nil {
		return err
	}

	if cmd.NArg() < 1 {
		return fmt.Errorf("list ID required")
	}

	id, err := strconv.Atoi(cmd.Args().First())
	if err != nil {
		return fmt.Errorf("invalid list ID: %w", err)
	}

	resp, err := app.client.GetList(ctx, id)
	if err != nil {
		return err
	}

	if err := app.formatter.FormatList(resp); err != nil {
		return err
	}

	if cmd.Bool("watch") {
		if resp.Status == nil || (resp.Status.State != "COMPLETED" && resp.Status.State != "FAILED") {
			return watchList(ctx, app, id)
		}
	}

	return nil
}

func listsDownloadCmd() *cli.Command {
	return &cli.Command{
		Name:      "download",
		Usage:     "Download results of a completed spreadsheet job",
		ArgsUsage: "<list-id>",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "output",
				Aliases: []string{"o"},
				Usage:   "Output file path",
			},
		},
		Action: listsDownloadAction,
	}
}

func listsDownloadAction(ctx context.Context, cmd *cli.Command) error {
	app, err := newApp(cmd)
	if err != nil {
		return err
	}

	if cmd.NArg() < 1 {
		return fmt.Errorf("list ID required")
	}

	id, err := strconv.Atoi(cmd.Args().First())
	if err != nil {
		return fmt.Errorf("invalid list ID: %w", err)
	}

	data, err := ui.WithSpinner(app.stderr, "Downloading results...", func() ([]byte, error) {
		return app.client.DownloadList(ctx, id)
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

func listsDeleteCmd() *cli.Command {
	return &cli.Command{
		Name:      "delete",
		Usage:     "Delete a spreadsheet job",
		ArgsUsage: "<list-id>",
		Action:    listsDeleteAction,
	}
}

func listsDeleteAction(ctx context.Context, cmd *cli.Command) error {
	app, err := newApp(cmd)
	if err != nil {
		return err
	}

	if cmd.NArg() < 1 {
		return fmt.Errorf("list ID required")
	}

	id, err := strconv.Atoi(cmd.Args().First())
	if err != nil {
		return fmt.Errorf("invalid list ID: %w", err)
	}

	if err := app.client.DeleteList(ctx, id); err != nil {
		return err
	}

	return app.formatter.FormatMessage(fmt.Sprintf("Deleted list %d", id))
}

func watchList(ctx context.Context, app *App, id int) error {
	display := ui.NewWatchDisplay(app.stderr)

	resp, err := app.client.PollList(ctx, id, func(list *api.ListResponse) {
		if list.Status == nil {
			return
		}
		display.Update(ui.WatchUpdate{
			Progress: list.Status.Progress,
			Status:   list.Status.State,
			TimeLeft: list.Status.TimeLeftDescription,
		})
	})

	display.Done()

	if err != nil {
		return err
	}

	return app.formatter.FormatList(resp)
}
