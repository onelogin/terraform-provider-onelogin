package app_resource

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &appResource{}
var _ resource.ResourceWithImportState = &appResource{}

func NewAppResource() resource.Resource {
	return &appResource{}
}

// appResource defines the resource implementation.
type appResource struct {
	client *http.Client
}

// appResourceModel describes the resource data model.
type appResourceModel struct {
	Name               types.String              `tfsdk:"name"`
	Visible            types.Bool                `tfsdk:"visible"`
	Description        types.String              `tfsdk:"description"`
	ConnectorID        types.Int64               `tfsdk:"connector_id"`
	ID                 types.Int64               `tfsdk:"id"`
	Notes              types.String              `tfsdk:"notes"`
	PolicyID           types.Int64               `tfsdk:"policy_id"`
	BrandID            types.Int64               `tfsdk:"brand_id"`
	IconURL            types.String              `tfsdk:"icon_url"`
	AuthMethod         types.Int64               `tfsdk:"auth_method"`
	TabID              types.Int64               `tfsdk:"tab_id"`
	CreatedAt          types.String              `tfsdk:"created_at"`
	UpdatedAt          types.String              `tfsdk:"updated_at"`
	RoleIDs            types.List                `tfsdk:"role_ids"`
	AllowAssumedSignin types.Bool                `tfsdk:"allow_assumed_signin"`
	Provisioning       ProvisioningModel         `tfsdk:"provisioning"`
	SSO                interface{}               `tfsdk:"sso"`
	Configuration      interface{}               `tfsdk:"configuration"`
	Parameters         map[string]ParameterModel `tfsdk:"parameters"`
	EnforcementPoint   EnforcementPointModel     `tfsdk:"enforcement_point"`
}

type ParameterModel struct {
	Values                    interface{}  `tfsdk:"values"`
	UserAttributeMappings     interface{}  `tfsdk:"user_attribute_mappings"`
	ProvisionedEntitlements   types.Bool   `tfsdk:"provisioned_entitlements"`
	SkipIfBlank               types.Bool   `tfsdk:"skip_if_blank"`
	ID                        types.Int64  `tfsdk:"id"`
	DefaultValues             interface{}  `tfsdk:"default_values"`
	AttributesTransformations interface{}  `tfsdk:"attributes_transformations"`
	Label                     types.String `tfsdk:"label"`
	UserAttributeMacros       interface{}  `tfsdk:"user_attribute_macros"`
	IncludeInSamlAssertion    types.Bool   `tfsdk:"include_in_saml_assertion"`
}

type ConfigurationOpenIdModel struct {
	RedirectURI             types.String `tfsdk:"redirect_uri"`
	LoginURL                types.String `tfsdk:"login_url"`
	OidcApplicationType     types.Int64  `tfsdk:"oidc_application_type"`
	TokenEndpointAuthMethod types.Int64  `tfsdk:"token_endpoint_auth_method"`
}

type ConfigurationSAMLModel struct {
	ProviderArn        interface{}  `tfsdk:"provider_arn"`
	SignatureAlgorithm types.String `tfsdk:"signature_algorithm"`
	CertificateID      types.Int64  `tfsdk:"certificate_id"`
}

type ProvisioningModel struct {
	Enabled types.Bool `tfsdk:"enabled"`
}

type SSOOpenIdModel struct {
	ClientID types.String `tfsdk:"client_id"`
}

type SSOSAMLModel struct {
	MetadataURL types.String     `tfsdk:"metadata_url"`
	AcsURL      types.String     `tfsdk:"acs_url"`
	SlsURL      types.String     `tfsdk:"sls_url"`
	Issuer      types.String     `tfsdk:"issuer"`
	Certificate CertificateModel `tfsdk:"certificate"`
}

type CertificateModel struct {
	ID    types.Int64  `tfsdk:"id"`
	Name  types.String `tfsdk:"name"`
	Value types.String `tfsdk:"value"`
}

type EnforcementPointModel struct {
	RequireSitewideAuthentication types.Bool                 `tfsdk:"require_sitewide_authentication"`
	Conditions                    *ConditionsModel           `tfsdk:"conditions,omitempty"`
	SessionExpiryFixed            DurationModel              `tfsdk:"session_expiry_fixed"`
	SessionExpiryInactivity       DurationModel              `tfsdk:"session_expiry_inactivity"`
	Permissions                   types.String               `tfsdk:"permissions"`
	Token                         types.String               `tfsdk:"token,omitempty"`
	Target                        types.String               `tfsdk:"target"`
	Resources                     []EnforcementResourceModel `tfsdk:"resources"`
	ContextRoot                   types.String               `tfsdk:"context_root"`
	UseTargetHostHeader           types.Bool                 `tfsdk:"use_target_host_header"`
	Vhost                         types.String               `tfsdk:"vhost"`
	LandingPage                   types.String               `tfsdk:"landing_page"`
	CaseSensitive                 types.Bool                 `tfsdk:"case_sensitive"`
}

type ConditionsModel struct {
	Type  types.String `tfsdk:"type"`
	Roles types.List   `tfsdk:"roles"`
}

type DurationModel struct {
	Value types.Int64 `tfsdk:"value"`
	Unit  types.Int64 `tfsdk:"unit"`
}

