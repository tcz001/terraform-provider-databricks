package token

import (
	"encoding/json"

	"github.com/tcz001/databricks-sdk-go/client"
	"github.com/tcz001/databricks-sdk-go/models"
)

type Endpoint struct {
	Client *client.Client
}

func (c *Endpoint) Create(request *models.TokenCreateRequest) (*models.TokenCreateReponse, error) {
	bytes, err := c.Client.Query("POST", "token/create", request)
	if err != nil {
		return nil, err
	}

	resp := models.TokenCreateReponse{}
	err = json.Unmarshal(bytes, &resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

func (c *Endpoint) List() (*models.TokenListResponse, error) {
	bytes, err := c.Client.Query("GET", "token/list", nil)
	if err != nil {
		return nil, err
	}

	resp := models.TokenListResponse{}
	err = json.Unmarshal(bytes, &resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

func (c *Endpoint) Revoke(request *models.TokenRevokeRequest) error {
	_, err := c.Client.Query("POST", "token/delete", request)
	if err != nil {
		return err
	}

	return nil
}
