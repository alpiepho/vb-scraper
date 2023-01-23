package main

// To run:
// go run main.go > results.txt 2>&1

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
	"golang.org/x/exp/slices"
)

type Configuration struct {
	StatesList               []string `json:"stateslist"`
	OpenChromedp             bool     `json:"open_chromedp"`
	ParseMap                 bool     `json:"parse_map"`
	ParseStates              bool     `json:"parse_states"`
	ParseColleges            bool     `json:"parse_colleges"`
	ExportColleges           bool     `json:"export_colleges"`
	ExportCollegesFile       string   `json:"export_colleges_file"`
	ImportColleges           bool     `json:"import_colleges"`
	ImportCollegesFile       string   `json:"import_colleges_file"`
	CollegeList              []string `json:"collegelist"`
	ParseCollegePages        bool     `json:"parse_college_pages"`
	ExportCollegeDetails     bool     `json:"export_college_details"`
	ExportCollegeDetailsFile string   `json:"export_college_details_file"`

	//DEBUG
	DumpStates   bool `json:"dump_states"`
	DumpColleges bool `json:"dump_colleges"`
}

var appConfig Configuration

const MAP_URL = "https://sportsrecruits.com/athletic-scholarships/womens-volleyball"

const NAV_TIME_MAX_STATE = 3
const NAV_TIME_MAX_COLLEGE = 3

type State struct {
	name      string
	stateLink string
}

type College struct {
	Name        string `json:"name"`
	State       string `json:"state"`
	City        string `json:"city"`
	Level       string `json:"level"`
	CollegeLink string `json:"college_link"`
	StateLink   string `json:"state_link"`
}

type CollegeDetail struct {
	Name                string `json:"name"`
	State               string `json:"state"`
	City                string `json:"city"`
	Level               string `json:"level"`
	CollegeLink         string `json:"college_link"`
	StateLink           string `json:"state_link"`
	Conference          string `json:"conference"`
	AcademicSelectivity string `json:"academic_selectivity"`
	UndergradEnrollment string `json:"undergrad_enrollment"`
	ControlAffiliation  string `json:"control_affilication"`
	Overview            string `json:"overview"`
	// TODO
	// Quick Facts
	// - [x]   conference
	// - [x]   academic selectivity
	// - [x]   undergrad enrollment (filter comma)
	// - [x]   control/affilication
	// Overview
	// - [x]   paragraph (filter stock sentences)
	// Athletics
	// - [ ]   division (already have this from state lists)
	// - [ ]   conference (already have this with quick facts)
	// - [ ]   head coach name (cross with other site)
	// - [ ]   assistant coach name (cross with other site)
	// School Profile
	// - [ ]   student ratio
	// - [ ]   graduation rate
	// - [ ]   enrollment by gender
	// - [ ]   calendar system
	// - [ ]   retention rate
	// - [ ]   on campus housing
	// Admissions
	// - [ ]   acceptance rate
	// - [ ]   total applicants
	// - [ ] Total Scores (25th-75th percentile)
	// - [ ]   SAT (should list max and max per section)
	// - [ ]     students submitting scores
	// - [ ]     reading range
	// - [ ]     math range
	// - [ ]     writing range
	// - [ ]   ACT (should list max and max per section)
	// - [ ]     students submitting scores
	// - [ ]     reading range
	// - [ ]     math range
	// - [ ]     writing range
	// - [ ]   Requirements
	// - [ ]     open admission policy
	// - [ ]     application fee
	// - [ ]     recomendations
	// - [ ]     secondary school record
	// - [ ]     secondary school rank
	// - [ ]     secondary school gpa
	// Cost (determine per ?)
	// - [ ]   In-State
	// - [ ]     total cost
	// - [ ]     tuition
	// - [ ]     fee
	// - [ ]     on campus room and board
	// - [ ]   Out-State
	// - [ ]     total cost
	// - [ ]     tuition
	// - [ ]     fee
	// - [ ]     on campus room and board
	// - [ ]   Financial Aid
	// - [ ]     percent undegrad receiving aid
	// Majors
	// - [ ]     college (multiple sections)
	// - [ ]       major category (multiple sections)
	// - [ ]         major category (multiple)
	// Other (from other source)
	// - [ ] mascot
	// - [ ] link to roster/coaches
	// - [ ] link to wikipedia for school
	// - [ ] geological coordinate
	// - [ ] ???
}

