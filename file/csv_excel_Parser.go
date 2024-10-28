package fileparser

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"mime/multipart"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
)

func ParseFile(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		log.Printf("Error retrieving file: %v", err)
		c.JSON(400, gin.H{"error": "Error getting file from request"})
		return
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))

	headerLabelsJSON := c.PostForm("headerLabels")
	var headerLabels map[string]string
	if err := json.Unmarshal([]byte(headerLabelsJSON), &headerLabels); err != nil {
		c.JSON(400, gin.H{"error": "Invalid JSON format for headerLabels"})
		return
	}

	considerFirstAsData := c.PostForm("considerFirstAsData") == "true"
	trimData := c.PostForm("trimData") == "true"

	if ext == ".xlam" || ext == ".xlsm" || ext == ".xlsx" || ext == ".xltm" || ext == ".xltx" {
		c.JSON(200, ReadFile(file, headerLabels, considerFirstAsData, trimData))
	} else if ext == ".csv" {
		c.JSON(200, ReadcsvFile(file, headerLabels, considerFirstAsData, trimData))
	} else {
		c.JSON(400, gin.H{"error": "Unsupported file type"})
	}
}

func ReadFile(file *multipart.FileHeader, headerLabels map[string]string, considerFirstAsData, trimData bool) []map[string]string {
	startTime := time.Now()
	defer func() {
		log.Println(time.Since(startTime))
	}()

	var output []map[string]string
	uploadedFile, err := file.Open()
	if err != nil {
		return nil
	}
	defer uploadedFile.Close()

	excelFile, err := excelize.OpenReader(uploadedFile)
	if err != nil {
		return nil
	}
	rows, err := excelFile.GetRows(excelFile.GetSheetName(0))
	if err != nil {
		return nil
	}

	if len(rows) > 0 {
		if considerFirstAsData {
			rows = append([][]string{make([]string, len(rows[0]))}, rows...)
		}
		cleanedHeaders, cleanedRows := CleanParsedData(rows[0], rows[1:], headerLabels, trimData)

		for i := 0; i < len(cleanedRows); i++ {
			row := cleanedRows[i]
			entry := make(map[string]string)
			for j := 0; j < len(cleanedHeaders); j++ {
				if j < len(row) {
					entry[cleanedHeaders[j]] = row[j]
				} else {
					entry[cleanedHeaders[j]] = ""
				}
			}
			output = append(output, entry)
		}
	}
	return output
}

func ReadcsvFile(file *multipart.FileHeader, headerLabels map[string]string, considerFirstAsData, trimData bool) []map[string]string {
	startTime := time.Now()
	defer func() {
		log.Println(time.Since(startTime))
	}()

	var output []map[string]string
	uploadedFile, err := file.Open()
	if err != nil {
		return nil
	}
	defer uploadedFile.Close()

	reader := csv.NewReader(uploadedFile)

	rows, err := reader.ReadAll()
	if err != nil {
		return nil
	}

	if len(rows) > 0 {
		if considerFirstAsData {
			rows = append([][]string{make([]string, len(rows[0]))}, rows...)
		}
		cleanedHeaders, cleanedRows := CleanParsedData(rows[0], rows[1:], headerLabels, trimData)

		for i := 0; i < len(cleanedRows); i++ {
			row := cleanedRows[i]
			entry := make(map[string]string)
			for j := 0; j < len(cleanedHeaders); j++ {
				if j < len(row) {
					entry[cleanedHeaders[j]] = row[j]
				} else {
					entry[cleanedHeaders[j]] = ""
				}
			}
			output = append(output, entry)
		}
	}
	return output
}

