package main

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strconv"

	"github.com/matzhouse/theauthapi-go-sdk"
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

	log.Println("found API key", apiKey)

	// do auth request
	data, err := p.authAPIClient.ApiKeys.GetValidKey(req.Context(), apiKey)
	if err != nil {
		log.Println("key not valid", err)
		return
	}

	headerData, err := json.Marshal(data)
	if err != nil {
		log.Println("data not valid")
		return
	}

	// Add custom headers - this can be ex
	req.Header.Del("X-Api-Key")
	req.Header.Set("x-auth-response", strconv.Itoa(data.HTTPResponse.StatusCode))
	req.Header.Set("x-auth-data", string(headerData))

	// Optionally log the modified request
	log.Printf("Proxying request to: %s with headers: %v\n", req.URL, req.Header)
}

// ServeHTTP implements the http.Handler interface
func (p *CustomProxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	p.proxy.ServeHTTP(w, r)
}

func main() {
	// Example target host - replace with your actual target
	targetHost := "http://localhost:5000"

	accessToken := os.Getenv("ACCESS_TOKEN")
	log.Println("access token:", accessToken)
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
