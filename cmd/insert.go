/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"phone_book_json/lib"
	"time"

	"github.com/spf13/cobra"
)

// insertCmd represents the insert command
var insertCmd = &cobra.Command{
	Use:   "insert",
	Short: "execute insert <first_name> <last_name> <phone_number>",
	Long:  ``,

	Args: cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		n, err := lib.FormatNumber(args[2])
		if err != nil {
			fmt.Println(err)
			return
		}
		c := http.Client{
			Timeout: 15 * time.Second,
		}
		request, err := http.NewRequest(http.MethodPost, fmt.Sprintf("http://localhost:1234/insert/%s/%s/%d", args[0], args[1], n), nil)
		if err != nil {
			fmt.Println("Get insert err:", err)
			return
		}

		httpData, err := c.Do(request)
		if err != nil {
			fmt.Println("Do() insert err:", err)
			return
		}
		_, err = io.Copy(os.Stdout, httpData.Body)
		fmt.Println("")
		if err != nil {
			fmt.Println("io.Copy insert err:", err)
		}

	},
}

func init() {
	rootCmd.AddCommand(insertCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// insertCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// insertCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
