package main

// structs for JSON responses
// generated via https://mholt.github.io/json-to-go/

// GetIspStruct used for JSON response returned by API on requesting e.g. https://portal-api.chinacache.com:444/api/public/statistics/getIsp.do?userName=%s&apiPasswd=%s&channelIds=%s&startTime=%s&endTime=%s
type GetIspStruct struct {
	TotalFlux int64 `json:"totalFlux"`
	Isps      []struct {
		Isp       string  `json:"isp"`
		FluxRatio float64 `json:"fluxRatio"`
		HitCount  int     `json:"hitCount"`
		HitRatio  float64 `json:"hitRatio"`
	} `json:"isps"`
	Code int `json:"code"`
}

// GetHitRateStruct used for JSON response returned by API on requesting e.g. https://portal-api.chinacache.com:444/api/public/statistics/getHitRate.do?userName=%s&apiPasswd=%s&channelIds=%s&startTime=%s&endTime=%s
type GetHitRateStruct struct {
	HitPercent  float64 `json:"HitPercent"`
	Hit         int     `json:"Hit"`
	Code        int     `json:"code"`
	MissPercent float64 `json:"MissPercent"`
	Miss        int     `json:"Miss"`
}

// GetRegionStruct used for JSON response returned by API on requesting e.g. https://portal-api.chinacache.com:444/api/public/statistics/getRegion.do?userName=%s&apiPasswd=%s&channelIds=%s&startTime=%s&endTime=%s
type GetRegionStruct struct {
	TotalFlux int64 `json:"totalFlux"`
	Provinces []struct {
		ProvinceName string  `json:"provinceName"`
		FluxRatio    float64 `json:"fluxRatio"`
		HitCount     int     `json:"hitCount"`
		HitRatio     float64 `json:"hitRatio"`
	} `json:"provinces"`
	States []struct {
		StateOrRegionName string  `json:"stateOrRegionName"`
		FluxRatio         float64 `json:"fluxRatio"`
		HitCount          int     `json:"hitCount"`
		HitRatio          float64 `json:"hitRatio"`
	} `json:"states"`
	Code int `json:"code"`
}

// GetStatusCodesStruct used for JSON response returned by API on requesting e.g. https://portal-api.chinacache.com:444/rest-api/public/statistics/http_code?api_user=%s&api_key=%s&start_time=%s&end_time=%s&channel_id=%s
type GetStatusCodesStruct struct {
	Data []struct {
		FluxPercent    string `json:"flux_percent"`
		HTTPCode       string `json:"http_code"`
		RequestCount   int    `json:"request_count"`
		RequestPercent string `json:"request_percent"`
	} `json:"data"`
	Msg     string `json:"msg"`
	Status  int    `json:"status"`
	Success bool   `json:"success"`
}
