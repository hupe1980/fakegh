package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/hupe1980/fakegh/pkg/git"
	"github.com/hupe1980/fakegh/pkg/github"
	"github.com/spf13/cobra"
)

type activityOptions struct {
	startdate string
	enddate   string
	url       string
	username  string
	email     string
	filename  string
	message   string
}

func newActivityCmd(globalOpts *globalOptions) *cobra.Command {
	opts := &activityOptions{}

	// Get today's date
	today := time.Now().Format("2006-01-02")

	// Get the first day of the current year
	firstDayOfYear := time.Date(time.Now().Year(), 1, 1, 0, 0, 0, 0, time.UTC).Format("2006-01-02")

	cmd := &cobra.Command{
		Use:           "activity",
		Short:         "Simulate fake activity on GitHub repositories",
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			startDate, err := time.Parse("2006-01-02", opts.startdate)
			if err != nil {
				return fmt.Errorf("error parsing startdate: %v", err)
			}

			endDate, err := time.Parse("2006-01-02", opts.enddate)
			if err != nil {
				return fmt.Errorf("error parsing enddate: %v", err)
			}

			if opts.email == "" {
				client := github.New(globalOpts.token)

				email, err := client.GetEmailByUsername(context.Background(), opts.username)
				if err != nil {
					return err
				}

				opts.email = *email
			}

			repository, err := git.New(globalOpts.token, opts.url)
			if err != nil {
				return err
			}

			// Loop from start date to end date
			for currentDate := startDate; currentDate.Before(endDate) || currentDate.Equal(endDate); currentDate = currentDate.AddDate(0, 0, 1) {
				if err := repository.AddFile(opts.filename, []byte(currentDate.Format("2006-01-02"))); err != nil {
					return err
				}

				if err := repository.Commit(opts.message, &object.Signature{
					Name:  opts.username,
					Email: opts.email,
					When:  currentDate,
				}); err != nil {
					return err
				}
			}

			return repository.Push()
		},
	}

	cmd.Flags().StringVarP(&opts.startdate, "startdate", "", firstDayOfYear, "startdate")
	cmd.Flags().StringVarP(&opts.enddate, "enddate", "", today, "enddate")
	cmd.Flags().StringVarP(&opts.url, "url", "", "", "github repo url")
	cmd.Flags().StringVarP(&opts.username, "username", "u", "", "github username")
	cmd.Flags().StringVarP(&opts.email, "email", "", "", "github email")
	cmd.Flags().StringVarP(&opts.filename, "filename", "", ".release", "filename")
	cmd.Flags().StringVarP(&opts.message, "message", "m", "Initial commit", "commit message")

	return cmd
}
