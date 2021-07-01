package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize/v2"
	_ "github.com/go-sql-driver/mysql"
)

type StateStatistics struct {
	State_Measure_ID  int
	ShipInfo_ID       int
	datetime          string
	NumofProcess      int
	NumofMeasurePoint int
	MeasurePointData  string
}

func main() {
	start := time.Now()

	db, err := sql.Open("mysql", "root:123456789@tcp(127.0.0.1:3306)/jmu")
	defer db.Close()

	if err != nil {
		log.Fatal(err)
	}

	res, err := db.Query("SELECT * FROM statistics.state_statistics WHERE ShipInfo_ID='1'")

	defer res.Close()

	if err != nil {
		log.Fatal(err)
	}

	alphaMap := make(map[int]string)
	for i := 1; i <= 26; i++ {
		alphaMap[i] = string(rune(i + 64))
		alphaMap[i+26] = fmt.Sprintf("%v%v", "A", string(rune(i+64)))
	}
	menrList := [][]float64{}
	devlList := [][]float64{}
	numofMeasurePoint := 0

	for res.Next() {

		var state_statistics StateStatistics
		err := res.Scan(&state_statistics.State_Measure_ID,
			&state_statistics.ShipInfo_ID,
			&state_statistics.datetime,
			&state_statistics.NumofProcess,
			&state_statistics.NumofMeasurePoint,
			&state_statistics.MeasurePointData)

		if err != nil {
			log.Fatal(err)
		}

		numofMeasurePoint = state_statistics.NumofMeasurePoint
		var m1 []float64
		json.Unmarshal([]byte(state_statistics.MeasurePointData), &m1)
		menrList = append(menrList, m1[0:(numofMeasurePoint+2)])
		devlList = append(devlList, m1[(numofMeasurePoint+2):(2*numofMeasurePoint+4)])
	}

	f, err := excelize.OpenFile("Stress_Acc_template.xlsx")
	if err != nil {
		log.Fatal(err)
	}

	for i := 1; i <= len(menrList); i++ {
		for j := 1; j <= numofMeasurePoint+2; j++ {
			cell := fmt.Sprintf("%v%v", alphaMap[j], i+1)
			f.SetCellValue("MENR", cell, menrList[i-1][j-1])
		}
	}

	for i := 1; i <= len(devlList); i++ {
		for j := 1; j <= numofMeasurePoint+2; j++ {
			cell := fmt.Sprintf("%v%v", alphaMap[j], i+1)
			f.SetCellValue("DEVL", cell, devlList[i-1][j-1])
		}
	}

	if err := f.SaveAs("Stress_Acc.xlsx"); err != nil {
		log.Fatal(err)
	}

	elapsed := time.Since(start)
	log.Printf("Exporting excels took %s", elapsed)
}
