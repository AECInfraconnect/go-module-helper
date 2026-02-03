package helper

import (
	"bytes"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"reflect"
	"strconv"
	"strings"

	"github.com/jung-kurt/gofpdf"
	"github.com/xuri/excelize/v2"
)

// ============================================================================
// EXCEL FUNCTIONS
// ============================================================================

// ExcelExportOptions configures Excel export behavior
type ExcelExportOptions struct {
	SheetName   string            // Sheet name (default: "Sheet1")
	Headers     []string          // Custom headers (if empty, uses struct field names or tags)
	HeaderStyle *excelize.Style   // Custom header style
	DataStyle   *excelize.Style   // Custom data style
	ColumnWidth map[string]float64 // Column widths (e.g., {"A": 20, "B": 30})
}

// DefaultExcelExportOptions returns default export options
func DefaultExcelExportOptions() ExcelExportOptions {
	return ExcelExportOptions{
		SheetName: "Sheet1",
	}
}

// ExportToExcel exports data to Excel bytes
// Supports []map[string]any or slice of structs
//
// Example with struct:
//
//	type User struct {
//	    Name  string `excel:"ชื่อ"`
//	    Email string `excel:"อีเมล"`
//	    Age   int    `excel:"อายุ"`
//	}
//	users := []User{{Name: "John", Email: "john@example.com", Age: 30}}
//	data, err := helper.ExportToExcel(users, helper.DefaultExcelExportOptions())
//
// Example with map:
//
//	data := []map[string]any{{"name": "John", "age": 30}}
//	bytes, err := helper.ExportToExcel(data, helper.DefaultExcelExportOptions())
func ExportToExcel(data any, opts ExcelExportOptions) ([]byte, error) {
	if opts.SheetName == "" {
		opts.SheetName = "Sheet1"
	}

	f := excelize.NewFile()
	defer f.Close()

	// Create sheet
	index, err := f.NewSheet(opts.SheetName)
	if err != nil {
		return nil, fmt.Errorf("failed to create sheet: %w", err)
	}
	f.SetActiveSheet(index)

	// Delete default sheet if different
	if opts.SheetName != "Sheet1" {
		f.DeleteSheet("Sheet1")
	}

	// Convert data to rows
	headers, rows, err := convertToRows(data, opts.Headers)
	if err != nil {
		return nil, err
	}

	// Write headers
	for col, header := range headers {
		cell, _ := excelize.CoordinatesToCellName(col+1, 1)
		f.SetCellValue(opts.SheetName, cell, header)
	}

	// Apply header style if provided
	if opts.HeaderStyle != nil {
		styleID, _ := f.NewStyle(opts.HeaderStyle)
		endCell, _ := excelize.CoordinatesToCellName(len(headers), 1)
		f.SetCellStyle(opts.SheetName, "A1", endCell, styleID)
	}

	// Write data rows
	for rowIdx, row := range rows {
		for colIdx, value := range row {
			cell, _ := excelize.CoordinatesToCellName(colIdx+1, rowIdx+2)
			f.SetCellValue(opts.SheetName, cell, value)
		}
	}

	// Apply column widths
	for col, width := range opts.ColumnWidth {
		f.SetColWidth(opts.SheetName, col, col, width)
	}

	// Write to buffer
	buf, err := f.WriteToBuffer()
	if err != nil {
		return nil, fmt.Errorf("failed to write excel: %w", err)
	}

	return buf.Bytes(), nil
}

// ImportFromExcel imports data from Excel reader
// Returns headers and rows as [][]string
//
// Example:
//
//	file, _ := os.Open("data.xlsx")
//	headers, rows, err := helper.ImportFromExcel(file, "Sheet1")
//	for _, row := range rows {
//	    fmt.Println(row)
//	}
func ImportFromExcel(reader io.Reader, sheetName string) (headers []string, rows [][]string, err error) {
	f, err := excelize.OpenReader(reader)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to open excel: %w", err)
	}
	defer f.Close()

	if sheetName == "" {
		sheetName = f.GetSheetName(0)
	}

	allRows, err := f.GetRows(sheetName)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read sheet: %w", err)
	}

	if len(allRows) == 0 {
		return nil, nil, nil
	}

	headers = allRows[0]
	if len(allRows) > 1 {
		rows = allRows[1:]
	}

	return headers, rows, nil
}

