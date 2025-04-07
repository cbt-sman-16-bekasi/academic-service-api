package tests

import (
	"fmt"
	"github.com/xuri/excelize/v2"
	"github.com/yon-module/yon-framework/logger"
	"os"
	"testing"
)

func TestUploadExcelUseImage(t *testing.T) {
	f, err := excelize.OpenFile("Sample Photo.xlsx")
	if err != nil {
		logger.Log.Fatal().Msgf("OpenFile err: %v", err)
	}
	defer f.Close()

	sheetName := "Sheet1"
	// Ambil semua gambar dari sheet
	pictures, err := f.GetPictures(sheetName, "A2")
	if err != nil {
		logger.Log.Fatal().Msgf("GetPictures err: %v", err)
	}

	// Loop melalui setiap gambar yang ditemukan
	for i, pic := range pictures {
		fmt.Printf("Gambar %d:\n", i+1)
		//fmt.Printf(" - Cell: %s\n", pic.Cell)
		fmt.Printf(" - Format: %s\n", pic.Extension)
		fmt.Printf(" - Nama File: %v\n", pic.InsertType)
		fmt.Printf(" - Nama File: %v\n", pic.InsertType)
		err := os.WriteFile(fmt.Sprintf("image-%v%s", i+1, pic.Extension), pic.File, 0644)
		if err != nil {
			fmt.Println("Gagal menyimpan gambar:", err)
			return
		}
		//fmt.Printf(" - Dimensi: %dx%d px\n", pic.Width, pic.Height)
	}

	fmt.Println("Selesai membaca gambar dalam Excel.")
}

func TestUploadExcelUseImage2(t *testing.T) {
}
