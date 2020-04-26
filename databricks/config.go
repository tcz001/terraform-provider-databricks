package databricks

import (
	"log"
	"time"

	"github.com/Azure/go-autorest/autorest/adal"
	azAuth "github.com/Azure/go-autorest/autorest/azure/auth"
	"github.com/tcz001/databricks-sdk-go/api/clusters"
	scim "github.com/tcz001/databricks-sdk-go/api/scim"
	secrets "github.com/tcz001/databricks-sdk-go/api/secrets"
	token "github.com/tcz001/databricks-sdk-go/api/token"
	"github.com/tcz001/databricks-sdk-go/api/workspace"
	apiClient "github.com/tcz001/databricks-sdk-go/client"
)

const (
	maxRetries = 3
	retryDelay = 5 * time.Second
)

type Config struct {
	Domain                            *string
	Token                             *string
	WorkspaceId                       *string
	XDatabricksAzureSPManagementToken *string

	AzCCC *azAuth.ClientCredentialsConfig
	DBCCC *azAuth.ClientCredentialsConfig
}

type Client struct {
	clusters  *clusters.Endpoint
	workspace *workspace.Endpoint
	secrets   *secrets.Endpoint
	token     *token.Endpoint
	scim      *scim.Endpoint
}

func ServicePrincipalToken(ccc *azAuth.ClientCredentialsConfig) (*adal.ServicePrincipalToken, error) {
	oauthConfig, err := adal.NewOAuthConfigWithAPIVersion(ccc.AADEndpoint, ccc.TenantID, nil)
	if err != nil {
		return nil, err
	}
	return adal.NewServicePrincipalToken(*oauthConfig, ccc.ClientID, ccc.ClientSecret, ccc.Resource)
}

func (c *Config) Client() (interface{}, error) {
	var client Client

	if c.DBCCC != nil && c.AzCCC != nil {
		// Get AZ Databricks Token
		log.Print("[DEBUG] Getting AZ Databricks Token")
		rsToken, err := ServicePrincipalToken(c.DBCCC)
		if err != nil {
			return nil, err
		}
		err = rsToken.EnsureFresh()
		if err != nil {
			return nil, err
		}
		rsOAuthToken := rsToken.OAuthToken()
		c.Token = &rsOAuthToken

		// Get Az Management SP Token
		log.Print("[DEBUG] Getting Az Management SP Token")
		azToken, err := ServicePrincipalToken(c.AzCCC)
		if err != nil {
			return nil, err
		}
		err = azToken.EnsureFresh()
		if err != nil {
			return nil, err
		}
		oauthToken := azToken.OAuthToken()
		c.XDatabricksAzureSPManagementToken = &oauthToken
	}

	if c.WorkspaceId != nil {
		log.Printf("[DEBUG] c.WorkspaceId: %s", *c.WorkspaceId)
	}
	if c.Token != nil {
		log.Printf("[DEBUG] c.Token: %s", *c.Token)
	}
	if c.XDatabricksAzureSPManagementToken != nil {
		log.Printf("[DEBUG] c.xDatabricksAzureSPManagementToken: %s", *c.XDatabricksAzureSPManagementToken)
	}

	opts := apiClient.Options{
		Domain:                              c.Domain,
		Token:                               c.Token,
		XDatabricksAzureSPManagementToken:   c.XDatabricksAzureSPManagementToken,
		XDatabricksAzureWorkspaceResourceId: c.WorkspaceId,
		MaxRetries:                          maxRetries,
		RetryDelay:                          retryDelay,
	}
	cl, err := apiClient.NewClient(opts)
	if err != nil {
		return nil, err
	}

	client.clusters = &clusters.Endpoint{Client: cl}
	client.workspace = &workspace.Endpoint{Client: cl}
	client.secrets = &secrets.Endpoint{Client: cl}
	client.token = &token.Endpoint{Client: cl}
	client.scim = &scim.Endpoint{Client: cl}

	return &client, nil
}
