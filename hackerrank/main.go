package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

func main() {
	reader := bufio.NewReaderSize(os.Stdin, 16*1024*1024)

	// stdout, err := os.Create(os.Getenv("OUTPUT_PATH"))
	stdout, err := os.Create("hello.txt")
	checkError(err)

	defer stdout.Close()

	writer := bufio.NewWriterSize(stdout, 16*1024*1024)

	actionsCount, err := strconv.ParseInt(strings.TrimSpace(readLine(reader)), 10, 64)
	checkError(err)

	fmt.Println(actionsCount)

	// http.HandleFunc("/get", getHandler)
	// http.HandleFunc("/post", postHandler)
	// http.HandleFunc("/delete", deleteHandler)
	// go http.ListenAndServe(portSuffix, nil)
	// time.Sleep(100 * time.Millisecond)

	// var actions []string

	// for i := 0; i < int(actionsCount); i++ {
	// 	actionsItem := readLine(reader)
	// 	actions = append(actions, actionsItem)
	// }

	// for _, actionStr := range actions {
	// 	var action Action
	// 	err := json.Unmarshal([]byte(actionStr), &action)
	// 	checkError(err)
	// 	switch action.Type {
	// 	case "post":
	// 		_, err := http.Post(address+"/post", "application/json", strings.NewReader(action.Payload))
	// 		checkError(err)
	// 	case "delete":
	// 		client := &http.Client{}
	// 		req, err := http.NewRequest("DELETE", address+"/delete?id="+action.Payload, nil)
	// 		checkError(err)
	// 		resp, err := client.Do(req)
	// 		checkError(err)
	// 		if resp.StatusCode != 200 {
	// 			fmt.Fprintf(writer, "%s\n", resp.Status)
	// 			continue
	// 		}
	// 	case "get":
	// 		resp, err := http.Get(address + "/get?id=" + action.Payload)
	// 		checkError(err)
	// 		if resp.StatusCode != 200 {
	// 			fmt.Fprintf(writer, "%s\n", resp.Status)
	// 			continue
	// 		}
	// 		var lake Lake
	// 		err = json.NewDecoder(resp.Body).Decode(&lake)
	// 		checkError(err)
	// 		fmt.Fprintf(writer, "%s\n", lake.Name)
	// 		fmt.Fprintf(writer, "%d\n", lake.Area)
	// 	}
	// }

	fmt.Fprintf(writer, "\n")

	writer.Flush()
}

func readLine(reader *bufio.Reader) string {
	str, _, err := reader.ReadLine()
	if err == io.EOF {
		return ""
	}

	return strings.TrimRight(string(str), "\r\n")
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

// func ModuloFibonacciSequence(requestChan chan bool, resultChan chan int) {
//     var skip, total int
//     fmt.Scanf("%s\n%s", skip, total)

//     result := fibonacci(skip, total)

//     for _, res:= range result {
//         select {
//             case request := <- requestChan:
//                 if request {
//                     time.Sleep(10*time.Millisecond)
//                     resultChan <- res
//                 }
//         }
//     }
// }
