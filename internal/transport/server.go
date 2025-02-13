package transport

import (
	"log"
	"net/http"
)

func Run() {
	MapRoutes()
	log.Fatal(http.ListenAndServe(":8080", nil))
}
