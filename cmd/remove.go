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

	"github.com/phone_book/internal/config"
	"github.com/phone_book/internal/lib"

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

		httpData, err := c.Do(request)
		if err != nil {
			fmt.Println("Do() remove err:", err)
			return
		}
		_, err = io.Copy(os.Stdout, httpData.Body)
		fmt.Println("")
		if err != nil {
			fmt.Println("io.Copy remove err:", err)
		}
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
