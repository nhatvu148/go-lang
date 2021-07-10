package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

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
	}
}