type EnforcementResourceModel struct {
	Path        types.String  `tfsdk:"path"`
	RequireAuth types.String  `tfsdk:"require_authentication"`
	Permissions types.String  `tfsdk:"permissions"`
	Conditions  *types.String `tfsdk:"conditions,omitempty"`
	IsPathRegex *types.Bool   `tfsdk:"is_path_regex,omitempty"`
	ResourceID  types.Int64   `tfsdk:"resource_id,omitempty"`
}

func (r *appResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_app"
}

func (r *appResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "app resource",

		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				MarkdownDescription: "app name",
				Required:            true,
			},
			"visible": schema.BoolAttribute{
				MarkdownDescription: "app visibility",
				Optional:            true,
				Default:             booldefault.StaticBool(true),
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "app description",
				Optional:            true,
				Default:             stringdefault.StaticString(""),
			},
			"notes": schema.StringAttribute{
				MarkdownDescription: "app notes",
				Optional:            true,
				Default:             stringdefault.StaticString(""),
			},
			"icon_url": schema.StringAttribute{
				MarkdownDescription: "app icon url",
				Computed:            true,
			},
			"auth_method": schema.Int64Attribute{
				MarkdownDescription: "app auth method",
				Computed:            true,
			},
			"connector_id": schema.Int64Attribute{
				MarkdownDescription: "app connector id",
				Required:            true,
			},
			"policy_id": schema.Int64Attribute{
				MarkdownDescription: "app policy id",
				Computed:            true,
			},
			"brand_id": schema.Int64Attribute{
				MarkdownDescription: "app brand id",
				Computed:            true,
			},
			"allow_assumed_signin": schema.BoolAttribute{
				MarkdownDescription: "app allow assumed signin",
				Optional:            true,
				Default:             booldefault.StaticBool(true),
			},
			"tab_id": schema.Int64Attribute{
				MarkdownDescription: "app tab id",
				Computed:            true,
			},
			"created_at": schema.StringAttribute{
				MarkdownDescription: "app created at",
				Computed:            true,
			},
			"updated_at": schema.StringAttribute{
				MarkdownDescription: "app updated at",
				Computed:            true,
			},
			"provisioning": schema.ListNestedAttribute{
				MarkdownDescription: "app provisioning",
				Optional:            true,
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"enabled": schema.BoolAttribute{
							MarkdownDescription: "app parameter key",
							Optional:            true,
							Default:             booldefault.StaticBool(false),
						},
					},
				},
			},
			"parameters": schema.ListNestedAttribute{
				MarkdownDescription: "app parameters",
				Optional:            true,
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"param_key": schema.StringAttribute{
							MarkdownDescription: "app parameter key",
							Required:            true,
						},
						"param_id": schema.Int64Attribute{
							MarkdownDescription: "app parameter id",
							Computed:            true,
						},
						"label": schema.StringAttribute{
							MarkdownDescription: "app parameter label",
							Optional:            true,
							Computed:            true,
						},
						"user_attribute_mappings": schema.StringAttribute{
							MarkdownDescription: "app parameter user attribute mappings",
							Optional:            true,
							Computed:            true,
						},
						"user_attribute_macros": schema.StringAttribute{
							MarkdownDescription: "app parameter user attribute macros",
							Optional:            true,
							Computed:            true,
						},
						"attributes_transformations": schema.StringAttribute{
							MarkdownDescription: "app parameter attributes transformations",
							Optional:            true,
							Computed:            true,
						},
						"default_value": schema.StringAttribute{
							MarkdownDescription: "app parameter default value",
							Optional:            true,
							Computed:            true,
						},
						"skip_if_blank": schema.BoolAttribute{
							MarkdownDescription: "app parameter skip if blank",
							Optional:            true,
							Computed:            true,
						},
						"values": schema.StringAttribute{
							MarkdownDescription: "app parameter values",
							Optional:            true,
							Computed:            true,
						},
						"provisioned_entitlements": schema.BoolAttribute{
							MarkdownDescription: "app parameter provisioned entitlements",
							Optional:            true,
							Computed:            true,
						},
						"safe_entitlements_enabled": schema.BoolAttribute{
							MarkdownDescription: "app parameter safe entitlements enabled",
							Optional:            true,
							Computed:            true,
						},
						"include_in_saml_assertion": schema.BoolAttribute{
							MarkdownDescription: "app parameter include in saml assertion",
							Optional:            true,
							Computed:            true,
						},
					},
				},
			},
		},
	}
}

func (r *appResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*http.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *http.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

func (r *appResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *appResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// If applicable, this is a great opportunity to initialize any necessary
	// provider client data and make a call using it.
	// httpResp, err := r.client.Do(httpReq)
	// if err != nil {
	//     resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create app, got error: %s", err))
	//     return
	// }

	// For the purposes of this app code, hardcoding a response value to
	// save into the Terraform state.

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, "created a resource")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *appResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *appResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// If applicable, this is a great opportunity to initialize any necessary
	// provider client data and make a call using it.
	// httpResp, err := r.client.Do(httpReq)
	// if err != nil {
	//     resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read app, got error: %s", err))
	//     return
	// }

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *appResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *appResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// If applicable, this is a great opportunity to initialize any necessary
	// provider client data and make a call using it.
	// httpResp, err := r.client.Do(httpReq)
	// if err != nil {
	//     resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update app, got error: %s", err))
	//     return
	// }

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *appResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *appResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// If applicable, this is a great opportunity to initialize any necessary
	// provider client data and make a call using it.
	// httpResp, err := r.client.Do(httpReq)
	// if err != nil {
	//     resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete app, got error: %s", err))
	//     return
	// }
}

func (r *appResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
