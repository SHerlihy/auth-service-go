package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
    fmt.Println("Hello, World!")

    routeHandler := http.NewServeMux()
    routeHandler.HandleFunc("GET /", func(w http.ResponseWriter, req *http.Request) {
        if req.URL.Path != "/" {
			http.NotFound(w, req)
			return
		}
    })

    routeHandler.HandleFunc("POST /", func(w http.ResponseWriter, req *http.Request) {
        if req.URL.Path != "/" {
			http.NotFound(w, req)
			return
		}
    })

    routeHandler.HandleFunc("PUT /", func(w http.ResponseWriter, req *http.Request) {
        if req.URL.Path != "/" {
			http.NotFound(w, req)
			return
		}
    })


    routeHandler.HandleFunc("DELETE /", func(w http.ResponseWriter, req *http.Request) {
        if req.URL.Path != "/" {
			http.NotFound(w, req)
			return
		}
    })

    server := &http.Server{
        Addr: ":8080",
        Handler: routeHandler,
    }

    log.Fatal(server.ListenAndServe())
}
