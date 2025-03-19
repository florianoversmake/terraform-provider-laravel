package provider

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"terraform-provider-laravel/internal/forge_client"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure ForgeSiteResource satisfies required interfaces.
var _ resource.Resource = &ForgeSiteResource{}
var _ resource.ResourceWithImportState = &ForgeSiteResource{}

// ForgeSiteResource implements a Terraform resource for a Forge site.
type ForgeSiteResource struct {
	client *forge_client.Client
}

// Resource model.
type ForgeSiteResourceModel struct {
	ID            types.Int64  `tfsdk:"id"`
	ServerID      types.Int64  `tfsdk:"server_id"`
	Domain        types.String `tfsdk:"domain"`
	ProjectType   types.String `tfsdk:"project_type"`
	Aliases       types.List   `tfsdk:"aliases"`
	Directory     types.String `tfsdk:"directory"`
	Isolated      types.Bool   `tfsdk:"isolated"`
	Username      types.String `tfsdk:"username"`
	Database      types.String `tfsdk:"database"`
	PHPVersion    types.String `tfsdk:"php_version"`
	NginxTemplate types.String `tfsdk:"nginx_template"`
	Wildcards     types.Bool   `tfsdk:"wildcards"`
	Status        types.String `tfsdk:"status"`
	CreatedAt     types.String `tfsdk:"created_at"`
	WebDirectory  types.String `tfsdk:"web_directory"`
}

// NewForgeSiteResource is a helper function to instantiate the resource.
func NewForgeSiteResource() resource.Resource {
	return &ForgeSiteResource{}
}

func (r *ForgeSiteResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_forge_site"
}

func (r *ForgeSiteResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Computed: true,
				// Use state for unknown IDs.
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"server_id": schema.Int64Attribute{
				Required: true,
			},
			"domain": schema.StringAttribute{
				Required: true,
			},
			"project_type": schema.StringAttribute{
				Required: true,
			},
			"aliases": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Computed:    true,
				Default:     listdefault.StaticValue(types.ListValueMust(types.StringType, []attr.Value{})),
			},
			"directory": schema.StringAttribute{
				Required: true,
			},
			"isolated": schema.BoolAttribute{
				Optional: true,
				Computed: true,
				Default:  booldefault.StaticBool(false),
			},
			"username": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString("forge"),
			},
			"database": schema.StringAttribute{
				Optional: true,
			},
			"php_version": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString("php82"),
			},
			"nginx_template": schema.StringAttribute{
				Optional: true,
			},
			"wildcards": schema.BoolAttribute{
				Optional: true,
				Computed: true,
			},
			"status": schema.StringAttribute{
				Computed: true,
			},
			"created_at": schema.StringAttribute{
				Computed: true,
			},
			"web_directory": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

