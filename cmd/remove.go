/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"phone_book_json/lib"

	"github.com/spf13/cobra"
)

// removeCmd represents the remove command
var removeCmd = &cobra.Command{
	Use:   "remove",
	Short: "",
	Long: ``,
	Args: cobra.ExactArgs(1),

	Run: func(cmd *cobra.Command, args []string) {
		db := getDB()
		number, err := lib.FormatNumber(args[0])
		if err != nil {
			fmt.Println(err)
		}
		removeErr := db.remove(number)
		if removeErr != nil {
			fmt.Println("removeErr:", removeErr)
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
