package main

import (
	"context"
	"fmt"
	"github.com/spf13/pflag"
	"github.com/the4thamigo-uk/workerpool"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
)

type (
	result struct {
		url  string
		urls []string
		err  error
	}
)

var (
	urlRegex = regexp.MustCompile(`https?://[A-Za-z0-9/\.]*`)
)

func main() {

	var workers *int = pflag.IntP("workers", "w", 10, "number of workers")

	// TODO read from stdin rather than a fixed set of urls in memory
	urls := []string{
		"http://www.google.com/index.html",
		"http://www.github.com/index.html",
		"http://www.bbc.co.uk/index.html",
		"http://uk.yahoo.com/index.html",
	}

	err := run(urls, *workers)
	if err != nil {
		fmt.Printf("%s", err)
		os.Exit(1)
	}
}

func run(urls []string, workers int) error {

	ctx := context.Background()

	wp, err := workerpool.New(workers, len(urls))
	if err != nil {
		return err
	}

	defer wp.Complete()

	rlts := make(chan result, len(urls))

	// TODO: when we read from stdin we will need to push work on a separate goroutine as
	// Add() can block if the queue size is exceeded.
	for _, url := range urls {
		err = wp.Add(ctx, extractLinksWork(url, rlts))
		if err != nil {
			return err
		}
	}

	var count int

	// TODO: to make this cancellable we will need to change this loop to wait on context
	for rlt := range rlts {
		count++
		if rlt.err != nil {
			fmt.Printf("%s: failed with error : %s\n", rlt.url, rlt.err)
			continue
		}
		fmt.Printf("%s: urls are :\n", rlt.url)
		for _, url := range rlt.urls {
			fmt.Printf("\t%s\n", url)
		}

		if count == len(urls) {
			break
		}
	}
	return nil
}

func extractLinksWork(url string, rlts chan result) func() {
	return func() {
		urls, err := extractLinks(url)
		if err != nil {
			rlts <- result{
				url: url,
				err: err,
			}
			return
		}
		rlts <- result{
			url:  url,
			urls: urls,
		}
	}
}

func extractLinks(url string) ([]string, error) {
	rsp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer rsp.Body.Close()

	body, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return nil, err
	}

	return urlRegex.FindAllString(string(body), -1), nil
}
