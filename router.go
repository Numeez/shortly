package shortly

import "net/http"

func InitRouter(app *Application, rateLimiterStore RateLimiterAllow) http.Handler {
	router := http.NewServeMux()
	router.HandleFunc("/health", app.HealthHandler.HandlerHealth)
	router.HandleFunc("/shorten", app.UrlShortenerHandler.HandlerShortener)
	router.HandleFunc("/{shortCode}", app.UrlShortenerHandler.HandlerGetURL)
	return RateLimiterMiddleware(rateLimiterStore)(router)
}
