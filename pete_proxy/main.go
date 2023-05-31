package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
)

func Proxy(targetHost string) (*httputil.ReverseProxy, error) {
	url, err := url.Parse(targetHost)
	if err != nil {
		return nil, err
	}

	return httputil.NewSingleHostReverseProxy(url), nil
}

func ProxyHandler(proxy *httputil.ReverseProxy) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		proxy.ServeHTTP(w, r)
	}
}

func FiveOhTwo(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusBadGateway)
}

func main() {
	proxy, err := Proxy(os.Getenv("TEST_PROXY_ADDR"))
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/r/insights/v1/systems/", FiveOhTwo)
	http.HandleFunc("/", ProxyHandler(proxy))
	log.Fatal(http.ListenAndServe(":8080", nil))
}