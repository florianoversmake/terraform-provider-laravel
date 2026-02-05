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
	ID               types.Int64  `tfsdk:"id"`
	ServerID         types.Int64  `tfsdk:"server_id"`
	Domain           types.String `tfsdk:"domain"`
	ProjectType      types.String `tfsdk:"project_type"`
	Aliases          types.List   `tfsdk:"aliases"`
	Directory        types.String `tfsdk:"directory"`
	Isolated         types.Bool   `tfsdk:"isolated"`
	Username         types.String `tfsdk:"username"`
	PHPVersion       types.String `tfsdk:"php_version"`
	Wildcards        types.Bool   `tfsdk:"wildcards"`
	Status           types.String `tfsdk:"status"`
	CreatedAt        types.String `tfsdk:"created_at"`
	WebDirectory     types.String `tfsdk:"web_directory"`
	DeleteProtection types.Bool   `tfsdk:"delete_protection"`
}

func NewForgeSiteResource() resource.Resource {
	return &ForgeSiteResource{}
}

func (r *ForgeSiteResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_forge_site"
}

func (r *ForgeSiteResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Forge site resource. This resource allows you to manage sites in Forge.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Computed: true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"server_id": schema.Int64Attribute{
				Required:            true,
				MarkdownDescription: "The ID of the server the site is on.",
			},
			"domain": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The domain name for the site.",
			},
			"project_type": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The type of project (e.g., 'php', 'html').",
			},
			"aliases": schema.ListAttribute{
				ElementType:         types.StringType,
				Optional:            true,
				Computed:            true,
				Default:             listdefault.StaticValue(types.ListValueMust(types.StringType, []attr.Value{})),
				MarkdownDescription: "List of additional domain names (aliases) for the site.",
			},
			"directory": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The web directory where the site files are located.",
			},
			"isolated": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
				MarkdownDescription: "Whether the site is isolated. If true, a username must be provided.",
			},
			"username": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("forge"),
				MarkdownDescription: "The username for the isolated site. Required if `isolated` is true. Default is 'forge'.",
			},
			"php_version": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString("php82"),
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
			"delete_protection": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
				MarkdownDescription: "This is a virtual attribute and not in the API. It is used to prevent accidental deletion of the site.",
			},
		},
	}
}

