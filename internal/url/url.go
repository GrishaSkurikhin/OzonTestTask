package url

import (
	"fmt"
	"math/rand"
	"net/url"
	"strings"
	"time"

	customerrors "github.com/GrishaSkurikhin/OzonTestTask/internal/custom-errors"
)

func GenerateToken() string {
	// Алгоритм:
	// 1. Взять текущее время в наносекундах и прибавить к нему случайное число от 0 до 1000
	// 2. Заполнить строку случайными символами из алфавита по следующему правилу:
	// 2.1. Взять остаток от деления на длину алфавита
	// 2.2. Добавить символ из алфавита с индексом, равным остатку
	// 2.3. Поделить текущее время на длину алфавита
	// Уникальность текущего времени гарантирует уникальность сгенерированной строки (в большинстве случаев)
	// Добавление случайного числа повышает вероятность уникальности сгенерированной строки в случае быстрого повторного вызова функции

	alphabet := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz_0123456789"
	urlLen := 10

	uniqueTime := time.Now().UnixNano() + rand.Int63n(1000)
	shortURL := strings.Builder{}

	for i := 0; i < urlLen; i++ {
		shortURL.WriteByte(alphabet[uniqueTime%int64(len(alphabet))])
		uniqueTime /= int64(len(alphabet))
	}

	return shortURL.String()
}

type URLSaver interface {
	SaveURL(longURL string, shortURL string) error
	IsShortURLExists(shortURL string) (bool, error)
}

func SaveURL(longURL string, host string, saver URLSaver) (string, error) {
	const op = "url.SaveURL"

	// Валидация URL
	u, err := url.Parse(longURL)
	if !(err == nil && u.Scheme != "" && u.Host != "") {
		return "", customerrors.WrongURL{Info: err.Error()}
	}

	token := GenerateToken()
	shortURL := host + "/" + token

	// Проверка на уникальность
	for {
		exists, err := saver.IsShortURLExists(shortURL)
		if err != nil {
			return "", fmt.Errorf("%s: failed to check is url exist: %v", op, err)
		}
		if !exists {
			break
		}
		token = GenerateToken()
		shortURL = host + "/" + token
	}

	err = saver.SaveURL(longURL, shortURL)
	if err != nil {
		return "", fmt.Errorf("%s: failed to save url: %v", op, err)
	}

	return shortURL, nil
}

type URLGetter interface {
	GetURL(shortURL string) (string, error)
}

func GetURL(shortURL string, getter URLGetter) (string, error) {
	const op = "url.GetURL"

	u, err := url.Parse(shortURL)
	if !(err == nil && u.Scheme != "" && u.Host != "") {
		return "", customerrors.WrongURL{Info: err.Error()}
	}

	longURL, err := getter.GetURL(shortURL)
	if err != nil {
		return "", fmt.Errorf("%s: failed to get url: %v", op, err)
	}

	if longURL == "" {
		return "", customerrors.URLNotFound{Info: fmt.Sprintf("url %s not found", shortURL)}
	}

	return longURL, nil
}
