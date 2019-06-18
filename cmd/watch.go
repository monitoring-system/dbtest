package cmd

import (
	"encoding/json"
	"fmt"
	tm "github.com/buger/goterm"
	"github.com/monitoring-system/dbtest/api/types"
	"github.com/spf13/cobra"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

var WatchTestCmd = &cobra.Command{
	Use:   "watch [watch tests]",
	Short: "watch test",
	Long:  "watch test",
	Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		watchTest()
	},
}

var server = "http://localhost:8080"
var failOnly = false

var query = "{}"

func init() {
	WatchTestCmd.Flags().StringVar(&server, "server", "http://localhost:8080", "db test server address")
	WatchTestCmd.Flags().BoolVar(&failOnly, "fail-only", false, "only list failed tests")
	WatchTestCmd.Flags().StringVar(&query, "query", "{}", "API raw rql query")
}

func watchTest() {
	for {
		// By moving cursor to top-left position we ensure that console output
		// will be overwritten each time, instead of adding new.
		tm.MoveCursor(1, 1)
		tm.Clear()

		results := getTestResult()
		totals := tm.NewTable(0, 4, 4, ' ', 0)
		fmt.Fprintf(totals, "TestID\tTestName\tLoop\tStatus\tFailed\tDuration\n")

		for _, result := range results {
			if failOnly && result.FailedLoopCount > 0 {
				continue
			}
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
	resp, err := http.Post(fmt.Sprintf("%s/results", server), "application/json", strings.NewReader(query))
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