func (r *ForgeSiteResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	providerConfig, ok := req.ProviderData.(*providerConfig)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Provider Configure Type",
			fmt.Sprintf("Expected *providerConfig, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	r.client = providerConfig.Forge
}

func (r *ForgeSiteResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan ForgeSiteResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Convert aliases list to []string.
	var aliases []string
	if !plan.Aliases.IsNull() {
		diags := plan.Aliases.ElementsAs(ctx, &aliases, false)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	if plan.Username.IsNull() && plan.Isolated.ValueBool() {
		resp.Diagnostics.AddError("Username required for isolated sites", "A username must be provided for isolated sites.")
		return
	}

	// Build the CreateSiteRequest payload.
	payload := forge_client.CreateSiteRequest{
		Domain:      plan.Domain.ValueString(),
		ProjectType: plan.ProjectType.ValueString(),
		Aliases:     aliases,
		Directory:   plan.Directory.ValueString(),
		Isolated:    plan.Isolated.ValueBool(),
		Username:    plan.Username.ValueString(),
		PHPVersion:  plan.PHPVersion.ValueString(),
	}
	// Optional fields.
	if !plan.Database.IsNull() && plan.Database.ValueString() != "" {
		payload.Database = plan.Database.ValueString()
	}
	if !plan.NginxTemplate.IsNull() && plan.NginxTemplate.ValueString() != "" {
		payload.NginxTemplate = plan.NginxTemplate.ValueString()
	}

	// Call CreateSite on the client using the provided server_id.
	site, err := r.client.CreateSite(ctx, int(plan.ServerID.ValueInt64()), payload)
	if err != nil {
		resp.Diagnostics.AddError("Error creating site", err.Error())
		return
	}

	// Update plan state with response values.
	plan.ID = types.Int64Value(site.ID)
	plan.Domain = types.StringValue(site.Name) // Assume site.Name equals the domain.
	plan.ProjectType = types.StringValue(site.ProjectType)
	listVal, diags := types.ListValueFrom(ctx, types.StringType, site.Aliases)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}
	plan.Aliases = listVal
	plan.Directory = types.StringValue(site.Directory)
	plan.Isolated = types.BoolValue(site.Isolated)
	plan.Username = types.StringValue(site.Username)
	plan.PHPVersion = types.StringValue(site.PHPVersion)
	plan.Wildcards = types.BoolValue(site.Wildcards)
	plan.Status = types.StringValue(site.Status)
	plan.CreatedAt = types.StringValue(site.CreatedAt)
	plan.WebDirectory = types.StringValue(site.WebDirectory)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *ForgeSiteResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state ForgeSiteResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	site, err := r.client.GetSite(ctx, int(state.ServerID.ValueInt64()), int(state.ID.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError("Error reading site", err.Error())
		return
	}

	state.Domain = types.StringValue(site.Name)
	state.ProjectType = types.StringValue(site.ProjectType)
	listVal, diags := types.ListValueFrom(ctx, types.StringType, site.Aliases)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}
	state.Aliases = listVal
	state.Directory = types.StringValue(site.Directory)
	state.Isolated = types.BoolValue(site.Isolated)
	state.Username = types.StringValue(site.Username)
	state.PHPVersion = types.StringValue(site.PHPVersion)
	state.Wildcards = types.BoolValue(site.Wildcards)
	state.Status = types.StringValue(site.Status)
	state.CreatedAt = types.StringValue(site.CreatedAt)
	state.WebDirectory = types.StringValue(site.WebDirectory)

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *ForgeSiteResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan ForgeSiteResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Convert aliases list to []string.
	var aliases []string
	if !plan.Aliases.IsNull() {
		diags := plan.Aliases.ElementsAs(ctx, &aliases, false)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	// Build the UpdateSiteRequest payload.
	updateReq := forge_client.UpdateSiteRequest{
		Name:       plan.Domain.ValueString(),
		Directory:  plan.Directory.ValueString(),
		PHPVersion: plan.PHPVersion.ValueString(),
		Aliases:    aliases,
		Wildcards:  plan.Wildcards.ValueBool(),
	}

	site, err := r.client.UpdateSite(ctx, int(plan.ServerID.ValueInt64()), int(plan.ID.ValueInt64()), updateReq)
	if err != nil {
		resp.Diagnostics.AddError("Error updating site", err.Error())
		return
	}

	// Update state with new values.
	plan.Domain = types.StringValue(site.Name)
	plan.Directory = types.StringValue(site.Directory)
	plan.PHPVersion = types.StringValue(site.PHPVersion)
	listVal, diags := types.ListValueFrom(ctx, types.StringType, site.Aliases)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}
	plan.Aliases = listVal
	plan.Wildcards = types.BoolValue(site.Wildcards)
	plan.Status = types.StringValue(site.Status)
	plan.CreatedAt = types.StringValue(site.CreatedAt)
	plan.WebDirectory = types.StringValue(site.WebDirectory)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *ForgeSiteResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state ForgeSiteResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	//err := r.client.DeleteSite(ctx, int(state.ServerID.ValueInt64()), int(state.ID.ValueInt64()))
	// if err != nil {
	// 	resp.Diagnostics.AddError("Error deleting site", err.Error())
	// 	return
	// }
}

func (r *ForgeSiteResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Expect import ID in format "server_id:site_id"
	parts := splitCompositeID(req.ID, 2)
	if parts == nil {
		resp.Diagnostics.AddError("Invalid import ID", "Expected format: server_id:site_id")
		return
	}
	serverID, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		resp.Diagnostics.AddError("Invalid server_id", err.Error())
		return
	}
	siteID, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		resp.Diagnostics.AddError("Invalid site_id", err.Error())
		return
	}

	site, err := r.client.GetSite(ctx, int(serverID), int(siteID))
	if err != nil {
		resp.Diagnostics.AddError("Error reading site", err.Error())
		return
	}

	var stateModel ForgeSiteResourceModel
	stateModel.ID = types.Int64Value(site.ID)
	stateModel.ServerID = types.Int64Value(serverID)
	stateModel.Domain = types.StringValue(site.Name)
	stateModel.ProjectType = types.StringValue(site.ProjectType)
	listVal, diags := types.ListValueFrom(ctx, types.StringType, site.Aliases)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}
	stateModel.Aliases = listVal
	stateModel.Directory = types.StringValue(site.Directory)
	stateModel.Isolated = types.BoolValue(site.Isolated)
	stateModel.Username = types.StringValue(site.Username)
	stateModel.PHPVersion = types.StringValue(site.PHPVersion)
	stateModel.Wildcards = types.BoolValue(site.Wildcards)
	stateModel.Status = types.StringValue(site.Status)
	stateModel.CreatedAt = types.StringValue(site.CreatedAt)
	stateModel.WebDirectory = types.StringValue(site.WebDirectory)

	diags = resp.State.Set(ctx, stateModel)
	resp.Diagnostics.Append(diags...)
}

// splitCompositeID splits an import id by colon and ensures the expected number of parts.
func splitCompositeID(id string, expectedParts int) []string {
	parts := strings.Split(id, ":")
	if len(parts) != expectedParts {
		return nil
	}
	return parts
}
