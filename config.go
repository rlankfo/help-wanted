package hw

import (
	"flag"
	"strings"
)

type Config struct {
	Hours   int
	Verbose bool

	Labels        listFlag
	Organizations listFlag
}

var defaultConfig = Config{
	Hours:   72,
	Verbose: false,

	Labels:        listFlag([]string{"help wanted", "good first issue"}),
	Organizations: listFlag([]string{"grafana", "ohmyzsh", "kubernetes", "homebrew"}),
}

func (c *Config) RegisterFlags(f *flag.FlagSet) {
	f.IntVar(&c.Hours, "hours", defaultConfig.Hours,
		"Hours since issue was created.")
	f.BoolVar(&c.Verbose, "verbose", defaultConfig.Verbose,
		"Prints Github search string associated with query..")
	f.Var(&c.Labels, "label", "Find issues with this label.")
	f.Var(&c.Organizations, "org", "Github organization to search.")
}

func (c *Config) ParseFlags() {
	flag.Parse()
	// set listFlag defaults
	if len(c.Labels) == 0 {
		c.Labels = defaultConfig.Labels
	}
	if len(c.Organizations) == 0 {
		c.Organizations = defaultConfig.Organizations
	}
}

type listFlag []string

func (f *listFlag) String() string {
	return strings.Join([]string(*f), ",")
}

func (f *listFlag) Set(v string) error {
	*f = append(*f, v)
	return nil
}
