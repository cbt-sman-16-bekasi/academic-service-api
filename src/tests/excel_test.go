package tests

import (
	"fmt"
	"github.com/xuri/excelize/v2"
	"github.com/yon-module/yon-framework/logger"
	_ "image/jpeg" // Jika pakai JPG
	_ "image/png"  // WAJIB untuk format PNG
	"os"
	"path/filepath"
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

func TestKopHeader(t *testing.T) {
	f := excelize.NewFile()

	// Pastikan file ada
	if _, err := os.Stat("./img.png"); os.IsNotExist(err) {
		fmt.Println("File img.png tidak ditemukan")
		return
	}

	file, _ := os.ReadFile(filepath.Clean("./img.png"))
	// Masukkan logo sekolah (PNG atau JPG)
	if err := f.AddPictureFromBytes("Sheet1", "A1", &excelize.Picture{
		Extension: ".png",
		File:      file,
		Format: &excelize.GraphicOptions{
			ScaleX: 0.5,
			ScaleY: 0.5,
		},
	}); err != nil {
		fmt.Println("Gagal menambahkan logo:", err)
		return
	}

	// Merge cells untuk nama dan alamat sekolah
	f.MergeCell("Sheet1", "B1", "G1")
	f.MergeCell("Sheet1", "B2", "G2")

	// Set isi dan style
	f.SetCellValue("Sheet1", "B1", "SDN MAJU BERSAMA")
	f.SetCellValue("Sheet1", "B2", "Jl. Pendidikan No. 123, Sukomulyo")

	// Style untuk tulisan kop
	styleKop, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold: true,
			Size: 14,
		},
		Alignment: &excelize.Alignment{
			Horizontal: "left",
			Vertical:   "left",
		},
	})
	f.SetCellStyle("Sheet1", "B1", "G1", styleKop)

	// Style untuk alamat
	styleAlamat, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Italic: true,
			Size:   11,
		},
		Alignment: &excelize.Alignment{
			Horizontal: "left",
			Vertical:   "left",
		},
	})
	f.SetCellStyle("Sheet1", "B2", "G2", styleAlamat)

	// Tambahkan garis bawah di seluruh baris ke-3 (sebagai underline kop)
	styleUnderline, _ := f.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "bottom", Color: "000000", Style: 2},
		},
	})
	f.SetCellStyle("Sheet1", "A3", "G3", styleUnderline)

	// Simpan file
	if err := f.SaveAs("kop_sekolah.xlsx"); err != nil {
		fmt.Println("Gagal menyimpan file:", err)
	}
}
