package main

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"
	"io"
	http "net/http"
	"time"
)

var oauthConfig oauth2.Config
var ctx = context.Background()
var provider *oidc.Provider

/* AuthSession is primary key, value is state */
var callbackData map[string]string

/* Primary key AuthSession, value is redirect address */
var callbackRedirect map[string]string

// Primary key Session, value is UserInfo struct.
var userData map[string]UserInfo

// UserInfo is a Struct of Keycloak / OpenID User Information standard data
type UserInfo struct {
	Sub               string   `json:"sub"`
	Groups            []string `json:"groups"`
	PreferredUsername string   `json:"preferred_username"`
	GivenName         string   `json:"given_name"`
	FamilyName        string   `json:"family_name"`
	Email             string   `json:"email"`
}

// Initialize authentication functions and processes
func initAuth() {
	callbackData = make(map[string]string)
	callbackRedirect = make(map[string]string)
	userData = make(map[string]UserInfo)
	provider, _ = oidc.NewProvider(ctx, Config.Issuer)

	oauthConfig = oauth2.Config{
		ClientID:     Config.ClientID,
		ClientSecret: Config.ClientSecret,
		RedirectURL:  "http://" + Config.Address + "/redirect",

		// Discovery returns the OAuth2 endpoints.
		Endpoint: provider.Endpoint(),

		// "openid" is a required scope for OpenID Connect flows.
		Scopes: []string{oidc.ScopeOpenID},
	}

	Mux.HandleFunc("/redirect", getRedirect)
	Mux.HandleFunc("/login", HttpLogger(getLogin))

}

func getLogin(w http.ResponseWriter, r *http.Request) {
	println("hi!")
	redirect := r.URL.Query().Get("redir")
	//generate secure random tags for state auth
	state, _ := randString(16)
	sess, _ := randString(8)

	//save state authentication data
	callbackData[sess] = state
	callbackRedirect[sess] = redirect

	setTimedCookie(w, r, "AuthSession", sess, int(time.Minute.Seconds()*5))
	//redirect to SSO
	http.Redirect(w, r, oauthConfig.AuthCodeURL(state), http.StatusFound)
}

/*Completes Authentication*/
func getRedirect(w http.ResponseWriter, r *http.Request) {
	//see if data is valid
	sess, err := r.Cookie("AuthSession")
	if err != nil {
		http.Error(w, "session not found", http.StatusBadRequest)
		return
	}
	sessId := sess.Value

	state := callbackData[sessId]
	if r.URL.Query().Get("state") != state {
		http.Error(w, "state did not match", http.StatusBadRequest)
		return
	}
	delete(callbackData, sessId)
	sess.Expires = time.Now().Add(-time.Minute * 5)
	sess.MaxAge = -1
	http.SetCookie(w, sess)

	oauth2Token, err := oauthConfig.Exchange(ctx, r.URL.Query().Get("code"))
	if err != nil {
		http.Error(w, "Token failed", http.StatusInternalServerError)
		return
	}

	userInfo, err := getUserInfo(oauth2Token)
	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, "Getting user info failed", http.StatusInternalServerError)
	}

	//create new Session cookie
	session, _ := randString(16)
	setTimedCookie(w, r, "Session", session, int(time.Hour.Seconds()*48))
	userData[session] = userInfo

	redir := callbackRedirect[sessId]
	if redir == "" {
		redir = "http://" + r.Host
	}

	if userInfo.PreferredUsername != "" {
		fmt.Println("User logged in:", userInfo.PreferredUsername)
		http.Redirect(w, r, redir, http.StatusFound)
	} else {
		http.Redirect(w, r, "https://csh.rit.edu", http.StatusFound)
	}

}

// getUserInfo returns UserInfo as described here https://openid.net/specs/openid-connect-core-1_0.html#UserInfo.
//
// This requires an OAuth2 Token, and uses the Config Issuer as the query destination.
func getUserInfo(token *oauth2.Token) (UserInfo, error) {
	req, err := http.NewRequest("GET", Config.Issuer+"/protocol/openid-connect/userinfo", nil)
	token.SetAuthHeader(req)

	//do http client request stuff
	res, err := http.DefaultClient.Do(req)
	jsonRaw, err := io.ReadAll(res.Body)
	if err != nil {
		return UserInfo{}, fmt.Errorf("reading JSON response failed, %s", err.Error())
	}

	var userInfo UserInfo
	err = json.Unmarshal(jsonRaw, &userInfo)
	if err != nil {
		return UserInfo{}, err
	}
	return userInfo, nil

}

// sets a cookie with an expiration time.
func setTimedCookie(w http.ResponseWriter, r *http.Request, name, value string, time int) {
	c := &http.Cookie{
		Name:     name,
		Value:    value,
		MaxAge:   time,
		Secure:   r.TLS != nil,
		SameSite: 3,
		HttpOnly: true,
	}
	http.SetCookie(w, c)
}

// Generates a random string using crypto rand for increased security
func randString(nByte int) (string, error) {
	b := make([]byte, nByte)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

// IsAuthenticated checks if their session cookie is valid and therefore if they have authentication.
func IsAuthenticated(w http.ResponseWriter, r *http.Request) bool {
	ses, err := r.Cookie("Session")
	if errors.Is(err, http.ErrNoCookie) {
		println("can't find a cookie")
		return false
	}
	_, ok := userData[ses.Value]
	if !ok {
		return false
	}
	ses.MaxAge = int(time.Hour.Seconds() * 24 * 7)
	http.SetCookie(w, ses)
	return true
}

// GetUser returns UserInfo for a user, given their session cookie.
func GetUser(r *http.Request) (UserInfo, error) {
	ses, err := r.Cookie("Session")
	if errors.Is(err, http.ErrNoCookie) {
		return UserInfo{}, errors.New("cookie not found")
	}
	session := ses.Value
	data, ok := userData[session]
	if ok {
		return data, nil
	}
	return UserInfo{}, nil

}
