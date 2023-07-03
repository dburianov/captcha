package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// go build -ldflags "-X main.appVersion=$HOSTNAME-$(date +%Y/%m/%d-%H:%M:%S%z)"
var appVersion string

func main() {
	welcomeString := "Starting SoftEther exporter Namecheap edition version: " + appVersion
	listenOn := "Listen on: " + bindAddress
	log.Println(welcomeString)
	log.Println(listenOn)

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", homeLink)
	router.HandleFunc("/api/v0/token/check/{token}/{hash}/{clientip}/{ttl}", tokenCheck).Methods("GET")
	router.HandleFunc("/api/v0/token/new/{clientip}", newToken).Methods("GET")
	router.HandleFunc("/api/v0/captcha/text", captchaTextHandle).Methods("GET")
	router.HandleFunc("/api/v0/captcha/math", captchaMathHandle).Methods("GET")
	router.HandleFunc("/api/v0/captcha/check", captchaCheckHandler).Methods("POST")
	router.HandleFunc("/api/v0/captcha", captchaIndexHandle).Methods("GET")

	log.Println(http.ListenAndServe(bindAddress, router))
}

func homeLink(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome home!")
}
