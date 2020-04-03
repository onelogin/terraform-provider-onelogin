package resources

import (
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/onelogin/onelogin-go-sdk/pkg/client"
	"github.com/onelogin/onelogin-go-sdk/pkg/models"
)

func OneloginApps() *schema.Resource {
	return &schema.Resource{
		Create: appCreate,
		Read:   appRead,
		Update: appUpdate,
		Delete: appDelete,
		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"visible": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"connector_id": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},
			"auth_method": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			},
			"parameters": &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"param_key_name": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
						"param_id": &schema.Schema{
							Type:     schema.TypeInt,
							Optional: true,
						},
						"label": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"user_attribute_mappings": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"user_attribute_macros": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"attributes_transformations": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"default_values": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"skip_if_blank": &schema.Schema{
							Type:     schema.TypeBool,
							Optional: true,
						},
						"values": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"provisioned_entitlements": &schema.Schema{
							Type:     schema.TypeBool,
							Optional: true,
						},
						"safe_entitlements_enabled": &schema.Schema{
							Type:     schema.TypeBool,
							Optional: true,
						},
					},
				},
			},
		},
	}
}

func appCreate(d *schema.ResourceData, m interface{}) error {
	paramsList := d.Get("parameters").(*schema.Set).List()

	appParams := make(map[string]models.AppParameters, len(paramsList))

	for _, paramI := range paramsList {
		param := paramI.(map[string]interface{})
		pid := int32(param["param_id"].(int))
		lbl := param["label"].(string)
		uam := param["user_attribute_mappings"].(string)
		uac := param["user_attribute_macros"].(string)
		atr := param["attributes_transformations"].(string)
		sib := param["skip_if_blank"].(bool)
		val := param["values"].(string)
		dfv := param["default_values"].(string)
		pet := param["provisioned_entitlements"].(bool)
		see := param["safe_entitlements_enabled"].(bool)
		key := param["param_key_name"].(string)

		appParams[key] = models.AppParameters{
			ID:                        &pid,
			Label:                     &lbl,
			UserAttributeMappings:     &uam,
			UserAttributeMacros:       &uac,
			AttributesTransformations: &atr,
			SkipIfBlank:               &sib,
			Values:                    &val,
			DefaultValues:             &dfv,
			ProvisionedEntitlements:   &pet,
			SafeEntitlementsEnabled:   &see,
		}
	}

	vis := d.Get("visible").(bool)
	nam := d.Get("name").(string)
	des := d.Get("description").(string)
	cid := int32(d.Get("connector_id").(int))
	aum := int32(d.Get("auth_method").(int))

	app := &models.App{
		Visible:     &vis,
		Name:        &nam,
		Description: &des,
		ConnectorID: &cid,
		AuthMethod:  &aum,
		Parameters:  appParams,
	}

	client := m.(*client.APIClient)
	resp, app, err := client.Services.AppsV2.CreateApp(app)
	if err != nil {
		log.Printf("[ERROR] There was a problem creating the app!")
		log.Println(err)
	}
	log.Printf("[CREATED] Created app with %d", *(app.ID))
	log.Println(resp)
	d.SetId(fmt.Sprintf("%d", *(app.ID)))
	return nil
}

func appRead(d *schema.ResourceData, m interface{}) error {
	return nil
}

func appUpdate(d *schema.ResourceData, m interface{}) error {
	paramsList := d.Get("parameters").(*schema.Set).List()

	appParams := make(map[string]models.AppParameters, len(paramsList))

	for _, paramI := range paramsList {
		param := paramI.(map[string]interface{})
		pid := int32(param["param_id"].(int))
		lbl := param["label"].(string)
		uam := param["user_attribute_mappings"].(string)
		uac := param["user_attribute_macros"].(string)
		atr := param["attributes_transformations"].(string)
		sib := param["skip_if_blank"].(bool)
		val := param["values"].(string)
		dfv := param["default_values"].(string)
		pet := param["provisioned_entitlements"].(bool)
		see := param["safe_entitlements_enabled"].(bool)
		key := param["param_key_name"].(string)

		appParams[key] = models.AppParameters{
			ID:                        &pid,
			Label:                     &lbl,
			UserAttributeMappings:     &uam,
			UserAttributeMacros:       &uac,
			AttributesTransformations: &atr,
			SkipIfBlank:               &sib,
			Values:                    &val,
			DefaultValues:             &dfv,
			ProvisionedEntitlements:   &pet,
			SafeEntitlementsEnabled:   &see,
		}
	}

	aid, _ := strconv.Atoi(d.Id())
	vis := d.Get("visible").(bool)
	nam := d.Get("name").(string)
	des := d.Get("description").(string)
	cid := int32(d.Get("connector_id").(int))
	aum := int32(d.Get("auth_method").(int))

	app := &models.App{
		Visible:     &vis,
		Name:        &nam,
		Description: &des,
		ConnectorID: &cid,
		AuthMethod:  &aum,
		Parameters:  appParams,
	}

	client := m.(*client.APIClient)
	resp, app, err := client.Services.AppsV2.UpdateAppByID(int32(aid), app)
	if err != nil {
		log.Printf("[ERROR] There was a problem creating the app!")
		log.Println(err)
	}
	log.Printf("[UPDATED] Updated app with %d", *(app.ID))
	log.Println(resp)
	d.SetId(fmt.Sprintf("%d", *(app.ID)))
	return nil
}

func appDelete(d *schema.ResourceData, m interface{}) error {
	aid, err := strconv.Atoi(d.Id())
	if err != nil {
		log.Printf("[ERROR] There was a problem reading the id!")
		log.Println(err)
		return nil
	}
	client := m.(*client.APIClient)
	resp, err := client.Services.AppsV2.DeleteApp(int32(aid))
	if err != nil {
		log.Printf("[ERROR] There was a problem creating the app!")
		log.Println(err)
	} else {
		log.Printf("[DELETED] Deleted app with %d", aid)
		log.Println(resp)
		d.SetId("")
	}

	return nil
}
