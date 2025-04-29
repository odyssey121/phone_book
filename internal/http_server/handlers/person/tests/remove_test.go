package person_test

import (
	"encoding/json"
	"errors"
	"fmt"
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

func TestRemoveHandler(t *testing.T) {
	reqParam := make(map[string]string)
	reqParam["number"] = "890158834343"
	mockSearchResult := store.Person{"test", "test", 890158834343, "null"}
	cases := []struct {
		name             string
		alias            string
		url              string
		respError        string
		respStatus       string
		respMsg          string
		reqParam         map[string]string
		jsonPayload      string
		mockError        error
		mockSearchResult *store.Person
		status           int
	}{
		{
			name:       "Error Not Found",
			respError:  fmt.Sprintf("Person with number %s not found", reqParam["number"]),
			status:     http.StatusOK,
			respStatus: apiResp.StatusError,
			reqParam:   reqParam,
		},
		{
			name:             "Success",
			status:           http.StatusOK,
			respStatus:       apiResp.StatusOK,
			respMsg:          fmt.Sprintf("Record with number %s deleted", reqParam["number"]),
			reqParam:         reqParam,
			mockSearchResult: &mockSearchResult,
		},
		{
			name:             "Error db",
			respError:        apiResp.InternalServerErrorMsg,
			mockError:        errors.New("internal db error"),
			status:           http.StatusOK,
			respStatus:       apiResp.StatusError,
			reqParam:         reqParam,
			mockSearchResult: &mockSearchResult,
		},
	}

	for _, tc := range cases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			mockDb := store.NewMockDB(t)

			numberInt, _ := strconv.Atoi(tc.reqParam["number"])
			if tc.mockSearchResult != nil {
				mockDb.On("Remove", numberInt).
					Return(tc.mockError).
					Once()
			}

			mockDb.On("Search", numberInt).
				Return(tc.mockSearchResult).
				Once()

			handler := person.Remove(slogdummy.NewDiscardLogger(), mockDb)

			req, err := http.NewRequest(http.MethodDelete, "/", nil)

			req = mux.SetURLVars(req, tc.reqParam)

			require.NoError(t, err)

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			require.Equal(t, rr.Code, tc.status)

			body := rr.Body.String()

			var resp apiResp.Response

			require.NoError(t, json.Unmarshal([]byte(body), &resp))

			require.Equal(t, tc.respError, resp.Error)
			require.Equal(t, resp.Message, tc.respMsg)
			require.Equal(t, tc.respStatus, resp.Status)

		})
	}
}