var STATE_NAMES = [51]string{
	"Alabama",
	"Alaska",
	"Arizona",
	"Arkansas",
	"California",
	"Colorado",
	"Connecticut",
	"Delaware",
	"Florida",
	"Georgia",
	"Hawaii",
	"Idaho",
	"Illinois",
	"Indiana",
	"Iowa",
	"Kansas",
	"Kentucky",
	"Louisiana",
	"Maine",
	"Maryland",
	"Massachusetts",
	"Michigan",
	"Minnesota",
	"Mississippi",
	"Missouri",
	"Montana",
	"Nebraska",
	"Nevada",
	"New Hampshire",
	"New Jersey",
	"New Mexico",
	"New York",
	"North Carolina",
	"North Dakota",
	"Ohio",
	"Oklahoma",
	"Oregon",
	"Pennsylvania",
	"Rhode Island",
	"South Carolina",
	"South Dakota",
	"Tennessee",
	"Texas",
	"Utah",
	"Vermont",
	"Virginia",
	"Washington",
	"West Virginia",
	"Wisconsin",
	"Wyoming",
	"District of Columbia",
}

func testStatesSkip(i int) bool {
	if slices.Contains(appConfig.StatesList, "All") {
		return false
	}
	if slices.Contains(appConfig.StatesList, "all") {
		return false
	}
	name := STATE_NAMES[i]
	if slices.Contains(appConfig.StatesList, name) {
		return false
	}
	return true
}

func testCollegesSkip(i int) bool {
	//DEBUG - limit products
	//return !(i == 1 || i == 2)
	//return i > 40
	//return i > 10
	return false
}

func parseForStates(ctx *context.Context, states *[]State) {
	var err error
	n := rand.Intn(NAV_TIME_MAX_STATE)
	random_delay := time.Duration(n) * time.Second

	err = chromedp.Run(*ctx,
		chromedp.Navigate(MAP_URL),
		chromedp.Sleep(1*time.Second),
		chromedp.Sleep(random_delay),
		//DEBUG: chromedp.Sleep(400*time.Second),
	)
	if err != nil {
		// ignore error
		//DEBUG:
		fmt.Println(err)
	}

	var stateLinkNodes []*cdp.Node
	err = chromedp.Run(*ctx,
		chromedp.Nodes(`.states li`, &stateLinkNodes, chromedp.ByQueryAll),
	)
	if err != nil {
		// ignore error
		//DEBUG:
		fmt.Println(err)
	}
	//DEBUG:
	//fmt.Println(len(stateLinkNodes))

	for _, n := range stateLinkNodes {
		var ok bool
		data := State{}
		err := chromedp.Run(*ctx,
			chromedp.Text(`a`, &data.name, chromedp.ByQuery, chromedp.FromNode(n)),
			chromedp.AttributeValue(`a`, "href", &data.stateLink, &ok, chromedp.NodeVisible, chromedp.ByQuery, chromedp.AtLeast(0), chromedp.FromNode(n)),
		)
		if err != nil {
			// ignore error
			//DEBUG:
			fmt.Println(err)
		}
		//data.pageLink = BASE_URL + data.pageLink
		data.name = strings.Replace(data.name, "\u0026", "&", 1)

		//DEBUG:
		//fmt.Println(i, ":\t", data.name, " ", data.pageLink)
		*states = append(*states, data)
	}
}

func parseForColleges(ctx *context.Context, colleges *[]College, state State) {
	var err error
	n := rand.Intn(NAV_TIME_MAX_COLLEGE)
	random_delay := time.Duration(n) * time.Second

	err = chromedp.Run(*ctx,
		chromedp.Navigate(state.stateLink),
		chromedp.Sleep(2*time.Second),
		chromedp.Sleep(random_delay),
		//DEBUG: chromedp.Sleep(400*time.Second),
	)
	if err != nil {
		// ignore error
		//DEBUG:
		fmt.Println(err)
	}

	var collegeLinkNodes []*cdp.Node
	err = chromedp.Run(*ctx,
		chromedp.Nodes(`.data-table a`, &collegeLinkNodes, chromedp.ByQueryAll),
	)
	if err != nil {
		// ignore error
		//DEBUG:
		fmt.Println(err)
	}
	//DEBUG:
	fmt.Println(state.name)
	fmt.Println(len(collegeLinkNodes))

	for i, n := range collegeLinkNodes {
		//var ok bool
		data := College{}
		data.State = state.name
		data.StateLink = state.stateLink
		data.CollegeLink = n.Attributes[3]

		if testCollegesSkip(i) {
			continue
		}

		err := chromedp.Run(*ctx,
			chromedp.Text(`.col-sm-5 p`, &data.Name, chromedp.ByQuery, chromedp.FromNode(n)),
			chromedp.Text(`.col-sm-4 p`, &data.City, chromedp.ByQuery, chromedp.FromNode(n)),
			chromedp.Text(`.col-sm-3 p`, &data.Level, chromedp.ByQuery, chromedp.FromNode(n)),
			//chromedp.AttributeValue(`.avatar img`, "src", &data.avatarLink, &ok, chromedp.NodeVisible, chromedp.ByQuery, chromedp.AtLeast(0), chromedp.FromNode(n)),
		)
		if err != nil {
			// ignore error
			//DEBUG:
			fmt.Println(err)
		}
		//data.pageLink = BASE_URL + data.pageLink

		//DEBUG:
		//fmt.Println(i, ":\t", data.name, " ", data.pageLink)
		*colleges = append(*colleges, data)
	}
	//DEBUG:
	//fmt.Println(len(stateLinkNodes))
}

