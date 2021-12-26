package main

import (
	"github.com/Jedsonofnel/hubdc-api/handlers"
	"log"
	"net/http"
	"os"
)

func main() {
	l := log.New(os.Stdout, "event-api", log.LstdFlags)
	hh := handlers.NewHello(l)

    sm := http.NewServeMux()
    sm.Handle("/", hh)

    http.ListenAndServe(":9090", sm)
}
