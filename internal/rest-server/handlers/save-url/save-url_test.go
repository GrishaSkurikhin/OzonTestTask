package saveurl_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	customerrors "github.com/GrishaSkurikhin/OzonTestTask/internal/custom-errors"
	saveurl "github.com/GrishaSkurikhin/OzonTestTask/internal/rest-server/handlers/save-url"
	"github.com/GrishaSkurikhin/OzonTestTask/internal/rest-server/handlers/save-url/mocks"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestSaveURLHandler(t *testing.T) {
	cases := []struct {
		name      string
		longURL   string
		respError string
		mockError error
	}{
		{
			name:      "Successfully",
			longURL:   "http://www.example.com",
			respError: "",
			mockError: nil,
		},
		{
			name:      "Empty request body",
			longURL:   "",
			respError: "empty request",
			mockError: nil,
		},
		{
			name:      "Wrong request body",
			longURL:   "",
			respError: "failed to decode request",
			mockError: nil,
		},
		{
			name:      "Empty URL",
			longURL:   "",
			respError: "url is required",
			mockError: nil,
		},
		{
			name:      "Wrong URL",
			longURL:   "example.com",
			respError: "wrong url",
			mockError: customerrors.WrongURL{},
		},
		{
			name: "Internal error",
			longURL:  "http://www.example.com",
			respError: "internal error",
			mockError: fmt.Errorf("internal error"),
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			serviceURLSaverMock := mocks.NewServiceURLSaver(t)

			if tc.name == "Successfully" || 
			tc.name == "Wrong URL" || 
			tc.name == "Internal error"{
				serviceURLSaverMock.On("SaveURL", mock.Anything, tc.longURL, mock.AnythingOfType("string"), mock.Anything).
					Return("string", tc.mockError).
					Once()
			}

			handler := saveurl.New(zerolog.Nop(), nil, "host", serviceURLSaverMock)

			var (
				req *http.Request
				err error
			)
			if tc.name == "Empty request body" {
				req, err = http.NewRequest(http.MethodPost, "/url/save", bytes.NewBufferString(""))
			} else if tc.name == "Wrong request body" {
				req, err = http.NewRequest(http.MethodPost, "/url/save", bytes.NewBufferString("wrong"))
			} else {
				input := fmt.Sprintf(`{"longURL": "%s"}`, tc.longURL)
				req, err = http.NewRequest(http.MethodPost, "/url/save", bytes.NewBufferString(input))
			}
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			require.Equal(t, rr.Code, http.StatusOK)

			body := rr.Body.String()

			var resp saveurl.Response

			require.NoError(t, json.Unmarshal([]byte(body), &resp))

			require.Equal(t, tc.respError, resp.Error)
		})
	}
}
