package inmemory

type InMemory struct {
	urls map[string]string
}

func New() *InMemory {
	return &InMemory{
		urls: make(map[string]string),
	}
}

func (i *InMemory) SaveURL(longURL string, shortURL string) error {
	i.urls[shortURL] = longURL
	return nil
}

func (i *InMemory) GetURL(shortURL string) (string, error) {
	if longURL, ok := i.urls[shortURL]; ok {
		return longURL, nil
	}
	return "", nil
}

func (i *InMemory) IsShortURLExists(shortURL string) (bool, error) {
	_, ok := i.urls[shortURL]
	return ok, nil
}