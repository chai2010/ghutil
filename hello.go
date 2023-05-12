//go:build ignore
// +build ignore

package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/chai2010/ghutil"
)

var (
	flagUserName = flag.String("u", "KusionStack", "set user name")
)

func main() {
	flag.Parse()

	c := ghutil.NewGHClient(&ghutil.GHOptions{
		Token: os.Getenv("GITHUB_ACCESS_TOKEN"),
	})
	repos, err := c.GetRepositories(context.Background(), *flagUserName)
	if err != nil {
		panic(err)
	}

	var starCount int
	var forkCount int
	for _, repo := range repos {
		starCount += *repo.StargazersCount
		forkCount += *repo.ForksCount
	}

	fmt.Printf("%s/{*}:StargazersCount: %d\n", *flagUserName, starCount)
	fmt.Printf("%s/{*}:ForksCount:      %d\n", *flagUserName, forkCount)
}
