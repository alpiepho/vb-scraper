

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

- [x] parse college details page
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
- [x] other stuff
    - [x] google link
    - [x] wikipedia link
    - [x] lat/long
    - [x] volleyball link
- [x] export details as text
- [x] export details as html
- [ ] export details as csv

- [x] filter while scraping details
- [-] filter after import details
- [x] filter details based on
    - [x] state
    - [x] college name
    - [ ] city
    - [x] level
    - [ ] conference
    - [ ] academic_selectivity
    - [ ] undergrad_enrollment
    - [ ] cost in/out state
- [ ] sort details based on
    - [ ] state
    - [ ] city
    - [ ] level
    - [ ] conference
    - [ ] academic_selectivity
    - [ ] undergrad_enrollment
    - [ ] cost in/out state

- [x] google maps link from radius?


- [x] build as web server
- [x] generate html output
- [ ] build map (scrape site to get lat/long)

- [x] alpine.js, expand all
- [x] alpine.js, filters: level
- [x] alpine.js, filters: state
- [-] alpine.js, filters: other?
- [x] cleaner buttons
- [x] labels as links (no url on page)
- [x] no logo link
- [x] remove label under scores
- [x] fonts?
- [-] mascot (link for mascot)
- [ ] link to google map? (ex. https://www.google.com/maps/search/?api=1&query=Northern+Arizona+University)
 
- [ ] fix lat bugs (automate around per state) (ex. https://www.google.com/search?q=latitude+Vassar+College)

- [ ] use aplpine.js and local storage to created todo list of majors
    - see https://stackoverflow.com/questions/70670336/how-to-integrate-alpinejs-persist-onto-a-toggle-using-localstoage
    - https://codewithhugo.com/alpinejs-localstorage-sessionstorage/



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
- https://developers.google.com/maps/documentation/javascript/reference/marker
- https://pkg.go.dev/github.com/chromedp/chromedp@v0.8.0#Text
- calc distance:
- https://gist.github.com/hotdang-ca/6c1ee75c48e515aec5bc6db6e3265e49



## APENDIX: pending kip2
<pre>
Las Vegas
Arizona Christian
Eastern Washington.
Georgia
George Fox
Kean 
Le moine
Linfield 
Mcla
Midwestern state 
Nevada 
New Mexico State 
Northern arizona 
Occidental 
Pacific
Pacific Lutheran
Providence. 
Puget sound
Southern Connecticut 
UNLV
Temple
Washington college 
Westminster. 
Whitman 
William Paterson



NIT
Arkansas state x2
Babson x2
Baker 
Belmont
Cameron 
Colorado College
Cornell college
Creighton x3
Elon x2
Emporia 
Georgia
Georgia tech
Gettysburg
Graceland
Grand Canyon x2
Grand valley
Harvard
Houston x2
Iowa x3
Iowa state x2 
Kansas 
Lehigh
Loyola 
Marist 
Mason.  Emerson 
Mesa
Michigan
Middlebury
Ohio state 
Oregon state 
Palm beach Atlantic
Rhodes 
Sac state 
Samford Christian school 
Seattle pacific
Sewanee
SFA x2
Southern Utah 
Texas x2
UCF x2
UNLV 
UTSA x2
Washington 
West Virginia x3 
Western Carolina.  Shorter athletic on move was
Westmont 
</pre>
 
## APENDIX: pending kip spokane tour

<pre>
Seattle U
UW
Seattle Pacific
Northwest Nazarene
Lewis and Clark (ID)
Pacific Lutheran
Pacific
Linfield
Geogre Fox
Eastern Washington
</pre>

## APENDIX: pending kip camp

<pre>
Southern Utah
Creighton
Lehigh
Westmont
Oklahoma Christian
Iowa
Mesa
Pacific
Emporia State
Friends
Cameron
Arkansas State
Bethel
Anderson
Temple
Lynn
Nova
Flagler
</pre>

## APENDIX - Trick for comparing

```
cat college_details_usavCoaches.html  | grep ">name:<" > junk2.txt
cat config_usavCoaches.json| grep "\"," > junk1.txt
# edit
meld junk1.txt junk2.txt
```
