package server

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"ws_app/user"
)

var routes = map[string]func(w http.ResponseWriter, r *http.Request){
	"/":      actionIndex(),
	"/login": actionLogin(),
	"/auth":  actionAuth(),
}

func InitServer() {
	for path, handler := range routes {
		http.HandleFunc(path, handler)
	}

	http.ListenAndServe(":8000", nil)
}

func actionIndex() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		checkAuth(w, r)
		render("frontend/index.html", w)
	}
}

func checkAuth(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(user.AuthCookieName)
	if err != nil {
		switch {
		case errors.Is(err, http.ErrNoCookie):
			log.Println(err)
		default:
			log.Println(err)
			http.Error(w, "server error", http.StatusInternalServerError)
		}

		http.Redirect(w, r, "/login", 302)
		return
	}

	data := strings.Split(cookie.Value, "|")

	userID, err := strconv.Atoi(data[0])
	if err != nil {
		log.Println(err)
		return
	}
	token := data[1]

	if !user.CheckToken(userID, token) {
		http.Redirect(w, r, "/login", 302)
	}
}

func actionLogin() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		render("frontend/login.html", w)
	}
}

func actionAuth() func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate;")

		r.ParseForm()

		login := r.Form.Get("login")
		pass := r.Form.Get("password")

		// check credentials and get user instance
		userIdentity, loginError := user.Login(login, pass)

		if loginError != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		userID := strconv.Itoa(userIdentity.Id)
		cookieVal := userID + "|" + userIdentity.Token.String

		cookie := http.Cookie{
			Name:     user.AuthCookieName,
			Value:    cookieVal,
			Path:     "/",
			MaxAge:   60,
			HttpOnly: true,
			Secure:   true,
		}

		http.SetCookie(w, &cookie)
		http.Redirect(w, r, "/", 302)
		fmt.Println("Set cookie")
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
