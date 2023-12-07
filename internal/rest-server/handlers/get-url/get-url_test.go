package geturl_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	customerrors "github.com/GrishaSkurikhin/OzonTestTask/internal/custom-errors"
	geturl "github.com/GrishaSkurikhin/OzonTestTask/internal/rest-server/handlers/get-url"
	"github.com/GrishaSkurikhin/OzonTestTask/internal/rest-server/handlers/get-url/mocks"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestGetURLHandler(t *testing.T) {
	cases := []struct {
		name      string
		shortURL  string
		respError string
		mockError error
	}{
		{
			name:      "Successfully",
			shortURL:  "example.com/ghJfkaisDn",
			respError: "",
			mockError: nil,
		},
		{
			name:      "Empty request body",
			shortURL:  "",
			respError: "empty request",
			mockError: nil,
		},
		{
			name:      "Wrong request body",
			shortURL:  "",
			respError: "failed to decode request",
			mockError: nil,
		},
		{
			name:      "Empty URL",
			shortURL:  "",
			respError: "url is required",
			mockError: nil,
		},
		{
			name:      "Wrong URL",
			shortURL:  "example.com",
			respError: "wrong url",
			mockError: customerrors.WrongURL{},
		},
		{
			name:      "URL not found",
			shortURL:  "example.com/ghJfkaisDn",
			respError: "url not found",
			mockError: customerrors.URLNotFound{},
		},
		{
			name: "Internal error",
			shortURL:  "example.com/ghJfkaisDn",
			respError: "internal error",
			mockError: fmt.Errorf("internal error"),
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			serviceURLGetterMock := mocks.NewServiceURLGetter(t)

			if tc.name == "Successfully" || 
			tc.name == "Wrong URL" || 
			tc.name == "URL not found" || 
			tc.name == "Internal error"{
				serviceURLGetterMock.On("GetURL", mock.Anything, tc.shortURL, mock.Anything).
					Return("string", tc.mockError).
					Once()
			}

			handler := geturl.New(zerolog.Nop(), nil, serviceURLGetterMock)

			var (
				req *http.Request
				err error
			)
			if tc.name == "Empty request body" {
				req, err = http.NewRequest(http.MethodGet, "/url", bytes.NewBufferString(""))
			} else if tc.name == "Wrong request body" {
				req, err = http.NewRequest(http.MethodGet, "/url", bytes.NewBufferString("wrong"))
			} else {
				input := fmt.Sprintf(`{"shortURL": "%s"}`, tc.shortURL)
				req, err = http.NewRequest(http.MethodGet, "/url", bytes.NewBufferString(input))
			}
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			require.Equal(t, rr.Code, http.StatusOK)

			body := rr.Body.String()

			var resp geturl.Response

			require.NoError(t, json.Unmarshal([]byte(body), &resp))

			require.Equal(t, tc.respError, resp.Error)
		})
	}
}
