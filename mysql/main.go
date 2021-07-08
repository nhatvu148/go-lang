package main

import (
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"image/color"
	"log"
	"os"
	"os/exec"
	"strconv"
	"sync"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize/v2"
	_ "github.com/go-sql-driver/mysql"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

type StateStatistics struct {
	MENR string
	DEVL string
}
type Gyro struct {
	datetime  string
	Roll_Max  float64
	Pitch_Max float64
	Yaw_Max   float64
}
type Wave struct {
	datetime                          string
	SignificantWaveHeight_Arbitrary00 float64
	WavePeriod_Arbitrary00            float64
}
type Map struct {
	datetime  string
	GPGGA_LAT float64
	GPGGA_LON float64
}
type timeTicks struct{}
type errPoints struct {
	plotter.XYs
	plotter.YErrors
	plotter.XErrors
}

func (timeTicks) Ticks(min, max float64) []plot.Tick {
	tks := plot.DefaultTicks{}.Ticks(min, max)
	for i, t := range tks {
		tks[i].Label = toDateTime(strconv.FormatFloat(t.Value, 'f', -1, 64))
	}
	return tks
}
func toDateTime(s string) string {
	var t time.Time
	if n, err := strconv.ParseInt(s, 10, 64); err == nil {
		t = time.Unix(n, 0)
	}
	return t.Format("2006-01-02")
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
	jpt_root := flag.String("jpt_root", "D:/AKIYAMA/Trunk_Rev56717_ForWeb/bin/Release/x64", "Jupiter Root Directory")
	outCsvDir := flag.String("outCsvDir", fmt.Sprintf("%s/TechnoStar/00", os.TempDir()), "Output CSV Directory")
	flag.Parse()

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
	wg.Add(5)

	go func() {
		defer wg.Done()
		args := []string{
			"cmd",
			"/C",
			fmt.Sprintf("%s/Start_It.bat -b -mcr runme.jpl", *jpt_root),
		}

		CmdExec(args...)
		Transform(*outCsvDir)
	}()

	go func() {
		defer wg.Done()
		sql := fmt.Sprintf("SELECT MENR, DEVL FROM statistics.state_statistics WHERE ShipInfo_ID='%d' ORDER BY datetime ASC", *shipInfoID)
		if *startTime != "" && *endTime != "" {
			sql = fmt.Sprintf("SELECT MENR, DEVL FROM statistics.state_statistics WHERE ShipInfo_ID='%d' AND datetime BETWEEN '%s' AND '%s' ORDER BY datetime ASC", *shipInfoID, *startTime, *endTime)
		}

		res, err := db.Query(sql)

		defer res.Close()

		if err != nil {
			log.Fatal(err)
		}

		menrList := [][]float64{}
		devlList := [][]float64{}
		stressNameSlice := []string{"DFP", "DFS", "SFP", "SFS", "DMP", "DMS", "SMP", "SMS", "DAP", "DAS", "SAP", "SAS", "L21", "L22", "L23", "L24", "L25", "L41", "L42", "L43", "L44", "L45", "L46", "L47", "L48", "L49", "L410", "L51", "L52", "L53", "L54", "L55", "L56", "L57", "AFx", "AFy", "AFz", "AAx", "AAy", "AAz"}
		numofMeasurePoint := 0

		for res.Next() {
			var state_statistics StateStatistics
			err := res.Scan(
				&state_statistics.MENR,
				&state_statistics.DEVL)

			if err != nil {
				log.Fatal(err)
			}

			var menr, devl []float64
			json.Unmarshal([]byte(state_statistics.MENR), &menr)
			json.Unmarshal([]byte(state_statistics.DEVL), &devl)
			menrList = append(menrList, menr)
			devlList = append(devlList, devl)
		}

		f, err := excelize.OpenFile("Stress_Acc_template.xlsx")
		if err != nil {
			log.Fatal(err)
		}

		numofMeasurePoint = len(menrList[0]) - 2
		// Export Excel
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

		// Draw Charts
		ptsSlice := []plotter.XYs{}
		for i := 0; i < numofMeasurePoint; i++ {
			n := len(menrList)
			pts := make(plotter.XYs, n)
			for j := 0; j < n; j++ {
				pts[j].X = (menrList[j][0] - 25569) * 86400
				pts[j].Y = menrList[j][i+2]
			}
			ptsSlice = append(ptsSlice, pts)
		}

		for i, pts := range ptsSlice {
			p := plot.New()
			p.Title.Text = stressNameSlice[i]
			p.Y.Label.Text = "Stress [MPa]"
			p.X.Tick.Marker = timeTicks{}
			p.Add(plotter.NewGrid())

			l, err := plotter.NewLine(pts)
			if err != nil {
				panic(err)
			}
			l.LineStyle.Width = vg.Points(1)
			l.LineStyle.Color = color.RGBA{R: 79, G: 129, B: 189, A: 255}

			// Add Error Bars
			yErr := make(plotter.Errors, len(pts))
			for j := range pts {
				yErr[j].Low = (-1) * devlList[j][i+2]
				yErr[j].High = devlList[j][i+2]
			}

			data := errPoints{
				XYs:     pts,
				YErrors: plotter.YErrors(yErr),
			}

			yerrs, err := plotter.NewYErrorBars(data)
			if err != nil {
				panic(err)
			}
			yerrs.LineStyle.Width = vg.Points(1)
			yerrs.Color = color.RGBA{R: 237, G: 242, B: 248, A: 255}

			p.Add(yerrs, l)

			outName := fmt.Sprintf("%s/images/graph-Stress-%d.png", *outDir, i+1)

			// Save the plot to a PNG file.
			if err := p.Save(10*vg.Inch, 5*vg.Inch, outName); err != nil {
				panic(err)
			}
		}

		fmt.Println("DONE Stress_Acc")
	}()

	go func() {
		defer wg.Done()
		sql := fmt.Sprintf("SELECT datetime, Roll_Max, Pitch_Max, Yaw_Max FROM statistics.gyro WHERE ShipInfo_ID='%d' ORDER BY datetime ASC", *shipInfoID)
		if *startTime != "" && *endTime != "" {
			sql = fmt.Sprintf("SELECT datetime, Roll_Max, Pitch_Max, Yaw_Max FROM statistics.gyro WHERE ShipInfo_ID='%d' AND datetime BETWEEN '%s' AND '%s' ORDER BY datetime ASC", *shipInfoID, *startTime, *endTime)
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
		gyroNameSlice := []string{"Roll Max", "Pitch Max", "Yaw Max"}

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

		// Draw Charts
		ptsSlice := []plotter.XYs{}
		n := len(rollList)
		ptsRoll := make(plotter.XYs, n)
		ptsPitch := make(plotter.XYs, n)
		ptsYaw := make(plotter.XYs, n)
		for j := 0; j < n; j++ {
			const layout = "2006-01-02 15:04:05"
			tm, _ := time.Parse(layout, dateList[j])
			t := float64(tm.Unix())
			ptsRoll[j].X = t
			ptsRoll[j].Y = rollList[j]

			ptsPitch[j].X = t
			ptsPitch[j].Y = pitchList[j]

			ptsYaw[j].X = t
			ptsYaw[j].Y = yawList[j]
		}
		ptsSlice = append(ptsSlice, ptsRoll, ptsPitch, ptsYaw)

		for i, pts := range ptsSlice {
			p := plot.New()

			p.Title.Text = gyroNameSlice[i]
			p.Y.Label.Text = "Angle [deg]"
			p.X.Tick.Marker = timeTicks{}
			p.Add(plotter.NewGrid())

			l, err := plotter.NewLine(pts)
			if err != nil {
				panic(err)
			}
			l.LineStyle.Width = vg.Points(1)
			l.LineStyle.Color = color.RGBA{R: 79, G: 129, B: 189, A: 255}

			p.Add(l)

			outName := fmt.Sprintf("%s/images/graph-Gyro-%d.png", *outDir, i+1)

			// Save the plot to a PNG file.
			if err := p.Save(10*vg.Inch, 5*vg.Inch, outName); err != nil {
				panic(err)
			}
		}

		fmt.Println("DONE Gyro")
	}()

	go func() {
		defer wg.Done()
		sql := fmt.Sprintf("SELECT datetime, SignificantWaveHeight_Arbitrary00, WavePeriod_Arbitrary00 FROM statistics.compweather WHERE ShipInfo_ID='%d' ORDER BY datetime ASC", *shipInfoID)
		if *startTime != "" && *endTime != "" {
			sql = fmt.Sprintf("SELECT datetime, SignificantWaveHeight_Arbitrary00, WavePeriod_Arbitrary00 FROM statistics.compweather WHERE ShipInfo_ID='%d' AND datetime BETWEEN '%s' AND '%s' ORDER BY datetime ASC", *shipInfoID, *startTime, *endTime)
		}

		res, err := db.Query(sql)

		defer res.Close()

		if err != nil {
			log.Fatal(err)
		}

		dateList2 := []string{}
		waveHList := []float64{}
		wavePList := []float64{}
		waveNameSlice := []string{"Wave Height", "Wave Period"}

		for res.Next() {
			var wave Wave
			err := res.Scan(
				&wave.datetime,
				&wave.SignificantWaveHeight_Arbitrary00,
				&wave.WavePeriod_Arbitrary00,
			)

			if err != nil {
				log.Fatal(err)
			}

			dateList2 = append(dateList2, wave.datetime)
			waveHList = append(waveHList, wave.SignificantWaveHeight_Arbitrary00)
			wavePList = append(wavePList, wave.WavePeriod_Arbitrary00)
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

		// Draw Charts
		ptsSlice := []plotter.XYs{}
		n := len(waveHList)
		ptsWaveH := make(plotter.XYs, n)
		ptsWaveP := make(plotter.XYs, n)
		for j := 0; j < n; j++ {
			const layout = "2006-01-02 15:04:05"
			tm, _ := time.Parse(layout, dateList2[j])
			t := float64(tm.Unix())
			ptsWaveH[j].X = t
			ptsWaveH[j].Y = waveHList[j]

			ptsWaveP[j].X = t
			ptsWaveP[j].Y = wavePList[j]
		}
		ptsSlice = append(ptsSlice, ptsWaveH, ptsWaveP)

		for i, pts := range ptsSlice {
			p := plot.New()

			p.Title.Text = waveNameSlice[i]
			if i == 0 {
				p.Y.Label.Text = "Significant Wave Height [m]"
			} else {
				p.Y.Label.Text = "Wave Period [sec]"
			}
			p.X.Tick.Marker = timeTicks{}
			p.Add(plotter.NewGrid())

			l, err := plotter.NewLine(pts)
			if err != nil {
				panic(err)
			}
			l.LineStyle.Width = vg.Points(1)
			l.LineStyle.Color = color.RGBA{R: 79, G: 129, B: 189, A: 255}

			p.Add(l)

			outName := fmt.Sprintf("%s/images/graph-Wave-%d.png", *outDir, i+1)

			// Save the plot to a PNG file.
			if err := p.Save(10*vg.Inch, 5*vg.Inch, outName); err != nil {
				panic(err)
			}
		}

		fmt.Println("DONE Wave")
	}()

	go func() {
		defer wg.Done()
		sql := fmt.Sprintf("SELECT datetime, GPGGA_LAT, GPGGA_LON FROM statistics.operation WHERE ShipInfo_ID='%d' ORDER BY datetime ASC", *shipInfoID)
		if *startTime != "" && *endTime != "" {
			sql = fmt.Sprintf("SELECT datetime, GPGGA_LAT, GPGGA_LON FROM statistics.operation WHERE ShipInfo_ID='%d' AND datetime BETWEEN '%s' AND '%s' ORDER BY datetime ASC", *shipInfoID, *startTime, *endTime)
		}

		res, err := db.Query(sql)

		defer res.Close()

		if err != nil {
			log.Fatal(err)
		}

		dateList3 := []string{}
		latList := []float64{}
		lonList := []float64{}

		for res.Next() {
			var mapL Map
			err := res.Scan(
				&mapL.datetime,
				&mapL.GPGGA_LAT,
				&mapL.GPGGA_LON,
			)

			if err != nil {
				log.Fatal(err)
			}

			dateList3 = append(dateList3, mapL.datetime)
			latList = append(latList, mapL.GPGGA_LAT)
			lonList = append(lonList, mapL.GPGGA_LON)
		}

		f, err := excelize.OpenFile("Map_template.xlsx")
		if err != nil {
			log.Fatal(err)
		}

		for i := 0; i < len(dateList3); i++ {
			cell1 := fmt.Sprintf("%v%v", alphaMap[1], i+2)
			cell2 := fmt.Sprintf("%v%v", alphaMap[2], i+2)
			cell3 := fmt.Sprintf("%v%v", alphaMap[3], i+2)
			cell4 := fmt.Sprintf("%v%v", alphaMap[4], i+2)
			cell5 := fmt.Sprintf("%v%v", alphaMap[5], i+2)
			cell6 := fmt.Sprintf("%v%v", alphaMap[6], i+2)
			f.SetCellValue("MapData", cell1, dateList3[i])
			f.SetCellValue("MapData", cell2, latList[i])
			f.SetCellValue("MapData", cell3, lonList[i])
			f.SetCellValue("MapData", cell4, 1300*lonList[i]/360)
			f.SetCellValue("MapData", cell5, 960*(latList[i]-19)/260)

			if i == 0 || i == len(dateList3)-1 {
				f.SetCellValue("MapData", cell6,
					fmt.Sprintf("Date: %v\nLat(deg): %v\nLon(deg): %v", dateList3[i], latList[i], lonList[i]))
			}
		}

		if err := f.SaveAs(fmt.Sprintf("%s/Map.xlsx", *outDir)); err != nil {
			log.Fatal(err)
		}

		fmt.Println("DONE Map")
	}()

	wg.Wait()
	elapsed := time.Since(start)
	log.Printf("Exporting excels took %s", elapsed)
}

func CmdExec(args ...string) (string, error) {

	baseCmd := args[0]
	cmdArgs := args[1:]

	fmt.Printf("Exec: %v", args)

	cmd := exec.Command(baseCmd, cmdArgs...)
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return string(out), nil
}
