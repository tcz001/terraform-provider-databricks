package secret

import (
	"encoding/json"

	"github.com/tcz001/databricks-sdk-go/client"
	"github.com/tcz001/databricks-sdk-go/models"
)

type Endpoint struct {
	Client *client.Client
}

func (c *Endpoint) Put(request *models.SecretsPutRequest) error {
	_, err := c.Client.Query("POST", "secrets/put", request)
	if err != nil {
		return err
	}

	return nil
}

func (c *Endpoint) List(request *models.SecretsListRequest) (*models.SecretsListResponse, error) {
	bytes, err := c.Client.Query("GET", "secrets/list", request)
	if err != nil {
		return nil, err
	}

	resp := models.SecretsListResponse{}
	err = json.Unmarshal(bytes, &resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

func (c *Endpoint) Delete(request *models.SecretsDeleteRequest) error {
	_, err := c.Client.Query("POST", "secrets/delete", request)
	if err != nil {
		return err
	}

	return nil
}

func (c *Endpoint) AddScope(request *models.SecretsScopesCreateRequest) error {
	_, err := c.Client.Query("POST", "secrets/scopes/create", request)
	if err != nil {
		return err
	}

	return nil
}

func (c *Endpoint) ListScopes() (*models.SecretsScopesListResponse, error) {
	bytes, err := c.Client.Query("GET", "secrets/scopes/list", nil)
	if err != nil {
		return nil, err
	}

	resp := models.SecretsScopesListResponse{}
	err = json.Unmarshal(bytes, &resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

func (c *Endpoint) DeleteScope(request *models.SecretsScopesDeleteRequest) error {
	_, err := c.Client.Query("POST", "secrets/scopes/delete", request)
	if err != nil {
		return err
	}

	return nil
}