// ImportExcelToMaps imports Excel to []map[string]string
// Uses first row as keys
//
// Example:
//
//	file, _ := os.Open("data.xlsx")
//	data, err := helper.ImportExcelToMaps(file, "Sheet1")
//	for _, row := range data {
//	    fmt.Println(row["name"], row["email"])
//	}
func ImportExcelToMaps(reader io.Reader, sheetName string) ([]map[string]string, error) {
	headers, rows, err := ImportFromExcel(reader, sheetName)
	if err != nil {
		return nil, err
	}

	var result []map[string]string
	for _, row := range rows {
		item := make(map[string]string)
		for i, header := range headers {
			if i < len(row) {
				item[header] = row[i]
			} else {
				item[header] = ""
			}
		}
		result = append(result, item)
	}

	return result, nil
}

// ImportExcelToStructs imports Excel to slice of structs
// Uses struct tags `excel:"column_name"` to map columns
//
// Example:
//
//	type User struct {
//	    Name  string `excel:"ชื่อ"`
//	    Email string `excel:"อีเมล"`
//	}
//	file, _ := os.Open("data.xlsx")
//	var users []User
//	err := helper.ImportExcelToStructs(file, "Sheet1", &users)
func ImportExcelToStructs(reader io.Reader, sheetName string, dest any) error {
	headers, rows, err := ImportFromExcel(reader, sheetName)
	if err != nil {
		return err
	}

	return mapRowsToStructs(headers, rows, dest)
}

// ============================================================================
// CSV FUNCTIONS
// ============================================================================

// CSVExportOptions configures CSV export behavior
type CSVExportOptions struct {
	Delimiter rune     // Field delimiter (default: ',')
	Headers   []string // Custom headers
	UseCRLF   bool     // Use \r\n as line ending
}

// DefaultCSVExportOptions returns default CSV options
func DefaultCSVExportOptions() CSVExportOptions {
	return CSVExportOptions{
		Delimiter: ',',
		UseCRLF:   false,
	}
}

// ExportToCSV exports data to CSV bytes
// Supports []map[string]any or slice of structs
//
// Example:
//
//	type User struct {
//	    Name  string `csv:"name"`
//	    Email string `csv:"email"`
//	}
//	users := []User{{Name: "John", Email: "john@example.com"}}
//	data, err := helper.ExportToCSV(users, helper.DefaultCSVExportOptions())
func ExportToCSV(data any, opts CSVExportOptions) ([]byte, error) {
	if opts.Delimiter == 0 {
		opts.Delimiter = ','
	}

	headers, rows, err := convertToRows(data, opts.Headers)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)
	writer.Comma = opts.Delimiter
	writer.UseCRLF = opts.UseCRLF

	// Write headers
	if err := writer.Write(headers); err != nil {
		return nil, fmt.Errorf("failed to write headers: %w", err)
	}

	// Write rows
	for _, row := range rows {
		strRow := make([]string, len(row))
		for i, v := range row {
			strRow[i] = fmt.Sprintf("%v", v)
		}
		if err := writer.Write(strRow); err != nil {
			return nil, fmt.Errorf("failed to write row: %w", err)
		}
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		return nil, fmt.Errorf("csv write error: %w", err)
	}

	return buf.Bytes(), nil
}

// ImportFromCSV imports data from CSV reader
// Returns headers and rows
//
// Example:
//
//	file, _ := os.Open("data.csv")
//	headers, rows, err := helper.ImportFromCSV(file, ',')
func ImportFromCSV(reader io.Reader, delimiter rune) (headers []string, rows [][]string, err error) {
	if delimiter == 0 {
		delimiter = ','
	}

	csvReader := csv.NewReader(reader)
	csvReader.Comma = delimiter
	csvReader.LazyQuotes = true
	csvReader.TrimLeadingSpace = true

	allRows, err := csvReader.ReadAll()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read csv: %w", err)
	}

	if len(allRows) == 0 {
		return nil, nil, nil
	}

	headers = allRows[0]
	if len(allRows) > 1 {
		rows = allRows[1:]
	}

	return headers, rows, nil
}

// ImportCSVToMaps imports CSV to []map[string]string
//
// Example:
//
//	file, _ := os.Open("data.csv")
//	data, err := helper.ImportCSVToMaps(file, ',')
//	for _, row := range data {
//	    fmt.Println(row["name"])
//	}
func ImportCSVToMaps(reader io.Reader, delimiter rune) ([]map[string]string, error) {
	headers, rows, err := ImportFromCSV(reader, delimiter)
	if err != nil {
		return nil, err
	}

	var result []map[string]string
	for _, row := range rows {
		item := make(map[string]string)
		for i, header := range headers {
			if i < len(row) {
				item[header] = row[i]
			} else {
				item[header] = ""
			}
		}
		result = append(result, item)
	}

	return result, nil
}

