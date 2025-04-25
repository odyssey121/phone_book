/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strings"
	"time"

	"github.com/phone_book/internal/lib"
	api "github.com/phone_book/internal/lib/api/response"

	"github.com/spf13/cobra"
)

// searchCmd represents the search command
var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "execute search <phone_number_for_search>",
	Long:  ``,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		n, err := lib.FormatNumber(args[0])
		if err != nil {
			fmt.Println(err)
			return
		}

		var paramStr string
		paramSlice := make([]string, 0)

		if cmd.Flags().Lookup("startWith").Changed {
			paramSlice = append(paramSlice, "?start_with=1")
		}

		c := http.Client{
			Timeout: cfg.HTTPServer.Timeout * time.Second,
		}

		if len(paramSlice) > 0 {
			paramStr = strings.Join(paramSlice, "&")
		}

		request, err := http.NewRequest(http.MethodGet, fmt.Sprintf("http://%s/search/%d%s", cfg.HTTPServer.Address, n, paramStr), nil)
		if err != nil {
			fmt.Println("Get search err:", err)
			return
		}

		resp, err := c.Do(request)
		if err != nil {
			fmt.Println("Do() search err:", err)
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

		var output string
		if responseJson.Data != nil && !reflect.ValueOf(responseJson.Data).IsNil() {
			output, _ = lib.PrettyPrintJSONstream(responseJson.Data)
		} else {
			output = "not found"
		}

		fmt.Println(output)
	},
}

func init() {
	rootCmd.AddCommand(searchCmd)
	searchCmd.Flags().Bool("startWith", false, "Поиск возвращает записи по первому вхождению")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// searchCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// searchCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
