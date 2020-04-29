package databricks

import (
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/tcz001/databricks-sdk-go/models"
)

func resourceDatabricksGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceDatabricksGroupCreate,
		Update: resourceDatabricksGroupUpdate,
		Read:   resourceDatabricksGroupRead,
		Delete: resourceDatabricksGroupDelete,

		Schema: map[string]*schema.Schema{
			"display_name": {
				Type:     schema.TypeString,
				Required: true,
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
			"members": {
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

func resourceDatabricksScimGroupExpandGroup(groups []interface{}) []models.Groups {
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

func resourceDatabricksScipGroupFlattenGroups(groups []models.Groups) []map[string]interface{} {
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

func resourceDatabricksScimGroupExpandEntitlement(entitlements []interface{}) []models.Entitlements {
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

func resourceDatabricksScipGroupFlattenEntitlements(entitlements []models.Entitlements) []map[string]interface{} {
	result := make([]map[string]interface{}, 0)
	for _, v := range entitlements {
		result = append(result, map[string]interface{}{
			"value": v.Value,
		})
	}
	return result
}

func resourceDatabricksScimGroupExpandMembers(members []interface{}) []models.ScimMember {
	memberList := make([]models.ScimMember, len(members))
	for _, v := range members {
		memberElem := v.(map[string]interface{})
		memberList = append(memberList,
			models.ScimMember{
				Value: memberElem["value"].(string),
			})
	}

	return memberList
}

func resourceDatabricksScipGroupFlattenMembers(members []models.ScimMember) []map[string]interface{} {
	result := make([]map[string]interface{}, 0)
	for _, v := range members {
		result = append(result, map[string]interface{}{
			"value": v.Value,
		})
	}
	return result
}

func resourceDatabricksGroupCreate(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*Client).scim

	log.Print("[DEBUG] Creating SCIM Group")

	request := models.ScimGroup{
		DisplayName: d.Get("display_name").(string),
	}

	if v, ok := d.GetOk("groups"); ok {
		request.Groups = resourceDatabricksScimGroupExpandGroup(v.(*schema.Set).List())
	}

	if v, ok := d.GetOk("entitlements"); ok {
		request.Entitlements = resourceDatabricksScimGroupExpandEntitlement(v.(*schema.Set).List())
	}

	if v, ok := d.GetOk("members"); ok {
		request.Members = resourceDatabricksScimGroupExpandMembers(v.(*schema.Set).List())
	}

	log.Print("[DEBUG] Creating SCIM Group")
	resp, err := apiClient.CreateGroup(&request)
	if err != nil {
		log.Printf("[DEBUG] err: %s", err.Error())
		return err
	}

	log.Print("[DEBUG] Created SCIM Group with response")
	//log.Printf("[DEBUG] Group Id: %s, Members: %s", d.Id(), d.Get("members"))

	d.SetId(resp.Id)
	d.Set("display_name", resp.DisplayName)

	return nil
}

func resourceDatabricksGroupRead(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*Client).scim

	resp, err := apiClient.GetGroup(d.Id())
	if err != nil {
		return err
	}
	if resourceDatabricksGroupNotExistError(err) {
		log.Printf("[WARN] Group (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}
	d.SetId(resp.Id)
	d.Set("entitlements", resourceDatabricksScipGroupFlattenEntitlements(resp.Entitlements))
	d.Set("groups", resourceDatabricksScipGroupFlattenGroups(resp.Groups))
	d.Set("members", resourceDatabricksScipGroupFlattenMembers(resp.Members))
	d.Set("display_name", resp.DisplayName)
	return nil
}

func resourceDatabricksGroupUpdate(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*Client).scim

	log.Print("[DEBUG] Updating Scim group")

	request := models.ScimGroup{
		Id:          d.Id(),
		DisplayName: d.Get("display_name").(string),
	}

	if v, ok := d.GetOk("groups"); ok {
		request.Groups = resourceDatabricksScimGroupExpandGroup(v.(*schema.Set).List())
	}

	log.Print("Display name set")

	if v, ok := d.GetOk("entitlements"); ok {
		request.Entitlements = resourceDatabricksScimGroupExpandEntitlement(v.(*schema.Set).List())
	}

	log.Print("Groups Set")
	if v, ok := d.GetOk("members"); ok {
		request.Members = resourceDatabricksScimGroupExpandMembers(v.(*schema.Set).List())
	}
	log.Print("Members Set")

	log.Print("Request :")
	log.Print(request)

	resp, err := apiClient.UpdateGroup(request.Id, request)
	if err != nil {
		log.Printf("[DEBUG] err: %s", err.Error())
		return err
	}

	d.SetId(resp.Id)
	d.Set("entitlements", resourceDatabricksScipGroupFlattenEntitlements(resp.Entitlements))
	d.Set("groups", resourceDatabricksScipGroupFlattenGroups(resp.Groups))
	d.Set("members", resourceDatabricksScipGroupFlattenMembers(resp.Members))
	d.Set("display_name", resp.DisplayName)

	//log.Printf("[DEBUG] Group Id: %s, Name: %s", d.Id(), d.Get("display_name").(string))

	return nil
}

func resourceDatabricksGroupDelete(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*Client).scim

	log.Printf("[DEBUG] Deleting group: %s", d.Id())

	err := apiClient.DeleteGroup(d.Id())
	if err != nil {
		return err
	}

	d.SetId("")

	return nil
}

//TODO
func resourceDatabricksGroupNotExistError(err error) bool {
	return false
}
