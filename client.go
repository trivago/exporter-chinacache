package main

/*
 * The Client is the part of the program that actually communicates with the ChinaCache API.
 * It uses user-defined environment-variables for logging in and dynamically builds the Query-URLs
 * to gather all needed metrics from ChinaCache CDN.
 * The return values of the exported functions are structures that are individually designed for the JSON-responses from ChinaCache.
 * Those structures are defined in types.go.
 */

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// ChinaCacheClient is used to request data from ChinaCache-API
type ChinaCacheClient struct {
	User      string
	Pass      string
	Channels  []string
	Querytime int64
}

const (
	// APIENDPOINT e.g. used for getHitRate, getRegion, getIsp
	APIENDPOINT = "https://portal-api.chinacache.com:444/api/public/statistics/%s.do?userName=%s&apiPasswd=%s&channelIds=%s&startTime=%s&endTime=%s"

	// RESTAPIENDPOINT e.g. used for StatusCodes
	RESTAPIENDPOINT = "https://portal-api.chinacache.com:444/rest-api/public/statistics/%s?api_user=%s&api_key=%s&start_time=%s&end_time=%s&channel_id=%s"

	// QUERYTIME : query data for a timeslot of QUERYTIME minutes
	QUERYTIME = 5

	// Methods on API
	getHitRate    = "getHitRate"
	getIspData    = "getIsp"
	getRegionData = "getRegion"

	// Methods on REST API
	getStatusCodes = "http_code"
)

var (
	apiMethods = map[string]bool{getHitRate: true, getIspData: true, getRegionData: true} // used for URL-building
	// restApiMethods = []string{getStatusCodes} // not used yet, as there's currently only one method

	logTemplate = "Method = %s, Err = %s, Returned Object = %+v"
)

// NewChinaCacheClient returns a new client that sends HTTP-Requests to the ChinaCache-API
// arguments given are the username, password and available channel-IDs for the API
// querytime specifies the length of the time interval that is to be requested
func NewChinaCacheClient(user string, pass string, channelIDs string, querytime string) *ChinaCacheClient {
	// remove whitespaces, which might be in there by natural typing and then split list by comma
	channels := strings.Split(strings.Replace(channelIDs, " ", "", -1), ",")
	if len(channels) == 0 {
		log.Fatal("error: empty list of channels")
	}
	var qt int64
	var err error
	qt = QUERYTIME
	if len(querytime) != 0 {
		if qt, err = strconv.ParseInt(querytime, 10, 64); err != nil {
			qt = QUERYTIME
			log.Println("warning: couldn't parse given QUERYTIME -> set to default\n", err)
		}
	}
	return &ChinaCacheClient{
		User:      user,
		Pass:      pass,
		Channels:  channels,
		Querytime: qt,
		// TODO: add user-defined timeouts and retries
	}
}

// GetHitRate requests Hit-Rate-Data from the API and returns the resulting JSON-response as bytes
func (c *ChinaCacheClient) GetHitRate(channelID string) (str *GetHitRateStruct, err error) {

	defer log.Println(fmt.Sprintf(logTemplate, "GetHitRate", err, str))

	body, err := c.request(getHitRate, channelID)
	if err != nil {
		log.Println("error: GetHitRate failed on channel ID", channelID)
		return nil, err
	}
	var data GetHitRateStruct
	err = json.Unmarshal(body, &data)
	return &data, err
}

// GetIsp requests ISP-Data from the API and returns the resulting JSON-response as bytes
func (c *ChinaCacheClient) GetIsp(channelID string) (*GetIspStruct, error) {
	body, err := c.request(getIspData, channelID)
	if err != nil {
		log.Println("error: GetIsp failed on channel ID", channelID)
		return nil, err
	}
	var data GetIspStruct
	err = json.Unmarshal(body, &data)
	return &data, err
}

// GetRegion requests Region-Data from the API and returns the resulting JSON-response as bytes
func (c *ChinaCacheClient) GetRegion(channelID string) (*GetRegionStruct, error) {
	body, err := c.request(getRegionData, channelID)
	if err != nil {
		log.Println("error: GetRegion failed on channel ID", channelID)
		return nil, err
	}
	var data GetRegionStruct
	err = json.Unmarshal(body, &data)
	return &data, err
}

// GetStatusCodes requests Status-Code-Data from the REST-API and returns the resulting JSON-response as bytes
func (c *ChinaCacheClient) GetStatusCodes(channelID string) (*GetStatusCodesStruct, error) {
	body, err := c.request(getStatusCodes, channelID)
	if err != nil {
		log.Println("error: GetStatusCodes failed on channel ID", channelID)
		return nil, err
	}
	var data GetStatusCodesStruct
	err = json.Unmarshal(body, &data)
	return &data, err
}

// request sends a GET-request to the specified endpoint and returns the response-data in a byte array
func (c *ChinaCacheClient) request(method, channel string) ([]byte, error) {
	url := c.buildURL(method, channel)

	// GET request to API-Endpoint
	resp, err := http.Get(url)
	if err != nil {
		log.Println("error: couldn't get", url)
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }() // TODO: should the error-value be checked here? ¯\_(ツ)_/¯

	// read and (hopefully) return response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("error: couldn't read response body queried from ", url)
		return nil, err
	}
	return body, nil
}

func (c *ChinaCacheClient) buildURL(method, channel string) string {
	var url string
	// get last interesting time slot and build URL
	end := time.Now()
	start := end.Add(-time.Duration(c.Querytime) * 60 * time.Second) // requested time slot
	if apiMethods[method] {
		url = fmt.Sprintf(APIENDPOINT, method, c.User, c.Pass, channel, start.Format("200601021504"), end.Format("200601021504"))
	} else {
		url = fmt.Sprintf(RESTAPIENDPOINT, method, c.User, c.Pass, start.Format("20060102"), end.Format("20060102"), channel)
	}
	return fmt.Sprint(url, "&timeZone=GMT%2B8")
}
