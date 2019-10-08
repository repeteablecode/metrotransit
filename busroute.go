package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strings"
)

// Route represents the data structure of availabe bus routes for today
type Route struct {
	Description string `json:"Description"`
	ProviderID  string `json:"ProviderID"`
	Route       string `json:"Route"`
}

// Stop represents data structure of a routes available stops
type Stop struct {
	Text  string `json:"Text"`
	Value string `json:"Value"`
}

// NexTripDeparture represents the stop details
type NexTripDeparture struct {
	Actual           bool    `json:"Actual"`
	BlockNumber      int64   `json:"BlockNumber"`
	DepartureText    string  `json:"DepartureText"`
	DepartureTime    string  `json:"DepartureTime"`
	Description      string  `json:"Description"`
	Gate             string  `json:"Gate"`
	Route            string  `json:"Route"`
	RouteDirection   string  `json:"RouteDirection"`
	Terminal         string  `json:"Terminal"`
	VehicleHeading   int64   `json:"VehicleHeading"`
	VehicleLatitude  float64 `json:"VehicleLatitude"`
	VehicleLongitude float64 `json:"VehicleLongitude"`
}

// errorMsg is the error message presented when detecting error
var errorMsg = "Arguments Error, Please Format Command as below. Goodbye\n\tCOMMAND: go run casestudy.go \"BUS ROUTE\" \"BUS STOP NAME\" \"DIRECTION\"\n\t*Note - OS may require escape characters in terminal, example \\\" and \\&"

// main
func main() {
	// test commands
	// go run casestudy.go \"Express - Target - Hwy 252 and 73rd Av P\&R - Mpls\" \"Target North Campus Building F\" \"south\"
	// go run casestudy.go \"METRO Blue Line\" \"Target Field Station Platform 1\" \"south\"

	// get arguments provided by user
	args := procCommands()

	// begin processing inputs to find the next departure time
	begin(args[0], args[1], args[2])
}

// procCommands returns "BUS ROUTE" "BUS STOP" "DIRECTION"
func procCommands() []string {
	//var test1	`\"METRO Blue Line\" \"Target Field Station Platform 1\" \"south\"`
	//var test2	`"Express - Target - Hwy 252 and 73rd Av P&R - Mpls" "Target North Campus Building F" "south"`

	// build strCommand containing the combined arguments inputed by user
	var strCommand = strings.Join(os.Args[1:], " ")

	// check for no double qoutes in string
	if strings.Index(strCommand, "\"") == -1 {
		// print error message
		fmt.Println(errorMsg)
		// return commands based on test string
		os.Exit(0)
		//return parseCommandDetails(`"Express - Target - Hwy 252 and 73rd Av P&R - Mpls" "Target North Campus Building F" "south"`)
	}

	// get command string array based on strCommand
	var commands = parseCommandDetails(strCommand)

	// test of number of arguments match
	if len(commands) != 3 {
		// print error message presenting command format
		fmt.Println(errorMsg)
		// return commands based on test string
		os.Exit(0)
		//return parseCommandDetails(`"Express - Target - Hwy 252 and 73rd Av P&R - Mpls" "Target North Campus Building F" "south"`)
	}

	return commands
}

// parseCommandDetails returns a string array of "BUS ROUTE NAME", "BUS STOP NAME", "DIRECTION"
// input from arguments string provided by command line
func parseCommandDetails(s string) []string {
	// create regexp to capture text inbetween quotes
	var r = regexp.MustCompile(`".*?"`)
	// find all
	ms := r.FindAllString(s, -1)
	// create string array to match ms
	ss := make([]string, len(ms))
	// iterate through strings in ms
	for i, m := range ms {
		// add  string m from ms minus the double quotes at the end
		ss[i] = m[1 : len(m)-1]
	}
	// return commands string array
	return ss
}

// begin starts processing provided BUS ROUTE NAME, BUS STOP NAME and DIRECTION
func begin(BUSROUTE string, BUSSTOPNAME string, DIRECTION string) {
	// display input to user
	//fmt.Println("Input\n  Bus Route: " + BUSROUTE + "\n  Bus Stop Name: " + BUSSTOPNAME + "\n  Direction: " + DIRECTION + "\n")

	// start with finding the route number as string
	var routeNumberString = findRouteNumber(BUSROUTE)

	// if empty string then route isn't running, inform user and exit
	if routeNumberString == "" {
		fmt.Println("Route not running today")
		os.Exit(1)
	}

	// Now to translate the directional to a number
	var routeDirectionalString = translateDirectional(DIRECTION)

	// Check for errors, inform user if necessary and exit
	if routeDirectionalString == "" {
		fmt.Println("Direction " + DIRECTION + " is malformed, please correct")
		os.Exit(1)
	}
	// Now to get the stop string
	routeStopString := findRoutesStopValue(routeNumberString, routeDirectionalString, BUSSTOPNAME)

	// print out next arrival time
	fmt.Println("Departure at " + findNextDepartureTime(routeNumberString, routeDirectionalString, routeStopString) + "\t*No time will show if last stop already departed")
}

