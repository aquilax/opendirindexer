package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gocolly/colly"
	"github.com/gocolly/colly/debug"
	"github.com/gocolly/colly/queue"
)

const appName = "opendirindexer"
const defaultUserAgent = appName + "/1.0"

var set map[string]bool

var ignoreRobots = flag.Bool("ignoreRobots", false, "Ignores robots.txt restrictions")
var insecure = flag.Bool("insecure", false, "Allow insecure tls connections")
var enableDebug = flag.Bool("debug", false, "Enable debugging")
var userAgent = flag.String("userAgent", defaultUserAgent, "set user agent")

// CommandLine is the default set of commands
var CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

// Usage prints a usage message documenting all defined command-line flags
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
	return len(link) > 0 && link[len(link)-1:] == "/"
}

func mustKeep(urlParam, link, linkFull string) bool {
	if !set[linkFull] {
		set[linkFull] = true
		return true
	}
	return false
}

func main() {
	output := os.Stdout
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

	if *insecure {
		c.WithTransport(&http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		})
	}

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
			fmt.Fprintln(output, linkFull)
		}
	})

	q.AddURL(urlParam)
	q.Run(c)
}
