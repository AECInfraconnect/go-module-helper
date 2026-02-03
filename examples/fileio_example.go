package main

import (
	"bytes"
	"fmt"
	"os"

	"github.com/AECInfraconnect/go-module-helper/helper"
)

// User represents a user for import/export example
type User struct {
	Name   string `excel:"Name" csv:"name"`
	Email  string `excel:"Email" csv:"email"`
	Age    int    `excel:"Age" csv:"age"`
	Active bool   `excel:"Active" csv:"active"`
}

func main() {
	fmt.Println("=== File Import/Export Examples ===\n")

	// ============================================
	// 1. CSV Export
	// ============================================
	fmt.Println("1. CSV Export Example")
	users := []User{
		{Name: "Somchai", Email: "somchai@example.com", Age: 30, Active: true},
		{Name: "Somying", Email: "somying@example.com", Age: 25, Active: true},
		{Name: "John", Email: "john@example.com", Age: 35, Active: false},
	}

	csvData, err := helper.ExportToCSV(users, helper.DefaultCSVExportOptions())
	if err != nil {
		fmt.Printf("   Error: %v\n", err)
		return
	}
	os.WriteFile("users.csv", csvData, 0644)
	fmt.Printf("   Exported to users.csv:\n%s\n", string(csvData))

	// ============================================
	// 2. CSV Import to Maps
	// ============================================
	fmt.Println("2. CSV Import to Maps")
	csvReader := bytes.NewReader(csvData)
	maps, err := helper.ImportCSVToMaps(csvReader, ',')
	if err != nil {
		fmt.Printf("   Error: %v\n", err)
		return
	}
	fmt.Printf("   Imported %d rows:\n", len(maps))
	for i, row := range maps {
		fmt.Printf("   Row %d: %v\n", i+1, row)
	}

	// ============================================
	// 3. CSV Import to Structs
	// ============================================
	fmt.Println("\n3. CSV Import to Structs")
	csvReader2 := bytes.NewReader(csvData)
	var importedUsers []User
	err = helper.ImportCSVToStructs(csvReader2, ',', &importedUsers)
	if err != nil {
		fmt.Printf("   Error: %v\n", err)
		return
	}
	fmt.Printf("   Imported %d users:\n", len(importedUsers))
	for _, u := range importedUsers {
		fmt.Printf("   - %s (%s), Age: %d, Active: %v\n", u.Name, u.Email, u.Age, u.Active)
	}

	// ============================================
	// 4. Import from existing CSV file
	// ============================================
	fmt.Println("\n4. Import from CSV file")
	if file, err := os.Open("users.csv"); err == nil {
		defer file.Close()
		headers, rows, err := helper.ImportFromCSV(file, ',')
		if err != nil {
			fmt.Printf("   Error: %v\n", err)
			return
		}
		fmt.Printf("   Headers: %v\n", headers)
		fmt.Printf("   Rows: %d\n", len(rows))
	}

	// ============================================
	// 5. Excel Export
	// ============================================
	fmt.Println("\n5. Excel Export Example")
	excelData, err := helper.ExportToExcel(users, helper.DefaultExcelExportOptions())
	if err != nil {
		fmt.Printf("   Error: %v\n", err)
		return
	}
	os.WriteFile("users.xlsx", excelData, 0644)
	fmt.Printf("   Exported to users.xlsx (%d bytes)\n", len(excelData))

	// ============================================
	// 6. Excel Import
	// ============================================
	fmt.Println("\n6. Excel Import Example")
	if file, err := os.Open("users.xlsx"); err == nil {
		defer file.Close()
		var excelUsers []User
		err = helper.ImportExcelToStructs(file, "Sheet1", &excelUsers)
		if err != nil {
			fmt.Printf("   Error: %v\n", err)
			return
		}
		fmt.Printf("   Imported %d users from Excel:\n", len(excelUsers))
		for _, u := range excelUsers {
			fmt.Printf("   - %s (%s), Age: %d, Active: %v\n", u.Name, u.Email, u.Age, u.Active)
		}
	}

	// ============================================
	// 7. PDF Export
	// ============================================
	fmt.Println("\n7. PDF Export Example")
	pdfOpts := helper.DefaultPDFExportOptions()
	pdfOpts.Title = "User Report"
	pdfData, err := helper.ExportToPDF(users, pdfOpts)
	if err != nil {
		fmt.Printf("   Error: %v\n", err)
		return
	}
	os.WriteFile("users.pdf", pdfData, 0644)
	fmt.Printf("   Exported to users.pdf (%d bytes)\n", len(pdfData))

	fmt.Println("\n=== Examples Complete ===")
}