// findNextDepartureTime returns the next estimated departure time for a given route number direction and stop
func findNextDepartureTime(routeNumber string, direction string, stop string) string {
	// sample url    https://svc.metrotransit.org/nextrip/901/1/TF12?format=JSON
	// build the get api from information provided
	var url = "https://svc.metrotransit.org/nextrip/" + routeNumber + "/" + direction + "/" + stop + "?format=JSON"
	// create NexTripDeparture array
	var departures []NexTripDeparture
	// download data via GET
	data := download(url)
	// convert raw json data to departures
	err := json.Unmarshal(data, &departures)
	// test if err isn't nil, and inform user of error and exit
	if err != nil {
		fmt.Println("findNextDepartureTime() err: ", err.Error())
		os.Exit(1)
	}
	// at last the departure time is returned
	return departures[0].DepartureText
}

// findRoutesStopValue returns a routes stop number for a give route number and direction
func findRoutesStopValue(routeNumber string, directional string, stopName string) string {
	// create stops struct array for give route number and direction
	stops := getRouteStops(routeNumber, directional)
	// initalize index
	var dex = 0
	// intialize stopValue to empty
	var stopValue = ""
	// go's while loop, just spelt for for some reason
	for dex < len(stops) {
		// if we get a match assign stop value then break the loop
		if strings.ToLower(stops[dex].Text) == strings.ToLower(stopName) {
			stopValue = stops[dex].Value
			break
		}
		// increment dex for go while loop, just spelt for for some reason
		dex++
	}
	// return stopValue
	return stopValue
}

// gerRouteStops returns an array of Stop arguments for a route and direction
func getRouteStops(routeNumber string, directional string) []Stop {
	// sample url    https://svc.metrotransit.org/nextrip/Stops/901/1?format=JSON
	var url = "https://svc.metrotransit.org/nextrip/Stops/" + routeNumber + "/" + directional + "?format=JSON"
	// create Stop struct array
	var stops []Stop
	// get data from url generated from inputs provided
	data := download(url)
	// covert raw json byte data into Stop array
	err := json.Unmarshal(data, &stops)
	// test if err isn't nil, and inform user of error and exit
	if err != nil {
		fmt.Println("getRouteStops() err: ", err.Error())
		os.Exit(1)
	}
	// return Stop struct array
	return stops
}

// translateDirectional returns the numerical value for a direction
// south = 1, east = 2, west = 3, and north = 4
func translateDirectional(direction string) string {
	// switch direction converted to lower case returning number as string
	switch strings.ToLower(direction) {
	case "south":
		return "1"
	case "east":
		return "2"
	case "west":
		return "3"
	case "north":
		return "4"
	default:
		return ""
	}
}

// findRouteNumber searches available routes for the day and returns the routes number as a string
func findRouteNumber(route string) string {
	// get the current routes for the day
	var currentRoutes = getRoutes()
	// initialize index
	var dex = 0
	var routeNumber = ""
	// go's while loop, just spelt for for some reason
	for dex < len(currentRoutes) {
		// if we match the route assign route number and break loop
		if currentRoutes[dex].Description == route {
			routeNumber = currentRoutes[dex].Route
			break
		}
		// increment dex for go's while loop, just spelt for for some reason
		dex++
	}
	// return the route number
	return routeNumber
}

// getRoutes returns an array of Route structures
func getRoutes() []Route {
	// create get url
	var url = "https://svc.metrotransit.org/NexTrip/Routes?format=JSON"
	// initialize Route array structure
	var routes []Route
	// get data from url generated from inputs provided
	data := download(url)
	// covert raw json byte data into routes array
	err := json.Unmarshal(data, &routes)
	// test if err isn't nil, and inform user of error and exit if so
	if err != nil {
		fmt.Println("getRoutes() err", err.Error())
		os.Exit(1)
	}
	// return routes
	return routes
}

// download returns byte array of Get Request
func download(url string) []byte {
	// create response for GET from url provided
	response, err := http.Get(url)
	// test if err isn't nil, and inform user of error and exit if so
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	// close immediately to prevent resource leak
	defer response.Body.Close()
	// create data populated with body of response
	data, err := ioutil.ReadAll(response.Body)
	// test if err isn't nil, and inform user of error and exit if so
	if err != nil {
		fmt.Println("iotil.ReadAll error", err.Error())
		os.Exit(1)
	}
	// return byte array data
	return data
}
