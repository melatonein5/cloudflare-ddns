package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/cloudflare/cloudflare-go"
)

//Global variables
var apiToken string
var ctx context.Context
var api *cloudflare.API
var zoneId string
var recordId string
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

//Get and Set
//GetAPIToken will search for a system variable and pull the API key to set in the program
func GetAPIToken() {
	apiToken = os.Getenv("CLOUDFLARE_API_TOKEN")
}

//SetAPIToken sets the token as a system variable
func SetAPIToken(token string) {
	os.Setenv("CLOUDFLARE_API_TOKEN", token)
}

func GetZoneID() {
	apiToken = os.Getenv("CLOUDFLARE_ZONE_ID")
}

func GetRecordID() {
	apiToken = os.Getenv("CLOUDFLARE_RECORD_ID")
}

//Main Functions
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

//Compare IP will check if the current IP matches the DNS records and return true if it does not
func CompareIPToRecord(zoneId, recordID string) bool {
	//Make an API request to cloudflare to get the IP address in the A record
	dnsRecord, err := api.DNSRecord(ctx, zoneId, recordID)
	if err != nil {
		log.Println("Could not fetch DNS Record! Is the API key valid? Does it have the right permissions?")
		return false
	}

	//Compare IPs
	if dnsRecord.Data == ip.Query {
		return false
	}
	return true
}

func main() {
	//Startup options

	//Set logging location

	//Declare the API
	GetAPIToken()
	api, _ = cloudflare.NewWithAPIToken(apiToken)

	//Set the context
	ctx = context.Background()
}
