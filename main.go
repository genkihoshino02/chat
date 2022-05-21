package main

import (
	"chat/trace"
	"flag"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"text/template"

	"github.com/stretchr/gomniauth"
	"github.com/stretchr/gomniauth/providers/google"
	"github.com/stretchr/objx"
)

const securityKey string = "securityKey"
const clientID string = "147770207720-53qfqmfjgack4hap4hiqbbtoof6g3b9g.apps.googleusercontent.com"
const privateKey string = "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"

type templateHandler struct {
	once     sync.Once
	filename string
	templ    *template.Template
}

func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		t.templ = template.Must(template.ParseFiles(filepath.Join("templates", t.filename)))
	})
	data := map[string]interface{}{
		"Host": r.Host,
	}
	if authCookie, err := r.Cookie("auth"); err == nil {
		data["UserData"] = objx.MustFromBase64(authCookie.Value)
	}
	t.templ.Execute(w, data)
}

func main() {
	var addr = flag.String("addr", ":8080", "Address of application")
	flag.Parse()

	gomniauth.SetSecurityKey(securityKey)
	gomniauth.WithProviders(google.New(clientID, privateKey, "http://localhost:8080/auth/callback/google"))

	r := newRoom()
	r.tracer = trace.New(os.Stdout)

	http.Handle("/chat", MustAuth(&templateHandler{filename: "chat.html"}))
	http.Handle("/login", &templateHandler{filename: "login.html"})
	http.HandleFunc("/auth/", loginHandler)
	http.Handle("/room", r)

	go r.run()

	log.Println("Launching web server. port: ", *addr)

	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("Listen and Serve :", err)
	}
}
