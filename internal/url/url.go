package url

func GenerateShortURL(longURL string) string {
	return "shortURL"
}

type URLSaver interface {
	SaveURL(longURL string) error
}

func SaveURL(longURL string, saver URLSaver) (string, error) {
	return "", nil
}

type URLGetter interface {
	GetURL(shortURL string) (string, error)
}

func GetURL(shortURL string, getter URLGetter) (string, error) {
	return "", nil
}