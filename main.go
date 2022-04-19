package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/google/go-github/v43/github"
)

var defaultConfig = config{
	hours:        72,
	label:        "help wanted",
	organization: "grafana",
}

type config struct {
	hours        int
	label        string
	organization string
}

func (c *config) registerFlags(f *flag.FlagSet) {
	f.IntVar(&c.hours, "hours", defaultConfig.hours,
		"Hours since issue was created.")
	f.StringVar(&c.label, "label", defaultConfig.label,
		"Find issues with this label.")
	f.StringVar(&c.organization, "org", defaultConfig.organization,
		"Github organization to search.")

}

func main() {
	var (
		issues = make(chan *github.Issue)
		wg     sync.WaitGroup
		cfg    config
	)
	cfg.registerFlags(flag.CommandLine)
	flag.Parse()

	d, err := time.ParseDuration(fmt.Sprintf("%dh", cfg.hours))
	if err != nil {
		fmt.Println("error parsing duration: ", err)
		os.Exit(-1)
	}
	query := fmt.Sprintf("is:issue is:open org:%s label:\"help wanted\" created:>=%s archived:false",
		cfg.organization, time.Now().Add(-d).Format(time.RFC3339))

	fmt.Println(query)

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