// ImportCSVToStructs imports CSV to slice of structs
// Uses struct tags `csv:"column_name"` to map columns
//
// Example:
//
//	type User struct {
//	    Name  string `csv:"name"`
//	    Email string `csv:"email"`
//	}
//	file, _ := os.Open("data.csv")
//	var users []User
//	err := helper.ImportCSVToStructs(file, ',', &users)
func ImportCSVToStructs(reader io.Reader, delimiter rune, dest any) error {
	headers, rows, err := ImportFromCSV(reader, delimiter)
	if err != nil {
		return err
	}

	return mapRowsToStructs(headers, rows, dest)
}

// ============================================================================
// PDF FUNCTIONS
// ============================================================================

// PDFExportOptions configures PDF export behavior
type PDFExportOptions struct {
	PageSize    string  // A4, Letter, Legal (default: A4)
	Orientation string  // P (portrait) or L (landscape) (default: P)
	FontSize    float64 // Font size (default: 12)
	FontFamily  string  // Font family (default: Arial)
	Title       string  // Document title
	MarginLeft  float64 // Left margin (default: 10)
	MarginTop   float64 // Top margin (default: 10)
	MarginRight float64 // Right margin (default: 10)
}

// DefaultPDFExportOptions returns default PDF options
func DefaultPDFExportOptions() PDFExportOptions {
	return PDFExportOptions{
		PageSize:    "A4",
		Orientation: "P",
		FontSize:    12,
		FontFamily:  "Arial",
		MarginLeft:  10,
		MarginTop:   10,
		MarginRight: 10,
	}
}

// PDFTable represents a table to be rendered in PDF
type PDFTable struct {
	Headers     []string   // Column headers
	Rows        [][]string // Data rows
	ColumnWidth []float64  // Column widths (optional)
}

// PDFTableFromData creates PDFTable from data (struct slice or []map[string]any)
//
// Example:
//
//	users := []User{{Name: "John", Age: 30}}
//	table, err := helper.PDFTableFromData(users, nil)
func PDFTableFromData(data any, customHeaders []string) (*PDFTable, error) {
	headers, rows, err := convertToRows(data, customHeaders)
	if err != nil {
		return nil, err
	}

	strRows := make([][]string, len(rows))
	for i, row := range rows {
		strRows[i] = make([]string, len(row))
		for j, v := range row {
			strRows[i][j] = fmt.Sprintf("%v", v)
		}
	}

	return &PDFTable{
		Headers: headers,
		Rows:    strRows,
	}, nil
}

// ExportToPDF exports data to PDF bytes with a table format
// Supports []map[string]any or slice of structs
//
// Example:
//
//	type User struct {
//	    Name  string `excel:"Name"`
//	    Email string `excel:"Email"`
//	}
//	users := []User{{Name: "John", Email: "john@example.com"}}
//	data, err := helper.ExportToPDF(users, helper.DefaultPDFExportOptions())
func ExportToPDF(data any, opts PDFExportOptions) ([]byte, error) {
	if opts.PageSize == "" {
		opts.PageSize = "A4"
	}
	if opts.Orientation == "" {
		opts.Orientation = "P"
	}
	if opts.FontSize == 0 {
		opts.FontSize = 10
	}
	if opts.MarginLeft == 0 {
		opts.MarginLeft = 10
	}
	if opts.MarginTop == 0 {
		opts.MarginTop = 10
	}
	if opts.MarginRight == 0 {
		opts.MarginRight = 10
	}

	// Create PDF
	pdf := gofpdf.New(opts.Orientation, "mm", opts.PageSize, "")
	pdf.SetMargins(opts.MarginLeft, opts.MarginTop, opts.MarginRight)
	pdf.AddPage()

	// Set font (use built-in font for basic ASCII)
	pdf.SetFont("Arial", "", opts.FontSize)

	// Add title if provided
	if opts.Title != "" {
		pdf.SetFont("Arial", "B", opts.FontSize+4)
		pdf.CellFormat(0, 10, opts.Title, "", 1, "C", false, 0, "")
		pdf.Ln(5)
		pdf.SetFont("Arial", "", opts.FontSize)
	}

	// Convert data to table
	table, err := PDFTableFromData(data, nil)
	if err != nil {
		return nil, err
	}

	// Calculate column widths
	pageWidth, _ := pdf.GetPageSize()
	marginLeft, _, marginRight, _ := pdf.GetMargins()
	availableWidth := pageWidth - marginLeft - marginRight
	colWidth := availableWidth / float64(len(table.Headers))

	// Draw header
	pdf.SetFont("Arial", "B", opts.FontSize)
	pdf.SetFillColor(220, 220, 220)
	for _, header := range table.Headers {
		pdf.CellFormat(colWidth, 8, header, "1", 0, "C", true, 0, "")
	}
	pdf.Ln(-1)

	// Draw rows
	pdf.SetFont("Arial", "", opts.FontSize)
	pdf.SetFillColor(255, 255, 255)
	for _, row := range table.Rows {
		for _, cell := range row {
			pdf.CellFormat(colWidth, 7, cell, "1", 0, "L", false, 0, "")
		}
		pdf.Ln(-1)
	}

	// Output to bytes
	var buf bytes.Buffer
	if err := pdf.Output(&buf); err != nil {
		return nil, fmt.Errorf("failed to generate PDF: %w", err)
	}

	return buf.Bytes(), nil
}

