package databricks

import (
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/tcz001/databricks-sdk-go/models"
)

func resourceDatabricksToken() *schema.Resource {
	return &schema.Resource{
		Create: resourceDatabricksTokenCreate,
		Read:   resourceDatabricksTokenRead,
		Delete: resourceDatabricksTokenRevoke,

		Schema: map[string]*schema.Schema{
			"lifetime_seconds": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"comment": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"token_value": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"creation_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"expiry_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func resourceDatabricksTokenCreate(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*Client).token

	log.Print("[DEBUG] Creating token")

	request := models.TokenCreateRequest{
		LifetimeSeconds: int32(d.Get("lifetime_seconds").(int)),
		Comment:         d.Get("comment").(string),
	}

	resp, err := apiClient.Create(&request)
	if err != nil {
		log.Printf("[DEBUG] err: %s", err.Error())
		return err
	}

	d.SetId(resp.TokenInfo.TokenId)
	d.Set("token_value", resp.TokenValue)
	d.Set("creation_time", resp.TokenInfo.CreationTime)
	d.Set("expiry_time", resp.TokenInfo.ExpiryTime)
	d.Set("comment", resp.TokenInfo.Comment)

	log.Printf("[DEBUG] Token Id: %s", d.Id())

	return nil
}

func resourceDatabricksTokenRead(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*Client).token

	resp, err := apiClient.List()
	if err != nil {
		return err
	}
	if !resourceDatabricksTokenNotExistsError(d.Id(), resp) {
		log.Printf("[WARN] Token (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}

	return nil
}

func resourceDatabricksTokenRevoke(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*Client).token

	log.Printf("[DEBUG] Deleting secret: %s", d.Id())

	request := models.TokenRevokeRequest{
		TokenId: d.Id(),
	}

	err := apiClient.Revoke(&request)
	if err != nil {
		return err
	}

	d.SetId("")
	d.SetId("")

	return nil
}

func resourceDatabricksTokenNotExistsError(tokenId string, resp *models.TokenListResponse) bool {
	tokens := resp.TokenInfos
	for _, token := range tokens {
		if token.TokenId == tokenId {
			return true
		}
	}
	return false
}
