package shortly

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

type UrlShortenerService struct {
	Db Database
}

func getUrlShortenerService(db *DB) *UrlShortenerService {
	return &UrlShortenerService{
		Db: db,
	}
}

func (us *UrlShortenerService) HandlerShortener(w http.ResponseWriter, r *http.Request) {
	var request UrlRequest
	if !strings.EqualFold(r.Method, "POST") {
		WriteErrorResponse(w, http.StatusMethodNotAllowed, "Invalid request method", nil)
		return
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		WriteErrorResponse(w, http.StatusBadRequest, "invalid request body", nil)
		return
	}
	if request.Url == "" {
		WriteErrorResponse(w, http.StatusBadRequest, "url should not be empty", nil)
		return
	}
	url, err := NormalizeURL(request.Url)
	if err != nil {
		WriteErrorResponse(w, http.StatusBadRequest, "", err)
		return

	}
	id, err := us.Db.InsertUrl(url)
	if err != nil {
		WriteErrorResponse(w, http.StatusInternalServerError, "", err)
		return
	}
	result := getShortUrl(id)
	WriteResponse(w, http.StatusOK, ShortenURL{Url: result})

}

func (us *UrlShortenerService) HandlerGetURL(w http.ResponseWriter, r *http.Request) {
	shortCode := r.PathValue("shortCode")
	code, err := DecodeBase62(shortCode)
	if err != nil {
		WriteErrorResponse(w, http.StatusBadRequest, "invalid short code url", nil)
		return
	}

	redirectURl, err := us.Db.GetUrl(code)
	if err != nil {
		WriteErrorResponse(w, http.StatusInternalServerError, "", err)
		return
	}
	http.Redirect(w, r, redirectURl, http.StatusFound)
}

func getShortUrl(id int64) string {
	serverURL := "http://localhost:9090/"
	return serverURL + EncodeBase62(id)
}

func NormalizeURL(rawURL string) (string, error) {
	rawURL = strings.TrimSpace(rawURL)

	if !strings.HasPrefix(rawURL, "http://") && !strings.HasPrefix(rawURL, "https://") {
		rawURL = "https://" + rawURL
	}

	parsed, err := url.Parse(rawURL)
	if err != nil {
		return "", fmt.Errorf("invalid url")
	}

	if parsed.Host == "" {
		return "", fmt.Errorf("invalid url: missing host")
	}

	parsed.Host = strings.ToLower(parsed.Host)

	parsed.Path = strings.TrimRight(parsed.Path, "/")

	host := parsed.Hostname()
	port := parsed.Port()
	if (port == "80" && parsed.Scheme == "http") || (port == "443" && parsed.Scheme == "https") {
		parsed.Host = host
	}

	parsed.Fragment = ""

	return parsed.String(), nil
}
