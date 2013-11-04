package mudlib

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"runtime"
	"runtime/debug"
)

var httpPort = flag.Int("httpPort", 8080, "Port for HTTP interface")

func init() {
	http.HandleFunc("/gc", gcHandler)
	http.HandleFunc("/mem", memHandler)
	http.HandleFunc("/errors", errorHandler)
	go startServing()
}

func startServing() {
	log.Printf("HTTP listening on port %d", *httpPort)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", *httpPort), nil); err != nil {
		errorLog.Printf("Failed to start HTTP server on port %d", *httpPort)
	}
}

// TODO: Parse the json and print in HTML
func gcHandler(w http.ResponseWriter, r *http.Request) {
	gcStats := new(debug.GCStats)
	debug.ReadGCStats(gcStats)
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, "%+v", *gcStats)
}

func memHandler(w http.ResponseWriter, r *http.Request) {
	memStats := new(runtime.MemStats)
	runtime.ReadMemStats(memStats)
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, "%+v", *memStats)
}

func errorHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: either read the last few lines of error log, or keep them in memory and write the
	// error log buffered.
}
