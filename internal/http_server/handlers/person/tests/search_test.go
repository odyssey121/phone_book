package person_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/gorilla/mux"
	"github.com/phone_book/internal/http_server/handlers/person"
	apiResp "github.com/phone_book/internal/lib/api/response"
	"github.com/phone_book/internal/loger/slogdummy"
	"github.com/phone_book/internal/store"
	"github.com/stretchr/testify/require"
)

func TestSearchHandler(t *testing.T) {
	reqParam := make(map[string]string)
	reqParam["number"] = "890134343331"
	searchPerson := store.Person{"i'm searched!", "wowo", 890134343331, ""}

	cases := []struct {
		name             string
		alias            string
		url              string
		respError        string
		mockSearchResult *store.Person
		respData         store.Person
		mockError        error
		status           int
		reqParam         map[string]string
	}{
		{
			name:             "Success",
			status:           http.StatusOK,
			mockSearchResult: &searchPerson,
			respData:         searchPerson,
			reqParam:         reqParam,
		},
		{
			name:     "Not found",
			status:   http.StatusOK,
			reqParam: reqParam,
		},
	}

	for _, tc := range cases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			mockDb := store.NewMockDB(t)
			if tc.respError == "" || tc.mockError != nil {
				phoneInt, _ := strconv.Atoi(reqParam["number"])
				mockDb.On("Search", phoneInt).
					Return(tc.mockSearchResult).
					Once()
			}

			handler := person.Search(slogdummy.NewDiscardLogger(), mockDb)

			req, err := http.NewRequest(http.MethodGet, "/", nil)
			req = mux.SetURLVars(req, tc.reqParam)
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			require.Equal(t, rr.Code, tc.status)

			body := rr.Body.String()

			var resp apiResp.ResponseDataPerson

			require.NoError(t, json.Unmarshal([]byte(body), &resp))

			require.Equal(t, tc.respError, resp.Error)
			require.Equal(t, tc.respData, resp.Data)

		})
	}
}
