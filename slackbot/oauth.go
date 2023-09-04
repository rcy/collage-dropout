package slackbot

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"os"

	slackapi "github.com/slack-go/slack"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/slack"
)

func Serve() {
	// Define the Slack app's client ID, client secret, and redirect URI
	clientID := os.Getenv("SLACK_CLIENT_ID")
	clientSecret := os.Getenv("SLACK_CLIENT_SECRET")
	redirectURI := os.Getenv("SLACK_REDIRECT_URI")

	// Set up the OAuth 2.0 config for Slack
	config := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Endpoint:     slack.Endpoint,
		RedirectURL:  redirectURI,
		Scopes:       []string{"channels:read"},
	}

	http.HandleFunc("/slack/oauth", func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get("code")

		values := url.Values{
			"code":          []string{code},
			"client_id":     []string{clientID},
			"client_secret": []string{clientSecret},
			"redirect_uri":  []string{redirectURI},
			"grant_type":    []string{"authorization_code"},
		}

		resp, err := http.PostForm("https://slack.com/api/oauth.v2.access", values)
		if err != nil {
			log.Printf("Error making token request: %v\n", err)
			http.Error(w, "Error making token request", http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			log.Printf("Unexpected status code: %d\n", resp.StatusCode)
			http.Error(w, "Unexpected status code", http.StatusInternalServerError)
			return
		}

		var tokenResponse struct {
			AccessToken string `json:"access_token"`
		}

		err = json.NewDecoder(resp.Body).Decode(&tokenResponse)
		if err != nil {
			log.Printf("Error decoding JSON response: %v\n", err)
			http.Error(w, "Error decoding JSON response", http.StatusInternalServerError)
			return
		}

		if tokenResponse.AccessToken == "" {
			log.Println("Server response missing access_token")
			http.Error(w, "Server response missing access_token", http.StatusInternalServerError)
			return
		}

		// Use the token to call Slack API methods with the appropriate scopes
		api := slackapi.New(tokenResponse.AccessToken)
		log.Printf("api: %v\n", api)

		// Redirect the user to a success page or do something else
		http.Redirect(w, r, "/success", http.StatusTemporaryRedirect)
	})
	// http.HandleFunc("/slack/oauth", func(w http.ResponseWriter, r *http.Request) {
	// 	code := r.URL.Query().Get("code")
	// 	token, err := config.Exchange(context.Background(), code, oauth2.SetAuthURLParam("client_id", clientID), oauth2.SetAuthURLParam("client_secret", clientSecret))
	// 	if err != nil {
	// 		log.Printf("Error exchanging code: %v\n", err)
	// 		http.Error(w, "Error exchanging code", http.StatusInternalServerError)
	// 		return
	// 	}

	// 	// Use the token to call Slack API methods with the appropriate scopes
	// 	api := slackapi.New(token.AccessToken)
	// 	log.Printf("api: %v\n", api)

	// 	// Redirect the user to a success page or do something else
	// 	http.Redirect(w, r, "/success", http.StatusTemporaryRedirect)
	// })

	// // Set up the HTTP request handler to handle the Slack OAuth callback
	// http.HandleFunc("/slack/oauth", func(w http.ResponseWriter, r *http.Request) {
	// 	code := r.URL.Query().Get("code")
	// 	token, err := config.Exchange(context.Background(), code)
	// 	if err != nil {
	// 		log.Printf("Error exchanging code: %v\n", err)
	// 		http.Error(w, "Error exchanging code", http.StatusInternalServerError)
	// 		return
	// 	}

	// 	// Use the token to call Slack API methods with the appropriate scopes
	// 	api := slackapi.New(token.AccessToken)
	// 	log.Printf("api: %v\n", api)
	// 	// ...

	// 	// Redirect the user to a success page or do something else
	// 	http.Redirect(w, r, "/success", http.StatusTemporaryRedirect)
	// })

	// Set up the HTTP request handler for the root path
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Redirect the user to the Slack authorization page
		authURL := config.AuthCodeURL("state", oauth2.AccessTypeOffline)
		http.Redirect(w, r, authURL, http.StatusTemporaryRedirect)
	})

	http.HandleFunc("/success", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("success"))
	})

	// Start the server on port 8080
	log.Fatal(http.ListenAndServe(":8080", nil))
}
