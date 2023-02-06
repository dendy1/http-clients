package httpclients

import (
	"context"
	"fmt"
	"golang.org/x/oauth2"
	"log"
	"net/http"
	"sync"
	"time"
)

var (
	once   sync.Once
	client *http.Client
)

func HttpClient() *http.Client {
	once.Do(func() {
		client = &http.Client{
			Timeout: time.Second * 30,
			Transport: &http.Transport{
				TLSHandshakeTimeout: 10 * time.Second,
			},
		}
	})
	return client
}

var (
	oauth2once   sync.Once
	oauth2client *http.Client
)

func OAuth2Client() *http.Client {
	oauth2once.Do(func() {
		ctx := context.Background()

		cfg := &oauth2.Config{
			ClientID:     "YOUR_CLIENT_ID",
			ClientSecret: "YOUR_CLIENT_SECRET",
			Scopes:       []string{"SCOPE1", "SCOPE2"},
			Endpoint: oauth2.Endpoint{
				AuthURL:  "https://provider.com/o/oauth2/auth",
				TokenURL: "https://provider.com/o/oauth2/token",
			},
		}
		// Redirect user to consent page to ask for permission
		// for the scopes specified above.
		url := cfg.AuthCodeURL("state", oauth2.AccessTypeOffline)
		fmt.Printf("Visit the URL for the auth dialog: %v", url)

		// Use the authorization code that is pushed to the redirect
		// URL. Exchange will do the handshake to retrieve the
		// initial access token. The HTTP Client returned by
		// conf.Client will refresh the token as necessary.
		var code string
		if _, err := fmt.Scan(&code); err != nil {
			log.Fatal(err)
		}
		tok, err := cfg.Exchange(ctx, code)
		if err != nil {
			log.Fatal(err)
		}

		oauth2client = cfg.Client(ctx, tok)
	})
	return oauth2client
}
