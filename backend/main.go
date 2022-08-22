package main

import (
	"net/http"
	"os"

	Router "github.com/sparkscience/logjam/backend/srv/router"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	// http.Handle("/room", http.FileServer(http.Dir("./public")))
	// http.Handle("/", http.FileServer(http.Dir("./public")))
	router := Router.NewRouter()
	http.Handle("/", router)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		panic(err)
	}
}
