package main

import (
	"log"
	"net/http"
	"net/url"
	"os"
)

func entry(client GoogleOauthClient, w http.ResponseWriter, r *http.Request) {
	v := client.GetValueForAuth()

	url, err := url.Parse(client.authEndpoint + "?" + v.Encode())
	if err != nil {
		log.Fatal(err)
	}
	http.Redirect(w, r, url.String(), http.StatusFound)
}

func redirectPoint(client GoogleOauthClient, w http.ResponseWriter, r *http.Request) {

	q := r.URL.Query()
	code := q.Get("code")
	if code == "" {
		log.Fatal("no code")
		return
	}
	// qscope := q.Get("scope")
	token, err := client.Exchange(code)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(
		token.AccessToken,
		token.RefreshToken,
		token.TokenType,
		token.ExpiresIn,
	)

	w.WriteHeader(http.StatusOK)
}

func oauthSample() {
	client := NewClient(
		os.Getenv("client_id"),
		os.Getenv("client_secret"),
		os.Getenv("scope"),
		os.Getenv("redirect_uri"),
	)
	var httpServer http.Server
	http.HandleFunc("/auth", func(w http.ResponseWriter, r *http.Request) {
		entry(client, w, r)
	})
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		redirectPoint(client, w, r)
	})

	log.Println("start http listen :18888")
	httpServer.Addr = ":8080"

	//http
	log.Println(httpServer.ListenAndServe())
}
