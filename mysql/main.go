package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"

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

	menrList := [][]float64{}
	devlList := [][]float64{}
	i := 1

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

		numMeasurePoint := state_statistics.NumofMeasurePoint
		var m1 []float64
		json.Unmarshal([]byte(state_statistics.MeasurePointData), &m1)
		menrList = append(menrList, m1[0:(numMeasurePoint+2)])
		devlList = append(devlList, m1[(numMeasurePoint+2):(2*numMeasurePoint+4)])

		// fmt.Printf("%v %v %v %v\n", i,
		// 	state_statistics.ShipInfo_ID,
		// 	state_statistics.NumofMeasurePoint)
		i++
	}

	fmt.Println(menrList[0])
	fmt.Println(devlList[0])
}
