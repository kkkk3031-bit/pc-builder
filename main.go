package main

import (
	"embed"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/khaledq84ever/pc-builder/build"
)

//go:embed templates/*
var tmplFS embed.FS

var indexTmpl = template.Must(template.ParseFS(tmplFS, "templates/index.html"))

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /", handleIndex)
	mux.HandleFunc("POST /recommend", handleRecommend)

	addr := ":8080"
	log.Printf("PC Builder running at http://localhost%s", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	render(w, nil)
}

type PageData struct {
	Result build.Build
}

func handleRecommend(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	raw := r.FormValue("budget")
	budget, err := strconv.ParseFloat(raw, 64)
	if err != nil {
		render(w, &PageData{Result: build.Build{ErrorMessage: "Invalid number. Please enter a valid amount."}})
		return
	}

	result := build.Recommend(budget)
	render(w, &PageData{Result: result})
}

func render(w http.ResponseWriter, data *PageData) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := indexTmpl.Execute(w, data); err != nil {
		log.Printf("template error: %v", err)
		http.Error(w, "Internal error", http.StatusInternalServerError)
	}
}