func (r *ForgeSiteResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	providerConfig, ok := req.ProviderData.(*providerConfig)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Provider Configure Type",
			fmt.Sprintf("Expected *providerConfig, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	if providerConfig.Forge == nil {
		resp.Diagnostics.AddError(
			"Forge Client Not Configured",
			"This resource requires the Forge API token to be configured in the provider. "+
				"Please set the 'forge_api_token' attribute in the provider configuration.",
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

	if plan.Username.IsNull() && plan.Isolated.ValueBool() {
		resp.Diagnostics.AddError("Username required for isolated sites", "A username must be provided for isolated sites.")
		return
	}

	webDir := plan.Directory.ValueString()

	// Build the CreateSiteRequest payload using the new API field names.
	payload := forge_client.CreateSiteRequest{
		Name:         plan.Domain.ValueString(),
		Type:         plan.ProjectType.ValueString(),
		WebDirectory: &webDir,
		IsIsolated:   plan.Isolated.ValueBool(),
		IsolatedUser: plan.Username.ValueString(),
		PHPVersion:   plan.PHPVersion.ValueString(),
	}

	// Call CreateSite on the client using the provided server_id.
	site, err := r.client.CreateSite(ctx, int(plan.ServerID.ValueInt64()), payload)
	if err != nil {
		resp.Diagnostics.AddError("Error creating site", err.Error())
		return
	}

	// Update plan state with response values.
	plan.ID = types.Int64Value(site.ID)
	plan.Domain = types.StringValue(site.Name)
	plan.ProjectType = types.StringValue(site.AppType)
	listVal, diags := types.ListValueFrom(ctx, types.StringType, site.Aliases)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}
	plan.Aliases = listVal
	plan.Directory = types.StringValue(site.WebDirectory)
	plan.Isolated = types.BoolValue(site.Isolated)
	plan.Username = types.StringValue(site.User)
	plan.PHPVersion = types.StringValue(site.PHPVersion)
	if site.Wildcards != nil {
		plan.Wildcards = types.BoolValue(*site.Wildcards)
	} else {
		plan.Wildcards = types.BoolValue(false)
	}
	plan.Status = types.StringValue(site.Status)
	if site.CreatedAt != nil {
		plan.CreatedAt = types.StringValue(*site.CreatedAt)
	} else {
		plan.CreatedAt = types.StringValue("")
	}
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
	state.ProjectType = types.StringValue(site.AppType)
	listVal, diags := types.ListValueFrom(ctx, types.StringType, site.Aliases)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}
	state.Aliases = listVal
	state.Directory = types.StringValue(site.WebDirectory)
	state.Isolated = types.BoolValue(site.Isolated)
	state.Username = types.StringValue(site.User)
	state.PHPVersion = types.StringValue(site.PHPVersion)
	if site.Wildcards != nil {
		state.Wildcards = types.BoolValue(*site.Wildcards)
	} else {
		state.Wildcards = types.BoolValue(false)
	}
	state.Status = types.StringValue(site.Status)
	if site.CreatedAt != nil {
		state.CreatedAt = types.StringValue(*site.CreatedAt)
	} else {
		state.CreatedAt = types.StringValue("")
	}
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

	dir := plan.Directory.ValueString()

	// Build the UpdateSiteRequest payload using the new API field names.
	updateReq := forge_client.UpdateSiteRequest{
		Directory:  &dir,
		PHPVersion: plan.PHPVersion.ValueString(),
		Type:       plan.ProjectType.ValueString(),
	}

	site, err := r.client.UpdateSite(ctx, int(plan.ServerID.ValueInt64()), int(plan.ID.ValueInt64()), updateReq)
	if err != nil {
		resp.Diagnostics.AddError("Error updating site", err.Error())
		return
	}

	// Update state with new values.
	plan.Domain = types.StringValue(site.Name)
	plan.Directory = types.StringValue(site.WebDirectory)
	plan.PHPVersion = types.StringValue(site.PHPVersion)
	listVal, diags := types.ListValueFrom(ctx, types.StringType, site.Aliases)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}
	plan.Aliases = listVal
	if site.Wildcards != nil {
		plan.Wildcards = types.BoolValue(*site.Wildcards)
	} else {
		plan.Wildcards = types.BoolValue(false)
	}
	plan.Status = types.StringValue(site.Status)
	if site.CreatedAt != nil {
		plan.CreatedAt = types.StringValue(*site.CreatedAt)
	} else {
		plan.CreatedAt = types.StringValue("")
	}
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

	// Check if delete protection is enabled.
	if state.DeleteProtection.ValueBool() {
		return
	}
	// Call the client to delete the site.
	err := r.client.DeleteSite(ctx, int(state.ServerID.ValueInt64()), int(state.ID.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError("Error deleting site", err.Error())
		return
	}
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
	stateModel.ProjectType = types.StringValue(site.AppType)
	listVal, diags := types.ListValueFrom(ctx, types.StringType, site.Aliases)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}
	stateModel.Aliases = listVal
	stateModel.Directory = types.StringValue(site.WebDirectory)
	stateModel.Isolated = types.BoolValue(site.Isolated)
	stateModel.Username = types.StringValue(site.User)
	stateModel.PHPVersion = types.StringValue(site.PHPVersion)
	if site.Wildcards != nil {
		stateModel.Wildcards = types.BoolValue(*site.Wildcards)
	} else {
		stateModel.Wildcards = types.BoolValue(false)
	}
	stateModel.Status = types.StringValue(site.Status)
	if site.CreatedAt != nil {
		stateModel.CreatedAt = types.StringValue(*site.CreatedAt)
	} else {
		stateModel.CreatedAt = types.StringValue("")
	}
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
