package main

import (
	"OAuth2/handlers"
	"net/http"

	"github.com/fatih/color"
)

const htmlIndex = `
<a href="/facebook">Facebook Login</a>
<a href="/github">Github Login</a>
<a href="/linkedin">LinkedIn Login</a>
<a href="/spotify">Spotify Login</a>
<a href="/azure">Azure Login</a>
`

func handleMain(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(htmlIndex))
}

func main() {
	http.HandleFunc("/", handleMain)
	http.HandleFunc("/facebook", handlers.HandleFacebookLogin)
	http.HandleFunc("/facebook/callback", handlers.HandleFacebookCallback)
	http.HandleFunc("/github", handlers.HandleGitLogin)
	http.HandleFunc("/github/callback", handlers.HandleGitCallback)
	http.HandleFunc("/linkedin", handlers.HandleLinkedLogin)
	http.HandleFunc("/linkedin/callback", handlers.HandleLinkedinCallback)
	http.HandleFunc("/spotify", handlers.HandleSpotifyLogin)
	http.HandleFunc("/spotify/callback", handlers.HandleSpotifyCallback)
	http.HandleFunc("/azure", handlers.HandleAzureLogin)
	http.HandleFunc("/azure/callback", handlers.HandleAzureCallback)

	color.Cyan("Launching Oauth2 Server")
	color.Magenta("Azure in testing")
	http.ListenAndServe(":8080", nil)

}
