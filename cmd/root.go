package cmd

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

func Execute(version string) {
	PrintLogo()

	rootCmd := newRootCmd(version)
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

type globalOptions struct {
	token string
}

func newRootCmd(version string) *cobra.Command {
	globalOpts := &globalOptions{}

	cmd := &cobra.Command{
		Use:           "fakegh",
		Version:       version,
		Short:         "Fake GitHub Activity Generator",
		SilenceErrors: true,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if globalOpts.token != "" {
				return nil
			}

			_ = godotenv.Load()
			globalOpts.token = os.Getenv("GITHUB_TOKEN")
			return nil
		},
	}

	cmd.PersistentFlags().StringVarP(&globalOpts.token, "token", "", "", "Github token")

	cmd.AddCommand(
		newAchievementCmd(globalOpts),
		newActivityCmd(globalOpts),
		newContributorCmd(globalOpts),
		newGetEmailCmd(globalOpts),
	)

	return cmd
}

func PrintLogo() {
	fmt.Fprint(os.Stderr, `   _____        __                 .__     
 _/ ____\____  |  | __ ____   ____ |  |__  
 \   __\\__  \ |  |/ // __ \ / ___\|  |  \ 
  |  |   / __ \|    <\  ___// /_/  >   Y  \
  |__|  (____  /__|_ \\___  >___  /|___|  /
 	     \/     \/    \/_____/      \/ `, "\n\n")
}
