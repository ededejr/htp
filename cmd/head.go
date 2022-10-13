package cmd

import (
	"fmt"
	"net/http"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(headCmd)
}

var headCmd = &cobra.Command{
	Use:   "head [url]",
	Short: "Make a HEAD request",
	Long:  "Make a HEAD request to a provided url",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		resp, err := http.Head(args[0])
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()

		fmt.Printf("Response status: %s", resp.Status)
		headersString := ""
		fmt.Printf("\nResponse headers:")
		for k, v := range resp.Header {
			headersString += fmt.Sprintf("\n'%s': %s", k, v)
		}
		fmt.Println(headersString)
	},
}
