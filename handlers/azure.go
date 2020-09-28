package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/davecgh/go-spew/spew"
	"golang.org/x/oauth2"
)

var (
	oauthConfAzure = &oauth2.Config{
		ClientID:     "",
		ClientSecret: "",
		RedirectURL:  "http://localhost:8080/azure/callback",
		Scopes:       []string{"Mail.Read"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://login.microsoftonline.com/common/oauth2/v2.0/authorize",
			TokenURL: "https://login.microsoftonline.com/common/oauth2/v2.0/token",
		},
	}
	oauthStateStringAz = "thisshouldberandom"
)

func HandleAzureLogin(w http.ResponseWriter, r *http.Request) {
	u := oauthConfAzure.AuthCodeURL(oauthStateStringAz)
	http.Redirect(w, r, u, http.StatusTemporaryRedirect)
}

func HandleAzureCallback(w http.ResponseWriter, r *http.Request) {
	state := r.FormValue("state")
	if state != oauthStateStringAz {
		fmt.Printf("invalid oauth state, expected %q got %q", oauthStateStringAz, state)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	code := r.FormValue("code")

	token, err := oauthConfAzure.Exchange(oauth2.NoContext, code)
	if err != nil {
		fmt.Printf("oauthConfAzure.Exchange() failed with %q", err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	bearer := "Bearer " + token.AccessToken
	req, err := http.NewRequest("GET", "https://graph.microsoft.com/beta/me/messages?$select=sender,subject&$top=5", nil) // Query the first 5 emails with the sender and the subject
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

	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}
