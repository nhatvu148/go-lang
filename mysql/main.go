package main

import (
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize/v2"
	_ "github.com/go-sql-driver/mysql"
)

type StateStatistics struct {
	NumofMeasurePoint int
	MeasurePointData  string
}
type Gyro struct {
	datetime  string
	Roll_Max  float64
	Pitch_Max float64
	Yaw_Max   float64
}
type Wave struct {
	datetime   string
	WaveHeight float64
	WavePeriod float64
}

func main() {
	start := time.Now()

	host := flag.String("host", "localhost", "Host name")
	user := flag.String("user", "root", "User name")
	password := flag.String("password", "123456789", "Password")
	database := flag.String("database", "jmu", "Database")
	shipInfoID := flag.Int("shipInfoID", 1, "Ship information ID")
	startTime := flag.String("startTime", "", "Start time")
	endTime := flag.String("endTime", "", "End time")
	outDir := flag.String("outDir", ".", "Output Directory")

	// var svar string
	// flag.StringVar(&svar, "svar", "bar", "a string var")

	flag.Parse()

	// fmt.Println("host:", *host)
	// fmt.Println("user:", *user)
	// fmt.Println("password:", *password)
	// fmt.Println("database:", *database)
	// fmt.Println("shipInfoID:", *shipInfoID)
	// fmt.Println("startTime:", *startTime)
	// fmt.Println("endTime:", *endTime)
	// fmt.Println("outDir:", *outDir)
	// // fmt.Println("svar:", svar)
	// fmt.Println("tail:", flag.Args())

	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:3306)/%s", *user, *password, *host, *database))
	defer db.Close()

	if err != nil {
		log.Fatal(err)
	}

	alphaMap := make(map[int]string)
	for i := 1; i <= 26; i++ {
		alphaMap[i] = string(rune(i + 64))
		alphaMap[i+26] = fmt.Sprintf("%v%v", "A", string(rune(i+64)))
	}
	var wg sync.WaitGroup
	wg.Add(3)

	go func() {
		defer wg.Done()
		sql := fmt.Sprintf("SELECT NumofMeasurePoint, MeasurePointData FROM statistics.state_statistics WHERE ShipInfo_ID='%d'", *shipInfoID)
		if *startTime != "" && *endTime != "" {
			sql = fmt.Sprintf("SELECT NumofMeasurePoint, MeasurePointData FROM statistics.state_statistics WHERE ShipInfo_ID='%d' AND datetime BETWEEN '%s' AND '%s'", *shipInfoID, *startTime, *endTime)
		}

		res, err := db.Query(sql)

		defer res.Close()

		if err != nil {
			log.Fatal(err)
		}

		menrList := [][]float64{}
		devlList := [][]float64{}
		numofMeasurePoint := 0

		for res.Next() {
			var state_statistics StateStatistics
			err := res.Scan(
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

		if err := f.SaveAs(fmt.Sprintf("%s/Stress_Acc.xlsx", *outDir)); err != nil {
			log.Fatal(err)
		}
		fmt.Println("DONE Stress_Acc")
	}()

	go func() {
		defer wg.Done()
		sql := fmt.Sprintf("SELECT datetime, Roll_Max, Pitch_Max, Yaw_Max FROM statistics.gyro WHERE ShipInfo_ID='%d'", *shipInfoID)
		if *startTime != "" && *endTime != "" {
			sql = fmt.Sprintf("SELECT datetime, Roll_Max, Pitch_Max, Yaw_Max FROM statistics.gyro WHERE ShipInfo_ID='%d' AND datetime BETWEEN '%s' AND '%s'", *shipInfoID, *startTime, *endTime)
		}

		res, err := db.Query(sql)

		defer res.Close()

		if err != nil {
			log.Fatal(err)
		}

		dateList := []string{}
		rollList := []float64{}
		pitchList := []float64{}
		yawList := []float64{}

		for res.Next() {
			var gyro Gyro
			err := res.Scan(
				&gyro.datetime,
				&gyro.Roll_Max,
				&gyro.Pitch_Max,
				&gyro.Yaw_Max,
			)

			if err != nil {
				log.Fatal(err)
			}

			dateList = append(dateList, gyro.datetime)
			rollList = append(rollList, gyro.Roll_Max)
			pitchList = append(pitchList, gyro.Pitch_Max)
			yawList = append(yawList, gyro.Yaw_Max)
		}

		f, err := excelize.OpenFile("Gyro_template.xlsx")
		if err != nil {
			log.Fatal(err)
		}

		for i := 0; i < len(dateList); i++ {
			cell1 := fmt.Sprintf("%v%v", alphaMap[1], i+2)
			cell2 := fmt.Sprintf("%v%v", alphaMap[2], i+2)
			cell3 := fmt.Sprintf("%v%v", alphaMap[3], i+2)
			cell4 := fmt.Sprintf("%v%v", alphaMap[4], i+2)
			f.SetCellValue("GYRO", cell1, dateList[i])
			f.SetCellValue("GYRO", cell2, rollList[i])
			f.SetCellValue("GYRO", cell3, pitchList[i])
			f.SetCellValue("GYRO", cell4, yawList[i])
		}

		if err := f.SaveAs(fmt.Sprintf("%s/Gyro.xlsx", *outDir)); err != nil {
			log.Fatal(err)
		}
		fmt.Println("DONE Gyro")
	}()

	go func() {
		defer wg.Done()
		sql := fmt.Sprintf("SELECT datetime, WaveHeight, WavePeriod FROM statistics.waves WHERE ShipInfo_ID='%d'", *shipInfoID)
		if *startTime != "" && *endTime != "" {
			sql = fmt.Sprintf("SELECT datetime, WaveHeight, WavePeriod FROM statistics.waves WHERE ShipInfo_ID='%d' AND datetime BETWEEN '%s' AND '%s'", *shipInfoID, *startTime, *endTime)
		}

		res, err := db.Query(sql)

		defer res.Close()

		if err != nil {
			log.Fatal(err)
		}

		dateList2 := []string{}
		waveHList := []float64{}
		wavePList := []float64{}

		for res.Next() {
			var wave Wave
			err := res.Scan(
				&wave.datetime,
				&wave.WaveHeight,
				&wave.WavePeriod,
			)

			if err != nil {
				log.Fatal(err)
			}

			dateList2 = append(dateList2, wave.datetime)
			waveHList = append(waveHList, wave.WaveHeight)
			wavePList = append(wavePList, wave.WavePeriod)
		}

		f, err := excelize.OpenFile("Wave_template.xlsx")
		if err != nil {
			log.Fatal(err)
		}

		for i := 0; i < len(dateList2); i++ {
			cell1 := fmt.Sprintf("%v%v", alphaMap[1], i+2)
			cell2 := fmt.Sprintf("%v%v", alphaMap[2], i+2)
			cell3 := fmt.Sprintf("%v%v", alphaMap[3], i+2)
			f.SetCellValue("WaveData", cell1, dateList2[i])
			f.SetCellValue("WaveData", cell2, waveHList[i])
			f.SetCellValue("WaveData", cell3, wavePList[i])
		}

		if err := f.SaveAs(fmt.Sprintf("%s/Wave.xlsx", *outDir)); err != nil {
			log.Fatal(err)
		}
		fmt.Println("DONE Wave")
	}()

	wg.Wait()
	elapsed := time.Since(start)
	log.Printf("Exporting excels took %s", elapsed)
}
