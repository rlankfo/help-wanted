package hw

import (
	"flag"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"gopkg.in/yaml.v2"
)

const (
	hours      = "hours"
	verbose    = "verbose"
	label      = "label"
	org        = "org"
	configFile = "config"
)

type Config struct {
	Hours   int  `yaml:"hours"`
	Verbose bool `yaml:"verbose"`

	labels listFlag
	Labels []string `yaml:"labels":`

	organizations listFlag
	Organizations []string `yaml:"orgs"`

	CfgFile string `yaml:"-"`
}

var defaultConfig = Config{
	Hours:   72,
	Verbose: false,

	Labels:        []string{"help wanted"},
	Organizations: []string{"grafana", "prometheus", "kubernetes"},
	CfgFile:       "~/.config/help-wanted/config.yml",
}

func (c *Config) RegisterFlags(f *flag.FlagSet) {
	f.IntVar(&c.Hours, hours, defaultConfig.Hours,
		"Hours since issue was created.")
	f.BoolVar(&c.Verbose, verbose, defaultConfig.Verbose,
		"Prints Github search string associated with query.")
	f.Var(&c.labels, label, "Find issues with this label.")
	f.Var(&c.organizations, org, "Github organization to search.")
	f.StringVar(&c.CfgFile, configFile, defaultConfig.CfgFile,
		"YAML config file. Command line flags will override values set in this file.")
}

func (c *Config) ParseFlags() {
	flag.Parse()

	// set listFlag defaults
	if len(c.labels) == 0 {
		c.Labels = defaultConfig.Labels
	} else {
		c.Labels = []string(c.labels)
	}
	if len(c.organizations) == 0 {
		c.Organizations = defaultConfig.Organizations
	} else {
		c.Organizations = []string(c.organizations)
	}
}

func (c *Config) LoadConfigFile(f *flag.FlagSet) {
	loadCfg := flagWasSet(f, configFile)
	bytes, err := ioutil.ReadFile(c.CfgFile)
	if err != nil {
		// the file doesn't exist, and flag wasn't passed. don't load it.
		if !loadCfg && os.IsNotExist(err) {
			return
		}
		log.Fatal("error loading config file: ", err)
	}

	var cfg Config
	err = yaml.Unmarshal(bytes, &cfg)
	if err != nil {
		log.Fatal("cannot unmarshal config: ", err)
	}

	// reset values on *c if !flagWasSet. This is to override
	// values set in the config yaml file with command line args
	if !flagWasSet(f, hours) {
		c.Hours = cfg.Hours
	}
	if !flagWasSet(f, verbose) {
		c.Verbose = cfg.Verbose
	}
	if !flagWasSet(f, org) {
		c.Organizations = cfg.Organizations
	}
	if !flagWasSet(f, label) {
		c.Labels = cfg.Labels
	}
}

func flagWasSet(f *flag.FlagSet, name string) bool {
	set := false
	f.Visit(func(fl *flag.Flag) {
		if fl.Name == name {
			set = true
		}
	})
	return set
}

type listFlag []string

func (f *listFlag) String() string {
	return strings.Join([]string(*f), ",")
}

func (f *listFlag) Set(v string) error {
	*f = append(*f, v)
	return nil
}
