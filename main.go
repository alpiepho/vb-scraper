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
	Name        string `json:"name"`
	State       string `json:"state"`
	City        string `json:"city"`
	Level       string `json:"level"`
	CollegeLink string `json:"college_link"`
	StateLink   string `json:"state_link"`

	Conference          string `json:"conference"`
	AcademicSelectivity string `json:"academic_selectivity"`
	UndergradEnrollment string `json:"undergrad_enrollment"`
	ControlAffiliation  string `json:"control_affilication"`

	Overview string `json:"overview"`

	// Division            string `json:"division"`   // already have this from state lists, Level
	// Conference          string `json:"conference"` // already have this with quick facts
	HeadCoach      string `json:"head_coach_name"`
	AssistantCoach string `json:"assistant_coach_name"`

	StudentRatio       string `json:"student_ratio"`
	GraduationRate     string `json:"garduation_rate"`
	EnrollmentByGender string `json:"enrollment_by_gender"`
	CalendarSystem     string `json:"calendar_system"`
	RetentionRate      string `json:"retention_rate"`
	OnCampusHousing    string `json:"on_campus_housing"`

	// AcceptanceRate                    string `json:"acceptance_rate"`
	// TotalAppicants                    string `json:"total_applicants"`
	// SatStudentsSubmitting             string `json:"sat_students_submitting"`
	// SatSReadingRange                  string `json:"sat_reading_range"`
	// SatMathRange                      string `json:"sat_math_range"`
	// SatWritingRange                   string `json:"sat_writing_range"`
	// ActStudentsSubmitting             string `json:"act_students_submitting"`
	// ActSReadingRange                  string `json:"act_reading_range"`
	// ActMathRange                      string `json:"act_math_range"`
	// ActWritingRange                   string `json:"act_writing_range"`
	// RequirementsOpenAdmissionPolicy   string `json:"requirements_open_admission_policy"`
	// RequirementsApplicationFee        string `json:"requirements_application_fee"`
	// RequirementsRecommendations       string `json:"requirements_recomedations"`
	// RequirementsSecondarySchoolRecord string `json:"requirements_secondary_school_record"`
	// RequirementsSecondarySchoolRank   string `json:"requirements_secondary_school_rank"`
	// RequirementsSecondarySchoolGpa    string `json:"requirements_secondary_school_gpa"`

	// CostInStateTotal         string `json:"cost_in_state_total"`
	// CostInStateTuition       string `json:"cost_in_state_tuition"`
	// CostInStateFee           string `json:"cost_in_state_fee"`
	// CostInStateOnCampusRoom  string `json:"cost_in_state_on_campus_room"`
	// CostOutStateTotal        string `json:"cost_out_state_total"`
	// CostOutStateTuition      string `json:"cost_out_state_tuition"`
	// CostOutStateFee          string `json:"cost_out_state_fee"`
	// CostOutStateOnCampusRoom string `json:"cost_out_state_on_campus_room"`
	// CostPercentUndergradAid  string `json:"cost_percent_undergrad_aid"`

	// Majors []string `json:"majors"`
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
	data.HeadCoach = ""
	data.AssistantCoach = ""

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
		chromedp.Text(`#athletics-section > div > div:nth-child(6) > div:nth-child(1) > p`, &data.HeadCoach, chromedp.ByQuery),
		chromedp.Text(`#athletics-section > div > div:nth-child(6) > div:nth-child(2) > p`, &data.AssistantCoach, chromedp.ByQuery),
		chromedp.Text(`#school-section > div > div:nth-child(1) > div:nth-child(1) > p`, &data.StudentRatio, chromedp.ByQuery),
		chromedp.Text(`#school-section > div > div:nth-child(2) > div:nth-child(1) > p`, &data.GraduationRate, chromedp.ByQuery),
		chromedp.Text(`#school-section > div > div:nth-child(3) > div:nth-child(1) > p`, &data.EnrollmentByGender, chromedp.ByQuery),
		chromedp.Text(`#school-section > div > div:nth-child(1) > div:nth-child(2) > p`, &data.CalendarSystem, chromedp.ByQuery),
		chromedp.Text(`#school-section > div > div:nth-child(2) > div:nth-child(2) > p`, &data.RetentionRate, chromedp.ByQuery),
		chromedp.Text(`#school-section > div > div:nth-child(3) > div:nth-child(2) > p`, &data.OnCampusHousing, chromedp.ByQuery),
	)
	if err != nil {
		// ignore error
		//DEBUG:
		fmt.Println(err)
	}
	data.UndergradEnrollment = strings.ReplaceAll(data.UndergradEnrollment, ",", "")
	data.HeadCoach = strings.ReplaceAll(data.HeadCoach, "\nSend Message", "")
	data.AssistantCoach = strings.ReplaceAll(data.AssistantCoach, "\nSend Message", "")

	*details = append(*details, data)
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
