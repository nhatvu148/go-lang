package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

type StateStatistics struct {
	State_Measure_ID  int
	ShipInfo_ID       int
	datetime          string
	NumofProcess      int
	NumofMeasurePoint int
	MENR              string
	DEVL              string
}
type MsmwAllowable struct {
	datetime    string
	ShipInfo_ID int
	X_From_AP   float64
	MsMw_HOG    float64
	MsMw_SAG    float64
}
type CompWeather struct {
	datetime                          string
	SignificantWaveHeight_Arbitrary00 float64
}
type Operation struct {
	Operation_ID         int
	ShipInfo_ID          int
	datetime             string
	NumofRow             int
	GPZDA_DATETIME       string
	GPGNS_LAT            float64
	GPGNS_LON            float64
	GPGGA_LAT            float64
	GPGGA_LON            float64
	VDVBW_SOG            float64
	GPVTG_COG            float64
	WIMWV_REL_ANGLE      float64
	WIMWV_REL_SPEED      float64
	HEHDT_HEADING        float64
	HETHS_HEADING        float64
	VDVBW_SOW            float64
	AGRSA_ANGLE_1        float64
	AGRSA_ANGLE_2        float64
	RCRPM_SHAFT_REV_E_1  float64
	RCRPM_SHAFT_REV_E_2  float64
	RCRPM_PROP_PITCH_E_1 float64
	RCRPM_PROP_PITCH_E_2 float64
	VDVBW_1              float64
	VDVBW_6              float64
	VDVBW_4              float64
	VDVBW_8              float64
	SDDPT                float64
	WIMWV_TRUE_ANGLE     float64
	WIMWV_TRUE_SPEED     float64
	GPVTG_2              float64
	GPVTG_4              float64
	GPVTG_6              float64
}

