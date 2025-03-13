package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// Enum for the flag values
type OrgType string

const (
	EUPROD OrgType = "EUPROD"
	USPROD OrgType = "USPROD"
	EUSBX  OrgType = "EUSBX"
)

var allowedOrgTypes = []OrgType{EUPROD, USPROD, EUSBX}

func (o *OrgType) Set(val string) error {
	// Check if the value is a valid OrgType
	switch OrgType(val) {
	case EUPROD, USPROD, EUSBX:
		*o = OrgType(val)
		return nil
	default:
		return fmt.Errorf("invalid OrgType: %s. Allowed values are: %v", val, allowedOrgTypes)
	}
}

func (o *OrgType) Type() string {
	return "OrgType"
}

func (o *OrgType) String() string {
	return string(*o)
}

// Global variables for flags
var clientId, clientSecret, tenantId, alias string
var target OrgType

func main() {

	var target OrgType

	var rootCmd = &cobra.Command{
		Use:   "zuo",
		Short: "Zuora CLI Tool",
	}

	var orgCmd = &cobra.Command{
		Use:   "org",
		Short: "Manage Zuora organizations",
	}

	// Define the "authorize" subcommand
	var authorizeCmd = &cobra.Command{
		Use:   "authorize",
		Short: "Authorize a Zuora organization",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("Authorizing Zuora organization with target: %s, clientId: %s, clientSecret: %s, tenantId: %s, alias: %s\n",
				target, clientId, clientSecret, tenantId, alias)
			// Add the actual authorization logic here
		},
	}

	// Define flags for the "authorize" command
	authorizeCmd.Flags().VarP(&target, "target", "", "Target environment (e.g., PROD)")
	authorizeCmd.Flags().StringVar(&clientId, "clientId", "", "Zuora client ID")
	authorizeCmd.Flags().StringVar(&clientSecret, "clientSecret", "", "Zuora client secret")
	authorizeCmd.Flags().StringVar(&tenantId, "tenantId", "", "Zuora tenant ID (optional)")
	authorizeCmd.Flags().StringVar(&alias, "alias", "", "Alias for the Zuora organization")

	// Define the "delete" subcommand
	var deleteCmd = &cobra.Command{
		Use:   "delete",
		Short: "Delete a Zuora organization",
		Run: func(cmd *cobra.Command, args []string) {
			// Handle deletion
			fmt.Printf("Deleting Zuora organization with alias: %s\n", alias)
			// Add the actual deletion logic here
		},
	}

	// Define flags for the "list" command
	deleteCmd.Flags().StringVar(&alias, "alias", "", "Alias for the Zuora organization")

	var listCmd = &cobra.Command{
		Use:   "list",
		Short: "list the stored Zuora organizations",
		Run: func(cmd *cobra.Command, args []string) {
			// Handle deletion
			fmt.Printf("listing the stored Zuora organizations\n")
			// Add the actual deletion logic here
		},
	}

	// Define the "delete" subcommand
	var viewOrgCmd = &cobra.Command{
		Use:   "view",
		Short: "view a Zuora organization",
		Run: func(cmd *cobra.Command, args []string) {
			// Handle deletion
			fmt.Printf("viewing Zuora organization with alias: %s\n", alias)
			// Add the actual deletion logic here
		},
	}

	// Define flags for the "list" command
	viewOrgCmd.Flags().StringVar(&alias, "alias", "a", "Alias for the Zuora organization")

	// Add the commands to the org command
	orgCmd.AddCommand(authorizeCmd)
	orgCmd.AddCommand(deleteCmd)
	orgCmd.AddCommand(listCmd)
	orgCmd.AddCommand(viewOrgCmd)

	// Add org command to the root command
	rootCmd.AddCommand(orgCmd)

	// Execute the root command
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
