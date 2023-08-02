package main

import (
	"log"
	"net/http"
	"time"
)

func main() {
	// Use the http.NewServerMux() function to create an empty servemux.
	mux := http.NewServeMux()

	// Use the http.RedirectHandler() function to create a handler which 307
	// redirects all requests it receives to http://example.org
	rh := http.RedirectHandler("http://example.org", 307)

	// Next we use the mux.Handle() function to register this with our new
	// servemux, so it acts as the handler for all incoming requests with the URL
	// path /foo.
	mux.Handle("/foo", rh)

	// Initialize the timeHandler in exactly the same way we would any normal
	// struct.
	th := timeHandler{format: time.RFC1123}

	// Like the above example, we use the mux.Handle() function to register
	// this with our ServeMux.
	mux.Handle("/time", th)

	// Convert the timeHandlerFunc function to a http.HandlerFunc type.
	th1 := http.HandlerFunc(timeHandlerFunc)

	// Add it to ServeMux.
	mux.Handle("/time1", th1)

	// Another way to add handler func
	mux.HandleFunc("/time2", timeHandlerFunc)

	// 3rd way to add handler func
	th2 := timeHandlerFunc1(time.RFC1123)
	mux.Handle("/time2", th2)

	log.Print("Listening...")

	// Then we create a new server and start listening for incoming requests
	// with the http.ListenAndServe() function, passing in our servemux for it to
	// match requests against as the second parameter.
	http.ListenAndServe(":3000", mux)

}

type timeHandler struct {
	format string
}

func (th timeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	tm := time.Now().Format(th.format)
	w.Write([]byte("The time is: " + tm))
}

func timeHandlerFunc(w http.ResponseWriter, r *http.Request) {
	tm := time.Now().Format(time.RFC1123)
	w.Write([]byte("The time is: " + tm))
}

func timeHandlerFunc1(format string) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		tm := time.Now().Format(format)
		w.Write([]byte("The time is: " + tm))
	}
	return http.HandlerFunc(fn)
}
