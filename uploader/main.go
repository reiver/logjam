package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func UploadFile(w http.ResponseWriter, r *http.Request) {
	file, handler, err := r.FormFile("file")
	fileName := r.FormValue("file_name")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	path := fmt.Sprintf("public/%s", handler.Filename)

	type writeError struct {
		Message string `json:"message"`
	}

	respondError := func() {
		w.WriteHeader(500)
		b, err := json.Marshal(writeError{"Failed to create file"})
		if err != nil {
			panic(err)
		}
		w.Write(b)
	}

	err = os.MkdirAll("public", os.ModePerm)
	if err != nil {
		respondError()
		return
	}
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		respondError()
		return
	}

	_, err = io.Copy(f, file)
	if err != nil {
		respondError()
	}

	type response struct {
		Path    string `json:"path"`
		Message string `json:"message"`
	}

	defer f.Close()

	b, err := json.Marshal(response{fmt.Sprintf("/%s", path), fmt.Sprintf("Sucessfully written %s!", fileName)})
	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
}

func OnlyPost(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			w.WriteHeader(404)
			w.Write([]byte("Not found"))
			return
		}

		handler.ServeHTTP(w, r)
	})
}

func homeLink(w http.ResponseWriter, r *http.Request) {
	_, _ = fmt.Fprintf(w, "Welcome home!")
}

func main() {
	router := mux.NewRouter()
	router.Handle("/file", http.HandlerFunc(UploadFile)).Methods("POST")
	router.PathPrefix("/public").Handler(http.StripPrefix("/public/", http.FileServer(http.Dir("./public"))))
	router.HandleFunc("/", homeLink)
	port := 8081
	log.Printf("Server listening on port %d\n", port)
	handler := cors.Default().Handler(router)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), handler))
}
