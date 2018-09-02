package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gomodule/redigo/redis"
)

var (
	name    string
	version string
	gitSHA  string
)

// positionURL is the URL to get the ISS position
const positionURL = "http://api.open-notify.org/iss-now.json"

// issPosition contains the latitude and longitude of the ISS
type issPosition struct {
	Latitude  string `json:"longitude"`
	Longitude string `json:"latitude"`
}

// response prepresents the returned JSON from the
// ISS call
type response struct {
	ISSPosition *issPosition `json:"iss_position"`
	Message     string       `json:"message"`
	Timestamp   int64        `json:"timestamp"`
}

// position gets the current position of the ISS
func position(hc *http.Client) (*issPosition, error) {
	req, err := http.NewRequest(http.MethodGet, positionURL, nil)
	if err != nil {
		return nil, err
	}
	res, err := hc.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	var r response
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return nil, err
	}
	return r.ISSPosition, nil
}

const usage = `version: %s - git: %s
Usage: %s [-p port] [-l domain]

Options:
  -h            this help
  -v            show version and exit
  -i            api query interval in seconds
  -s            size of geofence in meters
  -t            tile38 address & port. format: "127.0.0.1:9851"
  -l            location for geofence. format: 33.4484,112.0740

Examples: 
  %[3]s -t localhost:9851 -i 10 -s 5000       run the server
`

func main() {
	var vers bool
	var intervalFlag int
	var geofenseSizeFlag string
	var tile38AddressFlag string
	var locationFlag string

	flag.Usage = func() {
		w := os.Stderr
		for _, arg := range os.Args {
			if arg == "-h" {
				w = os.Stdout
				break
			}
		}
		fmt.Fprintf(w, usage, version, gitSHA, name)
	}

	flag.BoolVar(&vers, "v", false, "")
	flag.IntVar(&intervalFlag, "i", 60, "")
	flag.StringVar(&geofenseSizeFlag, "s", "", "")
	flag.StringVar(&tile38AddressFlag, "t", "", "")
	flag.StringVar(&locationFlag, "l", "", "")
	flag.Parse()

	if vers {
		fmt.Fprintf(os.Stdout, "version: %s\n", version)
		return
	}

	if tile38AddressFlag == "" {
		fmt.Println("error: tile38 address required")
		os.Exit(1)
	}

	// make sure the Tile38 address is formatted correcrtly
	if !strings.Contains(tile38AddressFlag, ":") {
		fmt.Println("error: tile38 address format incorrect")
		os.Exit(1)
	}

	// connect to Tile38
	c, err := redis.Dial("tcp", tile38AddressFlag)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer c.Close()

	// make sure the given location is formatted correctly
	if !strings.Contains(locationFlag, ",") {
		fmt.Println("error: location format requires a ',' seperating lat and lon")
		os.Exit(1)
	}

	loc := strings.Split(locationFlag, ",")

	// set the geofence
	_, err = c.Do("NEARBY", "earth-orbit", "FENCE", "DETECT", "enter", "POINT", loc[0], loc[1], geofenseSizeFlag)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	hc := http.Client{
		Timeout: 10 * time.Second,
	}

	ticker := time.NewTicker(time.Second * time.Duration(intervalFlag))
	defer ticker.Stop()

	for range ticker.C {
		pos, err := position(&hc)
		if err != nil {
			fmt.Println(err)
			continue
		}
		_, err = c.Do("SET", "earth-orbit", "iss", "POINT", pos.Latitude, pos.Longitude)
		if err != nil {
			fmt.Println(err)
			continue
		}
	}
}
