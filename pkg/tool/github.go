package tool

import (
	"context"
	"fmt"

	"github.com/google/go-github/v32/github"
)

func getSubscribeUrlFromGithub(repo string, username string) {
	client := github.NewClient(nil)

	opt := &github.RepositoryListOptions{Type: "public"}
	repos, _, err := client.Repositories.List(context.Background(), "gotodown", opt)
	if err != nil {
		panic(err)
	}
	fmt.Println(repos.List())
}
