package person

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"

	"github.com/go-chi/render"
	"github.com/gorilla/mux"
	"github.com/phone_book/internal/lib"
	apiResp "github.com/phone_book/internal/lib/api/response"
	"github.com/phone_book/internal/loger"
	"github.com/phone_book/internal/store"
)

func List(log *slog.Logger, db store.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "server.listHandler"
		log := log.With(slog.String("op", op))
		log.Info("new request", slog.String("Serving", r.URL.Path), slog.String("from", r.Host))

		data, err := db.List()
		if err != nil {
			log.Error(fmt.Sprintf("%s", err))
			render.JSON(w, r, apiResp.Error(apiResp.InternalServerErrorMsg))
			return

		}
		render.JSON(w, r, apiResp.WithData(data))
	}
}

type StatusResponse struct {
	apiResp.Response
	Total int `json:"total"`
}

func Status(log *slog.Logger, db store.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "server.statusHandler"
		log := log.With(slog.String("op", op))
		log.Info("new request", slog.String("Serving", r.URL.Path), slog.String("from", r.Host))

		total := db.CountRecords()
		render.JSON(w, r, StatusResponse{apiResp.OK(""), total})

	}
}

func Insert(log *slog.Logger, db store.DB) http.HandlerFunc {
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
			render.JSON(w, r, apiResp.Error(apiResp.InternalServerErrorMsg))
			return
		}
		render.JSON(w, r, apiResp.OK("new record added successfully"))
	}
}

func Search(log *slog.Logger, db store.DB) http.HandlerFunc {
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

func Remove(log *slog.Logger, db store.DB) http.HandlerFunc {
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
			render.JSON(w, r, apiResp.Error(apiResp.InternalServerErrorMsg))
			return
		}
		render.JSON(w, r, apiResp.OK(fmt.Sprintf("Record with number %d deleted", n)))

	}
}
