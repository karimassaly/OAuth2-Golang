package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"golang.org/x/oauth2"
)

var (
	oauthConfFB = &oauth2.Config{
		ClientID:     "",
		ClientSecret: "",
		RedirectURL:  "http://localhost:8080/facebook/callback",
		Scopes:       []string{"public_profile"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://www.facebook.com/v5.0/dialog/oauth",
			TokenURL: "https://graph.facebook.com/v5.0/oauth/access_token",
		},
	}
	oauthStateStringFB = "thisshouldberandom"
)

func HandleFacebookLogin(w http.ResponseWriter, r *http.Request) {
	fmt.Println(oauthConfFB)
	u := oauthConfFB.AuthCodeURL(oauthStateStringFB)
	http.Redirect(w, r, u, http.StatusTemporaryRedirect)
}

func HandleFacebookCallback(w http.ResponseWriter, r *http.Request) {
	state := r.FormValue("state")
	if state != oauthStateStringFB {
		fmt.Printf("invalid oauth state, expected %q got %q", oauthStateStringFB, state)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	code := r.FormValue("code")

	token, err := oauthConfFB.Exchange(oauth2.NoContext, code)
	if err != nil {
		fmt.Printf("oauthConfFB.Exchange() failed with %q", err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	info, err := http.Get(fmt.Sprintf("https://graph.facebook.com/me?fields=name,middle_name,first_name,last_name,email,address,age_range,gender&access_token=%s", url.QueryEscape(token.AccessToken)))
	if err != nil {
		fmt.Printf("Get: %q", err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	defer info.Body.Close()

	response, err := ioutil.ReadAll(info.Body)
	if err != nil {
		fmt.Printf("ReadAll: %q", err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	var m map[string]interface{}
	if err := json.Unmarshal(response, &m); err != nil {
		fmt.Printf("error unmarshalling response: %s", err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	}
	{
		fmt.Println("got response body", m)
		resp, err := http.Get(fmt.Sprintf("https://graph.facebook.com/v5.0/%s/picture?redirect=0&access_token=%s", m["id"], url.QueryEscape(token.AccessToken)))

		if err != nil {

			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		defer resp.Body.Close()

		var m map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&m); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		fmt.Println("get user profile", m)
	}
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}
