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
	t.templ.Execute(w, r)
}

func main() {
	//option
	var addr = flag.String("addr", ":8888", "アプリケーションのアドレス")
	flag.Parse()
	// WebSocket
	r := newRoom()
	r.tracer = trace.New(os.Stdout)
	/*
		templateHandler内部のServeHTTPメソッドは、http.Handlerのインタフェスに適合しているので直接渡すことができる
	*/
	http.Handle("/", MustAuth(&templateHandler{filename: "chat.html"}))
	http.Handle("/room", r)
	// チャットルームを開始します
	go r.run()
	// Webサーバー起動
	log.Println("[*] Web serverを開始します。 Port :", *addr)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
