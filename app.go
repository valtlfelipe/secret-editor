package main

import (
	"context"

	"github.com/valtlfelipe/secret-editor/backend/services"
	"github.com/valtlfelipe/secret-editor/backend/services/secrets"
	"github.com/valtlfelipe/secret-editor/backend/types"
)

// App struct
type App struct {
	ctx            context.Context
	PrefStore      services.PreferencesStorage
	SecretsHandler secrets.SecretsHandler
}

// NewApp creates a new App application struct
func NewApp(pref services.PreferencesStorage) *App {
	return &App{
		PrefStore:      pref,
		SecretsHandler: *secrets.NewSecretsHandler(pref.Preferences.Provider),
	}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *App) GetPreferences() types.Preferences {
	return a.PrefStore.Preferences
}

func (a *App) SetPreference(path string, value string) error {
	return a.PrefStore.SetPreference(path, value)
}

func (a *App) GetSecrets() types.ResultList {
	return a.SecretsHandler.GetSecrets()
}

func (a *App) CreateSecret(name *string, secret *string) types.Result {
	return a.SecretsHandler.Provider.CreateSecret(name, secret)
}

func (a *App) LoadSecret(arn *string) types.Result {
	return a.SecretsHandler.Provider.LoadSecret(arn)
}

func (a *App) SaveSecret(arn *string, secret *string) types.Result {
	return a.SecretsHandler.Provider.SaveSecret(arn, secret)
}
