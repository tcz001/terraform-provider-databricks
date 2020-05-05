package databricks

import (
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/tcz001/databricks-sdk-go/models"
)

func resourceDatabricksUser() *schema.Resource {
	return &schema.Resource{
		Create: resourceDatabricksUserCreate,
		Update: resourceDatabricksUserUpdate,
		Read:   resourceDatabricksUserRead,
		Delete: resourceDatabricksUserDelete,

		Schema: map[string]*schema.Schema{
			"display_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"username": {
				Type:     schema.TypeString,
				Required: true,
			},
			"active": {
				Type:     schema.TypeBool,
				Default:  true,
				Optional: true,
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
			"email": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"value": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"primary": {
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

func resourceDatabricksScimUserExpandGroup(groups []interface{}) []models.Groups {
	groupsList := make([]models.Groups, len(groups))
	for _, v := range groups {
		groupsElem := v.(map[string]interface{})
		groupsList = append(groupsList,
			models.Groups{
				Value:   groupsElem["value"].(string),
				Ref:     groupsElem["ref"].(string),
				Display: groupsElem["display"].(string),
			})
	}

	return groupsList
}

func resourceDatabricksScimUserFlattenGroups(groups []models.Groups) []map[string]interface{} {
	result := make([]map[string]interface{}, 0)
	for _, v := range groups {
		result = append(result, map[string]interface{}{
			"value":   v.Value,
			"ref":     v.Ref,
			"display": v.Display,
		})
	}
	return result
}

func resourceDatabricksScimUserExpandEntitlement(entitlements []interface{}) []models.Entitlements {
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

func resourceDatabricksScimUserFlattenEntitlements(entitlements []models.Entitlements) []map[string]interface{} {
	result := make([]map[string]interface{}, 0)
	for _, v := range entitlements {
		result = append(result, map[string]interface{}{
			"value": v.Value,
		})
	}
	return result
}

func resourceDatabricksScimUserExpandEmails(emails []interface{}) []models.Emails {
	emailList := make([]models.Emails, len(emails))
	for _, v := range emails {
		emailElem := v.(map[string]interface{})
		emailList = append(emailList,
			models.Emails{
				Type_:   emailElem["type"].(string),
				Value:   emailElem["value"].(string),
				Primary: emailElem["primary"].(bool),
			})
	}

	return emailList
}

func resourceDatabricksScimUserFlattenEmails(emails []models.Emails) []map[string]interface{} {
	result := make([]map[string]interface{}, 0)
	for _, v := range emails {
		result = append(result, map[string]interface{}{
			"value":   v.Value,
			"type":    v.Type_,
			"primary": v.Primary,
		})
	}
	return result
}

func resourceDatabricksUserCreate(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*Client).scim

	log.Print("[DEBUG] Creating SCIM User")

	request := models.ScimUser{
		DisplayName: d.Get("display_name").(string),
		UserName:    d.Get("username").(string),
	}

	if v, ok := d.GetOk("groups"); ok {
		request.Groups = resourceDatabricksScimUserExpandGroup(v.(*schema.Set).List())
	}

	if v, ok := d.GetOk("entitlements"); ok {
		request.Entitlements = resourceDatabricksScimGroupExpandEntitlement(v.(*schema.Set).List())
	}

	if v, ok := d.GetOk("emails"); ok {
		request.Emails = resourceDatabricksScimUserExpandEmails(v.(*schema.Set).List())
	}

	log.Print("[DEBUG] Creating SCIM User")
	resp, err := apiClient.CreateUser(request)
	if err != nil {
		log.Printf("[DEBUG] err: %s", err.Error())
		return err
	}

	log.Print("[DEBUG] Created SCIM User with response")

	d.SetId(resp.Id)

	return nil
}

func resourceDatabricksUserRead(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*Client).scim

	resp, err := apiClient.GetUser(d.Id())
	if err != nil {
		return err
	}
	if resourceDatabricksUserNotExistError(err) {
		log.Printf("[WARN] User (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}
	d.SetId(resp.Id)
	d.Set("entitlements", resourceDatabricksScimUserFlattenEntitlements(resp.Entitlements))
	d.Set("groups", resourceDatabricksScimUserFlattenGroups(resp.Groups))
	d.Set("members", resourceDatabricksScimUserFlattenEmails(resp.Emails))
	d.Set("display_name", resp.DisplayName)
	d.Set("name", resp.Name)
	d.Set("active", resp.Active)
	d.Set("username", resp.UserName)
	return nil
}

func resourceDatabricksUserUpdate(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*Client).scim

	log.Print("[DEBUG] Updating Scim user")

	request := models.ScimUser{
		Id:          d.Id(),
		DisplayName: d.Get("display_name").(string),
		Active:      d.Get("active").(bool),
	}

	if v, ok := d.GetOk("groups"); ok {
		request.Groups = resourceDatabricksScimUserExpandGroup(v.(*schema.Set).List())
	}

	log.Print("Display name set")

	if v, ok := d.GetOk("entitlements"); ok {
		request.Entitlements = resourceDatabricksScimUserExpandEntitlement(v.(*schema.Set).List())
	}

	log.Print("Groups Set")

	if v, ok := d.GetOk("emails"); ok {
		request.Emails = resourceDatabricksScimUserExpandEmails(v.(*schema.Set).List())
	}
	log.Print("Members Set")

	resp, err := apiClient.UpdateUser(request.Id, request)
	if err != nil {
		log.Printf("[DEBUG] err: %s", err.Error())
		return err
	}

	d.SetId(resp.Id)
	d.Set("entitlements", resourceDatabricksScimUserFlattenEntitlements(resp.Entitlements))
	d.Set("groups", resourceDatabricksScimUserFlattenGroups(resp.Groups))
	d.Set("members", resourceDatabricksScimUserFlattenEmails(resp.Emails))
	d.Set("display_name", resp.DisplayName)
	d.Set("active", resp.Active)

	return nil
}

func resourceDatabricksUserDelete(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*Client).scim

	log.Printf("[DEBUG] Deleting User: %s", d.Id())

	err := apiClient.DeleteUser(d.Id())
	if err != nil {
		return err
	}

	d.SetId("")

	return nil
}

func resourceDatabricksUserNotExistError(err error) bool {
	return false
}
