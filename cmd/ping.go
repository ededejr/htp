package cmd

import (
	"fmt"

	"github.com/ededejr/htp/lib"
	"github.com/spf13/cobra"
)

var pingCmdFlags struct {
	// the duration in seconds to ping for
	Duration int
	Interval int
}

func init() {
	pingCmd.Flags().IntVarP(&pingCmdFlags.Duration, "duration", "d", 10, "Duration in seconds to ping for")
	pingCmd.Flags().IntVarP(&pingCmdFlags.Interval, "interval", "i", 200, "Interval for pings in milliseconds")
	rootCmd.AddCommand(pingCmd)
}

var pingCmd = &cobra.Command{
	Use:   "ping [url]",
	Short: "Repeatedly make GET requests",
	Long:  "Repeatedly make GET requests to a given url",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		id := lib.SetInterval(func() {
			m, err := makeMeasuredGetRequest(args[0])
			if err != nil {
				panic(err)
			}
			fmt.Printf("%s - %s\n", m.Res.Status, m.FirstByte)
		}, pingCmdFlags.Interval)
		lib.Sleep(pingCmdFlags.Duration)
		lib.ClearInterval(id)
	},
}
