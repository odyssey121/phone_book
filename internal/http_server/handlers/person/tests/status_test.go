package person_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/phone_book/internal/http_server/handlers/person"
	"github.com/phone_book/internal/loger/slogdummy"
	"github.com/phone_book/internal/store"
	"github.com/stretchr/testify/require"
)

func TestStatusHandler(t *testing.T) {

	cases := []struct {
		name             string
		alias            string
		url              string
		respError        string
		mockSearchResult int
		respData         int
		mockError        error
		status           int
		reqParam         map[string]string
	}{
		{
			name:             "Success",
			status:           http.StatusOK,
			mockSearchResult: 1,
			respData:         1,
		},
	}

	for _, tc := range cases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			mockDb := store.NewMockDB(t)
			if tc.respError == "" || tc.mockError != nil {
				mockDb.On("CountRecords").
					Return(tc.mockSearchResult).
					Once()
			}

			handler := person.Status(slogdummy.NewDiscardLogger(), mockDb)

			req, err := http.NewRequest(http.MethodGet, "/", nil)
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			require.Equal(t, rr.Code, tc.status)

			body := rr.Body.String()

			var resp person.StatusResponse

			require.NoError(t, json.Unmarshal([]byte(body), &resp))

			require.Equal(t, tc.respError, resp.Error)
			require.Equal(t, tc.respData, resp.Total)

		})
	}
}
