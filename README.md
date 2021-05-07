# goteams

A quick utility for common administrative functions in MS teams.

***Note: This utility depends on github.com/yaegashi/msgraph.go for connecting to the MS graph API***

## Introduction

To start using the module, create a new teams client
```go
import "github.com/micaiahwallace/goteams"

const (
  tenantId = "XXXX"
  clientId = "XXXX"
  clientSecret = "XXXX"
)

func main() {

  // Create new client
  client := goteams.New(tenantID, clientID, clientSecret)
}
```

then you can use any of the commands currently added.

### Available Commands
```go
// Install an app on a team
func (client *TeamsClient) InstallTeamsApp(teamId, appId string) error

// Install an app on a team if it doesn't exist
func (client *TeamsClient) InstallNewTeamsApp(teamId, appId string) error

// Install app on list of teams (teams must be an msgraph.Team struct with an ID set)
// any installation errors will be sent through the returned error chan
func (client *TeamsClient) InstallNewAppOnTeams(teams []msgraph.Team, appId string) chan error

// Check if an app is installed on a team
func (client *TeamsClient) IsAppInstalled(teamId, appId string) (bool, error)

// Get a list of installed apps for a team
func (client *TeamsClient) GetTeamsApps(teamId string) ([]msgraph.TeamsAppInstallation, error)

// Get a list of all team enabled groups
func (client *TeamsClient) GetTeams() ([]msgraph.Group, error)
```

## Examples

### Install teams app on all teams
This example installs a given app by ID to every team in an organization that doesn't have it already installed.
```bash
$ install-app-on-teams -tenant $tenantId -client $clientId -secret $clientSecret -app $appId
```
Command usage:
```bash
usage: install-app-on-teams
  -app string
        Teams app ID
  -client string
        MS app client ID
  -secret string
        MS app client secret
  -tenant string
        MS tenant ID
```