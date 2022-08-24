package sniphandler

import (
	"fmt"
	"github.com/dgraph-io/badger/v3"
	"log"
	"net/http"
)

func MainHandler() http.HandlerFunc {
	return func(out http.ResponseWriter, in *http.Request) {
		if in.Method == "GET" {

			redirectHandler(out, in)
			return

		} else if in.Method == "POST" {

			createRedirectHandler(out, in)
			return

		} else {
			out.WriteHeader(404)
			return
		}
	}
}

func createRedirectHandler(out http.ResponseWriter, in *http.Request) {
	slug := in.URL.Path
	_, found := GetRedirectUrl(slug)

	url := in.Header.Get("X-Redirect-Target")

	if url == "" {
		out.WriteHeader(400)
		out.Write([]byte("Redirect Target Header `X-Redirect-Target` not set"))
		return
	}

	if found {
		out.WriteHeader(409)
		return
	}

	err := saveRedirectUrl(slug, url)
	if err != nil {
		log.Print(err)
		out.WriteHeader(500)
	}

	out.WriteHeader(200)
	return
}

func saveRedirectUrl(slug string, url string) error {
	db, err := badger.Open(badger.DefaultOptions("/tmp/badger"))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	txn := db.NewTransaction(true)

	txn.Set([]byte(slug), []byte(url))
	return txn.Commit()
}

func redirectHandler(out http.ResponseWriter, in *http.Request) {
	slug := in.URL.Path
	url, ok := GetRedirectUrl(slug)

	if !ok {
		log.Printf("URL for %s not found..", slug)
		out.WriteHeader(404)
		return
	}

	log.Printf("URL for %s found: %s", slug, url)

	out.Header().Add("Location", url)
	out.WriteHeader(301)
	return
}

func GetRedirectUrl(slug string) (url string, ok bool) {
	log.Printf("Attempting to get URL for %s", slug)
	db, err := badger.Open(badger.DefaultOptions("/tmp/badger"))
	if err != nil {
		return handleError(err)
	}
	defer db.Close()

	txn := db.NewTransaction(false)
	defer txn.Discard()

	item, err := txn.Get([]byte(slug))
	if err != nil {
		return handleError(err)
	}

	var value []byte
	value, err = item.ValueCopy(value)
	if err != nil {
		return handleError(err)
	}

	if value != nil {
		return fmt.Sprintf("%s", value), true

	}
	return "", false
}

func handleError(err error) (string, bool) {
	log.Print(err)
	return "", false
}
