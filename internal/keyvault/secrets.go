package keyvault

import (
	"context"
	"errors"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
)

type Secret struct {
	Name    string `json:"name"`
	Enabled bool   `json:"enabled"`
}

func (c *Client) ListSecrets(ctx context.Context, vaultName string) ([]Secret, string, error) {

	secretsClient, err := c.NewSecretsClient(vaultName)
	if err != nil {
		return nil, "", err
	}

	var secrets []Secret
	pager := secretsClient.NewListSecretPropertiesPager(nil)

	for pager.More() {
		page, err := pager.NextPage(ctx)
		if err != nil {
			var respError *azcore.ResponseError
			if errors.As(err, &respError) {
				switch respError.StatusCode {
				case 403:
					return nil, "access_denied", nil
				case 404:
					return nil, "not_found", nil
				case 429:
					return nil, "too_many_requests", nil
				default:
					return nil, "error", fmt.Errorf("HTTP %d: %s", respError.StatusCode, respError.ErrorCode)
				}
			}
			return nil, "error", err
		}
		for _, secretProp := range page.Value {
			enabled := secretProp.Attributes.Enabled != nil && *secretProp.Attributes.Enabled

			secrets = append(secrets, Secret{
				Name:    secretProp.ID.Name(),
				Enabled: enabled,
			})
		}
	}
	return secrets, "success", nil
}
