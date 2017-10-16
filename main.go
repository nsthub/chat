package main

import (
	"flag"
	"github.com/nsthub/tracer"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"text/template"
)

// template represents a single template
type templateHandler struct {
	once     sync.Once
	filename string
	template *template.Template
}

// ServeHTTP handles the HTTP request.
func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		t.template =
			template.Must(template.ParseFiles(filepath.Join("templates", t.filename)))
	})
	t.template.Execute(w, r)
}

func main() {
	var addr = flag.String("addr", ":8080", "The addr of the application.")
	flag.Parse() // parse the flags

	r := newRoom()
	r.tracer = trace.New(os.Stdout)
	// root
	http.Handle("/", &templateHandler{filename: "chat.html"})
	http.Handle("/room", r)
	// get the room going
	go r.run()
	// start the web server
	log.Println("Starting web server on", *addr)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
