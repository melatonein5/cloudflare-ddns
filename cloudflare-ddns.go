package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/cloudflare/cloudflare-go"
)

//Global variables
var apiToken string
var api *cloudflare.API
var ip IPQuery

//Structures
type IPQuery struct {
	Status      string  `json:"status"`
	Country     string  `json:"country"`
	CountryCode string  `json:"countryCode"`
	Region      string  `json:"region"`
	RegionName  string  `json:"regionName"`
	City        string  `json:"city"`
	Zip         string  `json:"zip"`
	Lat         float64 `json:"lat"`
	Lon         float64 `json:"lon"`
	Timezone    string  `json:"timezone"`
	Isp         string  `json:"isp"`
	Org         string  `json:"org"`
	As          string  `json:"as"`
	Query       string  `json:"query"`
}

//Functions

//GetAPIToken will search for a system variable and pull the API key to set in the program
func GetAPIToken() {
	apiToken = os.Getenv("CLOUDFLARE_API_TOKEN")
}

//SetAPIToken sets the token as a system variable
func SetAPIToken(token string) {
	os.Setenv("CLOUDFLARE_API_TOKEN", token)
}

//fetchIP will get the users current public IP address
func fetchIP() error {
	//Make a request to an endpoint which will return a user's public IP
	req, err := http.Get("http://ip-api.com/json/")
	if err != nil {
		return err
	}
	defer req.Body.Close()

	//Read the JSON body and unmarshal into an IPQuery
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return err
	}
	json.Unmarshal(body, &ip)

	//No errors, return nil
	return nil
}

func main() {
	//Startup options

	//Set logging location

	//Declare the API
	GetAPIToken()
	api, _ = cloudflare.NewWithAPIToken(apiToken)
}
