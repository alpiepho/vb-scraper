package main

// To run:
// go run main.go > results.txt 2>&1

import (
	"context"
	"fmt"
	"time"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
)

const MAP_URL = "https://sportsrecruits.com/athletic-scholarships/womens-volleyball"

// const BASE_URL = "https://www.keysight.com"
// const PARSEDINFO = "parsedinfo_working.json"
// const INSTRUMENTINFO = "instrumentinfo_working.json"

type State struct {
	name     string
	pageLink string
}

type College struct {
	name       string
	state      string
	city       string
	level      string
	avatarLink string
	pageLink   string
}

func dumpStates(states *[]State) {
	for _, state := range *states {
		fmt.Println("")
		fmt.Println("name:     " + state.name)
		fmt.Println("pageLink: " + state.pageLink)
	}
	fmt.Println(len(*states))
}

func dumpColleges(colleges *[]College) {
	for _, college := range *colleges {
		fmt.Println("")
		fmt.Println("name:       " + college.name)
		fmt.Println("state:      " + college.state)
		fmt.Println("city:       " + college.city)
		fmt.Println("avatarLink: " + college.avatarLink)
		fmt.Println("pageLink:   " + college.pageLink)
	}
	fmt.Println(len(*colleges))
}

func testStatesContinue(i int) bool {
	//DEBUG - limit products
	//return !(i == 1 || i == 2)
	//return i > 40
	return i > 0
	//return false
}

func testCollegesContinue(i int) bool {
	//DEBUG - limit products
	//return !(i == 1 || i == 2)
	//return i > 40
	return i > 0
	//return false
}

