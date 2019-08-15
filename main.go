package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
)

var (
	endpoint      = "https://accounts.google.com/o/oauth2/v2/auth"
	tokenEndoiint = "https://www.googleapis.com/oauth2/v4/token"

	clientID     = os.Getenv("client_id")
	clientSecret = os.Getenv("client_secret")
	scope        = os.Getenv("scope")
	redirectURI  = os.Getenv("redirect_uri")

	responseType = "code"

	grantType = "authorization_code"
)

func entry(w http.ResponseWriter, r *http.Request) {
	v := url.Values{
		"client_id":     {clientID},
		"redirect_uri":  {redirectURI},
		"scope":         {scope},
		"response_type": {responseType},
	}
	url, err := url.Parse(endpoint + "?" + v.Encode())
	if err != nil {
		log.Fatal(err)
	}
	http.Redirect(w, r, url.String(), http.StatusFound)
}

func redirectPoint(w http.ResponseWriter, r *http.Request) {

	q := r.URL.Query()
	code := q.Get("code")
	if code == "" {
		log.Fatal("no code")
		return
	}
	// qscope := q.Get("scope")

	v := url.Values{
		"code":          {code},
		"client_id":     {clientID},
		"client_secret": {clientSecret},
		"redirect_uri":  {redirectURI},
		// "scope":         {scope},
		// "response_type": {responseType},
		"grant_type": {grantType},
	}

	resp, err := http.PostForm(tokenEndoiint, v)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	dump, _ := httputil.DumpResponse(resp, true)

	log.Println(string(dump))

	w.WriteHeader(http.StatusOK)
	io.Copy(w, resp.Body)
}

func token(w http.ResponseWriter, r *http.Request) {
	dump, _ := httputil.DumpRequest(r, true)
	r.ParseForm() // PostForm で値を取得するために必要
	fmt.Println(string(dump))
	w.WriteHeader(http.StatusOK)
}

func main() {
	var httpServer http.Server
	http.HandleFunc("/auth", entry)
	http.HandleFunc("/", redirectPoint)
	http.HandleFunc("/token", token)

	log.Println("start http listen :18888")
	httpServer.Addr = ":8080"

	//http
	log.Println(httpServer.ListenAndServe())
}
