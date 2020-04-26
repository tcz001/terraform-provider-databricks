package databricks

import (
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/tcz001/databricks-sdk-go/models"
)

func resourceDatabricksServicePrincipal() *schema.Resource {
	return &schema.Resource{
		Create: resourceDatabricksServicePrincipalCreate,
		Update: resourceDatabricksServicePrincipalUpdate,
		Read:   resourceDatabricksServicePrincipalRead,
		Delete: resourceDatabricksServicePrincipalDelete,

		Schema: map[string]*schema.Schema{
			"application_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"display_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"active": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"groups": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"value": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"display": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ref": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"entitlements": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"value": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
		},
	}
}

func resourceDatabricksServicePrincipalUpdate(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*Client).scim

	log.Print("[DEBUG] Creating servicePrincipal")

	//TODO support Entitlements and Groups

	request := models.ServicePrincipal{
		Id:            d.Id(),
		ApplicationId: d.Get("application_id").(string),
		DisplayName:   d.Get("display_name").(string),
	}

	if v, ok := d.GetOk("groups"); ok {
		request.Groups = resourceDatabricksServicePrincipalExpandGroup(v.(*schema.Set).List())
	}

	if v, ok := d.GetOk("entitlements"); ok {
		request.Entitlements = resourceDatabricksServicePrincipalExpandEntitlement(v.(*schema.Set).List())
	}

	resp, err := apiClient.UpdateServicePrincipal(&request)
	if err != nil {
		log.Printf("[DEBUG] err: %s", err.Error())
		return err
	}

	//TODO support Entitlements and Groups
	d.SetId(resp.Id)
	d.Set("application_id", resp.ApplicationId)
	d.Set("active", resp.Active)
	d.Set("display_name", resp.DisplayName)

	log.Printf("[DEBUG] ServicePrincipal Id: %s, ApplicationId: %s", d.Id(), d.Get("application_id").(string))

	return nil
}

func resourceDatabricksServicePrincipalExpandGroup(groups []interface{}) []models.Groups {
	groupsList := make([]models.Groups, len(groups))
	for _, v := range groups {
		groupsElem := v.(map[string]interface{})
		groupsList = append(groupsList,
			models.Groups{
				Value: groupsElem["value"].(string),
			})
	}

	return groupsList
}

func resourceDatabricksServicePrincipalExpandEntitlement(entitlements []interface{}) []models.Entitlements {
	entitlementsList := make([]models.Entitlements, len(entitlements))
	for _, v := range entitlements {
		entitlementsElem := v.(map[string]interface{})
		entitlementsList = append(entitlementsList,
			models.Entitlements{
				Value: entitlementsElem["value"].(string),
			})
	}

	return entitlementsList
}

func resourceDatabricksServicePrincipalCreate(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*Client).scim

	log.Print("[DEBUG] Creating servicePrincipal")

	//TODO support Entitlements and Groups
	request := models.ServicePrincipalCreateRequest{
		ApplicationId: d.Get("application_id").(string),
		DisplayName:   d.Get("display_name").(string),
	}

	if v, ok := d.GetOk("groups"); ok {
		request.Groups = resourceDatabricksServicePrincipalExpandGroup(v.(*schema.Set).List())
	}

	if v, ok := d.GetOk("entitlements"); ok {
		request.Entitlements = resourceDatabricksServicePrincipalExpandEntitlement(v.(*schema.Set).List())
	}

	resp, err := apiClient.CreateServicePrincipal(&request)
	if err != nil {
		log.Printf("[DEBUG] err: %s", err.Error())
		return err
	}

	//TODO support Entitlements and Groups
	d.SetId(resp.Id)
	d.Set("application_id", resp.ApplicationId)
	d.Set("active", resp.Active)
	d.Set("display_name", resp.DisplayName)

	log.Printf("[DEBUG] ServicePrincipal Id: %s, ApplicationId: %s", d.Id(), d.Get("application_id").(string))

	return nil
}

func resourceDatabricksServicePrincipalRead(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*Client).scim

	resp, err := apiClient.GetServicePrincipal(d.Id())
	if err != nil {
		return err
	}
	//TODO Handle Not Exist Error
	if resourceDatabricksServicePrincipalNotExistError(err) {
		log.Printf("[WARN] ServicePrincipal (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}
	//TODO support Entitlements and Groups
	d.SetId(resp.Id)
	d.Set("application_id", resp.ApplicationId)
	d.Set("active", resp.Active)
	d.Set("display_name", resp.DisplayName)

	return nil
}

func resourceDatabricksServicePrincipalDelete(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*Client).scim

	log.Printf("[DEBUG] Deleting secret: %s", d.Id())

	err := apiClient.DeleteServicePrincipal(d.Id())
	if err != nil {
		return err
	}

	d.SetId("")

	return nil
}

func resourceDatabricksServicePrincipalNotExistError(err error) bool {
	return false
}
