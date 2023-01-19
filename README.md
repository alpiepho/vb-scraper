

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
    - [ ] conference
    - [ ] academic selectivity
    - [ ] enroollment (remove comma)
    - [ ] Private/Public
    - [ ] coaches names (might be dated, check google -> volleybll roster)
    - [ ] student-to-faculty ratio
    - [ ] student-to-faculty ratio
    - [ ] calendar system
    - [ ] graduation rate
    - [ ] retention rate
    - [ ] enrollment by gender
    - [ ] on-campus housing
    - [ ] Admissions..
    - [ ] Cost...
    - [ ] Major...



- [ ] export/import colleges list to separate processes
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

