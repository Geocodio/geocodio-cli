package cli

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/geocodio/geocodio-cli/internal/api"
	"github.com/geocodio/geocodio-cli/internal/ui"
	"github.com/urfave/cli/v3"
)

// enqueuedHintDelay is how long a list must sit in ENQUEUED before watchList
// explains why. It's a var (not const) so tests can shorten it.
var enqueuedHintDelay = 30 * time.Second

// concurrencyGuideURL is the canonical source of truth for the actual
// concurrency limits per plan. Keep the numbers out of the CLI and this doc
// link current -- they live in one place instead of drifting across the
// CLI, website, agent guide, and API docs.
const concurrencyGuideURL = "https://www.geocod.io/guides/spreadsheet-concurrency"

// concurrencyExplanation describes why a list stays ENQUEUED: Geocodio caps
// how many spreadsheet jobs one billing owner can process at the same time.
// The API never rejects an upload over this limit -- it just delays it, so
// this is purely explanatory, not an error.
func concurrencyExplanation() string {
	return "Geocodio limits how many spreadsheet jobs one account can process at once. " +
		"The limit depends on your plan, and uploads started from the dashboard count toward it too. " +
		"No action needed -- it will start automatically once a slot frees up. " +
		"See " + concurrencyGuideURL + " for your plan's limit before running multiple lists."
}

// concurrencyExplanationBrief is a shorter version of concurrencyExplanation
// for a single (non-watching) status check.
func concurrencyExplanationBrief() string {
	return "Geocodio limits how many spreadsheet jobs one account can process at once. " +
		"The limit depends on your plan, and dashboard uploads count toward it too. " +
		"It will start automatically -- see " + concurrencyGuideURL + " before running multiple lists."
}

// shouldShowEnqueuedHint reports whether watchList should print the
// concurrency explanation for a list that has been ENQUEUED for elapsed.
func shouldShowEnqueuedHint(state string, elapsed time.Duration) bool {
	return state == "ENQUEUED" && elapsed >= enqueuedHintDelay
}

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
		Name:  "upload",
		Usage: "Upload a spreadsheet for geocoding",
		Description: "Uploads are queued and processed asynchronously. How many run at once for your " +
			"account is plan-dependent, and dashboard uploads share the same limit -- see " +
			concurrencyGuideURL + " before running multiple lists.",
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
		Name:  "status",
		Usage: "Get status of a spreadsheet job",
		Description: "A list can sit in ENQUEUED longer than expected while it waits for a free " +
			"processing slot -- how many run concurrently is plan-dependent. See " +
			concurrencyGuideURL + " before running multiple lists.",
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

	if resp.Status != nil && resp.Status.State == "ENQUEUED" {
		fmt.Fprintf(app.stderr, "List %d is queued. %s\n", id, concurrencyExplanationBrief())
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

	start := time.Now()
	hintShown := false

	resp, err := app.client.PollList(ctx, id, func(list *api.ListResponse) {
		if list.Status == nil {
			return
		}

		if !hintShown && shouldShowEnqueuedHint(list.Status.State, time.Since(start)) {
			fmt.Fprintf(app.stderr, "\nList %d is still queued. %s\n", id, concurrencyExplanation())
			hintShown = true
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
