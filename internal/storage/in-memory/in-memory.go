package inmemory

import (
	"context"
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

func (i *InMemory) SaveURL(ctx context.Context, longURL string, shortURL string) error {
	i.urls.Store(shortURL, longURL)
	return nil
}

func (i *InMemory) GetURL(ctx context.Context, shortURL string) (string, error) {
	longURL, ok := i.urls.Load(shortURL)
	if !ok {
		return "", nil
	}
	return longURL.(string), nil
}

func (i *InMemory) IsShortURLExists(ctx context.Context, shortURL string) (bool, error) {
	_, ok := i.urls.Load(shortURL)
	return ok, nil
}
