package cmd

import (
	"fmt"
	"math"
	"sync"
	"time"

	"github.com/spf13/cobra"
)

var stressCmdFlags struct {
	Limit    int
	Duration int
	Verbose  bool
	Workers  int
}

func init() {
	stressCmd.Flags().IntVarP(&stressCmdFlags.Limit, "limit", "l", int(math.Inf(1)), "Limit the number of requests to make")
	stressCmd.Flags().IntVarP(&stressCmdFlags.Duration, "duration", "d", 10, "Duration of the test in seconds")
	stressCmd.Flags().IntVarP(&stressCmdFlags.Workers, "workers", "w", 5, "Number of workers to use")
	stressCmd.Flags().BoolVarP(&stressCmdFlags.Verbose, "verbose", "v", false, "Print more verbose logs")
	rootCmd.AddCommand(stressCmd)
}

var stressCmd = &cobra.Command{
	Use:   "stress [url]",
	Short: "Stress test an endpoint with GET requests",
	Long:  "Stress test an endpoint with GET requests",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		start := time.Now()
		results := make(map[int]*measuredResponse)
		mut := sync.Mutex{}

		done := make(chan bool)
		canExit := make(chan bool)

		finish := func() {
			done <- true
			avgResponseTime := time.Duration(0)
			slowestRequest := time.Duration(0)
			fastestRequest := time.Duration(stressCmdFlags.Duration) * time.Second
			totalRequests := 0
			end := time.Now()
			totalTime := end.Sub(start)

			for _, v := range results {
				if v.FirstByte > slowestRequest {
					slowestRequest = v.FirstByte
				}

				if v.FirstByte < fastestRequest {
					fastestRequest = v.FirstByte
				}

				avgResponseTime += v.FirstByte

				totalRequests += 1
			}

			avgResponseTime /= time.Duration(float64(totalRequests))

			fmt.Printf("\nSummary:\n")
			fmt.Println("--------------------")
			fmt.Printf("Total requests: %d\n", totalRequests)
			fmt.Printf("Total time: %s\n", totalTime)
			fmt.Println("--------------------")
			fmt.Printf("Average response time: %s\n", avgResponseTime)
			fmt.Printf("Fastest request: %s\n", fastestRequest)
			fmt.Printf("Slowest request: %s\n", slowestRequest)

			canExit <- true
		}

		time.AfterFunc(time.Duration(stressCmdFlags.Duration)*time.Second, func() {
			finish()
		})

		makeRequest := func(wid int) {
			m, err := makeMeasuredGetRequest(args[0])
			if err != nil {
				panic(err)
			}

			mut.Lock()
			results[stressCmdFlags.Limit] = m
			stressCmdFlags.Limit--
			mut.Unlock()

			if stressCmdFlags.Verbose {
				fmt.Printf("[w%d|%s] %s - %s\n", wid, time.Now().Format(time.Stamp), m.Res.Status, m.FirstByte.Round(time.Microsecond))
			}
		}

		for id := 0; id < stressCmdFlags.Workers; id++ {

			worker := func(wid int) {
				for {
					select {
					case <-done:
						return
					default:
						makeRequest(wid + 1)
					}

					if !(stressCmdFlags.Limit > 0) {
						finish()
						return
					}
				}
			}

			go worker(id)
		}

		// requests are complete
		<-done
		// program is ready to exit
		<-canExit
	},
}
