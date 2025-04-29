/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/phone_book/internal/config"
	"github.com/phone_book/internal/http_server/handlers"
	personHandlers "github.com/phone_book/internal/http_server/handlers/person"
	"github.com/phone_book/internal/loger"
	"github.com/phone_book/internal/store"
	"github.com/spf13/cobra"
)

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

		getMux.HandleFunc("/", handlers.DefaultHeandler)
		getMux.HandleFunc("/list", personHandlers.List(log, db))
		getMux.HandleFunc("/status", personHandlers.Status(log, db))
		getMux.HandleFunc("/search", personHandlers.Search(log, db))
		getMux.HandleFunc("/search/{number}", personHandlers.Search(log, db))

		postMux := gMux.Methods(http.MethodPost).Subrouter()

		postMux.HandleFunc("/insert", personHandlers.Insert(log, db))

		deleteMux := gMux.Methods(http.MethodDelete).Subrouter()
		deleteMux.HandleFunc("/remove", personHandlers.Remove(log, db))
		deleteMux.HandleFunc("/remove/{number}", personHandlers.Remove(log, db))

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