func parseForStates(ctx *context.Context, states *[]State) {
	var err error
	err = chromedp.Run(*ctx,
		chromedp.Navigate(MAP_URL),
		chromedp.Sleep(1*time.Second),
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
			chromedp.AttributeValue(`a`, "href", &data.pageLink, &ok, chromedp.NodeVisible, chromedp.ByQuery, chromedp.AtLeast(0), chromedp.FromNode(n)),
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
	err = chromedp.Run(*ctx,
		chromedp.Navigate(state.pageLink),
		chromedp.Sleep(5*time.Second),
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
	fmt.Println(len(collegeLinkNodes))

	for i, n := range collegeLinkNodes {
		//var ok bool
		data := College{}
		data.state = state.name
		fmt.Println(n.Attributes)

		if testCollegesContinue(i) {
			continue
		}

		err := chromedp.Run(*ctx,

			//#main > div.wrapper_row > section:nth-child(4) > div > div > a:nth-child(2) > div.col-xs-12.col-sm-5 > div > p
			//#main > div.wrapper_row > section:nth-child(4) > div > div > a:nth-child(2) > div.col-xs-12.col-sm-4 > p
			//#main > div.wrapper_row > section:nth-child(4) > div > div > a:nth-child(2) > div.col-xs-12.col-sm-3 > p
			chromedp.Text(`.col-sm-5`, &data.name, chromedp.ByQuery, chromedp.FromNode(n)),
			chromedp.Text(`.col-sm-4`, &data.city, chromedp.ByQuery, chromedp.FromNode(n)),
			chromedp.Text(`.col-sm-2`, &data.level, chromedp.ByQuery, chromedp.FromNode(n)),
			//chromedp.AttributeValue(`a`, "href", &data.pageLink, &ok, chromedp.NodeVisible, chromedp.ByQuery, chromedp.AtLeast(0), chromedp.FromNode(n)),
			// chromedp.AttributeValue(`.avatar img`, "src", &data.pageLink, &ok, chromedp.NodeVisible, chromedp.ByQuery, chromedp.AtLeast(0), chromedp.FromNode(n)),
		)
		if err != nil {
			// ignore error
			//DEBUG:
			fmt.Println(err)
		}
		//data.pageLink = BASE_URL + data.pageLink

		//DEBUG:
		fmt.Println(i, ":\t", data.name, " ", data.pageLink)
		*colleges = append(*colleges, data)
	}
	//DEBUG:
	//fmt.Println(len(stateLinkNodes))

}

//#main > div.wrapper_row > section:nth-child(4) > div > div > a:nth-child(2)

// type Product struct {
// 	description string
// 	imageLink   string
// }

// type InfoElement struct {
// 	category    string
// 	productType string
// 	description string
// 	imageLink   string
// }

// func parseForCategories(ctx *context.Context, categories *[]Category) {
// 	var err error
// 	err = chromedp.Run(*ctx,
// 		chromedp.Navigate(BASE_URL+"/us/en/products.html"),
// 		chromedp.Sleep(1*time.Second),
// 		//DEBUG: chromedp.Sleep(400*time.Second),
// 	)
// 	if err != nil {
// 		// ignore error
// 		//DEBUG:
// 		fmt.Println(err)
// 	}

// 	fmt.Println("parse product categories...")
// 	var categoryNodes []*cdp.Node
// 	err = chromedp.Run(*ctx,
// 		chromedp.Nodes(`.link-red`, &categoryNodes, chromedp.ByQueryAll),
// 	)
// 	if err != nil {
// 		// ignore error
// 		//DEBUG:
// 		fmt.Println(err)
// 	}

// 	for i, n := range categoryNodes {
// 		var ok bool
// 		data := Category{}
// 		err := chromedp.Run(*ctx,
// 			chromedp.Text(`p a`, &data.name, chromedp.ByQuery, chromedp.FromNode(n)),
// 			chromedp.AttributeValue(`p a`, "href", &data.pageLink, &ok, chromedp.NodeVisible, chromedp.ByQuery, chromedp.AtLeast(0), chromedp.FromNode(n)),
// 		)
// 		if err != nil {
// 			// ignore error
// 			//DEBUG:
// 			fmt.Println(err)
// 		}
// 		data.pageLink = BASE_URL + data.pageLink

// 		//DEBUG:
// 		fmt.Println(i, ":\t", data.name, " ", data.pageLink)
// 		*categories = append(*categories, data)
// 	}
// 	//DEBUG:
// 	fmt.Println(len(categoryNodes))
// }

// func parseForProductPages(ctx *context.Context, categories *[]Category) {
// 	for i, category := range *categories {
// 		if testContunue(i) {
// 			continue
// 		}
// 		//DEBUG: fmt.Println("parse product page for " + category.name + "...")
// 		var ok bool
// 		var allLink string
// 		err := chromedp.Run(*ctx,
// 			chromedp.Navigate(category.pageLink),
// 			chromedp.Sleep(1*time.Second),
// 			chromedp.AttributeValue(`a.view-catalog-page`, "href", &allLink, &ok, chromedp.NodeVisible, chromedp.ByQuery, chromedp.AtLeast(0)),
// 		)
// 		if err != nil {
// 			// ignore error
// 			//DEBUG: fmt.Println("allLink error1:")
// 			//DEBUG: fmt.Println(err)
// 		}
// 		if len(allLink) > 0 {
// 			(*categories)[i].allLink = BASE_URL + allLink
// 			//DEBUG: fmt.Println(i, ":\t", category.name, " ", category.allLink)
// 		}
// 		if len((*categories)[i].allLink) == 0 {
// 			err = chromedp.Run(*ctx,
// 				chromedp.Navigate(category.pageLink),
// 				chromedp.Sleep(1*time.Second),
// 				chromedp.AttributeValue(`a.btn-view-all`, "href", &allLink, &ok, chromedp.NodeVisible, chromedp.ByQuery, chromedp.AtLeast(0)),
// 			)
// 			if err != nil {
// 				// ignore error
// 				//DEBUG: fmt.Println("allLink error2:")
// 				//DEBUG: fmt.Println(err)
// 			}
// 			if len(allLink) > 0 {
// 				(*categories)[i].allLink = BASE_URL + allLink
// 				//DEBUG: fmt.Println(i, ":\t", category.name, " ", category.allLink)
// 			}
// 		}
// 	}
// }

// func autoScroll(ctx *context.Context) {
// 	var found = true
// 	// until "Load More" button is missing
// 	for found {
// 		//DEBUG
// 		//fmt.Println("click button")
// 		timeoutContext, cancel := context.WithTimeout(*ctx, 5*time.Second)
// 		defer cancel()

// 		var totalString string
// 		err := chromedp.Run(timeoutContext,
// 			chromedp.Text(`.catalog__pagination--total`, &totalString, chromedp.ByQuery),
// 			chromedp.Click(`.pagination-next`, chromedp.ByQuery),
// 			chromedp.Sleep(2*time.Second),
// 		)
// 		//DEBUG
// 		fmt.Println(totalString)
// 		if err != nil {
// 			// ignore error
// 			//fmt.Println(err)
// 			found = false
// 		}
// 	}
// }

// func parseCatalogPage(
// 	ctx *context.Context,
// 	productType string,
// 	i int, category Category,
// 	categories *[]Category,
// 	infoElements *[]InfoElement) {

// 	timeoutContext, cancel := context.WithTimeout(*ctx, 30*time.Second)
// 	defer cancel()
// 	var totalProductsString string
// 	var totalProducts int64
// 	var err error

// 	if productType == "ACTIVE" {
// 		err = chromedp.Run(timeoutContext,
// 			chromedp.Navigate(category.allLink),
// 			chromedp.Sleep(1*time.Second),
// 			chromedp.Text(`.catalog__pagination--total`, &totalProductsString, chromedp.ByQuery),
// 		)
// 	}
// 	if productType == "DISCONTINUED" {
// 		err = chromedp.Run(timeoutContext,
// 			chromedp.Navigate(category.allLink),
// 			chromedp.Sleep(1*time.Second),
// 			chromedp.Click(`catalog__filters--link catalog__disco--link`),
// 			chromedp.Sleep(1*time.Second),
// 			chromedp.Text(`.catalog__pagination--total`, &totalProductsString, chromedp.ByQuery),
// 		)
// 	}
// 	if err != nil {
// 		// ignore error
// 		//DEBUG:
// 		fmt.Println(err)
// 	}

// 	if len(totalProductsString) > 0 {
// 		// Showing 1 - 68 of 68
// 		//                   ^
// 		//DEBUG: fmt.Println("totalString")
// 		//DEBUG: fmt.Println(totalProductsString)
// 		s := strings.Fields(totalProductsString)
// 		totalProducts, _ = strconv.ParseInt(s[5], 0, 32)
// 	}
// 	if productType == "ACTIVE" {
// 		(*categories)[i].totalActive = totalProducts
// 	}
// 	if productType == "DISCONTINUED" {
// 		(*categories)[i].totalDiscontinued = totalProducts
// 	}

// 	//DEBUG:
// 	fmt.Println(totalProducts)

// 	// find all active products
// 	if totalProducts > 0 {

// 		// load all product cards
// 		//DEBUG:  - comment out to go faster
// 		autoScroll(ctx)

// 		//DEBUG:fmt.Println("parse products on category page...")
// 		var productNodes []*cdp.Node
// 		err = chromedp.Run(*ctx,
// 			chromedp.Nodes(`.catalog__products--item`, &productNodes, chromedp.ByQueryAll),
// 		)
// 		if err != nil {
// 			// ignore error
// 			fmt.Println(err)
// 		}

// 		for _, n := range productNodes {
// 			var ok bool
// 			data := Product{}
// 			err := chromedp.Run(*ctx,
// 				chromedp.Text(`.catalog__products--title`, &data.description, chromedp.ByQuery, chromedp.FromNode(n)),
// 				chromedp.AttributeValue(`img`, "src", &data.imageLink, &ok, chromedp.NodeVisible, chromedp.ByQuery, chromedp.AtLeast(0), chromedp.FromNode(n)),
// 			)
// 			if err != nil {
// 				// ignore error
// 				//DEBUG:
// 				fmt.Println(err)
// 			}

// 			//DEBUG: fmt.Println(i, ":\t"+productType+" PRODUCT": ", data.description, "\t", data.imageLink)
// 			infoElement := InfoElement{}
// 			infoElement.category = category.name
// 			infoElement.productType = productType
// 			data.description = strings.Replace(data.description, "\"", "\\\"", -1)
// 			infoElement.description = data.description
// 			infoElement.imageLink = data.imageLink
// 			*infoElements = append(*infoElements, infoElement)

// 		}
// 		//DEBUG:
// 		fmt.Println(len(productNodes))
// 	}

// }

// func dumpAsParsedInfo(infoElements *[]InfoElement) {
// 	fmt.Println("dump " + PARSEDINFO)
// 	file, err := os.Create(PARSEDINFO)
// 	if err != nil {
// 		//fmt.Println(err)
// 	} else {
// 		file.WriteString("{\n")
// 		file.WriteString("  \"RawData\": [\n")
// 		for i, infoElement := range *infoElements {
// 			file.WriteString("    {\n")
// 			file.WriteString("      \"description\": \"" + infoElement.description + "\",\n")
// 			file.WriteString("      \"category\": \"" + infoElement.category + "\",\n")
// 			file.WriteString("      \"productType\": \"" + infoElement.productType + "\",\n")
// 			file.WriteString("      \"imageLink\": \"" + infoElement.imageLink + "\"\n")
// 			file.WriteString("    }")
// 			if i+1 < len(*infoElements) {
// 				file.WriteString(",")
// 			}
// 			file.WriteString("\n")
// 		}
// 		file.WriteString("  ]\n")
// 		file.WriteString("}\n")
// 	}
// 	file.Close()
// }

// func dumpAsInstrumentInfo(infoElements *[]InfoElement) {
// 	fmt.Println("dump " + INSTRUMENTINFO)
// 	file, err := os.Create(INSTRUMENTINFO)
// 	if err != nil {
// 		//fmt.Println(err)
// 	} else {
// 		imagePrefix := "https://keysight-h.assetsadobe.com/is/image/content/dam/keysight/en/img/"
// 		webPrefix := "https://www.keysight.com/us/en/product/"
// 		file.WriteString("{\n")
// 		file.WriteString("  \"RawData\": [\n")
// 		for i, infoElement := range *infoElements {
// 			model := strings.Split(infoElement.description, " ")[0]
// 			imagePath := strings.Replace(infoElement.imageLink, imagePrefix, "", 1)
// 			file.WriteString("    {\n")
// 			file.WriteString("      \"Model\": \"" + model + "\",\n")
// 			file.WriteString("      \"Title\": \"" + infoElement.description + "\",\n")

// 			file.WriteString("      \"Id1\": \"" + model + "\",\n")
// 			file.WriteString("      \"Id2\": \"" + imagePath + "\",\n")
// 			file.WriteString("      \"Web\": \"" + "{T1}" + "\",\n")
// 			file.WriteString("      \"Image\": \"" + "{T2}" + "\",\n")
// 			file.WriteString("      \"Class\": \"" + infoElement.category + "\",\n")
// 			file.WriteString("      \"Width\": \"" + "400" + "\",\n")
// 			file.WriteString("      \"ImageTS\": \"" + "1.41E+12" + "\"\n")
// 			file.WriteString("    }")
// 			if i+1 < len(*infoElements) {
// 				file.WriteString(",")
// 			}
// 			file.WriteString("\n")
// 		}
// 		file.WriteString("  ],\n")

// 		segment := `
//   "SchemaVersion": "1",
//   "T1": "` + webPrefix + `/{Id1}",
//   "T2": "` + imagePrefix + `/{Id2}",
//   "Mitigations": [
//     {
//         "Id": "GpibAutoPollDisable",
//         "Parameters": [
//         {
//             "Name": "Model",
//             "Value": "34410A"
//         }
//       ]
//     },
//     {
//         "Id": "GpibAutoPollDisable",
//         "Parameters": [
//         {
//             "Name": "Model",
//             "Value": "34411A"
//         }
//       ]
//     },
//     {
//         "Id": "GpibAutoPollDisable",
//         "Parameters": [
//         {
//             "Name": "Model",
//             "Value": "34980A"
//         }
//       ]
//     }
//   ]
// `
// 		file.WriteString(segment)
// 		file.WriteString("}\n")
// 	}
// 	file.Close()
// }

func main() {
	//headless := flag.Bool("headless", false, "a bool")

	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.DisableGPU,
		chromedp.Flag("headless", false),
	)
	//var infoElements []InfoElement

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
		if testStatesContinue(i) {
			continue
		}

		fmt.Println("parse colleges...")
		parseForColleges(&ctx, &colleges, state)
	}

	//dumpColleges(&colleges)

	// fmt.Println("parse product pages...")
	// parseForProductPages(&ctx, &categories)

	// for i, category := range categories {
	// 	if testContunue(i) {
	// 		continue
	// 	}

	// 	if len(categories[i].allLink) == 0 {
	// 		continue
	// 	}

	// 	fmt.Println("parse ACTIVE      catalog page for " + category.name + "...")
	// 	// DEBUG:
	// 	fmt.Println(i, ":\t", category.name, " ", category.allLink)
	// 	parseCatalogPage(
	// 		&ctx,
	// 		"ACTIVE",
	// 		i, category,
	// 		&categories,
	// 		&infoElements)

	// 	fmt.Println("parse DISCONTINUE catalog page for " + category.name + "...")
	// 	//DEBUG:
	// 	fmt.Println(i, ":\t", category.name, " ", category.allLink)
	// 	parseCatalogPage(
	// 		&ctx,
	// 		"DISCONTINUED",
	// 		i, category,
	// 		&categories,
	// 		&infoElements)

	// }
	cancel()

	// dumpAsParsedInfo(&infoElements)
	// dumpAsInstrumentInfo(&infoElements)

	fmt.Println("done.")
}
