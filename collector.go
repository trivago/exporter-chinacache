package main

import (
	"log"

	"sync"

	"strconv"

	"github.com/prometheus/client_golang/prometheus"
)

const (
	// NAMESPACE for all created/fetched metrics that are exposed to Prometheus
	NAMESPACE = "ChinaCache"
)

// ChinaCacheCollector only needs the ChinaCacheClient which executes the HTTP-Requests to the ChinaCache-API
type ChinaCacheCollector struct {
	client ChinaCacheClient
}

// NewChinaCacheCollector returns a new collector given an ChinaCacheClient
func NewChinaCacheCollector(client *ChinaCacheClient) *ChinaCacheCollector {
	return &ChinaCacheCollector{
		client: *client,
	}
}

var (
	// TODO: reduce amount of different Descriptions by adding more labels? e.g. region-label = "all" or "<specific region>" instead of one for "all" and one for "specifics"
	// descriptions for metrics fetched from "normal" API
	ispTotalDesc = prometheus.NewDesc(
		prometheus.BuildFQName(NAMESPACE, "metrics", "isp_total_bytes"), "Total Traffic in Bytes.", []string{"channel"}, nil,
	)
	ispSpecificDesc = prometheus.NewDesc(
		prometheus.BuildFQName(NAMESPACE, "metrics", "isp_specific_flux_ratio"), "Flow Rate for a single ISP.", []string{"channel", "isp"}, nil,
	)
	hitRateDesc = prometheus.NewDesc(
		prometheus.BuildFQName(NAMESPACE, "metrics", "hit_miss_total"), "Hits and Misses total values.", []string{"channel", "HitOrMiss"}, nil,
	)
	regionTotalDesc = prometheus.NewDesc(
		prometheus.BuildFQName(NAMESPACE, "metrics", "region_total_bytes"), "Total Traffic in Bytes.", []string{"channel"}, nil,
	)
	regionSpecificDesc = prometheus.NewDesc(
		prometheus.BuildFQName(NAMESPACE, "metrics", "region_specific_flux_ratio"), "Flow Rate for a single region.", []string{"channel", "region", "name"}, nil,
	)

	// descriptions for metrics fetched from REST-API
	statusCodesRequestPercentDesc = prometheus.NewDesc(
		prometheus.BuildFQName(NAMESPACE, "metrics", "statuscodes_request_percent"), "Percentage of Requests that result in the given StatusCode.", []string{"channel", "StatusCode"}, nil,
	)
	statusCodesRequestCountDesc = prometheus.NewDesc(
		prometheus.BuildFQName(NAMESPACE, "metrics", "statuscodes_request_count"), "Number of Requests that result in the given StatusCode.", []string{"channel", "StatusCode"}, nil,
	)
)

// Describe implements Collector Interface Function defined in Prometheus
func (c *ChinaCacheCollector) Describe(ch chan<- *prometheus.Desc) {
	// descriptions for metrics fetched from "normal" API
	ch <- ispTotalDesc
	ch <- ispSpecificDesc
	ch <- hitRateDesc
	ch <- regionTotalDesc
	ch <- regionSpecificDesc

	// descriptions for metrics fetched from REST-API
	ch <- statusCodesRequestPercentDesc
}

// Collect implements Collector Interface Function defined in Prometheus
// is either called on scrape by prometheus server or in-process from pushing service
func (c *ChinaCacheCollector) Collect(ch chan<- prometheus.Metric) {
	var collectGroup sync.WaitGroup
	for _, channel := range c.client.Channels {
		collectGroup.Add(1)
		go c.metrics(ch, &collectGroup, channel)
	}
	collectGroup.Wait()
}

// concurrently run collect-functions for all currently implemented possible metrics
func (c *ChinaCacheCollector) metrics(ch chan<- prometheus.Metric, collectGroup *sync.WaitGroup, channel string) {
	defer collectGroup.Done()
	var metricsGroup sync.WaitGroup
	metricsGroup.Add(4) // 4 = number of goroutines started below
	go c.collectHitRate(ch, &metricsGroup, channel)
	go c.collectRegion(ch, &metricsGroup, channel)
	go c.collectIsp(ch, &metricsGroup, channel)
	go c.collectStatusCodes(ch, &metricsGroup, channel)
	metricsGroup.Wait()
}

