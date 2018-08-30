package sso

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"github.com/antihax/goesi"
	"github.com/gorilla/sessions"
	"log"
	"net/http"
)

var Client *http.Client
var SSOAuthenticator *goesi.SSOAuthenticator
var store = sessions.NewCookieStore([]byte("something-very-secret"))
var scopes = []string{
	"publicData",
	"esi-location.read_location.v1",
	"esi-location.read_ship_type.v1",
	"esi-characters.read_notifications.v1",
}

func Load(id string, key string) {

	SSOAuthenticator = goesi.NewSSOAuthenticator(Client, id, key, "http://localhost:8080/callback", scopes)
	if SSOAuthenticator == nil {
		log.Fatal("Failed to create SSO Authenticator")
	}
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "<a href=\"/login\">login to fifth-bot</a>")
	})
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/callback", callbackHandler)
	log.Printf(`SSO server running`)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	s, _ := store.Get(r, "session-name")

	// Generate a random state string
	b := make([]byte, 16)
	rand.Read(b)
	state := base64.URLEncoding.EncodeToString(b)

	// Save the state on the session
	s.Values["state"] = state
	err := s.Save(r, w)
	if err != nil {
		w.WriteHeader(500)
		return
		//return http.StatusInternalServerError, err
	}

	// Generate the SSO URL with the state string
	url := SSOAuthenticator.AuthorizeURL(state, true, scopes)

	// Send the user to the URL
	http.Redirect(w, r, url, 302)
	return
	//return http.StatusMovedPermanently, nil

}

func callbackHandler(w http.ResponseWriter, r *http.Request) {

	s, _ := store.Get(r, "session-name")

	// get our code and state
	code := r.FormValue("code")
	state := r.FormValue("state")

	// Verify the state matches our randomly generated string from earlier.
	if s.Values["state"] != state {
		//return http.StatusInternalServerError, errors.New("Invalid State.")
		return
	}

	// Exchange the code for an Access and Refresh token.
	token, err := SSOAuthenticator.TokenExchange(code)
	if err != nil {
		//return http.StatusInternalServerError, err
		return
	}

	//token.

	// Obtain a token source (automaticlly pulls refresh as needed)
	tokSrc, err := SSOAuthenticator.TokenSource(token)
	if err != nil {
		return
		//return http.StatusInternalServerError, err
	}

	tokSrc.Token()

	// Assign an auth context to the calls
	//auth := context.WithValue(context.TODO(), goesi.ContextOAuth2, tokSrc.Token)

	// Verify the client (returns clientID)
	v, err := SSOAuthenticator.Verify(tokSrc)
	if err != nil {
		return
		//return http.StatusInternalServerError, err
	}

	if err != nil {
		return
		//return http.StatusInternalServerError, err
	}

	// Save the verification structure on the session for quick access.
	s.Values["character"] = v
	err = s.Save(r, w)
	if err != nil {
		return
		//return http.StatusInternalServerError, err
	}

	// Redirect to the account page.
	//http.Redirect(w, r, "/account", 302)
	return
	//return http.StatusMovedPermanently, nil
}
