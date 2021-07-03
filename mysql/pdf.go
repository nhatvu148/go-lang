package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/jung-kurt/gofpdf"
)

func main() {
	pdf := gofpdf.New("P", "mm", "A4", "")
	type countryType struct {
		nameStr, capitalStr, areaStr, popStr string
	}
	countryList := make([]countryType, 0, 8)
	header := []string{"Country", "Capital", "Area (sq km)", "Pop. (thousands)"}
	loadData := func(fileStr string) {
		fl, err := os.Open(fileStr)
		if err == nil {
			scanner := bufio.NewScanner(fl)
			var c countryType
			for scanner.Scan() {
				// Austria;Vienna;83859;8075
				lineStr := scanner.Text()
				list := strings.Split(lineStr, ";")
				if len(list) == 4 {
					c.nameStr = list[0]
					c.capitalStr = list[1]
					c.areaStr = list[2]
					c.popStr = list[3]
					countryList = append(countryList, c)
				} else {
					err = fmt.Errorf("error tokenizing %s", lineStr)
				}
			}
			fl.Close()
			if len(countryList) == 0 {
				err = fmt.Errorf("error loading data from %s", fileStr)
			}
		}
		if err != nil {
			pdf.SetError(err)
		}
	}
	// Simple table
	// basicTable := func() {
	// 	left := (210.0 - 4*40) / 2
	// 	pdf.SetX(left)
	// 	for _, str := range header {
	// 		pdf.CellFormat(40, 7, str, "1", 0, "", false, 0, "")
	// 	}
	// 	pdf.Ln(-1)
	// 	for _, c := range countryList {
	// 		pdf.SetX(left)
	// 		pdf.CellFormat(40, 6, c.nameStr, "1", 0, "", false, 0, "")
	// 		pdf.CellFormat(40, 6, c.capitalStr, "1", 0, "", false, 0, "")
	// 		pdf.CellFormat(40, 6, c.areaStr, "1", 0, "", false, 0, "")
	// 		pdf.CellFormat(40, 6, c.popStr, "1", 0, "", false, 0, "")
	// 		pdf.Ln(-1)
	// 	}
	// }
	// Better table
	improvedTable := func() {
		// Column widths
		w := []float64{40.0, 35.0, 40.0, 45.0}
		wSum := 0.0
		for _, v := range w {
			wSum += v
		}
		left := (210 - wSum) / 2
		// 	Header
		pdf.SetX(left)
		for j, str := range header {
			pdf.CellFormat(w[j], 7, str, "1", 0, "C", false, 0, "")
		}
		pdf.Ln(-1)
		// Data
		for _, c := range countryList {
			pdf.SetX(left)
			pdf.CellFormat(w[0], 6, c.nameStr, "LR", 0, "", false, 0, "")
			pdf.CellFormat(w[1], 6, c.capitalStr, "LR", 0, "", false, 0, "")
			pdf.CellFormat(w[2], 6, strDelimit(c.areaStr, ",", 3),
				"LR", 0, "R", false, 0, "")
			pdf.CellFormat(w[3], 6, strDelimit(c.popStr, ",", 3),
				"LR", 0, "R", false, 0, "")
			pdf.Ln(-1)
		}
		pdf.SetX(left)
		pdf.CellFormat(wSum, 0, "", "T", 0, "", false, 0, "")
	}
	fancyTable := func() {
		// Colors, line width and bold font
		pdf.SetFillColor(255, 0, 0)
		pdf.SetTextColor(255, 255, 255)
		pdf.SetDrawColor(128, 0, 0)
		pdf.SetLineWidth(.3)
		pdf.SetFont("", "B", 0)
		// 	Header
		w := []float64{40, 35, 40, 45}
		wSum := 0.0
		for _, v := range w {
			wSum += v
		}
		left := (210 - wSum) / 2
		pdf.SetX(left)
		for j, str := range header {
			pdf.CellFormat(w[j], 7, str, "1", 0, "C", true, 0, "")
		}
		pdf.Ln(-1)
		// Color and font restoration
		pdf.SetFillColor(224, 235, 255)
		pdf.SetTextColor(0, 0, 0)
		pdf.SetFont("", "", 0)
		// 	Data
		fill := false
		for _, c := range countryList {
			pdf.SetX(left)
			pdf.CellFormat(w[0], 6, c.nameStr, "LR", 0, "", fill, 0, "")
			pdf.CellFormat(w[1], 6, c.capitalStr, "LR", 0, "", fill, 0, "")
			pdf.CellFormat(w[2], 6, strDelimit(c.areaStr, ",", 3),
				"LR", 0, "R", fill, 0, "")
			pdf.CellFormat(w[3], 6, strDelimit(c.popStr, ",", 3),
				"LR", 0, "R", fill, 0, "")
			pdf.Ln(-1)
			fill = !fill
		}
		pdf.SetX(left)
		pdf.CellFormat(wSum, 0, "", "T", 0, "", false, 0, "")
	}
	loadData(TextFile("countries.txt"))
	pdf.SetFont("Arial", "", 14)
	pdf.AddPage()
	// pdf.SetFont("Arial", "B", 16)
	pdf.Cell(80, 0, "Hello, world")
	// basicTable()
	pdf.AddPage()
	improvedTable()
	pdf.AddPage()
	fancyTable()
	fileStr := Filename("Fpdf_CellFormat_tables")
	err := pdf.OutputFileAndClose(fileStr)
	Summary(err, fileStr)
}

func PdfDir() string {
	return filepath.Join(".", "pdf")
}

func PdfFile(fileStr string) string {
	return filepath.Join(PdfDir(), fileStr)
}

func Filename(baseStr string) string {
	return PdfFile(baseStr + ".pdf")
}

func Summary(err error, fileStr string) {
	if err == nil {
		fileStr = filepath.ToSlash(fileStr)
		fmt.Printf("Successfully generated %s\n", fileStr)
	} else {
		fmt.Println(err)
	}
}

func strDelimit(str string, sepstr string, sepcount int) string {
	pos := len(str) - sepcount
	for pos > 0 {
		str = str[:pos] + sepstr + str[pos:]
		pos = pos - sepcount
	}
	return str
}

func TextFile(fileStr string) string {
	return filepath.Join(".", "text", fileStr)
}