// ExportTableToPDF exports PDFTable to PDF bytes
//
// Example:
//
//	table := &helper.PDFTable{
//	    Headers: []string{"Name", "Age"},
//	    Rows:    [][]string{{"John", "30"}, {"Jane", "25"}},
//	}
//	data, err := helper.ExportTableToPDF(table, helper.DefaultPDFExportOptions())
func ExportTableToPDF(table *PDFTable, opts PDFExportOptions) ([]byte, error) {
	if opts.PageSize == "" {
		opts.PageSize = "A4"
	}
	if opts.Orientation == "" {
		opts.Orientation = "P"
	}
	if opts.FontSize == 0 {
		opts.FontSize = 10
	}
	if opts.MarginLeft == 0 {
		opts.MarginLeft = 10
	}
	if opts.MarginTop == 0 {
		opts.MarginTop = 10
	}
	if opts.MarginRight == 0 {
		opts.MarginRight = 10
	}

	// Create PDF
	pdf := gofpdf.New(opts.Orientation, "mm", opts.PageSize, "")
	pdf.SetMargins(opts.MarginLeft, opts.MarginTop, opts.MarginRight)
	pdf.AddPage()
	pdf.SetFont("Arial", "", opts.FontSize)

	// Add title if provided
	if opts.Title != "" {
		pdf.SetFont("Arial", "B", opts.FontSize+4)
		pdf.CellFormat(0, 10, opts.Title, "", 1, "C", false, 0, "")
		pdf.Ln(5)
		pdf.SetFont("Arial", "", opts.FontSize)
	}

	// Calculate column widths
	pageWidth, _ := pdf.GetPageSize()
	marginLeft, _, marginRight, _ := pdf.GetMargins()
	availableWidth := pageWidth - marginLeft - marginRight

	var colWidths []float64
	if len(table.ColumnWidth) > 0 {
		colWidths = table.ColumnWidth
	} else {
		colWidth := availableWidth / float64(len(table.Headers))
		colWidths = make([]float64, len(table.Headers))
		for i := range colWidths {
			colWidths[i] = colWidth
		}
	}

	// Draw header
	pdf.SetFont("Arial", "B", opts.FontSize)
	pdf.SetFillColor(220, 220, 220)
	for i, header := range table.Headers {
		pdf.CellFormat(colWidths[i], 8, header, "1", 0, "C", true, 0, "")
	}
	pdf.Ln(-1)

	// Draw rows
	pdf.SetFont("Arial", "", opts.FontSize)
	for _, row := range table.Rows {
		for i, cell := range row {
			if i < len(colWidths) {
				pdf.CellFormat(colWidths[i], 7, cell, "1", 0, "L", false, 0, "")
			}
		}
		pdf.Ln(-1)
	}

	// Output to bytes
	var buf bytes.Buffer
	if err := pdf.Output(&buf); err != nil {
		return nil, fmt.Errorf("failed to generate PDF: %w", err)
	}

	return buf.Bytes(), nil
}

// ============================================================================
// HELPER FUNCTIONS
// ============================================================================

