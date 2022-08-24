package sniphandler

import (
	"fmt"
	"github.com/dgraph-io/badger/v3"
	"log"
	"net/http"
)

func RedirectHandler() http.HandlerFunc {
	return func(out http.ResponseWriter, in *http.Request) {
		if in.Method != "GET" {
			log.Printf("%s is not supported, only GET", in.Method)
			out.WriteHeader(404)
			return
		}

		slug := in.URL.Path
		log.Printf("Attempting to get URL for %s", slug)
		url, ok := GetRedirectUrl(slug)

		if !ok {
			log.Printf("URL for %s not found..", slug)
			out.WriteHeader(404)
			return
		}

		log.Printf("URL for %s found: %s", slug, url)

		out.Header().Add("Location", url)
		out.WriteHeader(301)
	}
}

func GetRedirectUrl(slug string) (url string, ok bool) {
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
