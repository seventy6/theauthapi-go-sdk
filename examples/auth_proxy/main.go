package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"

	"github.com/seventy6/theauthapi-go-sdk"
)

// CustomProxy handles the proxying of requests with header insertion
type CustomProxy struct {
	targetHost string
	proxy      *httputil.ReverseProxy

	originalDirector func(*http.Request)
	authAPIClient    *theauthapi.Client
}

// NewCustomProxy creates a new proxy instance
func NewCustomProxy(targetHost string, accessToken string) (*CustomProxy, error) {
	targetURL, err := url.Parse(targetHost)
	if err != nil {
		return nil, err
	}

	proxy := httputil.NewSingleHostReverseProxy(targetURL)

	authAPIClient := theauthapi.NewClient(
		theauthapi.WithAccessToken(accessToken),
		theauthapi.WithDebug(),
	)

	cp := &CustomProxy{
		targetHost: targetHost,
		proxy:      proxy,

		originalDirector: proxy.Director,
		authAPIClient:    authAPIClient,
	}

	// Modify the default director to add our custom headers
	proxy.Director = cp.Director

	return cp, nil
}

func (p *CustomProxy) Director(req *http.Request) {
	// keep original director
	p.originalDirector(req)

	apiKey := req.Header.Get("x-api-key")
	if apiKey == "" {
		log.Println("api key not found in request")
		return
	}

	// do auth request
	data, err := p.authAPIClient.ApiKeys.GetValidKey(req.Context(), apiKey)
	if err != nil {
		log.Println("key not valid")
		return
	}

	// Add custom headers - this can be ex
	req.Header.Set("x-key-name", data.Name)

	// Optionally log the modified request
	log.Printf("Proxying request to: %s with headers: %v\n", req.URL, req.Header)
}

// ServeHTTP implements the http.Handler interface
func (p *CustomProxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	p.proxy.ServeHTTP(w, r)
}

func main() {
	// Example target host - replace with your actual target
	targetHost := "http://localhost:8000"

	accessToken := os.Getenv("ACCESS_TOKEN")
	proxy, err := NewCustomProxy(targetHost, accessToken)
	if err != nil {
		log.Fatal("Error creating proxy: ", err)
	}

	// Start the proxy server
	server := &http.Server{
		Addr:    ":3000",
		Handler: proxy,
	}

	log.Printf("Starting proxy server on :8080, targeting %s\n", targetHost)
	if err := server.ListenAndServe(); err != nil {
		log.Fatal("Proxy server error: ", err)
	}
}
