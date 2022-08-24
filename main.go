package main

import (
	"fmt"
	"net/http"
	"os"
	"snip/sniphandler"
	"snip/snipinit"
)

func main() {

	if len(os.Args) > 1 && os.Args[1] == "init" {
		snipinit.InitDb()
	} else {
		handler := sniphandler.RedirectHandler()

		fmt.Println("Listening on port 8080")
		http.ListenAndServe(":8080", handler)
	}
}
