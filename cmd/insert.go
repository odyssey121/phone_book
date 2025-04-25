/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/phone_book/internal/lib"
	api "github.com/phone_book/internal/lib/api/response"
	"github.com/phone_book/internal/store"
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

		// init client
		c := http.Client{
			Timeout: cfg.HTTPServer.Timeout * time.Second,
		}

		payoadEntity := store.Person{FirstName: args[0], LastName: args[1], Phone: n, LastAccess: "0"}

		buf := new(bytes.Buffer)

		lib.Serialize(payoadEntity, buf)

		request, err := http.NewRequest(http.MethodPost, fmt.Sprintf("http://%s/insert", cfg.HTTPServer.Address), buf)
		if err != nil {
			fmt.Println("Get insert err:", err)
			return
		}

		resp, err := c.Do(request)
		if err != nil {
			fmt.Println("Do() insert err:", err)
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
		output, _ := lib.PrettyPrintJSONstream(responseJson.Message)
		fmt.Println(output)

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
