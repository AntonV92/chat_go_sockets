package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
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
	cookie, err := r.Cookie("authCookie")
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

	//string json value
	cookieVal := cookie.Value
	var decodedVal user.CookieValue

	fmt.Println(cookieVal)

	if jsonErr := json.Unmarshal([]byte(cookieVal), &decodedVal); jsonErr != nil {
		fmt.Printf("Json unmarshal error: %v\n", jsonErr)

	}

	if !user.CheckToken(decodedVal.UserId, decodedVal.Token) {
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

		var cookieVal user.CookieValue
		cookieVal.UserId = userIdentity.Id
		cookieVal.Token = userIdentity.Token.String

		jsonValue, err := json.Marshal(cookieVal)

		if err != nil {
			fmt.Printf("Token marshal error: %v\n", err)
		}

		fmt.Printf("Json cookie value: %s\n", jsonValue)

		cookie := http.Cookie{
			Name:     user.AuthCookieName,
			Value:    string(jsonValue),
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
