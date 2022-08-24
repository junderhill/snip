package main

import (
	"net/http"
	"snip/sniphandler"
)

func main() {
	handler := sniphandler.RedirectHandler()
	http.ListenAndServe(":8080", handler)
}