// convertToRows converts data to headers and rows
// Supports []map[string]any or slice of structs
func convertToRows(data any, customHeaders []string) ([]string, [][]any, error) {
	v := reflect.ValueOf(data)
	if v.Kind() != reflect.Slice {
		return nil, nil, errors.New("data must be a slice")
	}

	if v.Len() == 0 {
		return customHeaders, nil, nil
	}

	firstElem := v.Index(0)
	if firstElem.Kind() == reflect.Interface {
		firstElem = firstElem.Elem()
	}

	var headers []string
	var rows [][]any

	switch firstElem.Kind() {
	case reflect.Map:
		// Handle []map[string]any
		if len(customHeaders) > 0 {
			headers = customHeaders
		} else {
			// Get keys from first map
			for _, key := range firstElem.MapKeys() {
				headers = append(headers, fmt.Sprintf("%v", key.Interface()))
			}
		}

		for i := 0; i < v.Len(); i++ {
			elem := v.Index(i)
			if elem.Kind() == reflect.Interface {
				elem = elem.Elem()
			}
			row := make([]any, len(headers))
			for j, header := range headers {
				val := elem.MapIndex(reflect.ValueOf(header))
				if val.IsValid() {
					row[j] = val.Interface()
				}
			}
			rows = append(rows, row)
		}

	case reflect.Struct:
		// Handle slice of structs
		t := firstElem.Type()

		if len(customHeaders) > 0 {
			headers = customHeaders
		} else {
			for i := 0; i < t.NumField(); i++ {
				field := t.Field(i)
				if !field.IsExported() {
					continue
				}
				// Check for excel/csv tags
				header := field.Tag.Get("excel")
				if header == "" {
					header = field.Tag.Get("csv")
				}
				if header == "" || header == "-" {
					header = field.Name
				}
				headers = append(headers, header)
			}
		}

		for i := 0; i < v.Len(); i++ {
			elem := v.Index(i)
			row := make([]any, 0)
			for j := 0; j < t.NumField(); j++ {
				if !t.Field(j).IsExported() {
					continue
				}
				row = append(row, elem.Field(j).Interface())
			}
			rows = append(rows, row)
		}

	case reflect.Ptr:
		// Handle pointer to struct
		if firstElem.Elem().Kind() == reflect.Struct {
			// Dereference and process
			derefData := reflect.MakeSlice(reflect.SliceOf(firstElem.Elem().Type()), v.Len(), v.Len())
			for i := 0; i < v.Len(); i++ {
				derefData.Index(i).Set(v.Index(i).Elem())
			}
			return convertToRows(derefData.Interface(), customHeaders)
		}
		return nil, nil, errors.New("unsupported data type")

	default:
		return nil, nil, fmt.Errorf("unsupported element type: %v", firstElem.Kind())
	}

	return headers, rows, nil
}

// mapRowsToStructs maps string rows to struct slice
func mapRowsToStructs(headers []string, rows [][]string, dest any) error {
	v := reflect.ValueOf(dest)
	if v.Kind() != reflect.Ptr || v.Elem().Kind() != reflect.Slice {
		return errors.New("dest must be a pointer to slice")
	}

	sliceVal := v.Elem()
	elemType := sliceVal.Type().Elem()

	// Build header to field index map
	headerToField := make(map[string]int)
	for i := 0; i < elemType.NumField(); i++ {
		field := elemType.Field(i)
		if !field.IsExported() {
			continue
		}
		// Check for excel/csv tags
		tag := field.Tag.Get("excel")
		if tag == "" {
			tag = field.Tag.Get("csv")
		}
		if tag == "" || tag == "-" {
			tag = field.Name
		}
		for _, header := range headers {
			if strings.EqualFold(tag, header) || strings.EqualFold(field.Name, header) {
				headerToField[header] = i
				break
			}
		}
	}

	// Map rows to structs
	for _, row := range rows {
		newElem := reflect.New(elemType).Elem()
		for colIdx, header := range headers {
			fieldIdx, ok := headerToField[header]
			if !ok || colIdx >= len(row) {
				continue
			}

			field := newElem.Field(fieldIdx)
			if !field.CanSet() {
				continue
			}

			value := row[colIdx]
			if err := setFieldValue(field, value); err != nil {
				continue // Skip on error
			}
		}
		sliceVal.Set(reflect.Append(sliceVal, newElem))
	}

	return nil
}

// setFieldValue sets field value from string
func setFieldValue(field reflect.Value, value string) error {
	switch field.Kind() {
	case reflect.String:
		field.SetString(value)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if value == "" {
			field.SetInt(0)
			return nil
		}
		i, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return err
		}
		field.SetInt(i)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if value == "" {
			field.SetUint(0)
			return nil
		}
		u, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			return err
		}
		field.SetUint(u)
	case reflect.Float32, reflect.Float64:
		if value == "" {
			field.SetFloat(0)
			return nil
		}
		f, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return err
		}
		field.SetFloat(f)
	case reflect.Bool:
		b := strings.ToLower(value) == "true" || value == "1"
		field.SetBool(b)
	default:
		return fmt.Errorf("unsupported field type: %v", field.Kind())
	}
	return nil
}
