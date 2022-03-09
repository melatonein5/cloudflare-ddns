package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

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

//GetZoneID gets the zone ID set in the system variables
func GetZoneID() {
	zoneId = os.Getenv("CLOUDFLARE_ZONE_ID")
}

//SetZoneID takes in a Zone ID and sets it as a system variable
func SetZoneID(zoneID string) {
	os.Setenv("CLOUDFLARE_ZONE_ID", zoneID)
}

//GetRecordID gets the record ID set in the system variables
func GetRecordID() {
	recordId = os.Getenv("CLOUDFLARE_RECORD_ID")
}

//SetRecordID sets the record ID to the provided parameter
func SetRecordID(recordID string) {
	os.Setenv("CLOUDFLARE_ZONE_ID", recordID)
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
func CompareIPToRecord() (cloudflare.DNSRecord, bool) {
	var dnsRecord cloudflare.DNSRecord
	//Make an API request to cloudflare to get the IP address in the A record
	dnsRecord, err := api.DNSRecord(ctx, zoneId, recordId)
	if err != nil {
		log.Println("Could not fetch DNS Record! Is the API key valid? Does it have the right permissions?")
		return dnsRecord, false
	}

	//Compare IPs
	if dnsRecord.Data == ip.Query {
		return dnsRecord, false
	}

	return dnsRecord, true
}

//UpdateDNSRecords updates an A record with the current IP
func UpdateDNSRecord(oldRecord cloudflare.DNSRecord) error {
	//Change the old record IP and send it
	oldRecord.Content = ip.Query
	err := api.UpdateDNSRecord(ctx, zoneId, recordId, oldRecord)
	if err != nil {
		return err
	}
	return nil
}

//RecordUpdateWorker is the main loop which matches the IP and changes the DNS record when nesscary
func RecordUpdateWorker() {
	for {
		//Fetch the current IP
		fetchIP()
		//Compare to the current record
		oldRecord, changed := CompareIPToRecord()
		if changed {
			UpdateDNSRecord(oldRecord)
		}
		time.Sleep(time.Minute)
	}
}

//Setup runs set ups functions to help set enviroment variables
func Setup() {

}

func main() {
	//Set logging location
	f, err := os.OpenFile("cloudflare-ddns.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()

	log.SetOutput(f)

	//Choose the startup options
	Setup()

	//Declare the API
	GetAPIToken()
	api, err = cloudflare.NewWithAPIToken(apiToken)
	if err != nil {
		log.Fatal("Could not start Cloudflare API")
	}

	//Set the context
	ctx = context.Background()

	//Start the worker
	RecordUpdateWorker()
}
