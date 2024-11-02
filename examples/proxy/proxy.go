// This example implements a reverse proxy that sends requests to postman-echo.com.
//
// To run: `tinygo run -target=wasip2-http.json ./examples/proxy`
// Test GET: `curl -v 'http://0.0.0.0:8080/get'`
// Test POST: `curl -v -d hello 'http://0.0.0.0:8080/post'`

package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/ydnar/wasi-http-go/wasihttp"
)

func init() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		ctx, cancel := context.WithCancel(r.Context())
		defer cancel()

		r2 := r.Clone(ctx)
		r2.Host = "postman-echo.com"
		r2.URL.Host = "postman-echo.com"
		r2.URL.Scheme = "https"

		defer func() {
			dur := time.Since(start)
			log.Printf("proxied %s in %s\n", r2.URL.String(), dur.String())
		}()

		client := &http.Client{
			Transport: &wasihttp.Transport{},
		}
		res, err := client.Do(r2)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "error: %v", err)
			log.Printf("error: %v", err)
			return
		}

		for k, v := range res.Header {
			for _, vv := range v {
				w.Header().Add(k, vv)
			}
		}

		w.WriteHeader(res.StatusCode)
		if res.Body != nil {
			io.Copy(w, res.Body)
			res.Body.Close()
		}
	})
}

func main() {}
