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
	StatesList []string `json:"stateslist"`
	// flags
	ParseMap       bool `json:"parse_map"`
	ParseStates    bool `json:"parse_states"`
	ExportColleges bool `json:"export_colleges"`
	ImportColleges bool `json:"import_colleges"`
	// files
	ExportCollegesFile string `json:"export_colleges_file"`
	ImportCollegesFile string `json:"import_colleges_file"`
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
	if !appConfig.DumpStates {
		return
	}
	for _, state := range *states {
		fmt.Println("")
		fmt.Println("name:     " + state.name)
		fmt.Println("pageLink: " + state.stateLink)
	}
	fmt.Println(len(*states))
}

func dumpColleges(colleges *[]College) {
	if !appConfig.DumpColleges {
		return
	}
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
	if !appConfig.ExportColleges || len(appConfig.ExportCollegesFile) == 0 {
		return
	}
	fileName := appConfig.ExportCollegesFile
	file, err := os.Create(fileName)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	if strings.Contains(fileName, ".json") {
		b, err := json.MarshalIndent(*colleges, "", "  ")
		if err != nil {
			fmt.Println(err)
			return
		}
		file.WriteString(string(b))
	} else {
		for _, college := range *colleges {
			file.WriteString("")
			file.WriteString("name:       " + college.Name)
			file.WriteString("state:      " + college.State)
			file.WriteString("city:       " + college.City)
			file.WriteString("level:      " + college.Level)
			file.WriteString("pageLink:   " + college.CollegeLink)
		}
	}
}

func importColleges(colleges *[]College) {
	if !appConfig.ImportColleges || len(appConfig.ImportCollegesFile) == 0 {
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
	} else {
		lines := strings.Split(string(data), "\n")
		for _, line := range lines {
			var college College
			line := strings.TrimSpace(line)
			if len(line) == 0 {
				continue
			}
			value := strings.Split(line, ":")[1]
			value = strings.TrimSpace(value)

			if strings.Contains(line, "name:") {
				college = College{}
				college.Name = value
			}
			if strings.Contains(line, "state:") {
				college := College{}
				college.State = value
			}
			if strings.Contains(line, "city:") {
				college := College{}
				college.City = value
			}
			if strings.Contains(line, "level:") {
				college := College{}
				college.Level = value
			}
			if strings.Contains(line, "pageLink:") {
				college := College{}
				college.CollegeLink = value
				*colleges = append(*colleges, college)
			}
		}
	}
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
	defer cancel()
	ctx, cancel = chromedp.NewContext(ctx)
	defer cancel()

	fmt.Println("parse map...")
	var states []State
	parseForStates(&ctx, &states)

	fmt.Println("parse states...")
	var colleges []College
	for i, state := range states {
		if testStatesSkip(i) {
			continue
		}

		fmt.Println("parse colleges...")
		parseForColleges(&ctx, &colleges, state)
	}

	dumpStates(&states)
	dumpColleges(&colleges)
	exportColleges(&colleges)
	importColleges(&colleges)

	cancel()

	fmt.Println("done.")
}
