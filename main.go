package main

// To run:
// go run main.go > results.txt 2>&1

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"math"
	"math/rand"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
	"golang.org/x/exp/slices"
)

// const CONFIG_FILE = "config.json"
// const CONFIG_FILE = "config_all.json"
// const CONFIG_FILE = "config_spokane300.json"
// const CONFIG_FILE = "config_spokaneStates.json"
// const CONFIG_FILE = "config_spokaneCoaches.json"
// const CONFIG_FILE = "config_chicago300.json"
// const CONFIG_FILE = "config_chicago100.json"
// const CONFIG_FILE = "config_reno300.json"
// const CONFIG_FILE = "config_florida.json"
// const CONFIG_FILE = "config_kip.json"        // ???
// const CONFIG_FILE = "config_kip2.json"       // ???
// const CONFIG_FILE = "config_kipSpokane.json" // ???
// const CONFIG_FILE = "config_kipCamp.json"    // ???
// const CONFIG_FILE = "config_16.json"
const CONFIG_FILE = "config_16.json"

var global_configfile = ""

type Configuration struct {
	OpenChromedp  bool     `json:"open_chromedp"`
	ParseMap      bool     `json:"parse_map"`
	ParseStates   bool     `json:"parse_states"`
	StatesList    []string `json:"stateslist"`
	ParseColleges bool     `json:"parse_colleges"`

	ExportColleges     bool   `json:"export_colleges"`
	ExportCollegesFile string `json:"export_colleges_file"`
	ImportColleges     bool   `json:"import_colleges"`
	ImportCollegesFile string `json:"import_colleges_file"`

	ParseCollegePages bool     `json:"parse_college_pages"`
	CollegeList       []string `json:"collegelist"`
	LevelList         []string `json:"levellist"`

	ParseLocation                    bool   `json:"parse_location"`
	ParseLocationName                string `json:"parse_location_name"`
	ParseLocationLatitudeLogitude    string `json:"parse_location_latitude_logitude"`
	ParseLocationLatitudeRadiusMiles string `json:"parse_location_radius_miles"`

	ImportCollegeDetails     bool   `json:"import_college_details"`
	ImportCollegeDetailsFile string `json:"import_college_details_file"`

	ParseLatitudeLogitude bool     `json:"parse_lat_long"`
	LatitudeLogitudeList  []string `json:"lat_long_list"`

	ExportCollegeDetails     bool   `json:"export_college_details"`
	ExportCollegeDetailsFile string `json:"export_college_details_file"`

	ExportCollegeDetailsText     bool   `json:"export_college_details_text"`
	ExportCollegeDetailsTextFile string `json:"export_college_details_text_file"`

	ExportCollegeDetailsHtml     bool   `json:"export_college_details_html"`
	ExportCollegeDetailsHtmlFile string `json:"export_college_details_html_file"`

	//DEBUG
	DumpStates   bool `json:"dump_states"`
	DumpColleges bool `json:"dump_colleges"`
}

var appConfig Configuration

const MAP_URL = "https://sportsrecruits.com/athletic-scholarships/womens-volleyball"
const LATLONG_URL = "https://www.findlatitudeandlongitude.com"

const NAV_TIME_MAX_STATE = 3
const NAV_TIME_MAX_COLLEGE = 3

type State struct {
	name      string
	stateLink string
}

type College struct {
	Name              string `json:"name"`
	State             string `json:"state"`
	City              string `json:"city"`
	Level             string `json:"level"`
	CollegeLink       string `json:"college_link"`
	StateLink         string `json:"state_link"`
	LogoLink          string `json:"logo_link"`
	LatitudeLongitude string `json:"latitude_logitude"`
}

