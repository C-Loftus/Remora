package main

import (
	"context"

	oc "github.com/c-loftus/orca-controller"
	"github.com/ollama/ollama/api"
)

type Connection struct {
	ConnectedToOrca   bool
	ConnectionMessage string
	OrcaClient        *oc.OrcaClient
}

func (a *Connection) Reset() {
	a.ConnectedToOrca = false
	a.ConnectionMessage = ""
	if a.OrcaClient != nil {
		a.OrcaClient.Close()
	}
	a.OrcaClient = nil
}

// App struct
type App struct {
	ctx               context.Context
	visionModelPrompt string
	orcaConnection    Connection
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{orcaConnection: Connection{}}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *App) TryCreateClient() (success bool) {
	if a.orcaConnection.OrcaClient != nil {
		version, err := a.orcaConnection.OrcaClient.GetVersion()
		if err == nil {
			a.orcaConnection.ConnectedToOrca = true
			a.orcaConnection.ConnectionMessage = version
			return true
		}
	}

	a.orcaConnection.Reset()

	client, err := oc.NewOrcaClient()
	if err == nil {
		a.orcaConnection.OrcaClient = client
		_, err := a.orcaConnection.OrcaClient.GetVersion()
		if err == nil {
			return true
		} else {
			a.orcaConnection.ConnectionMessage = err.Error()
			return false
		}
	} else {
		a.orcaConnection.ConnectionMessage = err.Error()
		return false
	}
}

func (a *App) OrcaVersion(name string) (string, error) {
	if a.orcaConnection.OrcaClient == nil {
		return a.orcaConnection.OrcaClient.GetVersion()
	} else {
		return "Not connected to Orca", nil
	}
}

func (a *App) ConnectionStatus() Connection {
	return a.orcaConnection
}

func (a *App) GetHotKeys() []string {
	var hotkeys []string
	for _, hotkey := range hotkeyList {
		hotkeys = append(hotkeys, hotkey.ToString())
	}
	return hotkeys
}

func (a *App) GetDisplayServerType() DisplayServerType {
	return DetectDisplayServer()
}

func (a *App) OllamaConnectionStatus() string {
	client, err := api.ClientFromEnvironment()
	if err != nil {
		return "Error connecting to Ollama: " + err.Error()
	}
	_, err = client.Version(context.Background())
	if err != nil {
		return "Error connecting to Ollama: " + err.Error()
	}
	return "Connected"
}

func (a *App) SetPrompt(prompt string) {
	a.visionModelPrompt = prompt
}

func (a *App) GetPrompt() string {
	return a.visionModelPrompt
}

func (a *App) LastOllamaResponse() string {
	return mostRecentOllamaResponse
}
