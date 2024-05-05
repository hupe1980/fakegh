package cmd

import (
	"context"

	"github.com/hupe1980/fakegh/pkg/github"
	"github.com/spf13/cobra"
)

type achievementOptions struct {
	url string
}

func newAchievementCmd(globalOpts *globalOptions) *cobra.Command {
	opts := &achievementOptions{}

	cmd := &cobra.Command{
		Use:           "achievement",
		Short:         "Generate fake achievements to add to a GitHub profile",
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			client := github.New(globalOpts.token)

			ctx := context.Background()

			// https://github.com/drknzz/GitHub-Achievements?tab=readme-ov-file
			// Quickdraw (closed an issue / pull request within 5 minutes of opening)
			id, err := client.OpenIssue(ctx, opts.url)
			if err != nil {
				return err
			}

			return client.CloseIssue(ctx, opts.url, *id)
		},
	}

	cmd.Flags().StringVarP(&opts.url, "url", "", "", "github repo url")

	return cmd
}
