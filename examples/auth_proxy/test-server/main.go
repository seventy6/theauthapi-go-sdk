package main

import (
	"fmt"
	"log"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "X-Auth-Response header:", r.Header.Get("X-Auth-Response"))
	fmt.Fprintln(w, "X-Auth-Data header:", r.Header.Get("X-Auth-Data"))
	fmt.Fprintln(w, "response")
}

func main() {
	http.HandleFunc("/", handler)

	fmt.Printf("Starting server at port 5000\n")
	log.Fatal(http.ListenAndServe(":5000", nil))
}
