/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/phone_book/internal/lib"
	api "github.com/phone_book/internal/lib/api/response"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "list records in store",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		// init client
		c := http.Client{
			Timeout: cfg.HTTPServer.Timeout * time.Second,
		}
		request, err := http.NewRequest(http.MethodGet, fmt.Sprintf("http://%s/list", cfg.HTTPServer.Address), nil)
		if err != nil {
			fmt.Println("Get list err:", err)
			return
		}

		resp, err := c.Do(request)
		if err != nil {
			fmt.Println("Do() list err:", err)
			return
		}
		defer resp.Body.Close()

		body, _ := io.ReadAll(resp.Body)

		var responseJson api.ResponseData

		err = json.Unmarshal(body, &responseJson)

		if err != nil {
			fmt.Println("json.Unmarshal() err:", err)
			return
		}
		if responseJson.Error != "" {
			fmt.Println("response err:", responseJson.Error)
			return
		}
		output, _ := lib.PrettyPrintJSONstream(responseJson.Data)
		fmt.Println(output)
	},
}

func init() {
	rootCmd.AddCommand(listCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
