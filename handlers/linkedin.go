package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"golang.org/x/oauth2"
	"github.com/davecgh/go-spew/spew"
)

var (
	LinkedOauthConf = &oauth2.Config{
		RedirectURL:  "http://localhost:8080/linkedin/callback",
		ClientID:     "",
		ClientSecret: "",
		Scopes:       []string{"r_liteprofile"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://www.linkedin.com/oauth/v2/authorization",
			TokenURL: "https://www.linkedin.com/oauth/v2/accessToken",
		},
	}
)

func HandleLinkedLogin(w http.ResponseWriter, r *http.Request) {
	url := LinkedOauthConf.AuthCodeURL(oauthStateString)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func HandleLinkedinCallback(w http.ResponseWriter, r *http.Request) {
	state := r.FormValue("state")
	if state != oauthStateString {
		fmt.Printf("invalid oauth state, expected %q got %q", oauthStateString, state)
		return
	}

	code := r.FormValue("code")
	fmt.Printf("this is the code %s\n", code)
	token, err := LinkedOauthConf.Exchange(oauth2.NoContext, code)
	if err != nil {
		fmt.Printf("LinkedOauthConf.Exchange() failed with %q", err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	fmt.Printf("this is the token %s\n", token.AccessToken)
	var bearer = "Bearer " + token.AccessToken
	req, err := http.NewRequest("GET", "https://api.linkedin.com/v2/me", nil)
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
	spew.Dump(m)
	Fname := m["localizedFirstName"]
	Lname := m["localizedLastName"]
	fmt.Println(Fname, Lname)
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}
