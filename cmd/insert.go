/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"phone_book_json/lib"

	"github.com/spf13/cobra"
)

// insertCmd represents the insert command
var insertCmd = &cobra.Command{
	Use:   "insert",
	Short: "",
	Long: ``,

	Args: cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		db := getDB()
		n, err := lib.FormatNumber(args[2])
		if err != nil {
			fmt.Println(err)
		}

		if p := db.search(n); p != nil {
			fmt.Printf("\nPerson with number: %v already exsist!\n", n)
		}

		insertErr := db.insert(args[0], args[1], n)

		if insertErr != nil {
			fmt.Println("insertErr:", insertErr)
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
