package main

import (
	"crypto/tls"
	"crypto/x509"
	"log"
	"net/http"
	"os"
	"regexp"

	"gopkg.in/elazarl/goproxy.v1"
)

func FiveOhTwo(r *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
	if os.Getenv("first") == "true" {
		os.Setenv("first", "false")
		return r, nil
	} else {
		os.Setenv("first", "true")
		return r, goproxy.NewResponse(r, goproxy.ContentTypeText, http.StatusBadGateway, "Bad Gateway")
	}
}

func main() {
	err := setCA("cert.pem", "key.pem")
	if err != nil {
		log.Fatal(err.Error())
		return
	}
	os.Setenv("first", "true")
	proxy := goproxy.NewProxyHttpServer()
	proxy.Verbose = true

	proxy.OnRequest(goproxy.ReqHostMatches(regexp.MustCompile(".*redhat.*:443$"))).HandleConnect(goproxy.AlwaysMitm)

	proxy.OnRequest(goproxy.UrlMatches(regexp.MustCompile(".*/r/insights/v1/systems/.*"))).DoFunc(FiveOhTwo)

	log.Fatal(http.ListenAndServe(":3129", proxy))
}

func setCA(caCert, caKey string) error {
	goproxyCa, err := tls.LoadX509KeyPair(caCert, caKey)
	if err != nil {
		return err
	}
	if goproxyCa.Leaf, err = x509.ParseCertificate(goproxyCa.Certificate[0]); err != nil {
		return err
	}
	goproxy.GoproxyCa = goproxyCa
	goproxy.OkConnect = &goproxy.ConnectAction{Action: goproxy.ConnectAccept, TLSConfig: goproxy.TLSConfigFromCA(&goproxyCa)}
	goproxy.MitmConnect = &goproxy.ConnectAction{Action: goproxy.ConnectMitm, TLSConfig: goproxy.TLSConfigFromCA(&goproxyCa)}
	goproxy.HTTPMitmConnect = &goproxy.ConnectAction{Action: goproxy.ConnectHTTPMitm, TLSConfig: goproxy.TLSConfigFromCA(&goproxyCa)}
	goproxy.RejectConnect = &goproxy.ConnectAction{Action: goproxy.ConnectReject, TLSConfig: goproxy.TLSConfigFromCA(&goproxyCa)}
	return nil
}
