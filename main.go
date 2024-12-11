package main

import (
	"flag"
	"fmt"
	"log/slog"

	"github.com/BKadirkhodjaev/request-cli/orders"
)

const CommandName string = "main"

var (
	environment string
	threadCount int
	enableDebug bool
)

func main() {
	parseArgs()

	gatewayHostname := getGatewayHostname(environment)
	orders.ParseCsvAndOpenOrdersInBulk(gatewayHostname, enableDebug, threadCount)
}

func parseArgs() {
	flag.StringVar(&environment, "env", "okapi", "Environment to run on (okapi or eureka)")
	flag.IntVar(&threadCount, "threads", 50, "Persistent HTTP thread count")
	flag.BoolVar(&enableDebug, "debug", false, "Enable debug output of HTTP request and response")
	flag.Parse()

	slog.Info(CommandName, "Using arguments", fmt.Sprintf("gatewayHostname: %s, enableDebug: %t, threadCount: %d", environment, enableDebug, threadCount))
}

func getGatewayHostname(environment string) string {
	if environment == "okapi" {
		return "http://localhost:9130"
	}
	return "http://localhost:8000"
}