type CollegeDetail struct {
	Name               string `json:"name"`
	State              string `json:"state"`
	City               string `json:"city"`
	Level              string `json:"level"`
	CollegeLink        string `json:"college_link"`
	StateLink          string `json:"state_link"`
	LogoLink           string `json:"logo_link"`
	GoogleLink         string `json:"google_link"`
	GoogleMascotLink   string `json:"google_mascot_link"`
	GoogleRosterLink   string `json:"google_roster_link"`
	GoogleCoachesLink  string `json:"google_coaches_link"`
	GoogleScheduleLink string `json:"google_schedule_link"`
	WikipediaLink      string `json:"wikipedia_link"`

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
	GraduationRate     string `json:"graduation_rate"`
	EnrollmentByGender string `json:"enrollment_by_gender"`
	CalendarSystem     string `json:"calendar_system"`
	RetentionRate      string `json:"retention_rate"`
	OnCampusHousing    string `json:"on_campus_housing"`

	AcceptanceRate                    string `json:"acceptance_rate"`
	TotalAppicants                    string `json:"total_applicants"`
	SatStudentsSubmitting             string `json:"sat_students_submitting"`
	SatSReadingRange                  string `json:"sat_reading_range"`
	SatMathRange                      string `json:"sat_math_range"`
	SatWritingRange                   string `json:"sat_writing_range"`
	ActStudentsSubmitting             string `json:"act_students_submitting"`
	ActComposite                      string `json:"act_composite"`
	ActEnglish                        string `json:"act_english"`
	ActMath                           string `json:"act_math"`
	ActWriting                        string `json:"act_writing"`
	RequirementsOpenAdmissionPolicy   string `json:"requirements_open_admission_policy"`
	RequirementsApplicationFee        string `json:"requirements_application_fee"`
	RequirementsRecommendations       string `json:"requirements_recomedations"`
	RequirementsSecondarySchoolRecord string `json:"requirements_secondary_school_record"`
	RequirementsSecondarySchoolRank   string `json:"requirements_secondary_school_rank"`
	RequirementsSecondarySchoolGpa    string `json:"requirements_secondary_school_gpa"`

	CostInStateTotal         string `json:"cost_in_state_total"`
	CostInStateTuition       string `json:"cost_in_state_tuition"`
	CostInStateFee           string `json:"cost_in_state_fee"`
	CostInStateOnCampusRoom  string `json:"cost_in_state_on_campus_room"`
	CostOutStateTotal        string `json:"cost_out_state_total"`
	CostOutStateTuition      string `json:"cost_out_state_tuition"`
	CostOutStateFee          string `json:"cost_out_state_fee"`
	CostOutStateOnCampusRoom string `json:"cost_out_state_on_campus_room"`
	CostPercentUndergradAid  string `json:"cost_percent_undergrad_aid"`

	Majors []string `json:"majors"`

	LatitudeLongitude string `json:"latitude_logitude"`
	Mascot            string `json:"mascot"`
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

func testStatesSkip(name string) bool {
	if slices.Contains(appConfig.StatesList, "All") {
		return false
	}
	if slices.Contains(appConfig.StatesList, "all") {
		return false
	}
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
		var ok bool
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
			chromedp.AttributeValue(`.avatar img`, "src", &data.LogoLink, &ok, chromedp.NodeVisible, chromedp.ByQuery, chromedp.AtLeast(0), chromedp.FromNode(n)),
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

func stateIndex(name string) int {
	for i, state := range STATE_NAMES {
		if name == state {
			return i
		}
	}
	return 100
}

func distance(lat1 float64, lng1 float64, lat2 float64, lng2 float64) float64 {
	radlat1 := float64(math.Pi * lat1 / 180)
	radlat2 := float64(math.Pi * lat2 / 180)
	theta := float64(lng1 - lng2)
	radtheta := float64(math.Pi * theta / 180)
	dist := math.Sin(radlat1)*math.Sin(radlat2) + math.Cos(radlat1)*math.Cos(radlat2)*math.Cos(radtheta)
	if dist > 1 {
		dist = 1
	}
	dist = math.Acos(dist)
	dist = dist * 180 / math.Pi
	dist = dist * 60 * 1.1515
	return dist
}

func checkCollegeLat(colleges *[]College) {
	// gather average of lat, long per state
	var avgLats [52]float64
	var avgLngs [52]float64
	var avgCnts [52]int
	for i := 0; i < 52; i++ {
		avgLats[i] = 0.0
		avgLngs[i] = 0.0
		avgCnts[i] = 0
	}
	for _, college := range *colleges {
		index := stateIndex(college.State)
		parts := strings.Split(college.LatitudeLongitude, ",")
		latStr := strings.TrimSpace(parts[0])
		lngStr := strings.TrimSpace(parts[1])
		lat, _ := strconv.ParseFloat(latStr, 64)
		lng, _ := strconv.ParseFloat(lngStr, 64)
		// WARNING: assume US where lat > 0 and lng < 0
		if lat > 0 && lng < 0 {
			avgLats[index] += lat
			avgLngs[index] += lng
			avgCnts[index] += 1
		}
		// DEBUG
		// if college.State == "Alabama" {
		// 	fmt.Printf("%-40s: %s, %s; %f  %f  %d\n", college.Name, latStr, lngStr, avgLats[index], avgLngs[index], avgCnts[index])
		// }
	}
	for i := 0; i < 52; i++ {
		if avgCnts[i] > 0 {
			avgLats[i] = avgLats[i] / float64(avgCnts[i])
			avgLngs[i] = avgLngs[i] / float64(avgCnts[i])
		}
	}

	// look for outliers compared to averages
	// print error and google search (ex. https://www.google.com/search?q=latitude+Vassar+College)
	for _, college := range *colleges {
		index := stateIndex(college.State)
		parts := strings.Split(college.LatitudeLongitude, ",")
		latStr := strings.TrimSpace(parts[0])
		lngStr := strings.TrimSpace(parts[1])
		lat, _ := strconv.ParseFloat(latStr, 64)
		lng, _ := strconv.ParseFloat(lngStr, 64)
		distMiles := distance(lat, lng, avgLats[index], avgLngs[index])
		// WARNING: assume US where lat > 0 and lng < 0
		if distMiles > 500 || lat < 0 || lng > 0 {
			// DEBUG
			// if college.State != "Alabama" {
			// 	continue
			// }
			fmt.Printf("bad lat/lng:\n")
			fmt.Printf("  name      : %s\n", college.Name)
			fmt.Printf("  current   : %f, %f\n", lat, lng)
			fmt.Printf("  state avg : %f, %f\n", avgLats[index], avgLngs[index])
			fmt.Printf("  distance  : %f\n", distMiles)
			url := "https://www.google.com/search?q=latitude+" + strings.Replace(college.Name, " ", "+", -1)
			fmt.Printf("  google lat: %s\n", url)
		}
	}
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

func testLevelSkip(name string) bool {
	if slices.Contains(appConfig.LevelList, "All") {
		return false
	}
	if slices.Contains(appConfig.LevelList, "all") {
		return false
	}
	if slices.Contains(appConfig.LevelList, name) {
		return false
	}
	return true
}

func testLocationSkip(latlong string) bool {
	if appConfig.ParseLocation && len(latlong) > 0 {
		centerParts := strings.Split(appConfig.ParseLocationLatitudeLogitude, ",")
		centerLatStr := strings.TrimSpace(centerParts[0])
		centerLngStr := strings.TrimSpace(centerParts[1])
		centerLat, _ := strconv.ParseFloat(centerLatStr, 64)
		centerLng, _ := strconv.ParseFloat(centerLngStr, 64)
		radiusMiles, _ := strconv.ParseFloat(appConfig.ParseLocationLatitudeRadiusMiles, 64)

		givenParts := strings.Split(latlong, ",")
		givenLatStr := strings.TrimSpace(givenParts[0])
		givenLngStr := strings.TrimSpace(givenParts[1])
		givenLat, _ := strconv.ParseFloat(givenLatStr, 64)
		givenLng, _ := strconv.ParseFloat(givenLngStr, 64)
		distMiles := distance(centerLat, centerLng, givenLat, givenLng)
		fmt.Println(latlong)
		fmt.Println(distMiles)
		return distMiles > radiusMiles
	}
	return false
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
	data.LogoLink = college.LogoLink
	data.LatitudeLongitude = college.LatitudeLongitude
	temp := data.Name
	data.GoogleLink = "https://www.google.com/search?q=" + url.QueryEscape(temp)
	data.GoogleMascotLink = "https://www.google.com/search?q=" + url.QueryEscape(temp+" mascot")
	data.GoogleRosterLink = "https://www.google.com/search?q=" + url.QueryEscape(temp+" volleyball roster")
	data.GoogleCoachesLink = "https://www.google.com/search?q=" + url.QueryEscape(temp+" volleyball coaches")
	data.GoogleScheduleLink = "https://www.google.com/search?q=" + url.QueryEscape(temp+" volleyball schedule")
	temp = strings.ReplaceAll(temp, " ", "_")
	temp = strings.ReplaceAll(temp, "&", "%26")
	data.WikipediaLink = "https://en.wikipedia.org/wiki/" + temp

	var majorNodes []*cdp.Node

	// <img class="profile-photo-college" src="https://cdn2-sr-application.sportsrecruits.com/file/images/lacrosserecruits/2015/25_arizona_state_university.png" alt="">

	err = chromedp.Run(*ctx,
		chromedp.Navigate(college.CollegeLink),
		chromedp.Sleep(2*time.Second),
		chromedp.Sleep(random_delay),
		//DEBUG: chromedp.Sleep(400*time.Second),
		chromedp.Text(`#quick-facts-section > div > p:nth-child(2)`, &data.Conference, chromedp.ByQuery, chromedp.AtLeast(0)),
		chromedp.Text(`#quick-facts-section > div > p:nth-child(4)`, &data.AcademicSelectivity, chromedp.ByQuery, chromedp.AtLeast(0)),
		chromedp.Text(`#quick-facts-section > div > p:nth-child(6)`, &data.UndergradEnrollment, chromedp.ByQuery, chromedp.AtLeast(0)),
		chromedp.Text(`#quick-facts-section > div > p:nth-child(8)`, &data.ControlAffiliation, chromedp.ByQuery, chromedp.AtLeast(0)),
		// TODO: is overview useful?
		chromedp.Click(`.toggle-content`),
		chromedp.Sleep(1*time.Second),
		chromedp.Text(`.read-more-container`, &data.Overview, chromedp.ByQuery, chromedp.AtLeast(0)),
	)
	if err != nil {
		// ignore error
		//DEBUG:
		fmt.Println(err)
	}

	// found schools with only 1 coach listed, seperate Run so remaining fields are parsed
	err = chromedp.Run(*ctx,
		chromedp.Text(`#athletics-section > div > div:nth-child(6) > div:nth-child(1) > p`, &data.HeadCoach, chromedp.ByQuery, chromedp.AtLeast(0)),
		chromedp.Text(`#athletics-section > div > div:nth-child(6) > div:nth-child(2) > p`, &data.AssistantCoach, chromedp.ByQuery, chromedp.AtLeast(0)),
	)
	if err != nil {
		// ignore error
		//DEBUG:
		fmt.Println(err)
	}

	err = chromedp.Run(*ctx,
		chromedp.Text(`#school-section > div > div:nth-child(1) > div:nth-child(1) > p`, &data.StudentRatio, chromedp.ByQuery, chromedp.AtLeast(0)),
		chromedp.Text(`#school-section > div > div:nth-child(2) > div:nth-child(1) > p`, &data.GraduationRate, chromedp.ByQuery, chromedp.AtLeast(0)),
		chromedp.Text(`#school-section > div > div:nth-child(3) > div:nth-child(1) > p`, &data.EnrollmentByGender, chromedp.ByQuery, chromedp.AtLeast(0)),
		chromedp.Text(`#school-section > div > div:nth-child(1) > div:nth-child(2) > p`, &data.CalendarSystem, chromedp.ByQuery, chromedp.AtLeast(0)),
		chromedp.Text(`#school-section > div > div:nth-child(2) > div:nth-child(2) > p`, &data.RetentionRate, chromedp.ByQuery, chromedp.AtLeast(0)),
		chromedp.Text(`#school-section > div > div:nth-child(3) > div:nth-child(2) > p`, &data.OnCampusHousing, chromedp.ByQuery, chromedp.AtLeast(0)),

		chromedp.Text(`#admissions-section > div > div:nth-child(1) > div:nth-child(1) > p`, &data.AcceptanceRate, chromedp.ByQuery, chromedp.AtLeast(0)),
		chromedp.Text(`#admissions-section > div > div:nth-child(1) > div:nth-child(2) > p`, &data.TotalAppicants, chromedp.ByQuery, chromedp.AtLeast(0)),
		chromedp.Text(`#admissions-section > div > div:nth-child(5) > div > p`, &data.SatStudentsSubmitting, chromedp.ByQuery, chromedp.AtLeast(0)),
		chromedp.Text(`#admissions-section > div > div:nth-child(6) > div:nth-child(1) > p`, &data.SatSReadingRange, chromedp.ByQuery, chromedp.AtLeast(0)),
		chromedp.Text(`#admissions-section > div > div:nth-child(6) > div:nth-child(2) > p`, &data.SatMathRange, chromedp.ByQuery, chromedp.AtLeast(0)),
		chromedp.Text(`#admissions-section > div > div:nth-child(6) > div:nth-child(3) > p`, &data.SatWritingRange, chromedp.ByQuery, chromedp.AtLeast(0)),
		chromedp.Text(`#admissions-section > div > div:nth-child(8) > div > p`, &data.ActStudentsSubmitting, chromedp.ByQuery, chromedp.AtLeast(0)),

		chromedp.Text(`#admissions-section > div > div:nth-child(9) > div:nth-child(1) > p`, &data.ActComposite, chromedp.ByQuery, chromedp.AtLeast(0)),
		chromedp.Text(`#admissions-section > div > div:nth-child(9) > div:nth-child(2) > p`, &data.ActEnglish, chromedp.ByQuery, chromedp.AtLeast(0)),
		chromedp.Text(`#admissions-section > div > div:nth-child(9) > div:nth-child(3) > p`, &data.ActMath, chromedp.ByQuery, chromedp.AtLeast(0)),
		chromedp.Text(`#admissions-section > div > div:nth-child(9) > div:nth-child(4) > p`, &data.ActWriting, chromedp.ByQuery, chromedp.AtLeast(0)),
		chromedp.Text(`#admissions-section > div > div:nth-child(12) > div:nth-child(1) > p`, &data.RequirementsOpenAdmissionPolicy, chromedp.ByQuery, chromedp.AtLeast(0)),
		chromedp.Text(`#admissions-section > div > div:nth-child(12) > div:nth-child(2) > p`, &data.RequirementsApplicationFee, chromedp.ByQuery, chromedp.AtLeast(0)),
		chromedp.Text(`#admissions-section > div > div:nth-child(12) > div:nth-child(3) > p`, &data.RequirementsRecommendations, chromedp.ByQuery, chromedp.AtLeast(0)),
		chromedp.Text(`#admissions-section > div > div:nth-child(13) > div:nth-child(1) > p`, &data.RequirementsSecondarySchoolRecord, chromedp.ByQuery, chromedp.AtLeast(0)),
		chromedp.Text(`#admissions-section > div > div:nth-child(13) > div:nth-child(3) > p`, &data.RequirementsSecondarySchoolRank, chromedp.ByQuery, chromedp.AtLeast(0)),
		chromedp.Text(`#admissions-section > div > div:nth-child(13) > div:nth-child(2) > p`, &data.RequirementsSecondarySchoolGpa, chromedp.ByQuery, chromedp.AtLeast(0)),

		chromedp.Text(`#cost-section > div > div:nth-child(3) > div:nth-child(1) > p`, &data.CostInStateTotal, chromedp.ByQuery, chromedp.AtLeast(0)),
		chromedp.Text(`#cost-section > div > div:nth-child(3) > div:nth-child(2) > p`, &data.CostInStateTuition, chromedp.ByQuery, chromedp.AtLeast(0)),
		chromedp.Text(`#cost-section > div > div:nth-child(3) > div:nth-child(3) > p`, &data.CostInStateFee, chromedp.ByQuery, chromedp.AtLeast(0)),
		chromedp.Text(`#cost-section > div > div:nth-child(3) > div:nth-child(4) > p`, &data.CostInStateOnCampusRoom, chromedp.ByQuery, chromedp.AtLeast(0)),
		chromedp.Text(`#cost-section > div > div:nth-child(6) > div:nth-child(1) > p`, &data.CostOutStateTotal, chromedp.ByQuery, chromedp.AtLeast(0)),
		chromedp.Text(`#cost-section > div > div:nth-child(6) > div:nth-child(2) > p`, &data.CostOutStateTuition, chromedp.ByQuery, chromedp.AtLeast(0)),
		chromedp.Text(`#cost-section > div > div:nth-child(6) > div:nth-child(3) > p`, &data.CostOutStateFee, chromedp.ByQuery, chromedp.AtLeast(0)),
		chromedp.Text(`#cost-section > div > div:nth-child(6) > div:nth-child(4) > p`, &data.CostOutStateOnCampusRoom, chromedp.ByQuery, chromedp.AtLeast(0)),
		chromedp.Text(`#cost-section > div > div:nth-child(9) > div > p`, &data.CostPercentUndergradAid, chromedp.ByQuery, chromedp.AtLeast(0)),
		chromedp.Nodes(`p.major`, &majorNodes, chromedp.ByQueryAll, chromedp.AtLeast(0)),
	)
	if err != nil {
		// ignore error
		//DEBUG:
		fmt.Println(err)
	}

	for _, n := range majorNodes {
		value := strings.TrimSpace(n.Children[0].NodeValue)
		if len(value) > 0 {
			data.Majors = append(data.Majors, value)
		}
	}

	data.UndergradEnrollment = strings.ReplaceAll(data.UndergradEnrollment, ",", "")
	data.HeadCoach = strings.ReplaceAll(data.HeadCoach, "\nSend Message", "")
	data.AssistantCoach = strings.ReplaceAll(data.AssistantCoach, "\nSend Message", "")

	*details = append(*details, data)
}

func parseForCollegeLatitudeLogitude(ctx *context.Context, colleges *[]College, index int) {
	var err error
	n := rand.Intn(NAV_TIME_MAX_COLLEGE)
	random_delay := time.Duration(n) * time.Second

	name := (*colleges)[index].Name
	state := (*colleges)[index].State
	fmt.Println(name)
	location := ""

	if len((*colleges)[index].LatitudeLongitude) > 0 {
		fmt.Println("lat exists")
		return
	}

	err = chromedp.Run(*ctx,
		chromedp.Navigate(LATLONG_URL),
		chromedp.Sleep(5*time.Second),
		chromedp.Sleep(random_delay),
		chromedp.SendKeys(`#search_box > input.address`, name+", "+state+", US", chromedp.ByID),
		chromedp.Sleep(1*time.Second),
		chromedp.Click(`#search_box > input.big_button`),
		chromedp.Sleep(60*time.Second),
		chromedp.Text(`#search_results > span`, &location, chromedp.ByQuery),
	)
	if err != nil {
		// ignore error
		//DEBUG:
		fmt.Println(err)
	}
	if len(location) > 0 {
		parts := strings.Split(location, "\n")
		if len(parts) > 5 {
			location = parts[5]
			location = strings.ReplaceAll(location, "\u00B0", "")
		}
	}
	fmt.Println(location)
	(*colleges)[index].LatitudeLongitude = location
}

func parseForDetailLatitudeLogitude(ctx *context.Context, details *[]CollegeDetail, index int) {
	var err error
	n := rand.Intn(NAV_TIME_MAX_COLLEGE)
	random_delay := time.Duration(n) * time.Second

	name := (*details)[index].Name
	state := (*details)[index].State
	fmt.Println(name)
	location := ""

	if len((*details)[index].LatitudeLongitude) > 0 {
		fmt.Println("lat exists")
		return
	}

	err = chromedp.Run(*ctx,
		chromedp.Navigate(LATLONG_URL),
		chromedp.Sleep(5*time.Second),
		chromedp.Sleep(random_delay),
		chromedp.SendKeys(`#search_box > input.address`, name+", "+state+", US", chromedp.ByID),
		chromedp.Sleep(1*time.Second),
		chromedp.Click(`#search_box > input.big_button`),
		chromedp.Sleep(60*time.Second),
		chromedp.Text(`#search_results > span`, &location, chromedp.ByQuery),
	)
	if err != nil {
		// ignore error
		//DEBUG:
		fmt.Println(err)
	}
	if len(location) > 0 {
		parts := strings.Split(location, "\n")
		if len(parts) > 5 {
			location = parts[5]
			location = strings.ReplaceAll(location, "\u00B0", "")
		}
	}
	(*details)[index].LatitudeLongitude = location
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

func importCollegeDetails(details *[]CollegeDetail) {
	if len(appConfig.ImportCollegeDetailsFile) == 0 {
		return
	}
	fileName := appConfig.ImportCollegeDetailsFile
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
		err = json.Unmarshal(data, &details)
		if err != nil {
			panic(err)
		}
	}
	fmt.Println(len(*details))
}

func exportCollegeDetailsText(details *[]CollegeDetail) {
	if len(appConfig.ExportCollegeDetailsTextFile) == 0 {
		return
	}
	fileName := appConfig.ExportCollegeDetailsTextFile
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
	msg := string(b)
	msg = strings.ReplaceAll(msg, "    \"", "")
	//msg = strings.ReplaceAll(msg, "name\":", " - ")
	msg = strings.ReplaceAll(msg, "\",", "")
	msg = strings.ReplaceAll(msg, "\"", "")
	msg = strings.ReplaceAll(msg, "[", "")
	msg = strings.ReplaceAll(msg, "],", "")
	msg = strings.ReplaceAll(msg, "]", "")
	msg = strings.ReplaceAll(msg, "{", "")
	msg = strings.ReplaceAll(msg, "},", "")
	msg = strings.ReplaceAll(msg, "}", "")
	msg = strings.ReplaceAll(msg, "null,", "")

	file.WriteString(msg)
}

func detailFromName(details *[]CollegeDetail, name string) CollegeDetail {
	result := CollegeDetail{}
	for _, detail := range *details {
		if name == detail.Name {
			result = detail
		}
	}
	return result
}

func anyLevel(level string, details *[]CollegeDetail) bool {
	result := false
	for _, detail := range *details {
		if level == detail.Level {
			result = true
		}
	}
	return result
}

func unique(stringSlice []string) []string {
	keys := make(map[string]bool)
	list := []string{}
	for _, entry := range stringSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

func statesUsed(details *[]CollegeDetail) []string {
	var result []string
	for _, detail := range *details {
		result = append(result, detail.State)
	}
	return unique(result)
}

func exportCollegeDetailsHtml(details *[]CollegeDetail) {
	if len(appConfig.ExportCollegeDetailsHtmlFile) == 0 {
		return
	}
	fileName := appConfig.ExportCollegeDetailsHtmlFile
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
	msg := string(b)
	msg = strings.ReplaceAll(msg, "    \"", "")
	//msg = strings.ReplaceAll(msg, "name\":", " - ")
	msg = strings.ReplaceAll(msg, "\",", "")
	msg = strings.ReplaceAll(msg, "\"", "")
	msg = strings.ReplaceAll(msg, "[", "")
	msg = strings.ReplaceAll(msg, "],", "")
	msg = strings.ReplaceAll(msg, "]", "")
	msg = strings.ReplaceAll(msg, "{", "")
	msg = strings.ReplaceAll(msg, "},", "")
	msg = strings.ReplaceAll(msg, "}", "")
	msg = strings.ReplaceAll(msg, "null,", "")
	lines := strings.Split(msg, "\n")

	msg2 := ""
	msg2 += "<html>\n"
	msg2 += "  <head>\n"
	msg2 += "    <link rel=\"stylesheet\" href=\"style.css\" />\n"

	msg2 += "    <script defer src=\"https://cdn.jsdelivr.net/npm/alpinejs@3.x.x/dist/cdn.min.js\"></script>\n"
	msg2 += "  </head>\n"

	msg2 += "  <body>\n"
	title := global_configfile
	title = strings.Replace(title, ".json", "", -1)
	title = strings.Replace(title, "config_", "", -1)
	msg2 += "  <h2>" + title + "</h2>\n"
	title = strconv.Itoa(len(*details))
	msg2 += "  <div class=\"run-count\">(school count: " + title + ")</div>\n"
	currentTime := time.Now()
	title = currentTime.Format("2006-01-02 3:4 pm")
	msg2 += "  <div class=\"run-date\">(" + title + ")</div>\n"
	msg2 += "    <ul>\n"
	msg2 += "<div x-data=\"{ openall: false }\">\n"
	msg2 += "<div x-data=\"{ openDI: true }\">\n"
	msg2 += "<div x-data=\"{ openDII: true }\">\n"
	msg2 += "<div x-data=\"{ openDIII: true }\">\n"
	msg2 += "<div x-data=\"{ openNAIA: true }\">\n"
	msg2 += "<div x-data=\"{ openNJCAA: true }\">\n"

	states := statesUsed(details)
	for _, state := range states {
		s := strings.ReplaceAll(state, " ", "")
		msg2 += "<div x-data=\"{ open" + s + ": true }\">\n"
	}

	msg2 += "<button class=\"home-button\"><a href=\"index.html\">home</a></button>\n"
	msg2 += "<button class=\"expand-button\" @click=\"openall = !openall\" :class=\"{ 'active': openall }\">expand all...</button>\n"

	if anyLevel("NCAA DI", details) {
		msg2 += "<button class=\"menu-button\" @click=\"openDI = !openDI\" :class=\"{ 'active': openDI }\">NCAA DI...</button>\n"
	}
	if anyLevel("NCAA DII", details) {
		msg2 += "<button class=\"menu-button\" @click=\"openDII = !openDII\" :class=\"{ 'active': openDII }\">NCAA DII...</button>\n"
	}
	if anyLevel("NCAA DIII", details) {
		msg2 += "<button class=\"menu-button\" @click=\"openDIII = !openDIII\" :class=\"{ 'active': openDIII }\">NCAA DIII...</button>\n"
	}
	if anyLevel("NAIA", details) {
		msg2 += "<button class=\"menu-button\" @click=\"openNAIA = !openNAIA\" :class=\"{ 'active': openNAIA }\">NAIA...</button>\n"
	}
	if anyLevel("NJCAA", details) {
		msg2 += "<button class=\"menu-button\" @click=\"openNJCAA = !openNJCAA\" :class=\"{ 'active': openNJCAA }\">NJCAA...</button>\n"
	}
	for _, state := range states {
		s := strings.ReplaceAll(state, " ", "")
		msg2 += "<button class=\"menu-button\" @click=\"open" + s + " = !open" + s + "\" :class=\"{ 'active': open" + s + " }\">" + s + "...</button>\n"
	}

	indent := ""
	count := 0
	inDetails := false
	inMajors := false
	pendingName := ""
	pendingState := ""
	pendingCity := ""
	pendingLevel := ""
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if len(line) == 0 {
			continue
		}
		if strings.Contains(line, "logo_link:") {
			continue
		}
		if strings.HasPrefix(line, "name:") {
			value := strings.Replace(line, "name:", "", -1)
			value = strings.TrimSpace(value)
			pendingName = value

			inDetails = false
			// fmt.Println(line)
			count += 1
			if count > 1 {
				// close previous list
				msg2 += "    </ul>\n"
				// close previous level filter div
				msg2 += "</div>\n"
			}
		} else if strings.HasPrefix(line, "state:") {
			value := strings.Replace(line, "state:", "", -1)
			value = strings.TrimSpace(value)
			pendingState = value
		} else if strings.HasPrefix(line, "city:") {
			value := strings.Replace(line, "city:", "", -1)
			value = strings.TrimSpace(value)
			pendingCity = value
		} else if strings.HasPrefix(line, "level:") {
			value := strings.Replace(line, "level:", "", -1)
			value = strings.TrimSpace(value)
			pendingLevel = value

			// TODO build line with image and name/state/city/level to the right
			detail := detailFromName(details, pendingName)

			// level filter div
			openFlag := detail.Level
			openFlag = strings.Replace(openFlag, "NCAA", "", -1)
			openFlag = strings.TrimSpace(openFlag)
			s := strings.ReplaceAll(detail.State, " ", "")
			openFlag += " && open" + s
			msg2 += "<div x-show=\"open" + openFlag + "\">\n"
			msg2 += "    <ul>\n"

			msg2 += "    <div class=\"school-title-bar\">\n"
			msg2 += "    <div class=\"school-title-image\">\n"

			// add logo
			msg2 += "    <img src=\"" + detail.LogoLink + "\" alt=\"" + detail.Name + "\">\n"

			msg2 += "    </div>\n"
			msg2 += "    <div class=\"school-title-details\">\n"

			// add pending name/state/city
			msg2 += "    <div><b>name:</b>  " + pendingName + "</div>\n"
			msg2 += "    <div><b>state:</b>  " + pendingState + "</div>\n"
			msg2 += "    <div><b>city:</b>  " + pendingCity + "</div>\n"
			msg2 += "    <div><b>level:</b>  " + pendingLevel + "</div>\n"
			msg2 += "    </div>\n"
			msg2 += "    </div>\n"

			// clean up
			pendingName = ""
			pendingState = ""
			pendingCity = ""
		} else {
			parts := strings.Split(line, ":")
			label := parts[0]
			rest := strings.Join(parts[1:], ":")
			if strings.Contains(line, "_link:") {
				// <a href="aaa" target="_blank">aaa</a>
				label2 := label
				label2 = strings.Replace(label2, "_link", "", -1)
				label2 = strings.Replace(label2, "_", " ", -1)
				// msg2 += "<li><button class=\"menu-button\"><a href=\"" + rest + "\" target=\"_blank\">" + label2 + "^</a></button></li>\n"
				msg2 += "<button class=\"menu-button\"><a href=\"" + rest + "\" target=\"_blank\">" + label2 + "^</a></button>\n"

			} else if strings.Contains(line, ":") {
				if strings.Contains(line, "overview:") {
					// replace line breaks with <br>
					rest = strings.ReplaceAll(rest, "\\n", "<br>")
					rest = strings.ReplaceAll(rest, "<br><br>", "<br>")
					rest = strings.ReplaceAll(rest, "Hide Content", "")
					// only 2nd paragraph is useful
					parts := strings.Split(rest, "<br>")
					if len(parts) > 1 {
						rest = parts[1]
					}
				}

				if strings.Contains(line, "conference:") {
					// details button after wikipedia_link
					msg2 += "<div x-data=\"{ open: false }\">\n"
					msg2 += "<button class=\"menu-button\" @click=\"open = !open\">details...</button>\n"
					inDetails = true
				}

				// process line
				if strings.Contains(line, "majors:") {
					// finish majors list
					msg2 += "</div>\n"
					inDetails = false

					// start majors list
					msg2 += "<div x-data=\"{ open: false }\">\n"
					msg2 += "<button class=\"menu-button\" @click=\"open = !open\">majors...</button>\n"
					inMajors = true
				} else {
					if strings.Contains(line, "mascot:") {
						// skip
					} else if strings.Contains(line, "latitude_logitude:") {
						// skip
						// finish majors list
						msg2 += "</div>\n"
						inMajors = false
					} else {
						if inDetails {
							msg2 += "<li x-show=\"open | openall\"><b>" + label + ":</b> " + rest + "</li>\n"
						} else {
							msg2 += "<li><b>" + label + ":</b> " + rest + "</li>\n"
						}
					}
				}

			} else {
				if inMajors {
					msg2 += "<li x-show=\"open | openall\">" + indent + line + "</li>\n"
				} else {
					msg2 += "<li>" + indent + line + "</li>\n"
				}
			}
			if strings.Contains(line, "majors:") {
				indent = "&nbsp;&nbsp;&nbsp;&nbsp;"
			}
		}
	}

	for _, state := range states {
		s := strings.ReplaceAll(state, " ", "")
		msg2 += "</div> <!-- open" + s + " -->\n"
	}

	msg2 += "</div> <!-- openNJCAA -->\n"
	msg2 += "</div> <!-- openNAIA -->\n"
	msg2 += "</div> <!-- openDIII -->\n"
	msg2 += "</div> <!-- openDII -->\n"
	msg2 += "</div> <!-- openDI -->\n"
	msg2 += "</div> <!-- openall -->\n"

	msg2 += "    </ul>\n"
	msg2 += "    <br><br><br><br><br><br><br><br>\n"
	msg2 += "  </body>\n"
	msg2 += "</html>\n"

	file.WriteString(msg2)
}

func main() {
	// headless := flag.Bool("headless", false, "a bool")
	configfile := flag.String("configfile", CONFIG_FILE, "a string")

	flag.Parse()
	fmt.Println(*configfile)
	global_configfile = *configfile

	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.DisableGPU,
		chromedp.Flag("headless", false),
	)

	fmt.Println("open config...")
	confFile, err := os.Open(global_configfile)
	if err != nil {
		panic(err)
	}
	defer confFile.Close()
	conf, err := ioutil.ReadAll(confFile)
	if err != nil {
		panic(err)
	}
	fmt.Println("unmarshal config...")
	err = json.Unmarshal(conf, &appConfig)
	if err != nil {
		panic(err)
	}
	//DEBUG
	//fmt.Printf("%+v", appConfig)

	ctx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	if appConfig.OpenChromedp {
		fmt.Println("open chromedp...")
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
			name := STATE_NAMES[i]
			if testStatesSkip(name) {
				continue
			}
			if appConfig.ParseColleges {
				fmt.Println("parse colleges...")
				parseForColleges(&ctx, &colleges, state)
			}
		}
		// add seperate step for lat/long after import details
		if appConfig.ParseLatitudeLogitude {
			fmt.Println("parse colleges latitiude longitude...")
			for i, college := range colleges {
				if testLevelSkip(college.Level) {
					continue
				}
				parseForCollegeLatitudeLogitude(&ctx, &colleges, i)
			}
		}
	}

	if appConfig.DumpStates {
		fmt.Println("dump states...")
		dumpStates(&states)
	}
	if appConfig.DumpColleges {
		fmt.Println("dump colleges...")
		dumpColleges(&colleges)
	}
	if appConfig.ImportColleges {
		fmt.Println("import colleges...")
		importColleges(&colleges)
		checkCollegeLat(&colleges)
	}
	// add seperate step for lat/long after import details
	if appConfig.ParseLatitudeLogitude {
		fmt.Println("parse colleges latitiude longitude...")
		for i, college := range colleges {
			if testLevelSkip(college.Level) {
				continue
			}
			parseForCollegeLatitudeLogitude(&ctx, &colleges, i)
		}
	}
	if appConfig.ExportColleges {
		fmt.Println("export colleges...")
		exportColleges(&colleges)
	}
	var details []CollegeDetail
	if appConfig.ParseCollegePages {
		fmt.Println("parse college details...")
		for _, college := range colleges {
			if testStatesSkip(college.State) {
				continue
			}
			if testCollegeSkip(college.Name) {
				continue
			}
			if testLevelSkip(college.Level) {
				continue
			}
			if testLocationSkip(college.LatitudeLongitude) {
				continue
			}
			fmt.Println("details: " + college.Name + "...")
			parseForCollegePages(&ctx, &details, &college)
		}
	}

	if appConfig.ImportCollegeDetails {
		fmt.Println("import college details...")
		importCollegeDetails(&details)
	}

	// add seperate step for lat/long after import details
	if appConfig.ParseLatitudeLogitude {
		fmt.Println("parse college detail latitiude longitude...")
		for i, _ := range details {
			parseForDetailLatitudeLogitude(&ctx, &details, i)
		}
	}

	if appConfig.ExportCollegeDetails {
		fmt.Println("export college details...")
		exportCollegeDetails(&details)
	}

	if appConfig.ExportCollegeDetailsText {
		fmt.Println("export college details text...")
		exportCollegeDetailsText(&details)
	}
	if appConfig.ExportCollegeDetailsHtml {
		fmt.Println("export college details html...")
		exportCollegeDetailsHtml(&details)
	}

	if appConfig.OpenChromedp {
		fmt.Println("close chromedp...")
		cancel()
	}

	fmt.Println("done.")
}
