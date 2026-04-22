package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Numeez/shortly"
)

func main() {
	app := shortly.NewApplication()
	defer app.Db.Conn.Close()
	rateLimiterStore := shortly.RedisRateLimiterStore()
	router := shortly.InitRouter(app, rateLimiterStore)
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("Environment variable PORT is not set")
	}
	server := http.Server{
		ReadTimeout:       10 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
		Handler:           router,
		Addr:              ":" + port,
	}
	go func() {
		log.Println("Server running on port", server.Addr)
		if err := server.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop
	log.Println("Shutting down server........")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Println(err)
	}
	log.Println("Server exited successfully")
}
