package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"strings"
	"sync"
	"text/template"

	"github.com/joanassum/weather-web-app/infrastructure"
)

// templ represents a single template
type templateHandler struct {
	once     sync.Once
	filename string
	templ    *template.Template
}

// ServeHTTP handles the HTTP request.

func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		t.templ = template.Must(template.ParseFiles(filepath.Join("templates", t.filename)))
	})

	info, err := infrastructure.UseOpenWeatherAPI.GetWeatherInfo(strings.Split(r.URL.Path, "/")[1])
	if err != nil {
		http.Error(w, fmt.Sprintf("Error when trying to get weather info: %s", err), http.StatusBadRequest)
		return
	}

	data := map[string]interface{}{
		"Weather":  info,
		"TempUnit": "°C",
	}

	t.templ.Execute(w, data)
}

func main() {
	var addr = flag.String("addr", ":8080", "The addr of the application.")
	flag.Parse()

	http.Handle("/", &templateHandler{filename: "weather.html"})

	// start the web server
	log.Println("Starting web server on", *addr)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
