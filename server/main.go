package server

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"sync/atomic"
)

func MiddlewareHandler(capacity int64, next http.Handler) http.Handler {
	var counter int64

	h := func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt64(&counter, 1)

		if atomic.LoadInt64(&counter) > capacity {
			defer r.Body.Close()
			fmt.Fprint(w, http.StatusTooManyRequests, http.StatusText(http.StatusTooManyRequests), "\n")
		} else {
			next.ServeHTTP(w, r)
		}

		atomic.AddInt64(&counter, -1)
	}

	return http.HandlerFunc(h)
}

func NewServer(capacity int64, path, port string, finalHandler http.Handler) {
	http.Handle(path, MiddlewareHandler(capacity, finalHandler))

	l, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Listen: %v", err)
	}

	defer l.Close()

	log.Fatal(http.Serve(l, nil))
}
