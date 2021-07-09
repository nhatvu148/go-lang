package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"
)

type MeasureType struct {
	Datetime string  `json:"datetime"`
	FromAp   float64 `json:"fromAp"`
	AftData  float64 `json:"aftData"`
	MidData  float64 `json:"midData"`
	ForeData float64 `json:"foreData"`
}
type MwDevType struct {
	X    float64
	Calc float64
}

func CsvToJson(outDir string, csvChan chan [][]string) {
	file, err := os.Open(fmt.Sprintf("%s/Output.csv", outDir))
	exitOnError(err)

	reader := csv.NewReader(file)
	records, _ := reader.ReadAll()

OuterLoop:
	for {
		select {
		case csvChan <- records:
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
			break OuterLoop

		default:
		}
	}
}

func GetFromAp(startTime string, fromApChan chan []float64, csvChan chan [][]string) {
	fromApArr := <-fromApChan
	records := <-csvChan

	numOfSec, err := strconv.Atoi(records[0][0])
	exitOnError(err)

	numOfData, err := strconv.Atoi(records[0][(numOfSec*2 + 1)])
	exitOnError(err)

	resIndex := []int{}
	mwDev := []MwDevType{}
	fromApSrc := records[0][(numOfSec + 1):(numOfSec*2 + 1)]
	arr := records[0][(numOfSec*2 + 2):]

	targetArr := []float64{}
	for j := 0; j < numOfSec; j++ {
		target, err := strconv.ParseFloat(fromApSrc[j], 64)
		exitOnError(err)
		targetArr = append(targetArr, target)
	}

	for _, faa := range fromApArr {
		_, index := findClosest(targetArr, faa)
		resIndex = append(resIndex, index)
	}

	const layout = "2006-01-02 15:04:05"
	tm, _ := time.Parse(layout, startTime)

	dateArr := []float64{}
	for i := 0; i < numOfData; i++ {
		const layout = "2006/01/02 15:04:05"
		ts, _ := time.Parse(layout, arr[(numOfSec+1)*i])
		dateArr = append(dateArr, float64(ts.Unix()))
	}

	_, tmIndex := findClosest(dateArr, float64(tm.Unix()))

	for i, ri := range resIndex {
		calc, err := strconv.ParseFloat(arr[(numOfSec+1)*tmIndex+ri], 64)
		exitOnError(err)

		data := MwDevType{X: fromApArr[i], Calc: calc}

		mwDev = append(mwDev, data)
	}

	// fmt.Println(resIndex)
	// fmt.Println(mwDev)
	fmt.Println("DONE mwDev")
	bytes, err := json.Marshal(mwDev)
	exitOnError(err)

	file, err := os.Create("mwDev.json")
	exitOnError(err)

	file.WriteString(string(bytes))
}

func exitOnError(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// Method to compare which one
// is the more close We find the
// closest by taking the difference
// between the target and both
// values. It assumes that val2 is
// greater than val1 and target
// lies between these two.
func getClosest(arr []float64, mid int,
	target float64) (float64, int) {
	if target-arr[mid] >= arr[mid+1]-target {
		return arr[mid+1], mid + 1
	} else {
		return arr[mid], mid
	}
}

// Returns element closest
// to target in arr[]
func findClosest(arr []float64,
	target float64) (float64, int) {
	n := len(arr)

	// Corner cases
	if target <= arr[0] {
		return arr[0], 0
	}
	if target >= arr[n-1] {
		return arr[n-1], n - 1
	}

	// Doing binary search
	i := 0
	j := n
	mid := 0
	for i < j {
		mid = (i + j) / 2

		if arr[mid] == target {
			return arr[mid], mid
		}

		/* If target is less
		   than array element,
		   then search in left */
		if target < arr[mid] {

			// If target is greater
			// than previous to mid,
			// return closest of two
			if mid > 0 && target > arr[mid-1] {
				return getClosest(arr,
					mid-1, target)
			}

			/* Repeat for left half */
			j = mid
		} else {
			// If target is
			// greater than mid
			if mid < n-1 && target < arr[mid+1] {
				return getClosest(arr,
					mid, target)
			}
			i = mid + 1 // update i
		}
	}

	// Only single element
	// left after search
	return arr[mid], mid
}
