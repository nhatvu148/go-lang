package main

import (
	"log"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize/v2"
)

func main() {
	f, err := excelize.OpenFile("Stress_Acc_template.xlsx")
	// f := excelize.NewFile()

	if err != nil {
		log.Fatal(err)
	}

	// c1, err := f.GetCellValue("MENR", "A1")

	// if err != nil {
	// 	log.Fatal(err)
	// }

	// fmt.Println(c1)

	f.SetCellValue("MENR", "A2", 100)
	f.SetCellValue("MENR", "B2", 50)
	f.SetCellValue("MENR", "C2", 50)
	f.SetCellValue("MENR", "D2", 50)
	f.SetCellValue("MENR", "E2", 50)
	f.SetCellValue("MENR", "F2", 50)

	now := time.Now()

	f.SetCellValue("Sheet1", "A4", now.Format(time.ANSIC))

	if err := f.SaveAs("Stress_Acc.xlsx"); err != nil {
		log.Fatal(err)
	}
}
