package main

import (
	"fmt"
	"os"

	"github.com/abendris/zuora-cli/commands"
	"github.com/spf13/cobra"
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "zuo",
		Short: "Zuora CLI Tool",
	}

	var orgCmd = &cobra.Command{
		Use:   "org",
		Short: "Manage Zuora organizations",
	}

	// Add commands to org command
	orgCmd.AddCommand(commands.NewAuthorizeCommand())
	orgCmd.AddCommand(commands.NewDeleteCommand())
	orgCmd.AddCommand(commands.NewListCommand())
	orgCmd.AddCommand(commands.NewViewCommand())
	orgCmd.AddCommand(commands.NewLoginCommand())

	// Add org command to root command
	rootCmd.AddCommand(orgCmd)
	rootCmd.AddCommand(commands.NewNgrokCommand())

	// Execute the root command
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
