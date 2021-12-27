package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/Jedsonofnel/hubdc-api/handlers"
	"github.com/gorilla/mux"
)

func main() {
	l := log.New(os.Stdout, "event-api", log.LstdFlags)

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

	// sm.Handle("/events", eh)

	s := &http.Server{
		Addr:         ":9090",
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
