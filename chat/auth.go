package main

import (
	"net/http"
)

type authHandler struct {
	next http.Handler
}

func (h *authHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if _, err := r.Cookie("auth"); err == http.ErrNoCookie {
		// not auth
		w.Header().Set("location", "/login")
		w.WriteHeader(http.StatusTemporaryRedirect)
	} else if err != nil {
		// 何らかのエラー
		panic(err.Error())
	} else {
		// success
		h.next.ServeHTTP(w, r)
	}
}

// 認証が必須の際に利用するハンドラー
func MustAuth(handler http.Handler) http.Handler {
	return &authHandler{next: handler}
}
