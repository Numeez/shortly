package shortly

import "log"

type Application struct {
	HealthHandler       *HealthService
	UrlShortenerHandler *UrlShortenerService
	Db                  *DB
}

func NewApplication() *Application {
	db, err := connectDB()
	if err != nil {
		log.Fatal(err)
	}
	healthHandler := newHealthHandler()
	urlShortenerHandler := getUrlShortenerService(db)

	return &Application{
		HealthHandler:       healthHandler,
		UrlShortenerHandler: urlShortenerHandler,
		Db:                  db,
	}
}
