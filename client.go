package main

import (
	"io/ioutil"
	"net/http"
	"net/url"
)

// GoogleOauthClient is oauth client for Google
type GoogleOauthClient struct {
	authEndpoint  string
	tokenEndoiint string
	clientID      string
	clientSecret  string
	scope         string
	redirectURI   string
	responseType  string
	grantType     string
}

// GetValueForAuth return url.Values for OAuth 2.0 endpoint
func (c GoogleOauthClient) GetValueForAuth() url.Values {
	return url.Values{
		"client_id":     {c.clientID},
		"redirect_uri":  {c.redirectURI},
		"scope":         {c.scope},
		"response_type": {c.responseType},
		"access_type":   {"offline"},
		//sample のため state は指定しない
	}
}

// GetValueForToken return url.Values for token endpoint
func (c GoogleOauthClient) GetValueForToken(code string) url.Values {
	return url.Values{
		"code":          {code},
		"client_id":     {c.clientID},
		"client_secret": {c.clientSecret},
		"redirect_uri":  {c.redirectURI},
		"grant_type":    {c.grantType},
	}
}

// Exchange exchange an authorization code for an access token
func (c GoogleOauthClient) Exchange(code string) (*Token, error) {
	v := c.GetValueForToken(code)

	resp, err := http.PostForm(c.tokenEndoiint, v)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return ParseToken(data)
}

// NewClient return GoogleOauthClient instance
func NewClient(
	clientID,
	clientSecret,
	scope,
	redirectURI string,
) GoogleOauthClient {
	return GoogleOauthClient{
		authEndpoint:  "https://accounts.google.com/o/oauth2/v2/auth",
		tokenEndoiint: "https://www.googleapis.com/oauth2/v4/token",
		clientID:      clientID,
		clientSecret:  clientSecret,
		scope:         scope,
		redirectURI:   redirectURI,
		responseType:  "code",
		grantType:     "authorization_code",
	}
}
