package person_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/phone_book/internal/http_server/handlers/person"
	apiResp "github.com/phone_book/internal/lib/api/response"
	"github.com/phone_book/internal/loger/slogdummy"
	"github.com/phone_book/internal/store"
	"github.com/stretchr/testify/require"
)

func TestListHandler(t *testing.T) {
	cases := []struct {
		name       string
		alias      string
		url        string
		respError  string
		respData   []store.Person
		mockResult []store.Person
		mockError  error
		status     int
	}{
		{
			name:      "Error",
			mockError: errors.New("unexpected err"),
			respError: apiResp.InternalServerErrorMsg,
			status:    http.StatusOK,
		},
		{
			name:       "Success",
			status:     http.StatusOK,
			respData:   []store.Person{store.Person{"test", "test", 1, "null"}},
			mockResult: []store.Person{store.Person{"test", "test", 1, "null"}},
		},
	}

	for _, tc := range cases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			mockDb := store.NewMockDB(t)
			if tc.respError == "" || tc.mockError != nil {
				mockDb.On("List").
					Return(tc.mockResult, tc.mockError).
					Once()
			}

			handler := person.List(slogdummy.NewDiscardLogger(), mockDb)

			req, err := http.NewRequest(http.MethodGet, "/", nil)
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			require.Equal(t, rr.Code, tc.status)

			body := rr.Body.String()

			var resp apiResp.ResponseDataPersons

			require.NoError(t, json.Unmarshal([]byte(body), &resp))

			require.Equal(t, tc.respError, resp.Error)
			if tc.respData != nil {
				require.Equal(t, len(tc.respData), len(resp.Data))
				for i, d := range resp.Data {
					require.Equal(t, d.FirstName, resp.Data[i].FirstName)
				}

			}

		})
	}
}
