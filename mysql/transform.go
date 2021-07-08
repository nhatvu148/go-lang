package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
)

type MeasureType struct {
	Datetime string  `json:"datetime"`
	FromAp   float64 `json:"fromAp"`
	AftData  float64 `json:"aftData"`
	MidData  float64 `json:"midData"`
	ForeData float64 `json:"foreData"`
}

func Transform(outDir string) {
	file, err := os.Open(fmt.Sprintf("%s/Output.csv", outDir))
	exitOnError(err)

	reader := csv.NewReader(file)
	records, _ := reader.ReadAll()

	numOfSec, err := strconv.Atoi(records[0][0])
	exitOnError(err)

	numOfData, err := strconv.Atoi(records[0][(numOfSec*2 + 1)])
	exitOnError(err)

	result := new([]MeasureType)
	arr := records[0][(numOfSec*2 + 2):]

	for i := 0; i < numOfData; i++ {
		aftOrigin, err := strconv.ParseFloat(arr[(numOfSec+1)*i+27], 64)
		exitOnError(err)

		midOrigin, err := strconv.ParseFloat(arr[(numOfSec+1)*i+30], 64)
		exitOnError(err)

		foreOrigin, err := strconv.ParseFloat(arr[(numOfSec+1)*i+35], 64)
		exitOnError(err)

		data := MeasureType{Datetime: arr[(numOfSec+1)*i], AftData: aftOrigin, MidData: midOrigin, ForeData: foreOrigin}
		*result = append(*result, data)
	}

	bytes, err := json.Marshal(*result)
	exitOnError(err)

	file, err = os.Create("result.json")
	exitOnError(err)

	file.WriteString(string(bytes))
}

func exitOnError(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
