package sso

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/gob"
	"fmt"
	"github.com/antihax/goesi"
	"github.com/bwmarrin/discordgo"
	"github.com/eddbc/fifth-bot/db"
	"github.com/gorilla/sessions"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/oauth2"
	"log"
	"net/http"
)

var Client *http.Client
var Session *discordgo.Session
var SSOAuthenticator *goesi.SSOAuthenticator
var store = sessions.NewCookieStore([]byte("something-very-secret"))
var scopes = []string{
	"publicData",
	"esi-location.read_location.v1",
	"esi-location.read_ship_type.v1",
}

func Load(id string, key string) {

	gob.Register(goesi.VerifyResponse{})

	SSOAuthenticator = goesi.NewSSOAuthenticatorV2(Client, id, key, "http://localhost:8080/callback", scopes)
	if SSOAuthenticator == nil {
		log.Fatal("Failed to create SSO Authenticator")
	}
	http.HandleFunc("/", loginHandler)
	http.HandleFunc("/callback", callbackHandler)
	log.Printf(`SSO server running`)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	s, _ := store.Get(r, "session-name")

	discordId := r.URL.Query()["id"]
	if discordId == nil {
		return
	}

	s.Values["id"] = discordId[0]

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

	// Obtain a token source (automaticlly pulls refresh as needed)
	tokSrc := SSOAuthenticator.TokenSource(token)
	if err != nil {
		return
		//return http.StatusInternalServerError, err
	}

	//tokSrc.Token()

	// Assign an auth context to the calls
	//auth := context.WithValue(context.TODO(), goesi.ContextOAuth2, tokSrc.Token)

	// Verify the client (returns clientID)
	v, err := SSOAuthenticator.Verify(tokSrc)
	if err != nil {
		return
		//return http.StatusInternalServerError, err
	}

	// Save the verification structure on the session for quick access.
	s.Values["character"] = v
	err = s.Save(r, w)
	if err != nil {
		log.Print(err)
		return
		//return http.StatusInternalServerError, err
	}

	count, err := db.Db().Collection("tokens").CountDocuments(context.TODO(), bson.D{
		{"characterid",v.CharacterID},
	})
	if count == 0 {
		_,err = db.Db().Collection("tokens").InsertOne(context.TODO(), UserToken{v.CharacterID,token})
		if err != nil {
			log.Print(err)
		}
	} else {
		_,err = db.Db().Collection("tokens").UpdateOne(context.TODO(), bson.D{
			{"characterid",v.CharacterID},
		}, bson.D{{
			"$set", bson.D{{
				"token", token,
			}},
		}})
		if err != nil {
			log.Print(err)
		}
	}

	discordId := s.Values["id"].(string)

	_,err = db.Db().Collection("discord").InsertOne(context.TODO(), DiscordUser{v.CharacterID, discordId})
	if err != nil {
		log.Print(err)
	}

	usr, _ := Session.User(discordId)

	fmt.Fprintf(w, "discord user %+v signed in as eve user %v", usr.Username, v.CharacterName)

	// Redirect to the account page.
	//http.Redirect(w, r, "/account", 302)
	return
	//return http.StatusMovedPermanently, nil
}

func GetTokenForDiscordUser(userID string) (ut UserToken, err error) {
	var du DiscordUser

	err = db.Db().Collection("discord").FindOne(context.TODO(), bson.D{{"discordid", userID}}).Decode(&du)
	if err == nil {
		err = db.Db().Collection("tokens").FindOne(context.TODO(), bson.D{{"characterid", du.CharacterId}}).Decode(&ut)
	}

	if err != nil {
		log.Print(err)
		return
	}

	return
}

func AuthContext(ut UserToken) (c context.Context) {
	src := SSOAuthenticator.TokenSource(ut.Token)
	c = context.WithValue(context.Background(), goesi.ContextOAuth2, src)
	return
}

func UpdateToken(char int32, tkn *oauth2.Token) (err error) {
	_,err = db.Db().Collection("tokens").UpdateOne(context.TODO(), bson.D{
		{"characterid",char},
	}, bson.D{{
		"$set", bson.D{{
			"token", tkn,
		}},
	}})
	return
}

type UserToken struct {
	CharacterId int32 `json:"characterid"`
	Token *oauth2.Token `json:"token"`
}
type DiscordUser struct {
	CharacterId int32  `json:"characterid"`
	DiscordId string  `json:"discordid"`
}
