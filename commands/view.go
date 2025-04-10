package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

func NewViewCommand() *cobra.Command {
	var alias string

	cmd := &cobra.Command{
		Use:   "view",
		Short: "view a Zuora organization",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("viewing Zuora organization with alias: %s\n", alias)
		},
	}

	cmd.Flags().StringVar(&alias, "alias", "a", "Alias for the Zuora organization")

	return cmd
}
