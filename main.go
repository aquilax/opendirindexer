package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/gocolly/colly"
	"github.com/gocolly/colly/debug"
	"github.com/gocolly/colly/queue"
)

const appName = "opendirindexer"
const defaultUserAgent = appName + "/1.0"

var set map[string]bool

var ignoreRobots = flag.Bool("ingoreRobots", false, "Ignores robots.txt restrictions")
var enableDebug = flag.Bool("debug", false, "Enable debugging")
var userAgent = flag.String("userAgent", defaultUserAgent, "set user agent")

var CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
var Usage = func() {
	fmt.Fprintf(CommandLine.Output(), "Usage of %s:\n", appName)
	fmt.Fprintf(CommandLine.Output(), "%s [OPTIONS] URL\n", appName)
	flag.PrintDefaults()
}

func mustIgnore(urlParam, link, linkFull string) bool {
	// Going up
	if !strings.HasPrefix(linkFull, urlParam) {
		return true
	}
	// Dynamic urls
	if strings.Contains(link, "?") {
		return true
	}
	return false
}

func mustQueue(urlParam, link, linkFull string) bool {
	return link[len(link)-1:] == "/"
}

func mustKeep(urlParam, link, linkFull string) bool {
	if !set[linkFull] {
		set[linkFull] = true
		return true
	}
	return false
}

func main() {
	flag.Parse()
	urlParam := flag.Arg(0)
	if urlParam == "" {
		Usage()
	}
	set = make(map[string]bool)
	var options = []func(*colly.Collector){
		colly.UserAgent(*userAgent),
	}

	if *ignoreRobots {
		options = append(options, colly.IgnoreRobotsTxt())
	}
	if *enableDebug {
		options = append(options, colly.Debugger(&debug.LogDebugger{}))
	}

	c := colly.NewCollector(options...)
	q, _ := queue.New(
		1,
		&queue.InMemoryQueueStorage{MaxSize: 10000},
	)

	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		linkFull := e.Request.AbsoluteURL(link)
		if mustIgnore(urlParam, link, linkFull) {
			return
		}
		if mustQueue(urlParam, link, linkFull) {
			q.AddURL(linkFull)
			return
		}
		if mustKeep(urlParam, link, linkFull) {
			fmt.Println(linkFull)
		}
	})

	q.AddURL(urlParam)
	q.Run(c)
}
