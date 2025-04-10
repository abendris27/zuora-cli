package commands

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

func printJSONFilesAsTable(directory string, columnsToDisplay []string) error {
	files, err := os.ReadDir(directory)
	if err != nil {
		return fmt.Errorf("failed to read directory: %v", err)
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(append([]string{"File Name"}, columnsToDisplay...))

	for _, file := range files {
		if filepath.Ext(file.Name()) == ".json" {
			filePath := filepath.Join(directory, file.Name())
			fileContent, err := os.ReadFile(filePath)
			if err != nil {
				fmt.Printf("Error reading file %s: %v\n", file.Name(), err)
				continue
			}

			var data map[string]interface{}
			if err := json.Unmarshal(fileContent, &data); err != nil {
				fmt.Printf("Error unmarshalling file %s: %v\n", file.Name(), err)
				continue
			}

			row := []string{file.Name()}
			for _, col := range columnsToDisplay {
				if value, exists := data[col]; exists {
					row = append(row, fmt.Sprintf("%v", value))
				} else {
					row = append(row, "N/A")
				}
			}

			table.Append(row)
		}
	}

	table.Render()
	return nil
}

func NewListCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "list the stored Zuora organizations",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("listing the stored Zuora organizations\n")
			columnsToDisplay := []string{"alias", "clientId", "tenantId", "firstGenerated"}
			err := printJSONFilesAsTable(path, columnsToDisplay)
			if err != nil {
				fmt.Println("Error listing JSON files:", err)
				return
			}
		},
	}

	return cmd
}
