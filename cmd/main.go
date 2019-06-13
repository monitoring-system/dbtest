package main

import (
	"encoding/json"
	"fmt"
	tm "github.com/buger/goterm"
	"github.com/monitoring-system/dbtest/api/types"
	"io/ioutil"
	"net/http"
	"time"
)

func main() {
	tm.Clear() // Clear current screen

	for {
		// By moving cursor to top-left position we ensure that console output
		// will be overwritten each time, instead of adding new.
		tm.MoveCursor(1, 1)

		results := getTestResult()
		totals := tm.NewTable(0, 10, 5, ' ', 0)
		fmt.Fprintf(totals, "TestID\tTestName\tLoop\tStatus\n")

		for _, result := range results {
			var color string
			switch result.Status {
			case types.TestStatusRunning:
				color = tm.Color(result.Status, tm.BLUE)
			case types.TestStatusDone:
				color = tm.Color(result.Status, tm.GREEN)
			case types.TestStatusPending:
				color = tm.Color(result.Status, tm.YELLOW)
			case types.TestStatusFail:
				color = tm.Color(result.Status, tm.RED)
			}
			fmt.Fprintf(totals, "%d\t%s\t%d\t%s\n", result.TestID, result.Name, result.Loop, color)
		}
		tm.Println(totals)

		//tm.Println(tm.Background(tm.Color(tm.Bold("Important header"), tm.RED), tm.WHITE))

		tm.Flush() // Call it every time at the end of rendering
		time.Sleep(5 * time.Second)
	}
}

func getTestResult() []*types.TestResult {
	resp, err := http.Get("http://localhost:8080/results")
	if err != nil || resp == nil {
		return nil
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// handle error
		return nil
	}
	var data []*types.TestResult
	json.Unmarshal([]byte(body), &data)
	return data
}
