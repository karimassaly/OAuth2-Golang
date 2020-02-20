package handlers

import (
	"context"

	"github.com/google/go-github/github"
)

func HandleGitInit() {
	client := github.NewClient(nil)

	// list all organizations for user "willnorris"
	orgs, _, err := client.Organizations.List(context.Background(), "karimcorp", nil)
}
