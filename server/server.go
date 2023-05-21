package server

import (
	"log"
	"net/http"
	"os"
)

type httpHanlder func(w http.ResponseWriter, r *http.Request)

var routes = map[string]httpHanlder{
	"/":      authenticatedRequest(actionIndex()).check,
	"/login": actionLogin(),
	"/auth":  actionAuth(),
	"/ws":    actionGetWsConnection(),
}

func InitServer() {
	for path, handler := range routes {
		http.HandleFunc(path, handler)
	}

	curr, _ := os.Getwd()

	fs := http.FileServer(http.Dir(curr + "/frontend"))

	http.Handle("/frontend/", http.StripPrefix("/frontend", noCache(fs)))

	http.ListenAndServe(":8000", nil)
	//http.ListenAndServeTLS(":443", "cert.crt", "private.key", nil)
}

func noCache(f http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Cache-Control", "no-cache, no-store, must-revalidate;")
		f.ServeHTTP(w, r)
	})
}

func actionIndex() httpHanlder {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate;")
		render("frontend/index.html", w)
	}
}

func actionLogin() httpHanlder {
	return func(w http.ResponseWriter, r *http.Request) {
		render("frontend/login.html", w)
	}
}

func actionAuth() httpHanlder {
	return authenticateUser()
}

func actionGetWsConnection() httpHanlder {
	return getConnection()
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
