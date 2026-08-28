package onelogin

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/onelogin/onelogin-go-sdk/v4/pkg/onelogin"
	brandschema "github.com/onelogin/terraform-provider-onelogin/ol_schema/brand"
	"github.com/onelogin/terraform-provider-onelogin/utils"
)

// Brands returns the onelogin_brands resource: an account brand, which is the
// look of the login page and the portal for the users it applies to.
//
// A brand is attached to an app through the app's brand_id. The account's
// master brand is created with the account and cannot be created or destroyed
// here; import it if you want to manage it.
func Brands() *schema.Resource {
	return &schema.Resource{
		CreateContext: brandCreate,
		ReadContext:   brandRead,
		UpdateContext: brandUpdate,
		DeleteContext: brandDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: brandschema.Schema(),
	}
}

func brandCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*onelogin.OneloginSDK)

	body := brandschema.RequestBody(d, brandschema.ConfiguredKeys(d.GetRawConfig()))

	tflog.Info(ctx, "[CREATE] Creating brand", map[string]interface{}{"name": d.Get("name")})

	result, err := client.CreateBrand(body)
	if err != nil {
		return utils.HandleAPIError(ctx, err, utils.ErrorCategoryCreate, "Brand", "")
	}

	brand, ok := brandFromResponse(result)
	if !ok {
		return diag.Errorf("expected a brand object in the create response, got %T", result)
	}

	id, ok := brandschema.ID(brand)
	if !ok {
		return diag.Errorf("failed to extract brand ID from the create response")
	}
	d.SetId(id)

	tflog.Info(ctx, "[CREATED] Created brand", map[string]interface{}{"id": id})

	return brandRead(ctx, d, m)
}

func brandRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*onelogin.OneloginSDK)

	tflog.Info(ctx, "[READ] Reading brand", map[string]interface{}{"id": d.Id()})

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	result, err := client.GetBrandByID(id, nil)
	if err != nil {
		// A brand that is gone is not an error on a read. Somebody deleted it
		// outside Terraform; dropping it from state lets the next plan offer to
		// create it again rather than failing every run until it is imported or
		// removed by hand.
		if utils.IsNotFoundError(err) {
			tflog.Warn(ctx, "[READ] Brand is gone, removing it from state", map[string]interface{}{"id": d.Id()})
			d.SetId("")
			return nil
		}
		return utils.HandleAPIError(ctx, err, utils.ErrorCategoryRead, "Brand", d.Id())
	}

	brand, ok := brandFromResponse(result)
	if !ok {
		return diag.Errorf("expected a brand object in the read response, got %T", result)
	}

	if err := brandschema.Flatten(d, brand); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func brandUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*onelogin.OneloginSDK)

	body := brandschema.RequestBody(d, brandschema.ConfiguredKeys(d.GetRawConfig()))

	tflog.Info(ctx, "[UPDATE] Updating brand", map[string]interface{}{"id": d.Id()})

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	if _, err := client.UpdateBrand(id, body); err != nil {
		return utils.HandleAPIError(ctx, err, utils.ErrorCategoryUpdate, "Brand", d.Id())
	}

	return brandRead(ctx, d, m)
}

func brandDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// The master brand is the account's own branding, created with the account
	// and the fallback for every app that names no other brand. It is reachable
	// here only by importing it, and destroying it is not what anyone importing
	// it to change a colour is asking for.
	//
	// This is deliberately refused rather than attempted-and-reported. Whether
	// the endpoint would allow it was not tested: doing so would have meant
	// deleting the master brand of a shared tenant, and the answer does not
	// change what this resource should do. If the API refuses, this is a
	// clearer error than the one it would return; if it allows, this is the
	// difference between a destroy and an outage.
	if master, ok := d.Get("master").(bool); ok && master {
		return diag.Errorf(
			"brand %s is the account's master brand and will not be deleted: it is the branding every app falls back to. "+
				"Remove it from state with `terraform state rm` to stop managing it without destroying it",
			d.Id(),
		)
	}

	client := m.(*onelogin.OneloginSDK)

	return utils.StandardDeleteFunc(ctx, d, func(id string) (interface{}, error) {
		brandID, err := strconv.Atoi(id)
		if err != nil {
			return nil, err
		}
		return client.DeleteBrand(brandID)
	}, "Brand")
}

// brandFromResponse narrows the SDK's decoded body to the brand itself. The
// endpoint returns the brand at the top level rather than wrapped in a key, so
// the decoded body is the brand.
func brandFromResponse(result interface{}) (map[string]interface{}, bool) {
	brand, ok := result.(map[string]interface{})
	return brand, ok
}
