/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

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

		httpData, err := c.Do(request)
		if err != nil {
			fmt.Println("Do() list err:", err)
			return
		}
		_, err = io.Copy(os.Stdout, httpData.Body)
		fmt.Println("")
		if err != nil {
			fmt.Println("io.Copy list err:", err)
		}
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
