package mudlib

import (
	"bufio"
	"flag"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"runtime"
	"runtime/debug"
)

const (
	numErrorLines = 100
	gcTemplateContent =`
<html>
	<head>
		<title>GC</title>
		<meta http-equiv="refresh" content="3">
	</head>
	<body>
		<h1>GC</h1>
		<table>
			<tr><th>Last GC</th><td>{{.LastGC}}</td></tr>
			<tr><th>Num GC</th><td>{{.NumGC}}</td></tr>
			<tr><th>Pause Total</th><td>{{.PauseTotal}}</td></tr>
			<tr><th>Pause</th><td>{{.Pause}}</td></tr>
			<tr><th>Pause Quantiles</th><td>{{.PauseQuantiles}}</td></tr>
		</table>
	</body>
</html>`
	memTemplateContent =`
<html>
	<head>
		<title>Mem</title>
		<meta http-equiv="refresh" content="3">
	</head>
	<body>
		<h1>Mem</h1>
		<h2>General</h2>
		<table>
			<tr><th>Alloc</th><td>{{.Alloc}}</td></tr>
			<tr><th>Total alloc</th><td>{{.TotalAlloc}}</td></tr>
			<tr><th>Sys</th><td>{{.Sys}}</td></tr>
			<tr><th>Lookups</th><td>{{.Lookups}}</td></tr>
			<tr><th>Mallocs</th><td>{{.Mallocs}}</td></tr>
			<tr><th>Frees</th><td>{{.Frees}}</td></tr>
		</table>

		<h2>Heap</h2>
		<table>
			<tr><th>Alloc</th><td>{{.HeapAlloc}}</td></tr>
			<tr><th>Sys</th><td>{{.HeapSys}}</td></tr>
			<tr><th>Idle</th><td>{{.HeapIdle}}</td></tr>
			<tr><th>Inuse</th><td>{{.HeapInuse}}</td></tr>
			<tr><th>Released</th><td>{{.HeapReleased}}</td></tr>
			<tr><th>Objects</th><td>{{.HeapObjects}}</td></tr>
		</table>

		<h2>Low-level</h2>
		<table>
			<tr><th>Stack Inuse</th><td>{{.StackInuse}}</td></tr>
			<tr><th>Stack Sys</th><td>{{.StackSys}}</td></tr>
			<tr><th>MSpan Inuse</th><td>{{.MSpanInuse}}</td></tr>
			<tr><th>MSpan Sys</th><td>{{.MSpanSys}}</td></tr>
			<tr><th>MCache Inuse</th><td>{{.MCacheInuse}}</td></tr>
			<tr><th>MCache Sys</th><td>{{.MCacheSys}}</td></tr>
			<tr><th>Bucket Hash Sys</th><td>{{.BuckHashSys}}</td></tr>
		</table>
		<h2>Per-size</h2>
		<table>
			<tr><th>Size</th><th>Mallocs</th><th>Frees</th></tr>
			{{range .BySize}}
				<tr><td>{{.Size}}</td><td>{{.Mallocs}}</td><td>{{.Frees}}</td></tr>
			{{end}}
		</table>
	</body>
</html>
`
	errorTemplateContent =`
<html>
	<head>
		<title>Error</title>
		<meta http-equiv="refresh" content="5">
	</head>
	<body>
		<h1>Error</h1>
		<table>
			<tr><th>Log entry</th></tr>
			{{range .}}
				<tr><td>{{.}}</td></tr>
			{{end}}
		</table>
	</body>
</html>`
)

var (
	httpPort = flag.Int("httpPort", 8080, "Port for HTTP interface")
	gcTemplate = template.New("GC")
	memTemplate = template.New("Mem")
	errorTemplate = template.New("Error")
)

func init() {
	gcTemplate = template.Must(gcTemplate.Parse(gcTemplateContent))
	memTemplate = template.Must(memTemplate.Parse(memTemplateContent))
	errorTemplate = template.Must(errorTemplate.Parse(errorTemplateContent))

	http.HandleFunc("/gc", gcHandler)
	http.HandleFunc("/mem", memHandler)
	http.HandleFunc("/errors", errorHandler)
	http.HandleFunc("/", rootHandler)
	go startServing()
}

func startServing() {
	log.Printf("HTTP listening on port %d", *httpPort)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", *httpPort), nil); err != nil {
		errorLog.Printf("Failed to start HTTP server on port %d", *httpPort)
	}
}

func gcHandler(w http.ResponseWriter, r *http.Request) {
	gcStats := new(debug.GCStats)
	debug.ReadGCStats(gcStats)
	if err := gcTemplate.Execute(w, *gcStats); err != nil {
		errorLog.Printf("Failed to execute GC template: %+v", err)
		w.WriteHeader(500)
	}
}

func memHandler(w http.ResponseWriter, r *http.Request) {
	memStats := new(runtime.MemStats)
	runtime.ReadMemStats(memStats)
	if err := memTemplate.Execute(w, *memStats); err != nil {
		errorLog.Printf("Failed to execute Mem template: %+v", err)
		w.WriteHeader(500)
	}
}

func errorHandler(w http.ResponseWriter, r *http.Request) {
	log := make([]string, numErrorLines)

	logFile, err := os.Open("logs/error")
	if err != nil {
		errorLog.Printf("Unable to open error log for reading")
	}
	defer logFile.Close()

	reader := bufio.NewReader(logFile)
	line, err := reader.ReadString('\n')
	for err == nil {
		// TODO: maybe parse lines?
		log = log[1:]
		log = append(log, line)

		line, err = reader.ReadString('\n')
	}
	if err != io.EOF {
		errorLog.Printf("Error while reading error log: %+v", err)
		return
	}

	if err := errorTemplate.Execute(w, log); err != nil {
		errorLog.Printf("Failed to execute Error template: %+v", err)
		w.WriteHeader(500)
	}
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	http.NotFound(w, r)
}