// CleanParsedData function to clean and modify headers and rows based on the given options
func CleanParsedData(headers []string, rows [][]string, headerLabels map[string]string, trimData bool) ([]string, [][]string) {
	headerMap := make(map[string]int)
	cleanedHeaders := make([]string, 0, len(headers))
	cleanedRows := make([][]string, 0, len(rows))
	defaultHeaderCount := 1
	defaultAttriCount := 1

	// Apply header labels if provided
	for index, newHeader := range headerLabels {
		i, err := strconv.Atoi(index)
		if err == nil && i < len(headers) {
			headers[i] = newHeader
		}
	}

	// Clean headers
	for _, header := range headers {
		header = strings.TrimSpace(header)
		if header == "" {
			header = fmt.Sprintf("header%d", defaultHeaderCount)
			defaultHeaderCount++
		}
		if _, exists := headerMap[header]; exists {
			header = fmt.Sprintf("%s_dupli%d", header, headerMap[header])
			headerMap[header]++
		} else {
			headerMap[header] = 1
		}
		cleanedHeaders = append(cleanedHeaders, header)
	}

	// Adjust rows based on cleaned headers and apply trimming if needed
	for _, row := range rows {
		cleanedRow := make([]string, len(cleanedHeaders))

		for j := 0; j < len(cleanedHeaders); j++ {
			if j < len(row) {
				value := strings.TrimSpace(row[j])
				if trimData {
					value = trimCellData(value)
				}
				cleanedRow[j] = value
			} else {
				cleanedRow[j] = fmt.Sprintf("emptyattri%d", defaultAttriCount)
				defaultAttriCount++
			}
		}

		if len(row) > len(cleanedHeaders) {
			for j := len(cleanedHeaders); j < len(row); j++ {
				if row[j] != "" {
					headerName := fmt.Sprintf("header%d", defaultHeaderCount)
					cleanedHeaders = append(cleanedHeaders, headerName)
					defaultHeaderCount++
					cleanedRow = append(cleanedRow, strings.TrimSpace(row[j]))
				}
			}
		}

		cleanedRows = append(cleanedRows, cleanedRow)
	}

	return cleanedHeaders, cleanedRows
}

// trimCellData function to trim and format data in a cell
func trimCellData(data string) string {
	parts := strings.Split(data, ",")
	for i, part := range parts {
		parts[i] = strings.TrimSpace(part)
	}
	return strings.Join(parts, ",")
}

func HandleCSV_ExcelParsing(router *gin.Engine) {
	router.POST("/upload", ParseFile)
}

