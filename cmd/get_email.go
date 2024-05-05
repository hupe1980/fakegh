package cmd

import (
	"context"
	"fmt"

	"github.com/hupe1980/fakegh/pkg/github"
	"github.com/spf13/cobra"
)

type getEmailOptions struct {
	username string
}

func newGetEmailCmd(globalOpts *globalOptions) *cobra.Command {
	opts := &getEmailOptions{}

	cmd := &cobra.Command{
		Use:           "get-email",
		Short:         "Retrieve the email address associated with a GitHub user",
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			client := github.New(globalOpts.token)
			email, err := client.GetEmailByUsername(context.Background(), opts.username)
			if err != nil {
				return err
			}

			fmt.Println(*email)

			return nil
		},
	}

	cmd.Flags().StringVarP(&opts.username, "username", "u", "", "github username")

	return cmd
}
