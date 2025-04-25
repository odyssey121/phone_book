/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"time"

	"github.com/go-chi/render"
	"github.com/gorilla/mux"
	"github.com/phone_book/internal/config"
	"github.com/phone_book/internal/lib"
	apiResp "github.com/phone_book/internal/lib/api/response"
	"github.com/phone_book/internal/loger"
	"github.com/phone_book/internal/store"
	"github.com/spf13/cobra"
)

type HandlerFunc func(w http.ResponseWriter, r *http.Request)

func defaultHeandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	Body := "Server running!\n"
	fmt.Fprintf(w, "%s", Body)
}

func listHandler(log *slog.Logger, db store.DB) HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "server.listHandler"
		log := log.With(slog.String("op", op))
		log.Info("new request", slog.String("Serving", r.URL.Path), slog.String("from", r.Host))

		data, err := db.List()
		if err != nil {
			log.Error(fmt.Sprintf("%s", err))
			render.JSON(w, r, apiResp.Error("internal server error"))
			return

		}
		render.JSON(w, r, apiResp.WithData(data))
	}
}

type StatusResponse struct {
	apiResp.Response
	Total int `json:"total"`
}

func statusHandler(log *slog.Logger, db store.DB) HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "server.statusHandler"
		log := log.With(slog.String("op", op))
		log.Info("new request", slog.String("Serving", r.URL.Path), slog.String("from", r.Host))

		total := db.CountRecords()
		render.JSON(w, r, StatusResponse{apiResp.OK(""), total})

	}
}

func insertHandler(log *slog.Logger, db store.DB) HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "server.insertHandler"
		log := log.With(slog.String("op", op))
		log.Info("new request", slog.String("Serving", r.URL.Path), slog.String("from", r.Host))

		insertEntity := store.Person{}
		err := lib.DeSerialize(&insertEntity, r.Body)

		if err != nil {
			log.Error("de serialize error", loger.ErrLogFmt(err))
			render.JSON(w, r, apiResp.Error("de serialize error"))
			return
		}

		err = db.Insert(insertEntity.FirstName, insertEntity.LastName, insertEntity.Phone)

		if errors.Is(err, store.ErrPhoneExist) {
			render.JSON(w, r, apiResp.Error(fmt.Sprintf("person with number: %d already exsist", insertEntity.Phone)))
			return
		} else if err != nil {
			log.Error("db error", loger.ErrLogFmt(err))
			render.JSON(w, r, apiResp.Error("internal server error"))
			return
		}
		render.JSON(w, r, apiResp.OK("new record added successfully"))
	}
}

func searchHandler(log *slog.Logger, db store.DB) HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "server.searchHandler"
		log := log.With(slog.String("op", op))

		parsedUrl, _ := url.Parse(r.URL.String())
		params, _ := url.ParseQuery(parsedUrl.RawQuery)

		log.Info("new request", slog.String("Serving", r.URL.Path), slog.Any("params", params), slog.String("from", r.Host))

		vars := mux.Vars(r)
		phone, ok := vars["number"]

		if !ok {
			render.JSON(w, r, apiResp.Error("not enough arguments: "+r.URL.Path))
			return
		}
		n, err := lib.FormatNumber(phone)
		if err != nil {
			render.JSON(w, r, apiResp.Error(fmt.Sprint(err)))
			return
		}

		if pv := params.Get("start_with"); pv == "1" {
			ps, _ := db.SearchStartWith(n)
			render.JSON(w, r, apiResp.WithData(ps))

		} else {
			p := db.Search(n)
			render.JSON(w, r, apiResp.WithData(p))
		}
	}
}

func removeHandler(log *slog.Logger, db store.DB) HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "server.removeHandler"
		log := log.With(slog.String("op", op))
		log.Info("new request", slog.String("Serving", r.URL.Path), slog.String("from", r.Host))

		vars := mux.Vars(r)
		phone, ok := vars["number"]
		if !ok {
			render.JSON(w, r, apiResp.Error("not enough arguments: "+r.URL.Path))
			return
		}
		n, err := lib.FormatNumber(phone)
		if err != nil {
			render.JSON(w, r, apiResp.Error(fmt.Sprint(err)))
			return
		}

		if p := db.Search(n); p == nil {
			render.JSON(w, r, apiResp.Error(fmt.Sprintf("Person with number %d not found", n)))
			return
		}

		err = db.Remove(n)

		if err != nil {
			log.Error("db error", loger.ErrLogFmt(err))
			render.JSON(w, r, apiResp.Error("internal server error"))
			return
		}
		render.JSON(w, r, apiResp.OK(fmt.Sprintf("records with number %d deleted", n)))

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
		// init conf
		cfg := config.MustLoad()
		// init loger
		log := loger.GetLoger(cfg)
		log = log.With(slog.String("env", cfg.Env))
		log.Info("starting loger")
		// init db
		db, err := store.GetDB(cfg.Storage)
		if err != nil {
			log.Error("init db error", loger.ErrLogFmt(err))
		}
		// init server
		s := http.Server{
			Addr:         cfg.HTTPServer.Address,
			Handler:      gMux,
			IdleTimeout:  cfg.HTTPServer.IdleTimeout * time.Second,
			ReadTimeout:  cfg.HTTPServer.Timeout * time.Second,
			WriteTimeout: cfg.HTTPServer.Timeout * time.Second,
		}

		// init router
		getMux := gMux.Methods(http.MethodGet).Subrouter()

		getMux.HandleFunc("/", defaultHeandler)
		getMux.HandleFunc("/list", listHandler(log, db))
		getMux.HandleFunc("/status", statusHandler(log, db))
		getMux.HandleFunc("/search", searchHandler(log, db))
		getMux.HandleFunc("/search/{number}", searchHandler(log, db))

		postMux := gMux.Methods(http.MethodPost).Subrouter()

		postMux.HandleFunc("/insert", insertHandler(log, db))

		deleteMux := gMux.Methods(http.MethodDelete).Subrouter()
		deleteMux.HandleFunc("/remove", removeHandler(log, db))
		deleteMux.HandleFunc("/remove/{number}", removeHandler(log, db))

		log.Info("server started", slog.String("address", cfg.HTTPServer.Address))

		err = s.ListenAndServe()
		if err != nil {
			log.Error("http server listen err", loger.ErrLogFmt(err))
			return
		}

	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
}
