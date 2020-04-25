package databricks

import (
	"log"

	"github.com/hashicorp/terraform/helper/schema"

	azAuth "github.com/Azure/go-autorest/autorest/azure/auth"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"domain": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"token": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"client_id": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ARM_CLIENT_ID", ""),
				Description: "The Client ID which should be used.",
			},

			"tenant_id": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ARM_TENANT_ID", ""),
				Description: "The Tenant ID which should be used.",
			},
			// Client Secret specific fields
			"client_secret": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ARM_CLIENT_SECRET", ""),
				Description: "The Client Secret which should be used. For use When authenticating as a Service Principal using a Client Secret.",
			},

			"workspace_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The databricks workspace id",
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"databricks_cluster":      resourceDatabricksCluster(),
			"databricks_notebook":     resourceDatabricksNotebook(),
			"databricks_secret_scope": resourceDatabricksSecretScope(),
			"databricks_secret":       resourceDatabricksSecret(),
			"databricks_token":        resourceDatabricksToken(),
		},
		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	config := Config{}

	if domain, ok := d.GetOk("domain"); ok {
		s := domain.(string)
		config.Domain = &s
	}

	if token, ok := d.GetOk("token"); ok {
		// Token authentication
		log.Print("[DEBUG] Token authentication")
		s := token.(string)
		config.Token = &s
		config.DBCCC = nil
		config.AzCCC = nil
	} else if workspaceId, ok := d.GetOk("workspace_id"); ok {
		s := workspaceId.(string)
		config.WorkspaceId = &s
		// Azure Client Credentials authentication
		log.Print("[DEBUG] Azure Client Credentials authentication")
		var clientID, clientSecret, tenantID string

		if s, ok := d.GetOk("client_id"); ok {
			clientID = s.(string)
		}

		if s, ok := d.GetOk("client_secret"); ok {
			clientSecret = s.(string)
		}

		if s, ok := d.GetOk("tenant_id"); ok {
			tenantID = s.(string)
		}

		if clientID != "" && clientSecret != "" && tenantID != "" {
			dbccc := azAuth.NewClientCredentialsConfig(clientID, clientSecret, tenantID)
			dbccc.Resource = "2ff814a6-3304-4ab8-85cb-cd0e6f879c1d"
			azccc := azAuth.NewClientCredentialsConfig(clientID, clientSecret, tenantID)
			azccc.Resource = "https://management.core.windows.net/"
			config.DBCCC = &dbccc
			config.AzCCC = &azccc
		}
	}

	return config.Client()
}