// request HitRate-Data using the client and extract fields into new metrics
func (c *ChinaCacheCollector) collectHitRate(ch chan<- prometheus.Metric, metricsGroup *sync.WaitGroup, channel string) {
	defer metricsGroup.Done()
	hitRateData, err := c.client.GetHitRate(channel)
	if err != nil {
		log.Println("error: couldn't collect hit rate \n", err)
		return
	}
	ch <- prometheus.MustNewConstMetric(hitRateDesc, prometheus.GaugeValue, float64(hitRateData.Hit), []string{channel, "Hit"}...)
	ch <- prometheus.MustNewConstMetric(hitRateDesc, prometheus.GaugeValue, float64(hitRateData.Miss), []string{channel, "Miss"}...)
}

// request region-data using the client and extract fields into new metrics
func (c *ChinaCacheCollector) collectRegion(ch chan<- prometheus.Metric, metricsGroup *sync.WaitGroup, channel string) {
	defer metricsGroup.Done()
	regionData, err := c.client.GetRegion(channel)
	if err != nil {
		log.Println("error: couldn't collect region data \n", err)
		return
	}
	ch <- prometheus.MustNewConstMetric(regionTotalDesc, prometheus.GaugeValue, float64(regionData.TotalFlux), []string{channel}...)

	if len(regionData.Provinces) != 0 { // only happens when requested timeslot spans over two days (e.g. 12:58 - 00:03), as this data is only generated once every day
		for _, region := range regionData.Provinces {
			name := translations[region.ProvinceName]
			if name == "" { // no translation present in translations-map
				name = region.ProvinceName
			}
			ch <- prometheus.MustNewConstMetric(regionSpecificDesc, prometheus.GaugeValue, region.FluxRatio, []string{channel, "Province", name}...)
		}
	}
}

// request ISP-Data using the client and extract fields into new metrics
func (c *ChinaCacheCollector) collectIsp(ch chan<- prometheus.Metric, metricsGroup *sync.WaitGroup, channel string) {
	defer metricsGroup.Done()
	ispData, err := c.client.GetIsp(channel)
	if err != nil {
		log.Println("error: couldn't collect isp data \n", err)
		return
	}
	ch <- prometheus.MustNewConstMetric(ispTotalDesc, prometheus.GaugeValue, float64(ispData.TotalFlux), []string{channel}...)
	if len(ispData.Isps) != 0 { // only happens when requested timeslot spans over two days (e.g. 12:58 - 00:03), as this data is only generated once every day
		for _, isp := range ispData.Isps {
			name := translations[isp.Isp]
			if name == "" { // translation of isp-name not available in map
				name = isp.Isp
			}
			ch <- prometheus.MustNewConstMetric(ispSpecificDesc, prometheus.GaugeValue, isp.FluxRatio, []string{channel, name}...)
		}
	}
}

// request StatusCodes-Data using the client and extract fields into new metrics
func (c *ChinaCacheCollector) collectStatusCodes(ch chan<- prometheus.Metric, metricsGroup *sync.WaitGroup, channel string) {
	defer metricsGroup.Done()
	statusCodesData, err := c.client.GetStatusCodes(channel)
	if err != nil {
		log.Println("error: couldn't collect status codes \n", err)
		return
	}
	if statusCodesData.Success {
		var requestPercentageString string
		var requestPercentageFloat float64
		for _, data := range statusCodesData.Data {
			requestPercentageString = data.RequestPercent[:len(data.RequestPercent)-1] // cut off percentage sign
			requestPercentageFloat, _ = strconv.ParseFloat(requestPercentageString, 64)
			ch <- prometheus.MustNewConstMetric(statusCodesRequestPercentDesc, prometheus.GaugeValue, requestPercentageFloat, []string{channel, data.HTTPCode}...)
			ch <- prometheus.MustNewConstMetric(statusCodesRequestCountDesc, prometheus.GaugeValue, float64(data.RequestCount), []string{channel, data.HTTPCode}...)
		}
	} else { // TODO: move error-handling to client
		log.Printf("error: getStatusCodes didn't succeed and returned: %+v", statusCodesData)
	}

}
