package main

import (
	"log"
	"os"
	"time"

	"github.com/prometheus/client_golang/prometheus/push"
)

var (
	// User-defined Environment-Variables
	pushGateway   = os.Getenv("PUSHGATEWAY")
	user          = os.Getenv("CHINACACHE_USER")
	pass          = os.Getenv("CHINACACHE_PASS")
	channelIDs    = os.Getenv("CHINACACHE_CHANNEL_IDS") // given a comma-separated list of available channel-ids for the specified chinacache user
	queryInterval = os.Getenv("CHINACACHE_INTERVAL")
	querytime     = os.Getenv("QUERYTIME")
)

func main() {
	// check if all the necessary environment-variables were provided and are not empty
	if user == "" || pass == "" || pushGateway == "" || channelIDs == "" {
		log.Fatal("error: Please provide the environment-variables PUSHGATEWAY, CHINACACHE_USER, CHINACACHE_PASS, CHINACACHE_CHANNEL_IDS (comma-separated list) and QUERYTIME (minutes)")
		return
	}

	interval, err := time.ParseDuration(queryInterval)
	if err != nil {
		// Fetch values once and exit
		interval = 0
	}

	// create new client and hand it over to create a new collector
	client := NewChinaCacheClient(user, pass, channelIDs, querytime)
	collector := NewChinaCacheCollector(client)

	for {
		// push the metrics exposed by the collector
		err := push.AddCollectors(
			"ChinaCachePush",
			nil,
			pushGateway,
			collector,
		)
		if err != nil {
			log.Printf("Could not push metrics: %s\n", err)
		}

		if interval == 0 {
			break
		} else {
			time.Sleep(interval)
		}
	}

	// Exit with non-zero statuscode if last fetch resulted in error
	if err != nil {
		os.Exit(1)
	}
}
