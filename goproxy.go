package main

import (
	"flag"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
)

type Prox struct {
	target *url.URL
	proxy  *httputil.ReverseProxy
}

func New(target string) *Prox {
	url, _ := url.Parse(target)

	return &Prox{target: url, proxy: httputil.NewSingleHostReverseProxy(url)}
}

func (p *Prox) handle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("X-GoProxy", "GoProxy")

	p.proxy.ServeHTTP(w, r)
}

func main() {
	const (
		defaultPort        = ":8001"
		defaultPortUsage   = "default server port, ':80', ':8080'..."
		defaultTarget      = "http://127.0.0.1:8080"
		defaultTargetUsage = "default redirect url, 'http://127.0.0.1:8080'"
	)

	// flags
	port := flag.String("port", defaultPort, defaultPortUsage)
	url := flag.String("url", defaultTarget, defaultTargetUsage)

	flag.Parse()

	fmt.Println("server will run on : %s", *port)
	fmt.Println("redirecting to :%s", *url)

	// proxy
	proxy := New(*url)

	// server
	http.HandleFunc("/", proxy.handle)
	http.ListenAndServe(*port, nil)
}
