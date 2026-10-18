package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	addr := flag.String("listen", "127.0.0.1:9091", "health listen address")
	flag.Parse()
	m := http.NewServeMux()
	m.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		h, _ := os.Hostname()
		_ = json.NewEncoder(w).Encode(map[string]any{"ok": true, "service": "nek-agent", "hostname": h, "time": time.Now().UTC()})
	})
	log.Printf("NEK Agent scaffold listening on %s", *addr)
	log.Fatal(http.ListenAndServe(*addr, m))
}
