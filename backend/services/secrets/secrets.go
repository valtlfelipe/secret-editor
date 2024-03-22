package secrets

import (
	"log"

	"github.com/valtlfelipe/secret-editor/backend/types"
)

type SecretsProvider interface {
	GetSecrets() types.ResultList
	CreateSecret(name *string, secret *string) types.Result
	SaveSecret(arn *string, secret *string) types.Result
	LoadSecret(arn *string) types.Result
}

type SecretsHandler struct {
	Preferences types.PreferencesProvider
	Provider    SecretsProvider
}

func NewSecretsHandler(preferences types.PreferencesProvider) *SecretsHandler {
	switch preferences.Current {
	case "AWS":
		return &SecretsHandler{
			Preferences: preferences,
			Provider:    NewAWSSecretsProvider(preferences),
		}
	default:
		log.Printf("unsupported provider %s\n", preferences.Current)
	}

	return nil
}

func (h *SecretsHandler) GetSecrets() types.ResultList {
	return h.Provider.GetSecrets()
}
