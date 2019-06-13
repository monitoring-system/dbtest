package main

import (
	"encoding/json"
	"fmt"
	tm "github.com/buger/goterm"
	"github.com/monitoring-system/dbtest/api/types"
	"io/ioutil"
	"net/http"
	"time"
	//"github.com/spf13/cobra"
)

func main() {
	//var rootCmd = &cobra.Command{
	//	Use:   "hugo",
	//	Short: "Hugo is a very fast static site generator",
	//	Long: `A Fast and Flexible Static Site Generator built with
	//            love by spf13 and friends in Go.
	//            Complete documentation is available at http://hugo.spf13.com`,
	//	Run: func(cmd *cobra.Command, args []string) {
	//		// Do Stuff Here
	//	},
	//}

	tm.Clear() // Clear current screen

	for {
		// By moving cursor to top-left position we ensure that console output
		// will be overwritten each time, instead of adding new.
		tm.MoveCursor(1, 1)
		tm.Clear()

		results := getTestResult()
		totals := tm.NewTable(0, 4, 4, ' ', 0)
		fmt.Fprintf(totals, "TestID\tTestName\tLoop\tStatus\tFailed\tDuration\n")

		for _, result := range results {
			var color string
			var duration int64 = 0
			switch result.Status {
			case types.TestStatusRunning:
				color = tm.Color(result.Status, tm.BLUE)
				duration = time.Now().Unix() - result.Start
			case types.TestStatusDone:
				if result.FailedLoopCount > 0 {
					color = tm.Color(result.Status, tm.RED)
				} else {
					color = tm.Color(result.Status, tm.GREEN)
				}
				duration = result.End - result.Start
			case types.TestStatusPending:
				color = tm.Color(result.Status, tm.YELLOW)
			}
			fmt.Fprintf(totals, "%d\t%s\t%d\t%s\t\t%d\t%d\n", result.TestID, result.Name, result.Loop, color, result.FailedLoopCount, duration)
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
