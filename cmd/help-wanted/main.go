package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/google/go-github/v43/github"
	hw "github.com/rlankfo/help-wanted"
)

func main() {
	var (
		issues = make(chan *github.Issue)
		wg     sync.WaitGroup
		cfg    hw.Config
	)
	cfg.RegisterFlags(flag.CommandLine)
	cfg.ParseFlags()

	d, err := time.ParseDuration(fmt.Sprintf("%dh", cfg.Hours))
	if err != nil {
		fmt.Println("error parsing duration: ", err)
		os.Exit(-1)
	}
	orgs := ""
	for _, org := range cfg.Organizations {
		orgs = fmt.Sprintf("org:%s %s", org, orgs)
	}
	labels := ""
	for _, label := range cfg.Labels {
		labels = fmt.Sprintf("label:\"%s\" %s", label, labels)
	}
	query := fmt.Sprintf("is:issue is:open %s %s created:>=%s archived:false",
		orgs, labels, time.Now().Add(-d).Format(time.RFC3339))

	if cfg.Verbose {
		fmt.Println(query)
	}

	wg.Add(1)
	go func(issues <-chan *github.Issue) {
		defer wg.Done()
		for issue := range issues {
			fmt.Printf("%s - %s\n", issue.GetTitle(), issue.GetHTMLURL())
		}
	}(issues)

	func(issues chan<- *github.Issue) {
		err := search(query, issues)
		if err != nil {
			fmt.Println("error searching for issues: ", err)
			os.Exit(-1)
		}
	}(issues)
}

func search(query string, issues chan<- *github.Issue) error {
	var (
		ctx  = context.Background()
		gh   = github.NewClient(nil)
		next = 1
	)

	for next != 0 {
		result, resp, err := gh.Search.Issues(ctx, query, &github.SearchOptions{
			ListOptions: github.ListOptions{
				Page:    next,
				PerPage: 50,
			},
		})
		if err != nil {
			return err
		}
		next = resp.NextPage
		for _, iss := range result.Issues {
			issues <- iss
		}
	}

	return nil
}
