/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"phone_book_json/lib"
	"phone_book_json/store"
	"reflect"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

type HandlerFunc func(w http.ResponseWriter, r *http.Request)

func defaultHeandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Serving:", r.URL.Path, "from", r.Host)
	w.WriteHeader(http.StatusOK)
	Body := "Server running!\n"
	fmt.Fprintf(w, "%s", Body)
}

func listHandler(db store.DB) HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("Serving:", r.URL.Path, "from", r.Host)

		data, err := db.List()
		if err != nil {
			log.Println("listHandler Error: ", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return

		}
		w.WriteHeader(http.StatusOK)
		lib.Serialize(data, w)
	}
}

func statusHandler(db store.DB) HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("Serving:", r.URL.Path, "from", r.Host)

		res := fmt.Sprintf("Total record: %d\n", db.CountRecords())
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, res)
	}
}

func insertHandler(db store.DB) HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("Serving:", r.URL.Path, "from", r.Host)
		urlSplited := strings.Split(r.URL.Path, "/")
		if len(urlSplited) < 5 {
			w.WriteHeader(http.StatusNotFound)
			http.Error(w, "Not enough arguments: "+r.URL.Path, http.StatusNotFound)
			return
		}
		n, err := lib.FormatNumber(urlSplited[4])
		if err != nil {
			http.Error(w, fmt.Sprint(err), http.StatusBadRequest)
			return
		}

		err = db.Insert(urlSplited[2], urlSplited[3], n)

		if err != nil {
			http.Error(w, fmt.Sprint(err), http.StatusOK)
			return
		}
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "New record added successfully")
	}
}

func searchHandler(db store.DB) HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("Serving:", r.URL.Path, "from", r.Host)
		parsedUrl, _ := url.Parse(r.URL.String())
		params, _ := url.ParseQuery(parsedUrl.RawQuery)

		urlSplited := strings.Split(parsedUrl.Path, "/")
		if len(urlSplited) < 3 {
			w.WriteHeader(http.StatusNotFound)
			http.Error(w, "Not enough arguments: "+r.URL.Path, http.StatusNotFound)
			return
		}
		n, err := lib.FormatNumber(urlSplited[2])
		if err != nil {
			http.Error(w, fmt.Sprint(err), http.StatusBadRequest)
			return
		}

		if pv := params.Get("start_with"); pv == "1" {
			if ps := db.SearchStartWith(n); !reflect.ValueOf(ps).IsNil() {
				w.WriteHeader(http.StatusOK)
				lib.Serialize(ps, w)
			} else {
				w.WriteHeader(http.StatusNotFound)
				fmt.Fprintf(w, "Person with number %d not found", n)
			}

		} else if p := db.Search(n); p != nil {
			w.WriteHeader(http.StatusOK)
			lib.Serialize(p, w)
		} else {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintf(w, "Person with number %d not found", n)
		}
	}
}

func removeHandler(db store.DB) HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("Serving:", r.URL.Path, "from", r.Host)
		urlSplited := strings.Split(r.URL.Path, "/")
		if len(urlSplited) < 3 {
			w.WriteHeader(http.StatusNotFound)
			http.Error(w, "Not enough arguments: "+r.URL.Path, http.StatusNotFound)
			return
		}
		n, err := lib.FormatNumber(urlSplited[2])
		if err != nil {
			http.Error(w, fmt.Sprint(err), http.StatusBadRequest)
			return
		}
		log.Println("N:", n)
		err = db.Remove(n)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintf(w, "%s", err)
			return
		}
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Records with number %d deleted", n)
	}
}

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "run server for phone book",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		PORT := ":" + cmd.Flag("port").Value.String()
		mux := http.DefaultServeMux
		s := http.Server{
			Addr:         PORT,
			Handler:      mux,
			IdleTimeout:  5 * time.Second,
			ReadTimeout:  time.Second,
			WriteTimeout: time.Second,
		}

		db := store.GetDB()
		mux.Handle("/", http.HandlerFunc(defaultHeandler))
		mux.Handle("GET /list", http.HandlerFunc(listHandler(db)))
		mux.Handle("GET /status", http.HandlerFunc(statusHandler(db)))
		mux.Handle("POST /insert/", http.HandlerFunc(insertHandler(db)))
		mux.Handle("GET /search/", http.HandlerFunc(searchHandler(db)))
		mux.Handle("DELETE /remove/", http.HandlerFunc(removeHandler(db)))
		fmt.Println("Ready to serve at", PORT)
		err := s.ListenAndServe()
		if err != nil {
			fmt.Println(err)
			return
		}

	},
}

func init() {
	rootCmd.AddCommand(serverCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serverCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serverCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	serverCmd.Flags().StringP("port", "p", "1234", "the port on which the server will be launched")
}
