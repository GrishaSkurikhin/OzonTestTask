package shortlinks_test

import (
	"context"
	"errors"
	"math/rand"
	"strings"
	"testing"

	customerrors "github.com/GrishaSkurikhin/OzonTestTask/internal/custom-errors"
	"github.com/GrishaSkurikhin/OzonTestTask/internal/service/shortlinks"
	"github.com/GrishaSkurikhin/OzonTestTask/internal/service/shortlinks/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestGenerateToken(t *testing.T) {
	tokensNum := 1000000
	tokens := make(map[string]struct{}, tokensNum)
	alphabet := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz_0123456789"

	for i := 0; i < tokensNum; i++ {
		token := shortlinks.GenerateToken()
		if len(token) != 10 {
			t.Errorf("token %s has wrong length", token)
		}
		for _, c := range token {
			if !strings.Contains(alphabet, string(c)) {
				t.Errorf("token %s contains wrong character %c", token, c)
			}
		}

		if _, ok := tokens[token]; ok {
			t.Errorf("token %s is not unique", token)
		} else {
			tokens[token] = struct{}{}
		}
	}
}

func TestSaveURL(t *testing.T) {
	host := "localhost:8080"
	service := shortlinks.ShortlinksService{}

	cases := []struct {
		name             string
		longURL          string
		Error            error
		IsExistMockError error
		SaveMockError    error
	}{
		{
			name:             "Successfully",
			longURL:          "https://www.google.com",
			Error:            nil,
			IsExistMockError: nil,
			SaveMockError:    nil,
		},
		{
			name:             "Wrong URL",
			longURL:          "www.google",
			Error:            customerrors.WrongURL{Info: "wrong url"},
			IsExistMockError: nil,
			SaveMockError:    nil,
		},
		{
			name:             "Error while checking is url exist",
			longURL:          "https://www.google.com",
			Error:            errors.New("failed to check is url exist"),
			IsExistMockError: errors.New(""),
			SaveMockError:    nil,
		},
		{
			name:             "Error while saving url",
			longURL:          "https://www.google.com",
			Error:            errors.New("failed to save url"),
			IsExistMockError: nil,
			SaveMockError:    errors.New(""),
		},
		{
			name:             "Repeat tokens",
			longURL:          "https://www.google.com",
			Error:            nil,
			IsExistMockError: nil,
			SaveMockError:    nil,
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			urlSaverMock := mocks.NewURLSaver(t)

			if tc.name != "Wrong URL" {
				if tc.name == "Repeat tokens" {
					urlSaverMock.On("IsShortURLExists", mock.Anything, mock.AnythingOfType("string")).
						Return(true, tc.IsExistMockError).
						Times(2)

					urlSaverMock.On("IsShortURLExists", mock.Anything, mock.AnythingOfType("string")).
						Return(false, tc.IsExistMockError).
						Once()
				} else {
					urlSaverMock.On("IsShortURLExists", mock.Anything, mock.AnythingOfType("string")).
						Return(false, tc.IsExistMockError).
						Once()
				}

				if tc.IsExistMockError == nil {
					urlSaverMock.On("SaveURL", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("string")).
						Return(tc.SaveMockError).
						Once()
				}
			}

			_, err := service.SaveURL(context.Background(), tc.longURL, host, urlSaverMock)

			if tc.Error != nil {
				require.ErrorAs(t, err, &tc.Error)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestGetURL(t *testing.T) {
	host1 := "localhost:8080/"
	host2 := "http://example.com/"
	host3 := "example.com/"

	service := shortlinks.ShortlinksService{}

	cases := []struct {
		name      string
		shortURL  string
		longURL   string
		Error     error
		MockError error
	}{
		{
			name:      "Successfully host1",
			shortURL:  host1 + RandomString(10),
			longURL:   "https://www.google.com",
			Error:     nil,
			MockError: nil,
		},
		{
			name:      "Successfully host2",
			shortURL:  host2 + RandomString(10),
			longURL:   "https://www.google.com",
			Error:     nil,
			MockError: nil,
		},
		{
			name:      "Successfully host3",
			shortURL:  host3 + RandomString(10),
			longURL:   "https://www.google.com",
			Error:     nil,
			MockError: nil,
		},
		{
			name:      "Wrong URL",
			shortURL:  RandomString(20),
			longURL:   "",
			Error:     customerrors.WrongURL{Info: "wrong url"},
			MockError: nil,
		},
		{
			name:      "Wrong token",
			shortURL:  host1 + RandomString(5) + "!@#$%",
			longURL:   "",
			Error:     customerrors.WrongURL{Info: "wrong url"},
			MockError: nil,
		},
		{
			name:      "Wrong size of token",
			shortURL:  host1 + RandomString(20),
			longURL:   "",
			Error:     customerrors.WrongURL{Info: "wrong url"},
			MockError: nil,
		},
		{
			name:      "URL is not found",
			shortURL:  host1 + RandomString(10),
			longURL:   "",
			Error:     customerrors.URLNotFound{Info: "url not found"},
			MockError: nil,
		},
		{
			name:      "Error while getting url",
			shortURL:  host1 + RandomString(10),
			longURL:   "",
			Error:     errors.New("failed to get url"),
			MockError: errors.New(""),
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			urlGetterMock := mocks.NewURLGetter(t)

			if tc.name != "Wrong URL" && tc.name != "Wrong size of token" && tc.name != "Wrong token" {
				urlGetterMock.On("GetURL", mock.Anything, mock.AnythingOfType("string")).
					Return(tc.longURL, tc.MockError).
					Once()
			}

			_, err := service.GetURL(context.Background(), tc.shortURL, urlGetterMock)

			if tc.Error != nil {
				require.ErrorAs(t, err, &tc.Error)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func RandomString(n int) string {
	var letters = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz_0123456789")

	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}

	return string(b)
}
