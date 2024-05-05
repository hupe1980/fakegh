package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/google/uuid"
	"github.com/hupe1980/fakegh/pkg/git"
	"github.com/hupe1980/fakegh/pkg/github"
	"github.com/spf13/cobra"
)

type contributorOptions struct {
	url      string
	username string
	email    string
	date     string
	filename string
	message  string
}

func newContributorCmd(globalOpts *globalOptions) *cobra.Command {
	opts := &contributorOptions{}

	// Get today's date
	today := time.Now().Format("2006-01-02")

	cmd := &cobra.Command{
		Use:           "contributor",
		Short:         "Generate fake contributors to simulate community involvement in GitHub repositories",
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			date, err := time.Parse("2006-01-02", opts.date)
			if err != nil {
				return fmt.Errorf("error parsing date: %v", err)
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

			id := uuid.New()
			if err := repository.AddFile(opts.filename, []byte(id.String())); err != nil {
				return err
			}

			if err := repository.Commit(opts.message, &object.Signature{
				Name:  opts.username,
				Email: opts.email,
				When:  date,
			}); err != nil {
				return err
			}

			return repository.Push()
		},
	}

	cmd.Flags().StringVarP(&opts.date, "date", "", today, "date")
	cmd.Flags().StringVarP(&opts.url, "url", "", "", "github repo url")
	cmd.Flags().StringVarP(&opts.username, "username", "u", "", "github username")
	cmd.Flags().StringVarP(&opts.email, "email", "", "", "github email")
	cmd.Flags().StringVarP(&opts.filename, "filename", "", ".release", "filename")
	cmd.Flags().StringVarP(&opts.message, "message", "m", "Initial commit", "commit message")

	return cmd
}
