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

// TODO: Templates
func gcHandler(w http.ResponseWriter, r *http.Request) {
	gcStats := new(debug.GCStats)
	debug.ReadGCStats(gcStats)
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintln(w, "<html><head><title>GC</title></head>")
	fmt.Fprintln(w, "<body>")
	fmt.Fprintln(w, "<h1>GC</h1>")
	fmt.Fprintln(w, "<table>")
	fmt.Fprintf(w, "<tr><th>Last GC</th><td>%v</td></tr>\n", gcStats.LastGC)
	fmt.Fprintf(w, "<tr><th>Num GC</th><td>%v</td></tr>\n", gcStats.NumGC)
	fmt.Fprintf(w, "<tr><th>Pause Total</th><td>%v</td></tr>\n", gcStats.PauseTotal)
	fmt.Fprintf(w, "<tr><th>Pause</th><td>%v</td></tr>\n", gcStats.Pause)
	fmt.Fprintf(w, "<tr><th>Pause Quantiles</th><td>%v</td></tr>\n", gcStats.PauseQuantiles)
	fmt.Fprintln(w, "</table>")
	fmt.Fprintln(w, "</body>")
	fmt.Fprintln(w, "</html>")
}

func memHandler(w http.ResponseWriter, r *http.Request) {
	memStats := new(runtime.MemStats)
	runtime.ReadMemStats(memStats)
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintln(w, "<html><head><title>Mem</title></head>")
	fmt.Fprintln(w, "<body>")
	fmt.Fprintln(w, "<h1>Mem</h1>")
	fmt.Fprintln(w, "<h2>General</h2>")
	fmt.Fprintln(w, "<table>")
	fmt.Fprintf(w, "<tr><th>Alloc</th><td>%v</td></tr>\n", memStats.Alloc)
	fmt.Fprintf(w, "<tr><th>Total alloc</th><td>%v</td></tr>\n", memStats.TotalAlloc)
	fmt.Fprintf(w, "<tr><th>Sys</th><td>%v</td></tr>\n", memStats.Sys)
	fmt.Fprintf(w, "<tr><th>Lookups</th><td>%v</td></tr>\n", memStats.Lookups)
	fmt.Fprintf(w, "<tr><th>Mallocs</th><td>%v</td></tr>\n", memStats.Mallocs)
	fmt.Fprintf(w, "<tr><th>Frees</th><td>%v</td></tr>\n", memStats.Frees)
	fmt.Fprintln(w, "</table>")
	fmt.Fprintln(w, "<h2>Heap</h2>")
	fmt.Fprintln(w, "<table>")
	fmt.Fprintf(w, "<tr><th>Alloc</th><td>%v</td></tr>\n", memStats.HeapAlloc)
	fmt.Fprintf(w, "<tr><th>Sys</th><td>%v</td></tr>\n", memStats.HeapSys)
	fmt.Fprintf(w, "<tr><th>Idle</th><td>%v</td></tr>\n", memStats.HeapIdle)
	fmt.Fprintf(w, "<tr><th>Inuse</th><td>%v</td></tr>\n", memStats.HeapInuse)
	fmt.Fprintf(w, "<tr><th>Released</th><td>%v</td></tr>\n", memStats.HeapReleased)
	fmt.Fprintf(w, "<tr><th>Objects</th><td>%v</td></tr>\n", memStats.HeapObjects)
	fmt.Fprintln(w, "</table>")
	fmt.Fprintln(w, "<h2>Low-level</h2>")
	fmt.Fprintln(w, "<table>")
	fmt.Fprintf(w, "<tr><th>Stack Inuse</th><td>%v</td></tr>\n", memStats.StackInuse)
	fmt.Fprintf(w, "<tr><th>Stack Sys</th><td>%v</td></tr>\n", memStats.StackSys)
	fmt.Fprintf(w, "<tr><th>MSpan Inuse</th><td>%v</td></tr>\n", memStats.MSpanInuse)
	fmt.Fprintf(w, "<tr><th>MSpan Sys</th><td>%v</td></tr>\n", memStats.MSpanSys)
	fmt.Fprintf(w, "<tr><th>MCache Inuse</th><td>%v</td></tr>\n", memStats.MCacheInuse)
	fmt.Fprintf(w, "<tr><th>MCache Sys</th><td>%v</td></tr>\n", memStats.MCacheSys)
	fmt.Fprintf(w, "<tr><th>Bucket Hash Sys</th><td>%v</td></tr>\n", memStats.BuckHashSys)
	fmt.Fprintln(w, "</table>")
	fmt.Fprintln(w, "<h2>Per-size</h2>")
	fmt.Fprintln(w, "<table>")
	fmt.Fprintln(w, "<tr><th>Size</th><th>Mallocs</th><th>Frees</th></tr>")
	// TODO: histogram
	for _, bs := range memStats.BySize {
		fmt.Fprintf(w, "<tr><td>%v</td><td>%v</td><td>%v</td></tr>\n", bs.Size, bs.Mallocs, bs.Frees)
	}
	fmt.Fprintln(w, "</table>")
	fmt.Fprintln(w, "</body>")
	fmt.Fprintln(w, "</html>")
}

func errorHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: either read the last few lines of error log, or keep them in memory and write the
	// error log buffered.
}