func dumpStates(states *[]State) {
	for _, state := range *states {
		fmt.Println("")
		fmt.Println("name:     " + state.name)
		fmt.Println("pageLink: " + state.stateLink)
	}
	fmt.Println(len(*states))
}

func dumpColleges(colleges *[]College) {
	for _, college := range *colleges {
		fmt.Println("")
		fmt.Println("name:       " + college.Name)
		fmt.Println("state:      " + college.State)
		fmt.Println("city:       " + college.City)
		fmt.Println("level:      " + college.Level)
		fmt.Println("pageLink:   " + college.CollegeLink)
	}
	fmt.Println(len(*colleges))
}

func exportColleges(colleges *[]College) {
	if len(appConfig.ExportCollegesFile) == 0 {
		return
	}
	fileName := appConfig.ExportCollegesFile
	file, err := os.Create(fileName)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	b, err := json.MarshalIndent(*colleges, "", "  ")
	if err != nil {
		fmt.Println(err)
		return
	}
	file.WriteString(string(b))
}

func importColleges(colleges *[]College) {
	if len(appConfig.ImportCollegesFile) == 0 {
		return
	}
	fileName := appConfig.ImportCollegesFile
	file, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	data, err := ioutil.ReadAll(file)
	if err != nil {
		panic(err)
	}
	if strings.Contains(fileName, ".json") {
		err = json.Unmarshal(data, &colleges)
		if err != nil {
			panic(err)
		}
	}
	fmt.Println(len(*colleges))
}

func testCollegeSkip(name string) bool {
	if slices.Contains(appConfig.CollegeList, "All") {
		return false
	}
	if slices.Contains(appConfig.CollegeList, "all") {
		return false
	}
	if slices.Contains(appConfig.CollegeList, name) {
		return false
	}
	return true
}

