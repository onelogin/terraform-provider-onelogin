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
			"notes": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"icon_url": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"auth_method": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			},
			"policy_id": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			},
			"allow_assumed_signin": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
			},
			"tab_id": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			},
			"connector_id": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},
			"created_at": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"updated_at": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"provisioning": &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": &schema.Schema{
							Type:     schema.TypeBool,
							Required: true,
						},
					},
				},
			},
			"configuration": &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"redirect_uri": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"refresh_token_expiration_minutes": &schema.Schema{
							Type:     schema.TypeInt,
							Optional: true,
						},
						"login_url": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"oidc_application_type": &schema.Schema{
							Type:     schema.TypeInt,
							Optional: true,
						},
						"token_endpoint_auth_method": &schema.Schema{
							Type:     schema.TypeInt,
							Optional: true,
						},
						"access_token_expiration_minutes": &schema.Schema{
							Type:     schema.TypeInt,
							Optional: true,
						},
						"provider_arn": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"signature_algorithm": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
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

func inflateWithSchema(s map[string]interface{}, obj interface{}) {
	switch o := obj.(type) {
	case *models.AppProvisioning:
		e := s["enabled"].(bool)

		o.Enabled = &e
	case *models.AppParameters:
		pid := int32(s["param_id"].(int))
		lbl := s["label"].(string)
		uam := s["user_attribute_mappings"].(string)
		uac := s["user_attribute_macros"].(string)
		atr := s["attributes_transformations"].(string)
		sib := s["skip_if_blank"].(bool)
		val := s["values"].(string)
		dfv := s["default_values"].(string)
		pet := s["provisioned_entitlements"].(bool)
		see := s["safe_entitlements_enabled"].(bool)

		o.ID = &pid
		o.Label = &lbl
		o.UserAttributeMappings = &uam
		o.UserAttributeMacros = &uac
		o.AttributesTransformations = &atr
		o.SkipIfBlank = &sib
		o.Values = &val
		o.DefaultValues = &dfv
		o.ProvisionedEntitlements = &pet
		o.SafeEntitlementsEnabled = &see

	case *models.AppConfiguration:
		rui := s["redirect_uri"].(string)
		rte := int32(s["refresh_token_expiration_minutes"].(int))
		lur := s["login_url"].(string)
		oat := int32(s["oidc_application_type"].(int))
		tea := int32(s["token_endpoint_auth_method"].(int))
		ate := int32(s["access_token_expiration_minutes"].(int))
		par := s["provider_arn"].(string)
		sal := s["signature_algorithm"].(string)

		o.RedirectURI = &rui
		o.RefreshTokenExpirationMinutes = &rte
		o.LoginURL = &lur
		o.OidcApplicationType = &oat
		o.TokenEndpointAuthMethod = &tea
		o.AccessTokenExpirationMinutes = &ate
		o.ProviderArn = &par
		o.SignatureAlgorithm = &sal
	}
}

func buildAppObject(d *schema.ResourceData) *models.App {
	provisioningList := d.Get("provisioning").(*schema.Set).List()
	appProv := models.AppProvisioning{}
	for _, s := range provisioningList {
		sMap := s.(map[string]interface{})
		inflateWithSchema(sMap, &appProv)
	}

	paramsList := d.Get("parameters").(*schema.Set).List()
	appParams := make(map[string]models.AppParameters, len(paramsList))
	for _, s := range paramsList {
		sMap := s.(map[string]interface{})
		key := sMap["param_key_name"].(string)
		appParam := models.AppParameters{}
		inflateWithSchema(sMap, &appParam)
		appParams[key] = appParam
	}

	configList := d.Get("configuration").(*schema.Set).List()
	appConfig := models.AppConfiguration{}
	for _, s := range configList {
		sMap := s.(map[string]interface{})
		inflateWithSchema(sMap, appConfig)
	}

	nam := d.Get("name").(string)
	des := d.Get("description").(string)
	not := d.Get("notes").(string)
	iur := d.Get("icon_url").(string)

	app := models.App{
		Name:          &nam,
		Description:   &des,
		Notes:         &not,
		IconURL:       &iur,
		Parameters:    appParams,
		Provisioning:  &appProv,
		Configuration: &appConfig,
	}
	if vis, visSet := d.GetOk("visible"); visSet {
		vis := vis.(bool)
		app.Visible = &vis
	}
	if aas, aasSet := d.GetOk("allow_assumed_signin"); aasSet {
		aas := aas.(bool)
		app.AllowAssumedSignin = &aas
	}
	if cid, cidSet := d.GetOk("connector_id"); cidSet {
		cid := int32(cid.(int))
		app.ConnectorID = &cid
	}
	if aum, aumSet := d.GetOk("auth_method"); aumSet {
		aum := int32(aum.(int))
		app.AuthMethod = &aum
	}
	if pid, pidSet := d.GetOk("policy_id"); pidSet {
		pid := int32(pid.(int))
		app.PolicyID = &pid
	}
	if tid, tidSet := d.GetOk("tab_id"); tidSet {
		tid := int32(tid.(int))
		app.TabID = &tid
	}

	return &app
}

func appCreate(d *schema.ResourceData, m interface{}) error {
	app := buildAppObject(d)
	log.Println(app)
	client := m.(*client.APIClient)
	resp, app, err := client.Services.AppsV2.CreateApp(app)
	if err != nil {
		log.Printf("[ERROR] There was a problem creating the app!")
		log.Println(err)
	}
	log.Printf("[CREATED] Created app with %d", *(app.ID))
	log.Println(resp)
	d.SetId(fmt.Sprintf("%d", *(app.ID)))
	return appRead(d, m)
}

func appRead(d *schema.ResourceData, m interface{}) error {
	return nil
}

func appUpdate(d *schema.ResourceData, m interface{}) error {
	app := buildAppObject(d)

	aid, _ := strconv.Atoi(d.Id())

	client := m.(*client.APIClient)
	resp, app, err := client.Services.AppsV2.UpdateAppByID(int32(aid), app)
	if err != nil {
		log.Printf("[ERROR] There was a problem creating the app!")
		log.Println(err)
	}
	log.Printf("[UPDATED] Updated app with %d", *(app.ID))
	log.Println(resp)
	d.SetId(fmt.Sprintf("%d", *(app.ID)))
	return appRead(d, m)
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
