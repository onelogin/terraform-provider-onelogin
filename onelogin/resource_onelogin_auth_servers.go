package onelogin

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/onelogin/onelogin-go-sdk/v4/pkg/onelogin"
	authserverschema "github.com/onelogin/terraform-provider-onelogin/ol_schema/auth_server"
)

// AuthServers returns a resource with the CRUD methods and Terraform Schema defined
func AuthServers() *schema.Resource {
	return &schema.Resource{
		CreateContext: authServersCreate,
		ReadContext:   authServersRead,
		UpdateContext: authServersUpdate,
		DeleteContext: authServersDelete,
		Importer:      &schema.ResourceImporter{},
		Schema:        authserverschema.Schema(),
	}
}

func authServersCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*onelogin.OneloginSDK)

	authServer, _ := authserverschema.Inflate(map[string]interface{}{
		"name":          d.Get("name"),
		"description":   d.Get("description"),
		"configuration": d.Get("configuration"),
	})

	result, err := client.CreateAuthServer(&authServer)
	if err != nil {
		return diag.Errorf("error creating auth server: %v", err)
	}

	// Extract the auth server ID from the response
	authServerMap, ok := result.(map[string]interface{})
	if !ok || authServerMap["id"] == nil {
		return diag.Errorf("failed to parse auth server creation response or auth server ID not found in response")
	}

	authServerID := int(authServerMap["id"].(float64))
	d.SetId(fmt.Sprintf("%d", authServerID))
	log.Printf("[CREATED] Created auth server with id %d", authServerID)

	return authServersRead(ctx, d, m)
}

func authServersUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*onelogin.OneloginSDK)

	authServer, _ := authserverschema.Inflate(map[string]interface{}{
		"id":            d.Id(),
		"name":          d.Get("name"),
		"description":   d.Get("description"),
		"configuration": d.Get("configuration"),
	})

	authID, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.Errorf("error converting id to integer: %v", err)
	}

	_, err = client.UpdateAuthServer(authID, &authServer)
	if err != nil {
		return diag.Errorf("error updating auth server: %v", err)
	}

	log.Printf("[UPDATED] Updated auth server with id %d", authID)
	return authServersRead(ctx, d, m)
}

func authServersRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*onelogin.OneloginSDK)

	authID, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.Errorf("error converting id to integer: %v", err)
	}

	result, err := client.GetAuthServerByID(authID, nil)
	if err != nil {
		log.Printf("[ERROR] There was a problem reading the auth server: %v", err)
		return diag.FromErr(err)
	}

	// Check if the resource was not found
	if result == nil {
		d.SetId("")
		return nil
	}

	// Parse the response
	authServerMap, ok := result.(map[string]interface{})
	if !ok {
		return diag.Errorf("failed to parse auth server response")
	}

	log.Printf("[READ] Reading auth server with id %d", authID)

	if authServerMap["name"] != nil {
		d.Set("name", authServerMap["name"])
	}

	if authServerMap["description"] != nil {
		d.Set("description", authServerMap["description"])
	}

	// Handle configuration
	if authServerMap["configuration"] != nil {
		configMap, ok := authServerMap["configuration"].(map[string]interface{})
		if ok {
			// Convert the configuration back to a nested structure
			d.Set("configuration", []interface{}{
				map[string]interface{}{
					"resource_identifier": configMap["resource_identifier"],
					"audiences":           configMap["audiences"],
					"access_token_expiration_minutes": func() interface{} {
						if v, ok := configMap["access_token_expiration_minutes"]; ok {
							return int(v.(float64))
						}
						return nil
					}(),
					"refresh_token_expiration_minutes": func() interface{} {
						if v, ok := configMap["refresh_token_expiration_minutes"]; ok {
							return int(v.(float64))
						}
						return nil
					}(),
				},
			})
		}
	}

	return nil
}

func authServersDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*onelogin.OneloginSDK)

	authID, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.Errorf("error converting id to integer: %v", err)
	}

	_, err = client.DeleteAuthServer(authID)
	if err != nil {
		log.Printf("[ERROR] There was a problem deleting the auth server: %v", err)
		return diag.FromErr(err)
	}

	log.Printf("[DELETED] Deleted auth server with id %d", authID)
	d.SetId("")

	return nil
}
