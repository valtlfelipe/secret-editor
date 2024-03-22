package secrets

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/valtlfelipe/secret-editor/backend/types"
)

type AWSSecretsProvider struct {
	Preferences types.PreferencesProvider
	client      *secretsmanager.Client
}

func NewAWSSecretsProvider(preferences types.PreferencesProvider) SecretsProvider {
	return &AWSSecretsProvider{
		Preferences: preferences,
	}
}

func (p *AWSSecretsProvider) GetSecretsManagerClient() (client *secretsmanager.Client, err error) {
	if p.client != nil {
		return p.client, nil
	}

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithSharedConfigProfile(p.Preferences.AWSProfile),
	)
	if err != nil {
		return nil, err
	}

	p.client = secretsmanager.NewFromConfig(cfg)

	return p.client, nil
}

func (p *AWSSecretsProvider) GetSecrets() types.ResultList {

	client, err := p.GetSecretsManagerClient()
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

func (p *AWSSecretsProvider) CreateSecret(name *string, secret *string) types.Result {
	client, err := p.GetSecretsManagerClient()
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

func (p *AWSSecretsProvider) LoadSecret(arn *string) types.Result {
	client, err := p.GetSecretsManagerClient()
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

func (p *AWSSecretsProvider) SaveSecret(arn *string, secret *string) types.Result {
	client, err := p.GetSecretsManagerClient()
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
