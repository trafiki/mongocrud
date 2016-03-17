package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/jackdanger/collectlinks" // The collectlinks library is made for parsing links.
)

var transport = &http.Transport{
	TLSClientConfig: &tls.Config{
		InsecureSkipVerify: true,
	},
}

var client = http.Client{Transport: transport}

var deadlinks = []error{}

func main() {
	target := flag.String("url", "", "Please supply a url")
	crawlAll := flag.Bool("fullcrawl", true, "Option to crawl the domain only")
	flag.Parse()

	fmt.Println("Crawling Address: ", *target)

	targetURI, err := url.Parse(*target)
	if err != nil {
		fmt.Printf("Invalid Address[%s]\n", *target)
		return
	}

	queue := make(chan string)
	filteredQueue := make(chan string)

	var seen = make(map[string]bool)

	// This go routine sends the target url supplied at the command line to the queue channel.
	go func() { queue <- *target }()

	for {
		select {
		case val, ok := <-queue:
			if !ok {
				return
			}
			if !seen[val] {
				seen[val] = true
				go func() { filteredQueue <- val }()
				continue
			}

		case currentURI, ok := <-filteredQueue:
			if !ok {
				return
			}

			enqueue(currentURI, queue, targetURI.Host, *crawlAll)

		case <-time.After(1 * time.Minute):
			fmt.Printf("Expired beyond request timeout.\n")
			fmt.Println("---------------Dead Links-----------------------")
			for _, i := range deadlinks {
				fmt.Println(i)
			}
			fmt.Println("------------------------------------------------")
			return
		}
	}

}

// enqueue retrieves HTTP and parses the links, putting them into the same queue
// used by main.
func enqueue(uri string, queue chan string, targetHost string, crawlAll bool) {
	fmt.Println("Fetching: ", uri)

	currentURI, err := url.Parse(uri)
	if err != nil {

		return
	}

	if crawlAll && !strings.Contains(currentURI.Host, targetHost) {
		return
	}

	resp, err := client.Get(uri)
	if err != nil {
		fmt.Println("")
		fmt.Printf("Host: %s\n", currentURI.Host)
		fmt.Printf("URL: %s\n", currentURI.String())
		fmt.Printf("Error: %s\n", err.Error())
		fmt.Println("")
		deadlinks = append(deadlinks, err)
		return
	}

	defer resp.Body.Close()

	links := collectlinks.All(resp.Body)

	for _, link := range links {
		fixedLink, err := fixURL(link, uri)
		if err != nil {
			fmt.Println(err)
			continue
		}

		queue <- fixedLink
	}
}

func fixURL(href, base string) (string, error) {
	uri, err := url.Parse(href)
	if err != nil {
		return "", err
	}

	baseURL, err := url.Parse(base)
	if err != nil {
		return "", err
	}

	return baseURL.ResolveReference(uri).String(), nil
}
