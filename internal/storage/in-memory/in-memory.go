package inmemory

import (
	"sync"
)

type InMemory struct {
	urls *sync.Map
}

func New() *InMemory {
	return &InMemory{
		urls: &sync.Map{},
	}
}

func (i *InMemory) SaveURL(longURL string, shortURL string) error {
	i.urls.Store(shortURL, longURL)
	return nil
}

func (i *InMemory) GetURL(shortURL string) (string, error) {
	longURL, ok := i.urls.Load(shortURL)
	if !ok {
		return "", nil
	}
	return longURL.(string), nil
}

func (i *InMemory) IsShortURLExists(shortURL string) (bool, error) {
	_, ok := i.urls.Load(shortURL)
	return ok, nil
}
