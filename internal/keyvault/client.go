package keyvault

import (
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/keyvault/armkeyvault"
	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azsecrets"
	"github.com/zemelkajakub/lskv/internal/profile"
)

const (
	keyvaultUriFormat = "https://%s.vault.azure.net/"
)

type Client struct {
	credential     azcore.TokenCredential
	VaultsClient   *armkeyvault.VaultsClient
	SubscriptionID string
	Vaults         []string
}

func NewClient(p *profile.Profile) (*Client, error) {

	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		return nil, err
	}

	vaultsClient, err := armkeyvault.NewVaultsClient(p.SubscriptionID, cred, nil)
	if err != nil {
		return nil, err
	}

	return &Client{
		credential:     cred,
		VaultsClient:   vaultsClient,
		SubscriptionID: p.SubscriptionID,
		Vaults:         p.Vaults,
	}, nil
}

// NewSecretsClient creates a data plane client for accessing secrets in a vault
func (c *Client) NewSecretsClient(vaultName string) (*azsecrets.Client, error) {
	vaultURL := fmt.Sprintf(keyvaultUriFormat, vaultName)

	secretsClient, err := azsecrets.NewClient(vaultURL, c.credential, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create secrets client: %w", err)
	}

	return secretsClient, nil
}
