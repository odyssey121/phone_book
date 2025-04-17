/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/phone_book/lib"
	"github.com/phone_book/store"
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
			log.Println("listHandler DB Error: ", err)
			http.Error(w, "Server Error", http.StatusInternalServerError)
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

		insertEntity := store.Person{}
		err := lib.DeSerialize(&insertEntity, r.Body)

		if err != nil {
			http.Error(w, fmt.Sprint(err), http.StatusBadRequest)
			return
		}

		if p := db.Search(insertEntity.Phone); p != nil {
			http.Error(w, fmt.Sprintf("Person with number: %d exsist!", insertEntity.Phone), http.StatusBadRequest)
			return

		}

		err = db.Insert(insertEntity.FirstName, insertEntity.LastName, insertEntity.Phone)

		if err != nil {
			log.Println(err)
			http.Error(w, "Server Error", http.StatusInternalServerError)
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

		if p := db.Search(n); p == nil {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintf(w, "Person with number %d not found", n)
			return
		}

		err = db.Remove(n)

		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "Server Error", err)
			return
		}
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Records with number %d deleted", n)
	}
}

// Create a new ServeMux using Gorilla
var gMux = mux.NewRouter()

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "run server for phone book",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		PORT := ":" + cmd.Flag("port").Value.String()
		s := http.Server{
			Addr:         PORT,
			Handler:      gMux,
			IdleTimeout:  5 * time.Second,
			ReadTimeout:  time.Second,
			WriteTimeout: time.Second,
		}

		db := store.GetDB()

		getMux := gMux.Methods(http.MethodGet).Subrouter()

		getMux.HandleFunc("/", defaultHeandler)
		getMux.HandleFunc("/list", listHandler(db))
		getMux.HandleFunc("/status", statusHandler(db))
		getMux.HandleFunc("/search/{number}", searchHandler(db))

		postMux := gMux.Methods(http.MethodPost).Subrouter()

		postMux.HandleFunc("/insert", insertHandler(db))

		deleteMux := gMux.Methods(http.MethodDelete).Subrouter()

		deleteMux.HandleFunc("/remove/{number}", removeHandler(db))

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
