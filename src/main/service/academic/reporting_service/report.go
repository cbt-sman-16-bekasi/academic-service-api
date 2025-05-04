package reporting_service

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/helper/bucket"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/model/entity/school"
	"github.com/minio/minio-go/v7"
	"github.com/xuri/excelize/v2"
	"github.com/yon-module/yon-framework/logger"
	_ "image/jpeg" // Jika pakai JPG
	_ "image/png"  // WAJIB untuk format PNG
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Report struct {
	school    school.School
	data      []DataExamSession
	excelFile *excelize.File

	patternStyleBackgroundWhite int
	styleKopeHeader             int
	styleAddress                int
	styleUnderlineHeader        int
	styleOutline                int
	leftOutline                 int
	topOutline                  int
	bottomOutline               int
	rightOutline                int
	styleHeaderCenter           int
	styleMengetahui             int

	IsSuccess  bool
	ResultUrl  *string
	uploadInfo *minio.UploadInfo
	Error      []error
}

func NewReport(school school.School) *Report {
	f := excelize.NewFile()
	styleWhite, _ := f.NewStyle(&excelize.Style{
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{"#FFFFFF"},
			Pattern: 1,
		},
	})

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

	styleAddress, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Italic: true,
			Size:   11,
		},
		Alignment: &excelize.Alignment{
			Horizontal: "left",
			Vertical:   "left",
		},
	})

	styleUnderlineHeader, _ := f.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "bottom", Color: "000000", Style: 2},
		},
	})

	styleOutline, _ := f.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
		},
	})
	leftOutline, _ := f.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
		},
	})
	topOutline, _ := f.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "top", Color: "000000", Style: 1},
		},
	})
	bottomOutline, _ := f.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "bottom", Color: "000000", Style: 1},
		},
	})
	rightOutline, _ := f.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "right", Color: "000000", Style: 1},
		},
	})
	styleHeaderCenter, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold: true,
			Size: 11,
		},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
		},
	})
	styleMengetahui, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold: true,
			Size: 11,
		},
		Alignment: &excelize.Alignment{
			Horizontal: "left",
			Vertical:   "left",
		},
	})

	return &Report{
		school:                      school,
		excelFile:                   f,
		patternStyleBackgroundWhite: styleWhite,
		styleKopeHeader:             styleKop,
		styleAddress:                styleAddress,
		styleUnderlineHeader:        styleUnderlineHeader,
		styleOutline:                styleOutline,
		topOutline:                  topOutline,
		bottomOutline:               bottomOutline,
		rightOutline:                rightOutline,
		styleHeaderCenter:           styleHeaderCenter,
		styleMengetahui:             styleMengetahui,
		leftOutline:                 leftOutline,
	}
}

func (report *Report) SetData(data []DataExamSession) *Report {
	report.data = data
	return report
}

