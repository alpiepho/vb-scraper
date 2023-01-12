

# Overview

A web scraper application to gather a list of colleges from sportsrecruits.com, specifically, the page that allows picking colleges by a map of US states. Each link show a list of colleges for all levels (DI, DII, DIII etc), but does not allow filtering.  This app allows filtering based on various criteria:

- state
- level
- in-state, out-state
- etc.


## TODO

- [ ] start and open first page
- [ ] gather state links
- [ ] add pacing mechanism
- [ ] add config file
- [ ] output csv/tab table
- [ ] determine fields to output
- [ ] tbd

- [ ] build as web server
- [ ] generate html output
- [ ] more tbd



## Golang version

The utility was written in Go and chromedp.  Go is a little
more difficult than python, but it seems to run a little faster.  Also,
chromedp will allow automating web flow and can run headless.

Install golang [here](https://golang.org/doc/install)

[chromedp](https://github.com/chromedp/chromedp) library.

`go get -u github.com/chromedp/chromedp`

`go run main.go`

or to gather all the output:

`go run main.go > results.txt 2>&1`

## REFERENCE

