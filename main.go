package main

import (
	"context"
	"embed"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/khaledq84ever/pc-builder/build"
)

//go:embed templates/*
var tmplFS embed.FS

var indexTmpl = template.Must(template.New("index").Funcs(template.FuncMap{
	"dollar": func(v float64) string {
		if v == float64(int64(v)) {
			return "$" + strconv.FormatInt(int64(v), 10)
		}
		return "$" + strconv.FormatFloat(v, 'f', 2, 64)
	},
}).ParseFS(tmplFS, "templates/index.html"))

type PageData struct {
	Result build.Build
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /", handleIndex)
	mux.HandleFunc("POST /recommend", handleRecommend)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	server := &http.Server{
		Addr:         ":" + port,
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	go func() {
		log.Printf("PC Builder running on port %s", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %v", err)
		}
	}()

	<-ctx.Done()
	log.Println("shutting down...")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("shutdown: %v", err)
	}
	log.Println("stopped")
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	render(w, &PageData{})
}

func handleRecommend(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil || r.PostForm == nil {
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
