package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/micaiahwallace/goteams"
	msgraph "github.com/yaegashi/msgraph.go/v1.0"
)

func main() {

	// Parse variables from command line
	var tenantID, clientID, clientSecret, appId string
	flag.StringVar(&tenantID, "tenant", "", "MS tenant ID")
	flag.StringVar(&clientID, "client", "", "MS app client ID")
	flag.StringVar(&clientSecret, "secret", "", "MS app client secret")
	flag.StringVar(&appId, "app", "", "Teams app ID")
	flag.Parse()

	// Check for arg existence
	if tenantID == "" || clientID == "" || clientSecret == "" || appId == "" {
		fmt.Fprintf(os.Stderr, "usage: %s\n", os.Args[0])
		flag.PrintDefaults()
		os.Exit(1)
	}

	// Create new client
	client := goteams.New(tenantID, clientID, clientSecret)

	// Get all the teams
	// teams, getTeamsErr := client.GetTeams()
	// if getTeamsErr != nil {
	// 	log.Fatal(getTeamsErr)
	// }

	teamId := "XXX"

	// sample team
	teams := make([]msgraph.Team, 1)
	teams[0].ID = &teamId

	// Install app on all teams
	installErrors := client.InstallNewAppOnTeams(teams, appId)
	for err := range installErrors {
		log.Printf("install error: %v", err)
	}
	fmt.Println("App installation complete.")
}
