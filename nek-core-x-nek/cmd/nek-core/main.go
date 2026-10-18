package main

import (
	"github.com/qopix/nek/internal/app"
	"log"
	"net/http"
	"os"
)

func env(k, fallback string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return fallback
}
func main() {
	s, err := app.New(env("NEK_DATA", "data/nek.json"), os.Getenv("NEK_ADMIN_TOKEN"), env("NEK_PUBLIC_URL", "http://127.0.0.1:8080"))
	if err != nil {
		log.Fatal(err)
	}
	if s.GeneratedToken != "" {
		log.Printf("SAVE THIS ADMIN TOKEN: %s", s.GeneratedToken)
	}
	addr := env("NEK_ADDR", "127.0.0.1:8080")
	log.Printf("X-NEK listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, s.Handler()))
}
