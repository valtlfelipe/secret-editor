package main

import (
	"context"
	"log"

	"github.com/valtlfelipe/secret-editor/backend/services"
	"github.com/valtlfelipe/secret-editor/backend/types"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

// App struct
type App struct {
	ctx       context.Context
	prefStore services.PreferencesStorage
}

// NewApp creates a new App application struct
func NewApp(pref services.PreferencesStorage) *App {
	return &App{
		prefStore: pref,
	}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

func GetSecretsManagerClient() (client *secretsmanager.Client, err error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithSharedConfigProfile("valtlfelipe"),
	)
	if err != nil {
		return
	}

	client = secretsmanager.NewFromConfig(cfg)

	return
}

func (a *App) GetPreferences() types.Preferences {
	return a.prefStore.Preferences
}

func (a *App) SetPreference(path string, value string) error {
	return a.prefStore.SetPreference(path, value)
}

func (a *App) GetSecrets() types.ResultList {

	client, err := GetSecretsManagerClient()
	if err != nil {
		log.Printf("error: %v", err)
		return types.ResultList{
			Success: false,
			Error:   err.Error(),
		}
	}

	// TODO: handle pagination
	output, err := client.ListSecrets(context.TODO(), &secretsmanager.ListSecretsInput{
		MaxResults: aws.Int32(20),
	})
	if err != nil {
		log.Printf("error: %v", err)
		return types.ResultList{
			Success: false,
			Error:   err.Error(),
		}
	}

	result := []types.Secret{}

	for _, object := range output.SecretList {
		result = append(result, types.Secret{
			ARN:  *object.ARN,
			Name: *object.Name,
		})
	}

	return types.ResultList{
		Success: true,
		Result:  result,
	}
}

func (a *App) CreateSecret(name *string, secret *string) types.Result {
	client, err := GetSecretsManagerClient()
	if err != nil {
		log.Printf("error: %v", err)
		return types.Result{
			Success: false,
			Error:   err.Error(),
		}
	}

	output, err := client.CreateSecret(context.TODO(), &secretsmanager.CreateSecretInput{
		Name:         name,
		SecretString: secret,
	})
	if err != nil {
		log.Printf("error: %v", err)
		return types.Result{
			Success: false,
			Error:   err.Error(),
		}
	}

	return types.Result{
		Success: true,
		Result: types.Secret{
			ARN:  *output.ARN,
			Name: *output.Name,
		},
	}
}

func (a *App) LoadSecret(arn *string) types.Result {
	client, err := GetSecretsManagerClient()
	if err != nil {
		log.Printf("error: %v", err)
		return types.Result{
			Success: false,
			Error:   err.Error(),
		}
	}

	output, err := client.GetSecretValue(context.TODO(), &secretsmanager.GetSecretValueInput{
		SecretId: arn,
	})
	if err != nil {
		log.Printf("error: %v", err)
		return types.Result{
			Success: false,
			Error:   err.Error(),
		}
	}

	return types.Result{
		Success: true,
		Result: types.Secret{
			ARN:    *output.ARN,
			Name:   *output.Name,
			Secret: *output.SecretString,
		},
	}
}

func (a *App) SaveSecret(arn *string, secret *string) types.Result {
	client, err := GetSecretsManagerClient()
	if err != nil {
		log.Printf("error: %v", err)
		return types.Result{
			Success: false,
			Error:   err.Error(),
		}
	}

	output, err := client.UpdateSecret(context.TODO(), &secretsmanager.UpdateSecretInput{
		SecretId:     arn,
		SecretString: secret,
	})
	if err != nil {
		log.Printf("error: %v", err)
		return types.Result{
			Success: false,
			Error:   err.Error(),
		}
	}

	return types.Result{
		Success: true,
		Result: types.Secret{
			ARN:  *output.ARN,
			Name: *output.Name,
		},
	}
}
