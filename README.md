# Exporter-ChinaCache

A [Prometheus](https://prometheus.io/download/) Exporter/Collector for [ChinaCache CDN](https://en.chinacache.com/) API serving a [Prometheus Pushgateway](https://github.com/prometheus/pushgateway)

### What is this repository for?
* This repository packages three important modules with different capabilities:
    + [**client**.go](./client.go): is the communicating interface between an application and the ChinaCache-API
        - on demand, it executes HTTP-Requests to fetch metrics from ChinaCache
        - the resulting JSON-response is then read into a go-struct conforming the JSON structure
    + [**collector**.go](./collector.go): asks the client to get those metrics and accesses the resulting struct's fields
        - from these fields, it creates metrics conforming to the conventions of Prometheus
        - it also implements the Collector interface defined by Prometheus
        - the functions `Describe` and `Collect` are called passively by another
          application/module or on scrape by a prometheus server
    + [**main**.go](./main.go): initializes client and collector and pushes the resulting metrics to a prometheus pushgateway

### Package Management
* This project uses **dep** as package manager
* versions are tracked in `Gopkg.lock`
* dep settings are included in `Gopkg.toml`
* get dep here: [https://github.com/golang/dep](https://github.com/golang/dep)

### Test locally via Prometheus (on Linux)
1. define the following environment-variables to grant access to your ChinaCache-Account and to the PushGateway:
    + `CHINACACHE_USER`
    + `CHINACACHE_PASS`
    + `CHINACACHE_CHANNEL_IDS`
    + `CHINACACHE_INTERVAL` (optional query interval in [Golang duration format][1]. By default, fetch metrics once and exit)
    + `PUSHGATEWAY` (default for the docker image: http://localhost:9091)
    + `QUERYTIME` (optional)
2. have Prometheus [downloaded](https://prometheus.io/download/) and installed locally
3. start Prometheus using the config-file from this repo:
    + `sudo prometheus -config.file=prometheus.yml`
4. download the docker image of the pushgateway:
    + `docker pull prom/pushgateway`
5. start it:
    + `sudo docker run -d -p 9091:9091 prom/pushgateway`
6. build and run the repo:
    + `go build bin/cc`
    + `./bin/cc`

### Exposed Metrics
#### Metrics fetched from "normal" API:
**API-Endpoint:** `https://portal-api.chinacache.com:444/api/public/statistics/%s.do?userName=%s&apiPasswd=%s&channelIds=%s&startTime=%s&endTime=%s`


- `ChinaCache_metrics_isp_total_bytes`
    + HELP: Total Traffic in Bytes
    + TYPE: GaugeValue
    + Labels: channel = [one of the available channel IDs]

- `ChinaCache_metrics_isp_specific_flux_ratio`
    + HELP: Flow Rate for a single ISP.
    + TYPE: GaugeValue
    + Labels:
        * channel = [one of the available channel IDs]
        * isp = [name of the specific ISP translated from chinese]
    + *NOTE: These metrics specific to single regions are only available if the queried time-frame crosses two days, as they are only created once per day.*

- `ChinaCache_metrics_hit_miss_total`
    + HELP: Hits and Misses total values.
    + TYPE: GaugeValue
    + Labels:
        * channel = [one of the available channel IDs]
        * HitOrMiss = [Hit|Miss]

- `ChinaCache_metrics_region_total_bytes`
    + HELP: Total Traffic in Bytes.
    + TYPE: GaugeValue
    + Labels: channel = [one of the available channel IDs]

- `ChinaCache_metrics_region_specific_flux_ratio`
    + HELP: Flow Rate for a single region.
    + TYPE: GaugeValue
    + Labels:
        * channel = [one of the available channel IDs]
        * region = [Province]
        * name = [name of the specific region translated from chinese]
    + *NOTE: These metrics specific to single regions are only available if the queried time-frame crosses two days, as they are only created once per day.*

#### Metrics fetched from REST-API:

**REST-API-Endpoint:** `https://portal-api.chinacache.com:444/rest-api/public/statistics/%s?api_user=%s&api_key=%s&start_time=%s&end_time=%s&channel_id=%s`

- `ChinaCache_metrics_statuscodes_request_percent`
    + HELP: Percentage of Requests that result in the given StatusCode.
    + TYPE: GaugeValue
    + Labels:
        * channel = [one of the available channel IDs]
        * StatusCode = [200|404|...]

- `ChinaCache_metrics_statuscodes_request_count`
    + HELP: Number of Requests that result in the given StatusCode.
    + TYPE: GaugeValue
    + Labels:
        * channel = [one of the available channel IDs]
        * StatusCode = [200|404|...]

[1]: https://golang.org/pkg/time/#ParseDuration