/**package fileparser

import (
	"encoding/csv"
	"fmt"
	"log"
	"mime/multipart"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
)

func ParseFile(c *gin.Context) {
	log.Println("error file debug")
	file, err := c.FormFile("file")
	if err != nil {
		log.Printf("Error retrieving file: %v", err)
		c.JSON(400, gin.H{"error": "Error getting file from request"})
		return
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))

	headerLabels := c.PostFormArray("headerLabels")
	considerFirstAsData := c.PostForm("considerFirstAsData") == "true"
	trimData := c.PostForm("trimData") == "true"

	if ext == ".xlam" || ext == ".xlsm" || ext == ".xlsx" || ext == ".xltm" || ext == ".xltx" {
		c.JSON(200, ReadFile(file, headerLabels, considerFirstAsData, trimData))
	} else if ext == ".csv" {
		c.JSON(200, ReadcsvFile(file, headerLabels, considerFirstAsData, trimData))
	} else {
		c.JSON(400, gin.H{"error": "Unsupported file type"})
	}
}

func ReadFile(file *multipart.FileHeader, headerLabels []string, considerFirstAsData, trimData bool) []map[string]string {
	startTime := time.Now()
	defer func() {
		log.Println(time.Since(startTime))
	}()

	var output []map[string]string
	uploadedFile, err := file.Open()
	if err != nil {
		return nil
	}
	defer uploadedFile.Close()

	excelFile, err := excelize.OpenReader(uploadedFile)
	if err != nil {
		return nil
	}
	rows, err := excelFile.GetRows(excelFile.GetSheetName(0))
	if err != nil {
		return nil
	}

	if len(rows) > 0 {
		if considerFirstAsData {
			rows = append([][]string{make([]string, len(rows[0]))}, rows...)
		}
		cleanedHeaders, cleanedRows := CleanParsedData(rows[0], rows[1:], headerLabels, trimData)

		for i := 0; i < len(cleanedRows); i++ {
			row := cleanedRows[i]
			entry := make(map[string]string)
			for j := 0; j < len(cleanedHeaders); j++ {
				if j < len(row) {
					entry[cleanedHeaders[j]] = row[j]
				} else {
					entry[cleanedHeaders[j]] = ""
				}
			}
			output = append(output, entry)
		}
	}
	return output
}

func ReadcsvFile(file *multipart.FileHeader, headerLabels []string, considerFirstAsData, trimData bool) []map[string]string {
	startTime := time.Now()
	defer func() {
		log.Println(time.Since(startTime))
	}()

	var output []map[string]string
	uploadedFile, err := file.Open()
	if err != nil {
		return nil
	}
	defer uploadedFile.Close()

	reader := csv.NewReader(uploadedFile)

	rows, err := reader.ReadAll()
	if err != nil {
		return nil
	}

	if len(rows) > 0 {
		if considerFirstAsData {
			rows = append([][]string{make([]string, len(rows[0]))}, rows...)
		}
		cleanedHeaders, cleanedRows := CleanParsedData(rows[0], rows[1:], headerLabels, trimData)

		for i := 0; i < len(cleanedRows); i++ {
			row := cleanedRows[i]
			entry := make(map[string]string)
			for j := 0; j < len(cleanedHeaders); j++ {
				if j < len(row) {
					entry[cleanedHeaders[j]] = row[j]
				} else {
					entry[cleanedHeaders[j]] = ""
				}
			}
			output = append(output, entry)
		}
	}
	return output
}

// CleanParsedData function to clean and modify headers and rows based on the given options
func CleanParsedData(headers []string, rows [][]string, headerLabels []string, trimData bool) ([]string, [][]string) {
	headerMap := make(map[string]int)
	cleanedHeaders := make([]string, 0, len(headers))
	cleanedRows := make([][]string, 0, len(rows))
	defaultHeaderCount := 1
	defaultAttriCount := 1

	// Apply header labels if provided
	if len(headerLabels) > 0 {
		for i := 0; i < len(headers) && i < len(headerLabels); i++ {
			headers[i] = strings.TrimSpace(headerLabels[i])
		}
	}

	// Clean headers
	for _, header := range headers {
		header = strings.TrimSpace(header)
		if header == "" {
			header = fmt.Sprintf("header%d", defaultHeaderCount)
			defaultHeaderCount++
		}
		if _, exists := headerMap[header]; exists {
			header = fmt.Sprintf("%s_dupli%d", header, headerMap[header])
			headerMap[header]++
		} else {
			headerMap[header] = 1
		}
		cleanedHeaders = append(cleanedHeaders, header)
	}

	// Adjust rows based on cleaned headers and apply trimming if needed
	for _, row := range rows {
		cleanedRow := make([]string, len(cleanedHeaders))

		for j := 0; j < len(cleanedHeaders); j++ {
			if j < len(row) {
				value := strings.TrimSpace(row[j])
				if trimData {
					value = trimCellData(value)
				}
				cleanedRow[j] = value
			} else {
				cleanedRow[j] = fmt.Sprintf("emptyattri%d", defaultAttriCount)
				defaultAttriCount++
			}
		}

		if len(row) > len(cleanedHeaders) {
			for j := len(cleanedHeaders); j < len(row); j++ {
				if row[j] != "" {
					headerName := fmt.Sprintf("header%d", defaultHeaderCount)
					cleanedHeaders = append(cleanedHeaders, headerName)
					defaultHeaderCount++
					cleanedRow = append(cleanedRow, strings.TrimSpace(row[j]))
				}
			}
		}

		cleanedRows = append(cleanedRows, cleanedRow)
	}

	return cleanedHeaders, cleanedRows
}

// trimCellData function to trim and format data in a cell
func trimCellData(data string) string {
	parts := strings.Split(data, ",")
	for i, part := range parts {
		parts[i] = strings.TrimSpace(part)
	}
	return strings.Join(parts, ",")
}

func HandleCSV_ExcelParsing(router *gin.Engine) {
	router.POST("/upload", ParseFile)
}

/**package fileparser

import (
	"encoding/csv"
	"fmt"
	"log"
	"mime/multipart"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
)

func ParseFile(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		log.Printf("Error retrieving file: %v", err)

		c.JSON(400, gin.H{"error": "Error getting file from request"})
		return
	}
	ext := strings.ToLower(filepath.Ext(file.Filename))

	// Choose the appropriate parser based on the file extension
	if ext == ".xlam" || ext == ".xlsm" || ext == ".xlsx" || ext == ".xltm" || ext == ".xltx" {
		// Handle Excel file
		c.JSON(200, ReadFile(file))
	} else if ext == ".csv" {
		// Handle CSV file
		c.JSON(200, ReadcsvFile(file))
	} else {
		// Unsupported file type
		c.JSON(400, gin.H{"error": "Unsupported file type"})
	}
}

func ReadFile(file *multipart.FileHeader) []map[string]string {
	startTime := time.Now()
	defer func() {
		log.Println(time.Since(startTime))
	}()

	var output []map[string]string
	uploadedFile, err := file.Open()
	if err != nil {
		return nil
	}
	defer uploadedFile.Close()

	excelFile, err := excelize.OpenReader(uploadedFile)
	if err != nil {
		return nil
	}
	rows, err := excelFile.GetRows(excelFile.GetSheetName(0))
	if err != nil {
		return nil
	}

	// Clean headers and rows
	// Clean headers and rows
	if len(rows) > 0 {
		cleanedHeaders, cleanedRows := CleanParsedData(rows[0], rows[1:])

		for i := 0; i < len(cleanedRows); i++ {
			row := cleanedRows[i]
			entry := make(map[string]string)
			for j := 0; j < len(cleanedHeaders); j++ {
				if j < len(row) {
					entry[cleanedHeaders[j]] = row[j]
				} else {
					entry[cleanedHeaders[j]] = ""
				}
			}
			output = append(output, entry)
		}
	}
	return output
}

func ReadcsvFile(file *multipart.FileHeader) []map[string]string {
	startTime := time.Now()
	defer func() {
		log.Println(time.Since(startTime))
	}()
	var output []map[string]string
	uploadedFile, err := file.Open()
	if err != nil {
		return nil
	}
	defer uploadedFile.Close()

	reader := csv.NewReader(uploadedFile)

	rows, err := reader.ReadAll()
	if err != nil {
		return nil
	}

	// Clean headers and rows
	if len(rows) > 0 {
		cleanedHeaders, cleanedRows := CleanParsedData(rows[0], rows[1:])
		for i := 0; i < len(cleanedRows); i++ {
			row := cleanedRows[i]
			entry := make(map[string]string)
			for j := 0; j < len(cleanedHeaders); j++ {
				if j < len(row) {
					entry[cleanedHeaders[j]] = row[j]
				} else {
					entry[cleanedHeaders[j]] = ""
				}
			}
			output = append(output, entry)
		}
	}
	return output
}

// CleanParsedData function to remove empty headers, empty attributes, spaces in headers, duplicate headers, and merged cells
// CleanParsedData function to remove empty headers, empty attributes, spaces in headers, duplicate headers, and merged cells
func CleanParsedData(headers []string, rows [][]string) ([]string, [][]string) {
	headerMap := make(map[string]int)
	cleanedHeaders := make([]string, 0, len(headers))
	cleanedRows := make([][]string, 0, len(rows))
	defaultHeaderCount := 1
	defaultAttriCount := 1

	// Clean headers
	for _, header := range headers {
		header = strings.TrimSpace(header) // Remove leading/trailing spaces
		if header == "" {
			header = fmt.Sprintf("header%d", defaultHeaderCount)
			defaultHeaderCount++
		}
		if _, exists := headerMap[header]; exists {
			header = fmt.Sprintf("%s_dupli%d", header, headerMap[header])
			headerMap[header]++
		} else {
			headerMap[header] = 1
		}
		cleanedHeaders = append(cleanedHeaders, header)
	}

	// Adjust rows based on cleaned headers
	for _, row := range rows {
		cleanedRow := make([]string, len(cleanedHeaders))

		// Fill existing columns with data
		for j := 0; j < len(cleanedHeaders); j++ {
			if j < len(row) {
				cleanedRow[j] = strings.TrimSpace(row[j]) // Remove leading/trailing spaces from cell values
			} else {
				cleanedRow[j] = fmt.Sprintf("emptyattri%d", defaultAttriCount)
				defaultAttriCount++
			}
		}

		// Handle columns with data but missing headers
		if len(row) > len(cleanedHeaders) {
			for j := len(cleanedHeaders); j < len(row); j++ {
				if row[j] != "" {
					headerName := fmt.Sprintf("header%d", defaultHeaderCount)
					cleanedHeaders = append(cleanedHeaders, headerName)
					defaultHeaderCount++
					cleanedRow = append(cleanedRow, strings.TrimSpace(row[j]))
				}
			}
		}

		cleanedRows = append(cleanedRows, cleanedRow)
	}

	return cleanedHeaders, cleanedRows
}

func HandleCSV_ExcelParsing(router *gin.Engine) {
	router.POST("/upload", ParseFile)
}
**/
