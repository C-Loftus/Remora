package main

import (
	"context"
	"fmt"

	oc "github.com/c-loftus/orca-controller"
)

// App struct
type App struct {
	ctx        context.Context
	orcaClient *oc.OrcaClient
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *App) TryCreateClient() string {
	client, err := oc.NewOrcaClient()
	if err == nil {
		a.orcaClient = client
		version, err := a.orcaClient.GetVersion()
		if err == nil {
			return fmt.Sprintf("Connected to Orca %s", version)
		} else {
			a.orcaClient = nil
			return err.Error()
		}
	} else {
		a.orcaClient = nil
		return err.Error()
	}
}

// Greet returns a greeting for the given name
func (a *App) OrcaVersion(name string) (string, error) {
	return a.orcaClient.GetVersion()
}

// ClientConnected returns true if the client is connected to Orca
func (a *App) ConnectedToOrca() bool {
	return a.orcaClient != nil
}

func (a *App) TakeScreenshot() (string, error) {
	return takeScreenshot()
}
