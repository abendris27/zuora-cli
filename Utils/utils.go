package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/olekukonko/tablewriter"
)

// Function to create a JSON file from any data
func createJSONFile(filePath string, data interface{}) error {
	// Open or create the file for writing
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("error creating file: %v", err)
	}
	defer file.Close()

	// Create a JSON encoder and set it to pretty-print
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")

	// Encode the data and write it to the file
	err = encoder.Encode(data)
	if err != nil {
		return fmt.Errorf("error encoding JSON: %v", err)
	}

	// Successfully created JSON file
	return nil
}

// Function to print selected columns of JSON files as tables
func printJSONFilesAsTable(directory string, columnsToDisplay []string) error {
	// Open the directory
	files, err := os.ReadDir(directory)
	if err != nil {
		return fmt.Errorf("failed to read directory: %v", err)
	}

	// Initialize table
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(append([]string{"File Name"}, columnsToDisplay...))

	// Loop over files in the directory
	for _, file := range files {
		// Only process JSON files
		if filepath.Ext(file.Name()) == ".json" {
			// Open the JSON file
			filePath := filepath.Join(directory, file.Name())
			fileContent, err := os.ReadFile(filePath)
			if err != nil {
				fmt.Printf("Error reading file %s: %v\n", file.Name(), err)
				continue
			}

			// Unmarshal JSON into a generic map
			var data map[string]interface{}
			if err := json.Unmarshal(fileContent, &data); err != nil {
				fmt.Printf("Error unmarshalling file %s: %v\n", file.Name(), err)
				continue
			}

			// Create a slice for row data with only selected columns
			row := []string{file.Name()}

			// Add selected columns to the row
			for _, col := range columnsToDisplay {
				if value, exists := data[col]; exists {
					row = append(row, fmt.Sprintf("%v", value))
				} else {
					row = append(row, "N/A") // If column doesn't exist in the data, add "N/A"
				}
			}

			// Add the row to the table
			table.Append(row)
		}
	}

	// Render the table
	table.Render()

	return nil
}
