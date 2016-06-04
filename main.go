package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

func main() {
	go StartAPI()
	StartMonitoring()

}

func StartAPI() {
	router := NewRouter()
	log.Fatal(http.ListenAndServe(":8080", router))
}

func StartMonitoring() {
	go UuidProvider()

	// Websites = TestWebsites // to be removed
	WebsitesToMonitor()

	go StartWebsiteMonitors(Websites)

	for {
		status := <-statusChan
		b, _ := json.MarshalIndent(status, "", "  ")
		fmt.Println(string(b))
	}

}

func StartWebsiteMonitors([]Website) {
	for _, website := range Websites {
		time.Sleep(10 * time.Millisecond)
		go MonitorWebsite(website)
	}
}

func StopMonitoring() {

}

func Logger(inner http.Handler, name string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		inner.ServeHTTP(w, r)
		log.Printf(
			"%s\t%s\t%s\t%s",
			r.Method,
			r.RequestURI,
			name,
			time.Since(start),
		)
	})
}

// This slice of websites is for temporary development testing,
// and will be removed once data persistence is in place.
var TestWebsites = []Website{
	Website{Url: `http://www.google.com/robots.txt`, Interval: 5, Email: `alert@chrisbenson.com`},
	Website{Url: `http://www.microsoft.com/robots.txt`, Interval: 5, Email: `alert@chrisbenson.com`},
	Website{Url: `http://www.amazon.com/robots.txt`, Interval: 5, Email: `alert@chrisbenson.com`},
	Website{Url: `http://www.codeguard.com/this-webpage-will-fail.html`, Interval: 5, Email: `alert@chrisbenson.com`},
}
