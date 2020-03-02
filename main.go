package main

import (
	"OAuth2/handlers"
	"log"
	"net/http"
)

const htmlIndex = `
<a href="/facebook">Facebook Login</a>
<a href="/github">Github Login</a>
<a href="/linkedin">LinkedIn Login</a>
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
	log.Fatal(http.ListenAndServe(":8080", nil))

}