func (report *Report) Generate() error {
	var allError []error
	for _, dataExamSession := range report.data {
		sheetName := dataExamSession.ClassName
		activeSheet, _ := report.excelFile.NewSheet(sheetName)
		report.excelFile.SetActiveSheet(activeSheet)

		// Set background
		report.excelFile.SetCellStyle(sheetName, "A1", "Z100", report.styleKopeHeader)

		// Set width
		report.excelFile.SetColWidth(sheetName, "A", "B", 14.29)
		report.excelFile.SetColWidth(sheetName, "B", "C", 14.29)
		report.excelFile.SetColWidth(sheetName, "D", "E", 14.29)
		report.excelFile.SetColWidth(sheetName, "E", "F", 14.29)

		if err := report.addPictureFromURL(sheetName, "A1", report.school.Logo, 0.5); err != nil {
			logger.Log.Error().Msgf("AddPictureFromURL for class %s error: %v", dataExamSession.ClassName, err)
			allError = append(allError, fmt.Errorf("AddPictureFromURL for class %s error: %v", dataExamSession.ClassName, err))
			continue
		}
		// Merge cells untuk nama dan alamat sekolah
		report.excelFile.MergeCell(sheetName, "B1", "F1")
		report.excelFile.MergeCell(sheetName, "B2", "F2")

		report.excelFile.SetCellValue(sheetName, "B1", report.school.SchoolName)
		report.excelFile.SetCellValue(sheetName, "B2", report.school.Address)

		report.excelFile.SetCellStyle(sheetName, "B1", "F1", report.styleKopeHeader)
		report.excelFile.SetCellStyle(sheetName, "B2", "F2", report.styleAddress)
		report.excelFile.SetCellStyle(sheetName, "A3", "F3", report.styleUnderlineHeader)

		report.excelFile.SetCellValue(sheetName, "A5", "Jenis Ujian")
		report.excelFile.SetCellValue(sheetName, "B5", dataExamSession.TypeExam)

		report.excelFile.SetCellValue(sheetName, "A6", "Mata Pelajaran")
		report.excelFile.SetCellValue(sheetName, "B6", dataExamSession.Subject)

		report.excelFile.SetCellValue(sheetName, "A7", "Kelas")
		report.excelFile.SetCellValue(sheetName, "B7", dataExamSession.ClassName)

		report.excelFile.SetCellValue(sheetName, "A8", "Nama Sesi")
		report.excelFile.SetCellValue(sheetName, "B8", dataExamSession.SessionName)

		report.excelFile.SetCellValue(sheetName, "A9", "Mulai")
		report.excelFile.SetCellValue(sheetName, "B9", dataExamSession.SessionStart.Format("02/01/2006 15:04:05"))

		report.excelFile.SetCellValue(sheetName, "A10", "Selesai")
		report.excelFile.SetCellValue(sheetName, "B10", dataExamSession.SessionEnd.Format("02/01/2006 15:04:05"))

		report.excelFile.SetCellStyle(sheetName, "A5", "A10", report.leftOutline)
		report.excelFile.SetCellStyle(sheetName, "A5", "F5", report.topOutline)
		report.excelFile.SetCellStyle(sheetName, "A10", "F10", report.bottomOutline)
		report.excelFile.SetCellStyle(sheetName, "F5", "F10", report.rightOutline)
		report.excelFile.SetCellStyle(sheetName, "F5", "F5", report.topOutline)
		report.excelFile.SetCellStyle(sheetName, "F10", "F10", report.bottomOutline)

		report.excelFile.SetCellValue(sheetName, "A12", "NO")
		report.excelFile.SetCellValue(sheetName, "B12", "NISN")
		report.excelFile.SetCellValue(sheetName, "C12", "NAMA SISWA")
		report.excelFile.SetCellValue(sheetName, "D12", "KELAS")
		report.excelFile.SetCellValue(sheetName, "E12", "JENIS KELAMIN")
		report.excelFile.SetCellValue(sheetName, "F12", "NILAI")

		report.excelFile.SetCellStyle(sheetName, "A12", "F12", report.styleHeaderCenter)
		currentCellNumber := 13
		for i, data := range dataExamSession.ScoreData {
			currentCellNumber = currentCellNumber + 1
			report.excelFile.SetCellValue(sheetName, "A"+strconv.Itoa(currentCellNumber), i+1)
			report.excelFile.SetCellValue(sheetName, "B"+strconv.Itoa(currentCellNumber), data.NISN)
			report.excelFile.SetCellValue(sheetName, "C"+strconv.Itoa(currentCellNumber), data.Name)
			report.excelFile.SetCellValue(sheetName, "D"+strconv.Itoa(currentCellNumber), data.ClassName)
			report.excelFile.SetCellValue(sheetName, "E"+strconv.Itoa(currentCellNumber), data.Gender)
			report.excelFile.SetCellValue(sheetName, "F"+strconv.Itoa(currentCellNumber), data.Score)
		}

		report.excelFile.SetCellStyle(sheetName, "A12", "F"+strconv.Itoa(currentCellNumber), report.styleOutline)

		// TTD
		currentCellNumber = currentCellNumber + 2
		report.excelFile.SetCellValue(sheetName, "E"+strconv.Itoa(currentCellNumber+1), formatTanggalIndonesia(time.Now()))
		report.excelFile.SetCellValue(sheetName, "E"+strconv.Itoa(currentCellNumber+2), "Mengetahui")
		report.excelFile.SetCellStyle(sheetName, "E"+strconv.Itoa(currentCellNumber+2), "E17"+strconv.Itoa(currentCellNumber+2), report.styleMengetahui)
		report.excelFile.SetCellValue(sheetName, "B"+strconv.Itoa(currentCellNumber+3), "Kepala Sekolah")
		report.excelFile.SetCellValue(sheetName, "E"+strconv.Itoa(currentCellNumber+3), "Wakasek Kurikulum")
		report.excelFile.SetCellValue(sheetName, "B"+strconv.Itoa(currentCellNumber+7), report.school.PrincipalName)
		report.excelFile.SetCellValue(sheetName, "E"+strconv.Itoa(currentCellNumber+7), report.school.VicePrincipalName)
		report.excelFile.SetCellValue(sheetName, "B"+strconv.Itoa(currentCellNumber+8), fmt.Sprintf("(%s)", report.school.PrincipalNIP))
		report.excelFile.SetCellValue(sheetName, "E"+strconv.Itoa(currentCellNumber+8), fmt.Sprintf("(%s)", report.school.VicePrincipalNIP))
	}
	if len(report.data) == 0 {
		return errors.New("empty data")
	}

	err := report.saveReportToMinio()

	if err != nil {
		return err
	}

	report.IsSuccess = true
	report.Error = allError
	return nil
}

