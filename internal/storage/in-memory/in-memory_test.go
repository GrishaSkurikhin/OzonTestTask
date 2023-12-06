package inmemory_test

import (
	"strconv"
	"testing"

	inmemory "github.com/GrishaSkurikhin/OzonTestTask/internal/storage/in-memory"
	"github.com/stretchr/testify/require"
)

func TestInMemory(t *testing.T) {
	storage := inmemory.New()

	casesNum := 100

	type testCase struct {
		shortURL string
		longURL  string
	}

	cases := make([]testCase, 0, casesNum)
	for i := 0; i < casesNum; i++ {
		cases = append(cases, testCase{
			shortURL: "localhost:8080/" + strconv.Itoa(i),
			longURL:  "https://example.com/" + strconv.Itoa(i),
		})
	}

	for _, tc := range cases {
		tc := tc
		go func() {
			isExist, err := storage.IsShortURLExists(tc.shortURL)
			if err != nil {
				t.Errorf("IsShortURLExists() error = %v", err)
			}
			require.False(t, isExist)

			err = storage.SaveURL(tc.longURL, tc.shortURL)
			if err != nil {
				t.Errorf("SaveURL() error = %v", err)
			}

			isExist, err = storage.IsShortURLExists(tc.shortURL)
			if err != nil {
				t.Errorf("IsShortURLExists() error = %v", err)
			}
			require.True(t, isExist)

			longURL, err := storage.GetURL(tc.shortURL)
			if err != nil {
				t.Errorf("GetURL() error = %v", err)
			}
			require.Equal(t, tc.longURL, longURL)
		}()
	}
}
