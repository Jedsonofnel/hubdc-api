package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/Jedsonofnel/hubdc-api/handlers"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
    // create new logger using standard go logger
	l := log.New(os.Stdout, "event-api", log.LstdFlags)

    // godotenv for env files
    if os.Getenv("APP_ENV") != "production" {
        err := godotenv.Load()
        if err != nil {
            l.Fatal("Error loading .env file")
        }
    }

	eh := handlers.NewEvents(l)

    // create a new servemux using gorilla/mux
	sm := mux.NewRouter()

    getRouter := sm.Methods(http.MethodGet).Subrouter()
    getRouter.HandleFunc("/events", eh.GetEvents)

    putRouter := sm.Methods(http.MethodPut).Subrouter()
    putRouter.HandleFunc("/event/{id:[0-9]+}", eh.UpdateEvent)
    putRouter.Use(eh.MiddlewareEventValidation)

    postRouter := sm.Methods(http.MethodPost).Subrouter()
    postRouter.HandleFunc("/events", eh.AddEvent)
    postRouter.Use(eh.MiddlewareEventValidation)

	s := &http.Server{
        Addr:         fmt.Sprintf(":%s", os.Getenv("PORT")),
		Handler:      sm,
        ErrorLog:     l,
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}

	go func() {
		err := s.ListenAndServe()
		if err != nil {
			l.Fatal(err)
		}
	}()

	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, os.Kill)

	sig := <-sigChan
	l.Println("Received terminate, graceful shutdown", sig)

	tc, _ := context.WithTimeout(context.Background(), 30*time.Second)
	s.Shutdown(tc)
}
