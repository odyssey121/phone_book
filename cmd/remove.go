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

	"github.com/phone_book/internal/config"
	"github.com/phone_book/internal/lib"
	api "github.com/phone_book/internal/lib/api/response"

	"github.com/spf13/cobra"
)

// removeCmd represents the remove command
var removeCmd = &cobra.Command{
	Use:   "remove",
	Short: "execute remove <phone_number_for_remove>",
	Long:  ``,
	Args:  cobra.ExactArgs(1),

	Run: func(cmd *cobra.Command, args []string) {
		n, err := lib.FormatNumber(args[0])
		if err != nil {
			fmt.Println(err)
			return
		}
		// init conf
		cfg := config.MustLoad()
		// init client
		c := http.Client{
			Timeout: cfg.HTTPServer.Timeout * time.Second,
		}
		request, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("http://%s/remove/%d", cfg.HTTPServer.Address, n), nil)
		if err != nil {
			fmt.Println("Get remove err:", err)
			return
		}

		resp, err := c.Do(request)

		if err != nil {
			fmt.Println("Do() search err:", err)
			return
		}
		defer resp.Body.Close()

		body, _ := io.ReadAll(resp.Body)

		var responseJson api.Response
		err = json.Unmarshal(body, &responseJson)

		if err != nil {
			fmt.Println("json.Unmarshal() err:", err)
			return
		}

		if responseJson.Error != "" {
			fmt.Println("response err:", responseJson.Error)
			return
		}

		fmt.Println(responseJson.Message)
	},
}

func init() {
	rootCmd.AddCommand(removeCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// removeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// removeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
