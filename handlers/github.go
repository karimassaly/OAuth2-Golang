package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"golang.org/x/oauth2"
)

var (
	GitOauthConf = &oauth2.Config{
		RedirectURL:  "http://localhost:8080/github/callback",
		ClientID:     "",
		ClientSecret: "",
		Scopes:       []string{"user:email"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://github.com/login/oauth/authorize",
			TokenURL: "https://github.com/login/oauth/access_token",
		},
	}
)

func HandleGitLogin(w http.ResponseWriter, r *http.Request) {
	url := GitOauthConf.AuthCodeURL(oauthStateString)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func HandleGitCallback(w http.ResponseWriter, r *http.Request) {
	state := r.FormValue("state")
	if state != oauthStateString {
		fmt.Printf("invalid oauth state, expected %q got %q", oauthStateString, state)
		return
	}

	code := r.FormValue("code")
	fmt.Printf("this is the code %s\n", code)
	token, err := GitOauthConf.Exchange(oauth2.NoContext, code)
	if err != nil {
		fmt.Printf("GitOauthConf.Exchange() failed with %q", err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	fmt.Printf("this is the token %s\n", token.AccessToken)
	var bearer = "Bearer " + token.AccessToken
	req, err := http.NewRequest("GET", "https://api.github.com/user", nil)
	req.Header.Add("Authorization", bearer)
	client := &http.Client{}
	info, err := client.Do(req)
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
	fmt.Println("get user profile", m["login"])
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}
