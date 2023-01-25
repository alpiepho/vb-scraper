

# Overview

A web scraper application to gather a list of colleges from sportsrecruits.com, specifically, the page that allows picking colleges by a map of US states. Each link show a list of colleges for all levels (DI, DII, DIII etc), but does not allow filtering.  This app allows filtering based on various criteria:

- state
- level
- in-state, out-state
- etc.


## TODO

- [x] start and open first page
- [x] gather state links
- [x] add pacing mechanism
- [x] add config file

- [x] test dump_states
- [x] test dump_colleges
- [x] test parse_map
- [x] test parse_states
- [x] test export_colleges
- [-] test import_colleges txt
- [x] test import_colleges json

- [ ] parse college details page
    - [x] conference
    - [x] academic selectivity
    - [x] enrollment (remove comma)
    - [x] Private/Public
    - [x] coaches names (might be dated, check google -> volleybll roster)
    - [x] student-to-faculty ratio
    - [x] calendar system
    - [x] graduation rate
    - [x] retention rate
    - [x] enrollment by gender
    - [x] on-campus housing
    - [x] Admissions..
    - [x] Cost...
    - [x] Major...
- [ ] other stuff
    - [x] google link
    - [x] wikipedia link
    - [ ] lat/long
    - [ ] mascot
    - [ ] volleyball link



- [ ] export/import colleges list to separate processes
- [ ] output csv/tab table
- [ ] determine fields to output

- [ ] build as web server
- [ ] generate html output
- [ ] build map (scrape site to get lat/long)
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

- https://www.google.com/search?q=create+google+map+with+markers&rlz=1CDGOYI_enUS729US729&oq=create+google+map+&aqs=chrome.1.69i57j0i512l5.14447j0j7&hl=en-US&sourceid=chrome-mobile&ie=UTF-8
- https://www.google.com/maps/about/mymaps/
- https://developers.google.com/maps/documentation/javascript/examples/marker-simple
- https://www.latlong.net/