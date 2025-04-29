package person_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/phone_book/internal/http_server/handlers/person"
	"github.com/phone_book/internal/lib"
	apiResp "github.com/phone_book/internal/lib/api/response"
	"github.com/phone_book/internal/loger/slogdummy"
	"github.com/phone_book/internal/store"
	"github.com/stretchr/testify/require"
)

func TestInsertHandler(t *testing.T) {
	insertPayload := store.Person{"test", "test", 1, "null"}
	cases := []struct {
		name        string
		alias       string
		url         string
		respError   string
		respStatus  string
		respMsg     string
		reqData     store.Person
		jsonPayload string
		mockError   error
		status      int
	}{
		{
			name:       "Error",
			mockError:  errors.New("unexpected err"),
			respError:  apiResp.InternalServerErrorMsg,
			status:     http.StatusOK,
			respStatus: apiResp.StatusError,
			reqData:    insertPayload,
		},
		{
			name:       "Success",
			status:     http.StatusOK,
			respStatus: apiResp.StatusOK,
			respMsg:    "new record added successfully",
			reqData:    insertPayload,
		},
		{
			name:        "De serialize error",
			status:      http.StatusOK,
			respStatus:  apiResp.StatusError,
			respError:   "de serialize error",
			reqData:     insertPayload,
			jsonPayload: "de serialize error",
		},
	}

	for _, tc := range cases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			mockDb := store.NewMockDB(t)
			if tc.respError == "" || tc.mockError != nil {
				mockDb.On("Insert", tc.reqData.FirstName, tc.reqData.LastName, tc.reqData.Phone).
					Return(tc.mockError).
					Once()
			}

			handler := person.Insert(slogdummy.NewDiscardLogger(), mockDb)
			buf := new(bytes.Buffer)
			if tc.jsonPayload != "" {
				buf.Write([]byte(tc.jsonPayload))
			} else {
				lib.Serialize(tc.reqData, buf)
			}

			req, err := http.NewRequest(http.MethodPost, "/", buf)

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