func parseForCollegePages(ctx *context.Context, details *[]CollegeDetail, college *College) {
	var err error
	n := rand.Intn(NAV_TIME_MAX_COLLEGE)
	random_delay := time.Duration(n) * time.Second

	data := CollegeDetail{}
	data.Name = college.Name
	data.State = college.State
	data.City = college.City
	data.Level = college.Level
	data.CollegeLink = college.CollegeLink
	data.StateLink = college.StateLink
	data.Conference = ""
	data.AcademicSelectivity = ""
	data.UndergradEnrollment = ""
	data.ControlAffiliation = ""
	data.Overview = ""

	overview1 := ""
	overview2 := ""

	//#main > div.wrapper_row > div.container-fluid.ProfileMain-wrap > div > div.col-xs-12.col-sm-8.ProfileMain-content > div.ProgramProfile-overview.start-xs > div > div > div > p
	//#main > div.wrapper_row > div.container-fluid.ProfileMain-wrap > div > div.col-xs-12.col-sm-8.ProfileMain-content > div.ProgramProfile-overview.start-xs > div > div > div

	err = chromedp.Run(*ctx,
		chromedp.Navigate(college.CollegeLink),
		chromedp.Sleep(2*time.Second),
		chromedp.Sleep(random_delay),
		//DEBUG: chromedp.Sleep(400*time.Second),
		chromedp.Text(`#quick-facts-section > div > p:nth-child(2)`, &data.Conference, chromedp.ByQuery),
		chromedp.Text(`#quick-facts-section > div > p:nth-child(4)`, &data.AcademicSelectivity, chromedp.ByQuery),
		chromedp.Text(`#quick-facts-section > div > p:nth-child(6)`, &data.UndergradEnrollment, chromedp.ByQuery),
		chromedp.Text(`#quick-facts-section > div > p:nth-child(8)`, &data.ControlAffiliation, chromedp.ByQuery),
		// TODO: is overview useful?
		chromedp.Click(`.toggle-content`),
		chromedp.Sleep(1*time.Second),
		chromedp.Text(`.read-more-container`, &data.Overview, chromedp.ByQuery),
	)
	if err != nil {
		// ignore error
		//DEBUG:
		fmt.Println(err)
	}
	data.UndergradEnrollment = strings.ReplaceAll(data.UndergradEnrollment, ",", "")

	//DEBUG:
	//fmt.Println(i, ":\t", data.name, " ", data.pageLink)
	*details = append(*details, data)

	// Enrollment: #quick-facts-section > div > p:nth-child(6)

	// var collegeLinkNodes []*cdp.Node
	// err = chromedp.Run(*ctx,
	// 	chromedp.Nodes(`.data-table a`, &collegeLinkNodes, chromedp.ByQueryAll),
	// )
	// if err != nil {
	// 	// ignore error
	// 	//DEBUG:
	// 	fmt.Println(err)
	// }
	// //DEBUG:
	// fmt.Println(state.name)
	// fmt.Println(len(collegeLinkNodes))

	// for i, n := range collegeLinkNodes {
	// 	//var ok bool
	// 	data := College{}
	// 	data.State = state.name
	// 	data.StateLink = state.stateLink
	// 	data.CollegeLink = n.Attributes[3]

	// 	if testCollegesSkip(i) {
	// 		continue
	// 	}

	// 	err := chromedp.Run(*ctx,
	// 		chromedp.Text(`.col-sm-5 p`, &data.Name, chromedp.ByQuery, chromedp.FromNode(n)),
	// 		chromedp.Text(`.col-sm-4 p`, &data.City, chromedp.ByQuery, chromedp.FromNode(n)),
	// 		chromedp.Text(`.col-sm-3 p`, &data.Level, chromedp.ByQuery, chromedp.FromNode(n)),
	// 		//chromedp.AttributeValue(`.avatar img`, "src", &data.avatarLink, &ok, chromedp.NodeVisible, chromedp.ByQuery, chromedp.AtLeast(0), chromedp.FromNode(n)),
	// 	)
	// 	if err != nil {
	// 		// ignore error
	// 		//DEBUG:
	// 		fmt.Println(err)
	// 	}
	// 	//data.pageLink = BASE_URL + data.pageLink

	// 	//DEBUG:
	// 	//fmt.Println(i, ":\t", data.name, " ", data.pageLink)
	// 	*colleges = append(*colleges, data)
	// }
	//DEBUG:
	//fmt.Println(len(stateLinkNodes))
}

func exportCollegeDetails(details *[]CollegeDetail) {
	if len(appConfig.ExportCollegeDetailsFile) == 0 {
		return
	}
	fileName := appConfig.ExportCollegeDetailsFile
	file, err := os.Create(fileName)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	b, err := json.MarshalIndent(*details, "", "  ")
	if err != nil {
		fmt.Println(err)
		return
	}
	file.WriteString(string(b))
}

func main() {
	//headless := flag.Bool("headless", false, "a bool")

	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.DisableGPU,
		chromedp.Flag("headless", false),
	)

	confFile, err := os.Open("config.json")
	if err != nil {
		panic(err)
	}
	defer confFile.Close()
	conf, err := ioutil.ReadAll(confFile)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(conf, &appConfig)
	if err != nil {
		panic(err)
	}
	//DEBUG
	//fmt.Printf("%+v", appConfig)

	ctx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	if appConfig.OpenChromedp {
		defer cancel()
		ctx, cancel = chromedp.NewContext(ctx)
		defer cancel()
	}

	var states []State
	if appConfig.ParseMap {
		fmt.Println("parse map...")
		parseForStates(&ctx, &states)
	}

	var colleges []College
	if appConfig.ParseStates {
		fmt.Println("parse states...")
		for i, state := range states {
			if testStatesSkip(i) {
				continue
			}
			if appConfig.ParseColleges {
				fmt.Println("parse colleges...")
				parseForColleges(&ctx, &colleges, state)
			}
		}
	}

	if appConfig.DumpStates {
		dumpStates(&states)
	}
	if appConfig.DumpColleges {
		dumpColleges(&colleges)
	}
	if appConfig.ExportColleges {
		exportColleges(&colleges)
	}
	if appConfig.ImportColleges {
		importColleges(&colleges)
	}
	var details []CollegeDetail
	if appConfig.ParseCollegePages {
		fmt.Println("parse college details...")
		for _, college := range colleges {
			if testCollegeSkip(college.Name) {
				continue
			}
			parseForCollegePages(&ctx, &details, &college)
		}

	}
	if appConfig.ExportCollegeDetails {
		exportCollegeDetails(&details)
	}

	if appConfig.OpenChromedp {
		cancel()
	}

	fmt.Println("done.")
}
