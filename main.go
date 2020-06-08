package main

import (
	"flag"
	"os"
	"sort"
	"time"

	"github.com/golang/glog"

	"github.com/BurntSushi/toml"
)

// Config contains config data for a news fetch/render cycle
type Config struct {
	Title    string
	Days     int `toml:"daysOfNews"`
	Template string
	OutFile  string
	Feeds    []string
}

var (
	debug  = false
	dryrun = false
)

func init() {
	flag.Set("logtostderr", "true")
	flag.Set("v", "2")
}

func main() {
	cfgFileName := flag.String("config", "config.toml", "path to config file")
	linkStateFileName := flag.String("linkstate", "", "path to link state file")
	timeout := flag.Duration("timeout", 60*time.Minute, "timeout in seconds for fetch operations")
	flag.BoolVar(&debug, "debug", false, "enable debug output")
	flag.BoolVar(&dryrun, "dryrun", false, "disable mutation of files")
	flag.Parse()

	var conf Config
	if _, err := toml.DecodeFile(*cfgFileName, &conf); err != nil {
		glog.Fatalf("error reading config file %s: %s", *cfgFileName, err)
	}

	if debug {
		glog.Infof("config: %+v", conf)
	}

	links, err := fetchLinks(conf.Feeds, *timeout)
	if err != nil {
		glog.Fatalf("error fetching links: %s", err)
	}

	if debug {
		glog.Infof("found links: %+v", links)
	}

	// use link state file if name is set
	if *linkStateFileName != "" {

	}

	// Sort items chronologically descending
	sort.Slice(links, func(i, j int) bool {
		return links[i].Published.After(links[j].Published)
	})

	if debug {
		glog.Infof("sorted links: %+v", links)
	}

	// Filter by age
	notBefore := time.Now().AddDate(0, 0, -conf.Days)
	filteredLinks := []Link{}
	for _, l := range links {
		if l.Published.Before(notBefore) {
			break
		}
		filteredLinks = append(filteredLinks, l)
	}
	links = filteredLinks

	out := os.Stdout
	if !dryrun {
		out, err = os.Create(conf.OutFile)
		if err != nil {
			glog.Fatalf("could not open newsfeed file: %s", err)
		}
	}

	err = renderLinks(out, links, conf.Template)
	if err != nil {
		glog.Fatalf("could not generate newsfeed: %s", err)
	}
}
