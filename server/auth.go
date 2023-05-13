package server

import (
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"
	"ws_app/user"
)

type CheckAuth struct {
	handler httpHanlder
}

func authenticatedRequest(f httpHanlder) *CheckAuth {
	return &CheckAuth{handler: f}
}

func (check *CheckAuth) check(w http.ResponseWriter, r *http.Request) {
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

	check.handler(w, r)

	if !user.CheckToken(userID, token) {
		http.Redirect(w, r, "/login", 302)
	}
}

func authenticateUser() httpHanlder {
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
	}
}
