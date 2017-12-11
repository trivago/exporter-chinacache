package main

import (
	"log"
	"os"

	"github.com/prometheus/client_golang/prometheus/push"
)

var (
	// User-defined Environment-Variables
	pushGateway = os.Getenv("PUSHGATEWAY")
	user        = os.Getenv("CHINACACHE_USER")
	pass        = os.Getenv("CHINACACHE_PASS")
	channelIDs  = os.Getenv("CHINACACHE_CHANNEL_IDS") // given a comma-separated list of available channel-ids for the specified chinacache user
	querytime   = os.Getenv("QUERYTIME")
)

func main() {

	// check if all the necessary environment-variables were provided and are not empty
	if user == "" || pass == "" || pushGateway == "" || channelIDs == "" || querytime == "" {
		log.Fatal("error: Please provide the environment-variables PUSHGATEWAY, CHINACACHE_USER, CHINACACHE_PASS, CHINACACHE_CHANNEL_IDS (comma-separated list) and QUERYTIME (minutes)")
		return
	}

	// create new client and hand it over to create a new collector
	client := NewChinaCacheClient(user, pass, channelIDs, querytime)
	collector := NewChinaCacheCollector(client)

	// push the metrics exposed by the collector
	err := push.AddCollectors(
		"ChinaCachePush",
		nil,
		pushGateway,
		collector,
	)
	if err != nil {
		log.Fatal("error: couldn't add collector to push metrics!")
	}
}