func (report *Report) saveReportToMinio() error {
	firstData := report.data[0]
	filename := fmt.Sprintf("%s_%s_%s.xlsx",
		firstData.TypeExam,
		firstData.Subject,
		firstData.SessionEnd.Format("20060102"))

	// Simpan ke buffer (tidak ke file lokal)
	buf := new(bytes.Buffer)
	if err := report.excelFile.Write(buf); err != nil {
		return fmt.Errorf("failed to write excel to buffer: %w", err)
	}

	// Upload ke MinIO
	clientMinio := bucket.NewMinio()
	uploadInfo, publicUrl, err := clientMinio.UploadObject(buf, "report", filename, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	if err != nil {
		return fmt.Errorf("failed to upload to minio: %w", err)
	}

	report.uploadInfo = uploadInfo
	report.ResultUrl = publicUrl

	return nil
}

// formatTanggalIndonesia mengubah time.Time menjadi string tanggal dalam format Indonesia
func formatTanggalIndonesia(date time.Time) string {
	bulanIndonesia := [...]string{
		"Januari", "Februari", "Maret", "April", "Mei", "Juni",
		"Juli", "Agustus", "September", "Oktober", "November", "Desember",
	}

	return fmt.Sprintf("%02d %s %d", date.Day(), bulanIndonesia[date.Month()-1], date.Year())
}

func (report *Report) addPictureFromURL(sheet, cell, url string, scale float64) error {
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("gagal mengunduh gambar: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("gagal mengunduh gambar: status %s", resp.Status)
	}

	imgBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("gagal membaca isi gambar: %w", err)
	}

	ext := ".png"
	if strings.HasSuffix(strings.ToLower(url), ".jpg") || strings.HasSuffix(strings.ToLower(url), ".jpeg") {
		ext = ".jpg"
	}

	return report.excelFile.AddPictureFromBytes(sheet, cell, &excelize.Picture{
		Extension: ext,
		File:      imgBytes,
		Format: &excelize.GraphicOptions{
			ScaleX: scale,
			ScaleY: scale,
		},
	})
}

func (report *Report) GetResult() (*string, bool) {
	return report.ResultUrl, report.IsSuccess
}

func (report *Report) GetError() []error {
	return report.Error
}
