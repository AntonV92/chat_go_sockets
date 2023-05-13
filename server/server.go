package server

import (
	"log"
	"net/http"
	"os"
)

var routes = map[string]func(w http.ResponseWriter, r *http.Request){
	"/": actionIndex(),
}

func InitServer() {
	for path, handler := range routes {
		http.HandleFunc(path, handler)
	}

	http.ListenAndServe(":8000", nil)
}

func actionIndex() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		render("frontend/index.html", w)
	}
}

func render(fileName string, w http.ResponseWriter) {
	content, err := os.ReadFile(fileName)
	if err != nil {
		log.Fatal(err)
	}

	_, wErr := w.Write(content)
	if wErr != nil {
		log.Fatal(wErr)
	}
}
