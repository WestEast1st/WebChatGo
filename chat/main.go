package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"text/template"

	"../trace"
	"github.com/joho/godotenv"
	"github.com/stretchr/gomniauth"
	"github.com/stretchr/gomniauth/providers/github"
	"github.com/stretchr/objx"
)

//http ハンドラーのtemplate
type templateHandler struct {
	once     sync.Once
	filename string
	templ    *template.Template
}

func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		t.templ =
			template.Must(template.ParseFiles(filepath.Join("templates", t.filename)))
	})
	data := map[string]interface{}{
		"Host": r.Host,
	}
	if authCookie, err := r.Cookie("auth"); err == nil {
		data["UserData"] = objx.MustFromBase64(authCookie.Value)
	}
	t.templ.Execute(w, data)
}

func Env_load() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	Env_load()
	//option
	var addr = flag.String("addr", ":8888", "アプリケーションのアドレス")
	flag.Parse()
	// Gomniauth
	gomniauth.SetSecurityKey("セキュリティキー")
	gomniauth.WithProviders(
		github.New(os.Getenv("GITHUB_CLIENT_ID"), os.Getenv("GITHUB_SECRET"), "http://localhost:8888/auth/callback/github"),
	)
	// WebSocket
	r := newRoom(UseFileSystemAvatar)
	r.tracer = trace.New(os.Stdout)
	/*
		templateHandler内部のServeHTTPメソッドは、http.Handlerのインタフェスに適合しているので直接渡すことができる
	*/
	http.Handle("/", MustAuth(&templateHandler{filename: "chat.html"}))
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("./assets/"))))
	http.Handle("/avatars/", http.StripPrefix("/avatars/", http.FileServer(http.Dir("./avatars/"))))
	http.Handle("/login", &templateHandler{filename: "login.html"})
	http.HandleFunc("/auth/", loginHandler)
	http.Handle("/upload", &templateHandler{filename: "upload.html"})
	http.HandleFunc("/uploader", uploaderHandler)
	http.Handle("/room", r)
	http.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
		http.SetCookie(w, &http.Cookie{
			Name:   "auth",
			Value:  "",
			Path:   "/",
			MaxAge: -1,
		})
		w.Header()["Location"] = []string{"/chat"}
		w.WriteHeader(http.StatusTemporaryRedirect)
	})

	// チャットルームを開始します
	go r.run()
	// Webサーバー起動
	log.Println("[*] Web serverを開始します。 Port :", *addr)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
