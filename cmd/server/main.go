package main

import (
	"fmt"

	"github.com/dnsoftware/gophkeeper/internal/server/app"
)

var (
	buildVersion = "N/A"
	buildDate    = "N/A"
	buildCommit  = "N/A"
)

func main() {
	fmt.Printf("Build version: %s\n", buildVersion)
	fmt.Printf("Build date: %s\n", buildDate)
	fmt.Printf("Build commit: %s\n", buildCommit)

	err := app.ServerRun()
	if err != nil {
		panic(err)
	}
}
