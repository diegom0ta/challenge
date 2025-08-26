package cli

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"challenge/db"
	"challenge/models"
	"challenge/repository"
)

// CLI represents the command line interface
type CLI struct {
	args []string
}

// NewCLI creates a new CLI instance
func NewCLI() *CLI {
	return &CLI{
		args: os.Args[1:], // Skip program name
	}
}

// Run executes the CLI with the provided arguments
func (c *CLI) Run() error {
	if len(c.args) == 0 {
		return fmt.Errorf("usage: %s <csv-file>\nPlease provide a CSV file path as argument", os.Args[0])
	}

	if len(c.args) > 1 {
		return fmt.Errorf("too many arguments provided. Expected 1 CSV file, got %d arguments\nUsage: %s <csv-file>", len(c.args), os.Args[0])
	}

	csvFile := c.args[0]

	fmt.Printf("Processing CSV file: %s\n", csvFile)

	// Process the CSV file
	return c.processCSVFile(csvFile)
}

// processCSVFile handles the processing of CSV files
func (c *CLI) processCSVFile(filePath string) error {
	// Check if it's a help request
	if filePath == "help" || filePath == "--help" || filePath == "-h" {
		c.showHelp()
		return nil
	}

	if filePath == "version" || filePath == "--version" || filePath == "-v" {
		c.showVersion()
		return nil
	}

	// Validate file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return fmt.Errorf("CSV file does not exist: %s", filePath)
	}

	// Validate file extension
	ext := strings.ToLower(filepath.Ext(filePath))
	if ext != ".csv" {
		return fmt.Errorf("file must be a CSV file (*.csv), got: %s", ext)
	}

	// Open and read the CSV file
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open CSV file: %v", err)
	}
	defer file.Close()

	// Create CSV reader
	reader := csv.NewReader(file)

	// Read all records
	records, err := reader.ReadAll()
	if err != nil {
		return fmt.Errorf("failed to read CSV file: %v", err)
	}

	if len(records) == 0 {
		return fmt.Errorf("CSV file is empty")
	}

	fmt.Printf("Successfully loaded CSV file: %s\n", filePath)
	fmt.Printf("Found %d rows\n", len(records))

	// Show headers if available
	if len(records) > 0 {
		fmt.Printf("Headers: %v\n", records[0])
	}

	// Process the CSV data
	return c.processCSVData(records)
}

// processCSVData processes the CSV records
func (c *CLI) processCSVData(records [][]string) error {
	fmt.Printf("Processing %d records...\n", len(records))

	// Skip header row if it exists
	dataRows := records
	if len(records) > 1 {
		dataRows = records[1:]
		fmt.Printf("Processing %d data rows (excluding header)...\n", len(dataRows))
	}

	// Connect to database
	database, err := db.NewConnection()
	if err != nil {
		return fmt.Errorf("failed to connect to database: %v", err)
	}
	defer database.Close()

	// Create B3 repository
	b3Repo := repository.NewB3Repository(database.GetConnection())

	// Parse CSV rows to B3 models
	var b3Records []*models.B3
	for i, row := range dataRows {
		b3Record, err := repository.ParseCSVRowToB3(row)
		if err != nil {
			fmt.Printf("Warning: Failed to parse row %d: %v\n", i+1, err)
			fmt.Printf("Skipping row: %v\n", row)
			continue
		}
		b3Records = append(b3Records, b3Record)
	}

	if len(b3Records) == 0 {
		return fmt.Errorf("no valid records found in CSV file")
	}

	fmt.Printf("Successfully parsed %d valid records\n", len(b3Records))

	// Insert records into database
	err = b3Repo.InsertBatch(b3Records)
	if err != nil {
		return fmt.Errorf("failed to insert records: %v", err)
	}

	fmt.Printf("Successfully persisted %d records to database!\n", len(b3Records))
	return nil
}

// showHelp displays help information
func (c *CLI) showHelp() {
	fmt.Println("Challenge CLI - CSV File Processor")
	fmt.Println("Usage: challenge <csv-file>")
	fmt.Println()
	fmt.Println("Arguments:")
	fmt.Println("  <csv-file>           Path to a local CSV file to process")
	fmt.Println("  help, --help, -h     Show this help message")
	fmt.Println("  version, --version, -v Show version information")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  challenge data.csv")
	fmt.Println("  challenge /path/to/b3_data.csv")
}

// showVersion displays version information
func (c *CLI) showVersion() {
	fmt.Println("Challenge CLI - CSV Processor v1.0.0")
}
