package main

import (
	"regexp"
	"log"
	"net/http"
	"gopkg.in/elazarl/goproxy.v1"
)

func FiveOhTwo(r *http.Request,ctx *goproxy.ProxyCtx) (*http.Request,*http.Response) {
        return r,goproxy.NewResponse(r, goproxy.ContentTypeText, http.StatusBadGateway, "Bad Gateway")
}

func main() {
	proxy := goproxy.NewProxyHttpServer()
	proxy.Verbose = true
	
	proxy.OnRequest(goproxy.ReqHostMatches(regexp.MustCompile(".*redhat.*:443$"))).HandleConnect(goproxy.AlwaysMitm)
	
	proxy.OnRequest(goproxy.UrlMatches(regexp.MustCompile(".*/r/insights/v1/systems/.*"))).DoFunc(
		func(r *http.Request,ctx *goproxy.ProxyCtx)(*http.Request,*http.Response) {
        return r,goproxy.NewResponse(r, goproxy.ContentTypeText, http.StatusBadGateway,"Bad Gateway")
	})

	log.Fatal(http.ListenAndServe(":3129", proxy))
}
