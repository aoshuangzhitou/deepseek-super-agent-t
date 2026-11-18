// Command hello-web serves a single page that outputs hello.
package main

import (
	"flag"
	"log"
	"net/http"
	"time"
)

// helloBody is the entire content of the page.
const helloBody = "hello"

// newHandler returns the HTTP handler for the hello page.
func newHandler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/", handleHello)
	return mux
}

// handleHello writes the hello page for GET and HEAD requests to the root
// path. Other paths return 404 and other methods return 405.
func handleHello(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	if r.Method != http.MethodGet && r.Method != http.MethodHead {
		w.Header().Set("Allow", "GET, HEAD")
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("Cache-Control", "no-store")
	w.WriteHeader(http.StatusOK)
	if r.Method == http.MethodGet {
		_, _ = w.Write([]byte(helloBody))
	}
}

func main() {
	addr := flag.String("addr", "127.0.0.1:8123", "address to listen on")
	flag.Parse()

	log.Printf("hello-web listening on http://%s/", *addr)
	srv := &http.Server{
		Addr:              *addr,
		Handler:           newHandler(),
		ReadHeaderTimeout: 5 * time.Second,
	}
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
