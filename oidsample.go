package main

import (
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func entryForOid(client GoogleOauthClient, w http.ResponseWriter, r *http.Request) {
	v := client.GetValueForAuth()

	url, err := url.Parse(client.authEndpoint + "?" + v.Encode())
	if err != nil {
		log.Fatal(err)
	}
	http.Redirect(w, r, url.String(), http.StatusFound)
}

func redirectPointForOid(client GoogleOauthClient, w http.ResponseWriter, r *http.Request) {

	q := r.URL.Query()
	code := q.Get("code")
	if code == "" {
		log.Fatal("no code")
		return
	}
	// ここでstateの検証が必要
	token, err := client.ExchangeForOid(code)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(
		token.AccessToken,
		token.RefreshToken,
		token.TokenType,
		token.ExpiresIn,
		token.IDToken,
	)

	w.WriteHeader(http.StatusOK)

	fmt.Fprintln(w, token.IDToken)

	for _, tkn := range strings.Split(token.IDToken, ".") {
		str, err := base64.RawURLEncoding.DecodeString(tkn)
		if err != nil {
			log.Println(err)
		} else {
			log.Println(string(str))
			w.Write(str)
		}
	}
}

func oidSample() {
	client := NewClient(
		os.Getenv("client_id"),
		os.Getenv("client_secret"),
		"openid email profile",
		os.Getenv("redirect_uri"),
	)
	var httpServer http.Server
	http.HandleFunc("/auth", func(w http.ResponseWriter, r *http.Request) {
		entryForOid(client, w, r)
	})
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		redirectPointForOid(client, w, r)
	})

	log.Println("start http listen :18888")
	httpServer.Addr = ":8080"

	//http
	log.Println(httpServer.ListenAndServe())
}
