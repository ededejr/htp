package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(postCmd)
}

var postCmd = &cobra.Command{
	Use:   "post [url] [data]",
	Short: "Make a POST request",
	Long:  "Make a POST request to a provided url. Currently only supports JSON format.",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		url := args[0]
		data := args[1]

		// ensure data is valid JSON
		valid := json.Valid([]byte(data))

		if !valid {
			fmt.Println("string provided for \"data\" is Invalid JSON")
			return
		}

		payload := strings.NewReader(data)

		resp, err := http.Post(url, "application/json", payload)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()

		fmt.Println("Response status:", resp.Status)
		scanner := bufio.NewScanner(resp.Body)
		for scanner.Scan() {
			fmt.Println(scanner.Text())
		}
		if err := scanner.Err(); err != nil {
			panic(err)
		}
	},
}
