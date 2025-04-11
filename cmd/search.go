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
	"reflect"
	"strings"
	"time"

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

		var result any

		var paramStr string
		paramSlice := make([]string, 0)

		if cmd.Flags().Lookup("startWith").Changed {
			paramSlice = append(paramSlice, "?start_with=1")
		}

		c := http.Client{
			Timeout: 15 * time.Second,
		}

		if len(paramSlice) > 0 {
			paramStr = strings.Join(paramSlice, "&")
		}

		request, err := http.NewRequest(http.MethodGet, fmt.Sprintf("http://localhost:1234/search/%d%s", n, paramStr), nil)
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

		if result != nil && !reflect.ValueOf(result).IsNil() {
			res, _ := lib.PrettyPrintJSONstream(result)
			fmt.Println(res)

		}
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
