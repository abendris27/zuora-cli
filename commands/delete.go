package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

func NewDeleteCommand() *cobra.Command {
	var alias string

	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete a Zuora organization",
		Run: func(cmd *cobra.Command, args []string) {
			var jsURL = "https://one.zuora.com"
			err := openBrowser(jsURL)
			if err != nil {
				fmt.Printf("Failed to open browser: %v\n", err)
			}
			fmt.Printf("Deleting Zuora organization with alias: %s\n", alias)
		},
	}

	cmd.Flags().StringVar(&alias, "alias", "", "Alias for the Zuora organization")

	return cmd
}
