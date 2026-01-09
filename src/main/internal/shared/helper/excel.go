package helper

import (
	"fmt"
	"io"
	"mime/multipart"
	"strings"

	"github.com/xuri/excelize/v2"
)

// ReadAndValidateExcel reads an Excel file and returns rows (skipping header)
func ReadAndValidateExcel(file multipart.File) ([][]string, error) {
	content, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	f, err := excelize.OpenReader(strings.NewReader(string(content)))
	if err != nil {
		return nil, err
	}
	defer f.Close()

	sheetName := f.GetSheetName(0)
	if sheetName == "" {
		return nil, fmt.Errorf("first sheet undefined")
	}

	rows, err := f.GetRows(sheetName)
	if err != nil {
		return nil, err
	}

	if len(rows) < 2 {
		return nil, fmt.Errorf("file has no data rows")
	}

	// Skip header row
	return rows[1:], nil
}
