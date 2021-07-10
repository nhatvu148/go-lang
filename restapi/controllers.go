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
			switch dataName {
			case "state_statistics":
				QUERY =
					fmt.Sprintf(`SELECT State_Measure_ID,ShipInfo_ID,datetime,NumofProcess,NumofMeasurePoint,MENR,DEVL FROM %s.%s WHERE datetime BETWEEN '%s' AND '%s' AND ShipInfo_ID='%s' ORDER BY datetime ASC`, dbName, dataName, startTime, endTime, shipInfoID)
			default:
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
