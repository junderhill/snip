package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"os"
	"snip/sniphandler"
	"snip/snipinit"
	"sync"
)

func main() {

	if len(os.Args) > 1 && os.Args[1] == "init" {
		snipinit.InitDb()
	} else {

		r := mux.NewRouter()

		r.HandleFunc("/{slug}", sniphandler.RedirectHandler).Methods(http.MethodGet)
		r.HandleFunc("/{slug}", sniphandler.CreateRedirectHandler).Methods(http.MethodPost)
		r.Use(loggingMiddleware)

		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			http.ListenAndServe(":8080", r)
		}()

		fmt.Println("Listening on port 8080")
		wg.Wait()
	}
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("%s Request to %s\n", r.Method, r.RequestURI)
		next.ServeHTTP(w, r)
	})
}
