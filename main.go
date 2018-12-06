package main

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/rs/zerolog"
)

func main() {
	logger := zerolog.New(os.Stderr).With().Timestamp().Logger().Level(zerolog.InfoLevel)
	var err error

	routes := httprouter.New()
	routes.GET("/greeting", Greeting)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: routes,
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			logger.Error().Err(err).Msgf("http server error")
		}
	}()

	logger.Info().Msgf("server started")

	<-done
	logger.Info().Msgf("server stoped")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err = srv.Shutdown(ctx); err != nil {
		logger.Error().Err(err).Msgf("https server shutdown failed")
	}
}

// Greeting is for hello world la!
func Greeting(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	json.NewEncoder(w).Encode(struct {
		Greeting string `json:"greeting"`
	}{
		Greeting: "Hello world!",
	})
}
