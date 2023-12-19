package main

import (
	"fmt"
	"github.com/coooold/simpleotp_go/internal"
	"log"
	"net/http"
	"time"
)



func startHttpServer() {
	http.HandleFunc("/totp/check", internal.CheckAuthHandler)
	http.HandleFunc("/totp/login", internal.LoginHandler)

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", internal.Port),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	log.Printf("Server is listening on port %d", internal.Port)
	log.Fatal(server.ListenAndServe())
}

func main() {
	internal.InitParams()
	startHttpServer()
}
