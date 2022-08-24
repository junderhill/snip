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
		log.Fatal(err)
		return "", false
	}
	defer db.Close()

	var value []byte

	e := db.View(func(txn *badger.Txn) error {
		item, e := txn.Get([]byte(slug))

		if e != nil {
			log.Print(e)
			return nil
		}

		err := item.Value(func(val []byte) error {
			value = append([]byte{}, val...)
			return nil
		})
		if err != nil {
			log.Print(err)
			return nil
		}
		return nil
	})

	if e != nil {
		log.Print(e)
		return "", false
	}

	if value != nil {
		return fmt.Sprintf("%b", value), true
	}
	return "", false
}
