package goteams

import (
	"context"
	"fmt"
	"log"

	"github.com/yaegashi/msgraph.go/msauth"
	msgraph "github.com/yaegashi/msgraph.go/v1.0"
	"golang.org/x/oauth2"
)

type TeamsClient struct {

	// MS Graph API client
	Client *msgraph.GraphServiceRequestBuilder

	// Context for client requests
	Ctx context.Context
}

// Create a new TeamsClient
func New(tenantID, clientID, clientSecret string) *TeamsClient {

	ctx := context.Background()

	// Create authorization manager
	m := msauth.NewManager()
	scopes := []string{msauth.DefaultMSGraphScope}
	ts, err := m.ClientCredentialsGrant(ctx, tenantID, clientID, clientSecret, scopes)
	if err != nil {
		log.Fatal(err)
	}

	// Create request clients
	httpClient := oauth2.NewClient(ctx, ts)
	graphClient := msgraph.NewClient(httpClient)

	// Create new TeamsClient
	return &TeamsClient{
		Client: graphClient,
		Ctx:    ctx,
	}
}

// Install an app on a team
func (client *TeamsClient) InstallTeamsApp(teamId, appId string) error {

	// Create installation request body
	appInstall := map[string]string{
		"teamsApp@odata.bind": fmt.Sprintf("https://graph.microsoft.com/v1.0/appCatalogs/teamsApps/%s", appId),
	}

	// Attempt app installation
	req := client.Client.Teams().ID(teamId).Request()
	err := req.JSONRequest(client.Ctx, "POST", "/installedApps", appInstall, nil)

	if err != nil {
		return err
	}

	return nil
}

// Install an app on a team if it doesn't exist
func (client *TeamsClient) InstallNewTeamsApp(teamId, appId string) error {

	// Retrieve app installations
	appInstalled, checkErr := client.IsAppInstalled(teamId, appId)
	if checkErr != nil {
		return checkErr
	}

	// Install it if not already installed
	if !appInstalled {
		installErr := client.InstallTeamsApp(teamId, appId)
		if installErr != nil {
			return installErr
		}
	}

	return nil
}

// Install app on list of teams
func (client *TeamsClient) InstallNewAppOnTeams(teams []msgraph.Team, appId string) chan error {

	// Create new comm channel
	errors := make(chan error)

	// Process all teams in a go routine
	go func(errors chan error) {

		// Loop through teams
		for _, team := range teams {
			err := client.InstallNewTeamsApp(*team.ID, appId)
			if err != nil {
				errors <- err
			}
		}

		// Close channel to indicate done
		close(errors)

	}(errors)

	return errors
}

// Check if an app is installed on a team
func (client *TeamsClient) IsAppInstalled(teamId, appId string) (bool, error) {

	// Get app installations for this team
	teamApps, getAppErr := client.GetTeamsApps(teamId)
	if getAppErr != nil {
		return false, getAppErr
	}

	// Check if app is installed
	appInstalled := false
	for _, ai := range teamApps {
		if *ai.TeamsAppDefinition.TeamsAppID == appId {
			appInstalled = true
		}
	}

	return appInstalled, nil
}

// Get a list of installed apps for a team
func (client *TeamsClient) GetTeamsApps(teamId string) ([]msgraph.TeamsAppInstallation, error) {
	req := client.Client.Teams().ID(teamId).InstalledApps().Request()
	req.Expand("teamsAppDefinition")
	return req.Get(client.Ctx)
}

// Get a list of all team enabled groups
func (client *TeamsClient) GetTeams() ([]msgraph.Group, error) {

	// Final list of filtered groups
	teams := make([]msgraph.Group, 0)

	// Request list of groups
	req := client.Client.Groups().Request()
	req.Select("id,resourceProvisioningOptions,displayName")
	groups, err := req.Get(client.Ctx)
	if err != nil {
		return nil, err
	}

	// Loop over groups to find team enabled ones
	for _, group := range groups {

		// Retrive the rpo for the team
		rpoRaw, ok := group.GetAdditionalData("resourceProvisioningOptions")
		if !ok || rpoRaw == nil {
			log.Println("Unable to get teams resource provisioning options")
			continue
		}

		// Convert resourceProvisioningOptions to a []interface{}
		rpo, ok := rpoRaw.([]interface{})
		if !ok {
			log.Println("Unable to parse resource provisioning options for team")
			continue
		}

		// Check for a "Team" entry in the rpo slice
		isTeam := false
		for _, option := range rpo {
			opt := option.(string)
			if opt == "Team" {
				isTeam = true
			}
		}

		// Append group to return array if it has a team
		if isTeam {
			teams = append(teams, group)
		}
	}
	return teams, nil
}
