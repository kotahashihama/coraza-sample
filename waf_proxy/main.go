package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"

	coraza "github.com/corazawaf/coraza/v3"
	chttp "github.com/corazawaf/coraza/v3/http"
)

func main() {
	waf, err := coraza.NewWAF(
		coraza.NewWAFConfig().
			WithDirectivesFromFile("coraza.conf").
			WithDirectivesFromFile("coreruleset/crs-setup.conf.example").
			WithDirectivesFromFile("coreruleset/rules/*.conf"),
	)
	if err != nil {
		log.Fatalf("Failed to create WAF: %v", err)
	}

	targetURL, _ := url.Parse("http://localhost:8080")
	proxy := httputil.NewSingleHostReverseProxy(targetURL)
	wafProtectedProxy := chttp.WrapHandler(waf, proxy)

	log.Println("WAF Proxy running on http://localhost:8081")
	if err := http.ListenAndServe(":8081", wafProtectedProxy); err != nil {
		log.Fatal(err)
	}
}
