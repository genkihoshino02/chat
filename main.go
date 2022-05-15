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
)

type templateHandler struct{
	once sync.Once
	filename string
	templ 	*template.Template
}

func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request){
	t.once.Do(func() {
		t.templ = template.Must(template.ParseFiles(filepath.Join("templates",t.filename)))
	})
	t.templ.Execute(w,r)
}

func main(){
	var addr = flag.String("addr",":8080","Address of application")
	flag.Parse()
	r := newRoom()
	r.tracer = trace.New(os.Stdout)
	http.Handle("/",&templateHandler{filename:"chat.html"})
	http.Handle("/room",r)

	go r.run()

	log.Println("Launching web server. port: ",*addr)

	if err := http.ListenAndServe(*addr,nil);err != nil{
		log.Fatal("Listen and Serve :",err);
	}
}