func getJmuData(db *sql.DB, dbName string, dataName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		shipInfoID := params["shipInfoID"]
		startTime := params["startTime"]
		endTime := params["endTime"]
		var QUERY string
		if startTime != "" && endTime != "" {
			switch {
			case dataName == "state_statistics":
				QUERY =
					fmt.Sprintf(`SELECT State_Measure_ID,ShipInfo_ID,datetime,NumofProcess,NumofMeasurePoint,MENR,DEVL FROM %s.%s WHERE datetime BETWEEN '%s' AND '%s' AND ShipInfo_ID='%s' ORDER BY datetime ASC`, dbName, dataName, startTime, endTime, shipInfoID)

			case dataName == "msmwallowable":
				QUERY =
					fmt.Sprintf(`SELECT datetime,ShipInfo_ID,X_From_AP,MsMw_HOG,MsMw_SAG FROM %s.%s WHERE ShipInfo_ID='%s' ORDER BY datetime ASC`, dbName, dataName, shipInfoID)

			case dataName == "compweather":
				QUERY =
					fmt.Sprintf(`SELECT datetime,SignificantWaveHeight_Arbitrary00 FROM %s.%s WHERE datetime BETWEEN '%s' AND '%s' AND ShipInfo_ID='%s' ORDER BY datetime ASC`, dbName, dataName, startTime, endTime, shipInfoID)

			case dataName == "operation":
				QUERY =
					fmt.Sprintf(`SELECT Operation_ID,ShipInfo_ID,datetime,NumofRow,GPZDA_DATETIME,GPGNS_LAT,GPGNS_LON,GPGGA_LAT,GPGGA_LON,VDVBW_SOG,GPVTG_COG,WIMWV_REL_ANGLE,WIMWV_REL_SPEED,HEHDT_HEADING,HETHS_HEADING,VDVBW_SOW,AGRSA_ANGLE_1,AGRSA_ANGLE_2,RCRPM_SHAFT_REV_E_1,RCRPM_SHAFT_REV_E_2,RCRPM_PROP_PITCH_E_1,RCRPM_PROP_PITCH_E_2,VDVBW_1,VDVBW_6,VDVBW_4,VDVBW_8,SDDPT,WIMWV_TRUE_ANGLE,WIMWV_TRUE_SPEED,GPVTG_2,GPVTG_4,GPVTG_6 FROM %s.%s WHERE datetime BETWEEN '%s' AND '%s' AND ShipInfo_ID='%s' ORDER BY datetime ASC`, dbName, dataName, startTime, endTime, shipInfoID)

			case shipInfoID != "":
				QUERY =
					fmt.Sprintf(`SELECT * FROM %s.%s WHERE datetime BETWEEN '%s' AND '%s' AND ShipInfo_ID='%s' ORDER BY datetime ASC`, dbName, dataName, startTime, endTime, shipInfoID)

			default:
				QUERY = fmt.Sprintf(`SELECT * FROM %s.%s WHERE datetime BETWEEN '%s' AND '%s' ORDER BY datetime ASC`, dbName, dataName, startTime, endTime)
			}
		} else {
			if dataName == "shipmeasurepoint" {
				QUERY = fmt.Sprintf(`SELECT * FROM %s.%s WHERE ShipInfo_ID='%s'`, dbName, dataName, shipInfoID)
			} else {
				QUERY = fmt.Sprintf(`SELECT * FROM %s.%s`, dbName, dataName)
			}
		}

		res, err := db.Query(QUERY)

		defer res.Close()

		if err != nil {
			log.Fatal(err)
		}

		switch {
		case dataName == "compweather":
			dateList2 := []string{}
			waveHList := []float64{}
			// wavePList := []float64{}
			// waveNameSlice := []string{"Wave Height", "Wave Period"}

			for res.Next() {
				var stateStatistics StateStatistics
				var compWeather CompWeather
				var msmwallowable MsmwAllowable
				var operation Operation
				err := res.Scan(
					&stateStatistics.State_Measure_ID,
					&stateStatistics.ShipInfo_ID,
					&stateStatistics.datetime,
					&stateStatistics.NumofProcess,
					&stateStatistics.NumofMeasurePoint,
					&stateStatistics.MENR,
					&stateStatistics.DEVL,

					&msmwallowable.datetime,
					&msmwallowable.ShipInfo_ID,
					&msmwallowable.X_From_AP,
					&msmwallowable.MsMw_HOG,
					&msmwallowable.MsMw_SAG,

					&compWeather.datetime,
					&compWeather.SignificantWaveHeight_Arbitrary00,
					// &compWeather.WavePeriod_Arbitrary00,

					&operation.Operation_ID,
					&operation.ShipInfo_ID,
					&operation.datetime,
					&operation.NumofRow,
					&operation.GPZDA_DATETIME,
					&operation.GPGNS_LAT,
					&operation.GPGNS_LON,
					&operation.GPGGA_LAT,
					&operation.GPGGA_LON,
					&operation.VDVBW_SOG,
					&operation.GPVTG_COG,
					&operation.WIMWV_REL_ANGLE,
					&operation.WIMWV_REL_SPEED,
					&operation.HEHDT_HEADING,
					&operation.HETHS_HEADING,
					&operation.VDVBW_SOW,
					&operation.AGRSA_ANGLE_1,
					&operation.AGRSA_ANGLE_2,
					&operation.RCRPM_SHAFT_REV_E_1,
					&operation.RCRPM_SHAFT_REV_E_2,
					&operation.RCRPM_PROP_PITCH_E_1,
					&operation.RCRPM_PROP_PITCH_E_2,
					&operation.VDVBW_1,
					&operation.VDVBW_6,
					&operation.VDVBW_4,
					&operation.VDVBW_8,
					&operation.SDDPT,
					&operation.WIMWV_TRUE_ANGLE,
					&operation.WIMWV_TRUE_SPEED,
					&operation.GPVTG_2,
					&operation.GPVTG_4,
					&operation.GPVTG_6,
				)

				if err != nil {
					log.Fatal(err)
				}

				dateList2 = append(dateList2, compWeather.datetime)
				waveHList = append(waveHList, compWeather.SignificantWaveHeight_Arbitrary00)
				// wavePList = append(wavePList, compWeather.WavePeriod_Arbitrary00)
			}

			fmt.Println(dateList2)
			fmt.Println(waveHList)
		}
	}
}
