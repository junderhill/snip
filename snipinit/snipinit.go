package snipinit

import (
	"github.com/dgraph-io/badger/v3"
	"log"
)

type Url struct {
	Slug string
	Url  string
}

func InitDb() {

	records := []Url{
		Url{
			Slug: "/twitter",
			Url:  "https://www.twitter.com/underhillj",
		},
		Url{
			Slug: "/photography",
			Url:  "https://www.junderhill.com/",
		},
	}

	log.Printf("Attempt to init db with %d records", len(records))

	db, err := badger.Open(badger.DefaultOptions("/tmp/badger"))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	txn := db.NewTransaction(true)
	for _, v := range records {
		log.Printf("Inserting %s : %s", v.Slug, v.Url)

		if err := txn.Set([]byte(v.Slug), []byte(v.Url)); err == badger.ErrTxnTooBig {
			_ = txn.Commit()
			txn = db.NewTransaction(true)
			_ = txn.Set([]byte(v.Slug), []byte(v.Url))
		}
	}
	_ = txn.Commit()
}
