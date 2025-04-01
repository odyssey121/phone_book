/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"phone_book_json/lib"
	"phone_book_json/store"
	"reflect"

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
		db := store.GetDB()

		var result any

		if cmd.Flags().Lookup("startWith").Changed {
			list := db.SearchStartWith(n)
			if len(list) != 0 {
				result = list
			}
		} else {
			result = db.Search(n)
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
