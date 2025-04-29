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

	personHandlers "github.com/phone_book/internal/http_server/handlers/person"
	"github.com/spf13/cobra"
)

// statusCmd represents the status command
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "A brief description of your command",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		c := http.Client{
			Timeout: cfg.HTTPServer.Timeout * time.Second,
		}
		request, err := http.NewRequest(http.MethodGet, fmt.Sprintf("http://%s/status", cfg.HTTPServer.Address), nil)
		if err != nil {
			fmt.Println("Get status err:", err)
			return
		}

		resp, err := c.Do(request)
		if err != nil {
			fmt.Println("Do() status err:", err)
			return
		}
		defer resp.Body.Close()
		body, _ := io.ReadAll(resp.Body)

		var responseJson personHandlers.StatusResponse
		err = json.Unmarshal(body, &responseJson)

		if err != nil {
			fmt.Println("json.Unmarshal() err:", err)
			return
		}
		if responseJson.Error != "" {
			fmt.Println("response err:", responseJson.Error)
			return
		}

		fmt.Println("total: ", responseJson.Total)
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// statusCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// statusCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
