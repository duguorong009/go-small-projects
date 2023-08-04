package main

import (
	"log"
	"net"
	"net/http"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

// Create a custom visitor struct which holds the rate limiter for each
// visitor and the last time that the visitor was seen.
type visitor struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

// Create a map to hold the rate limiters for each visitor and a mutex.
var visitors = make(map[string]visitor)
var mu sync.Mutex

// Run a background goroutine to remove old entries from the visitors map.
func init() {
	go cleanupVisitors()
}

// Retrieve and return the rate limiter for the current visitor if it
// already exists. Otherwise create a new rate limiter and add it to
// the visitors map, using the IP address as the key.
func getVisitor(ip string) *rate.Limiter {
	mu.Lock()
	defer mu.Unlock()

	v, exists := visitors[ip]
	if !exists {
		limiter := rate.NewLimiter(1, 3)
		// Include the current time when creating a new visitor.
		newVisitor := visitor{limiter, time.Now()}
		visitors[ip] = newVisitor
		return limiter
	}

	// Update the last seen time for the visitor.
	v.lastSeen = time.Now()
	return v.limiter
}

// every minute check the map for visitors that haven't been seen for
// more than 3 minutes and delete the entries.
func cleanupVisitors() {
	for {
		time.Sleep(time.Minute)

		mu.Lock()
		for ip, v := range visitors {
			if time.Since(v.lastSeen) > 3*time.Minute {
				delete(visitors, ip)
			}
		}
		mu.Unlock()
	}
}

func limit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get the IP address for the current user.
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			log.Print(err.Error())
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// Call the getVisitor function to retrieve the rate limiter for
		// the current user.
		limiter := getVisitor(ip)
		if !limiter.Allow() {
			http.Error(w, http.StatusText(http.StatusTooManyRequests), http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}
