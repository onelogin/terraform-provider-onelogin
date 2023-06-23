package apps

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
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
	ConfigurableAttribute types.String `tfsdk:"configurable_attribute"`
	Defaulted             types.String `tfsdk:"defaulted"`
	Id                    types.String `tfsdk:"id"`
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
				Default:             
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "app description",
				Optional:            true,
				Default:             "",
			},
			"notes": schema.StringAttribute{
				MarkdownDescription: "app notes",
				Optional:            true,
				Default:             "",
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
				Default:             types.BoolDefault(false),
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
			"provisioning": schema.NestedAttribute{
				MarkdownDescription: "app provisioning",
				Optional:            true,
				Computed:            true,
				schema.NestedAttributeObject{
					Attribute: schema.StringAttribute{
						"enabled": schema.BoolAttribute{
							MarkdownDescription: "app provisioning enabled",
							Optional:            true,
							Default:             types.BoolDefault(false),
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
	data.Id = types.StringValue("app-id")

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
