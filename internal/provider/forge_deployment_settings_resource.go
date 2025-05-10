package provider

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"terraform-provider-laravel/internal/forge_client"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &ForgeDeploymentSettingsResource{}
var _ resource.ResourceWithImportState = &ForgeDeploymentSettingsResource{}

func NewForgeDeploymentSettingsResource() resource.Resource {
	return &ForgeDeploymentSettingsResource{}
}

// ForgeDeploymentSettingsResource defines the resource implementation.
type ForgeDeploymentSettingsResource struct {
	client *forge_client.Client
}

// ForgeDeploymentSettingsResourceModel describes the resource data model.
type ForgeDeploymentSettingsResourceModel struct {
	ServerID      types.Int64  `tfsdk:"server_id"`
	SiteID        types.Int64  `tfsdk:"site_id"`
	QuickDeploy   types.Bool   `tfsdk:"quick_deploy"`
	Script        types.String `tfsdk:"script"`
	AutoSourceEnv types.Bool   `tfsdk:"auto_source_env"`
}

func (r *ForgeDeploymentSettingsResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_forge_deployment_settings"
}

func (r *ForgeDeploymentSettingsResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages deployment settings for a Laravel Forge site, including quick deployment and the deployment script.",

		Attributes: map[string]schema.Attribute{
			"server_id": schema.Int64Attribute{
				Required:            true,
				MarkdownDescription: "The ID of the server the site belongs to.",
			},
			"site_id": schema.Int64Attribute{
				Required:            true,
				MarkdownDescription: "The ID of the site to manage deployment settings for.",
			},
			"quick_deploy": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
				MarkdownDescription: "Whether quick deployment is enabled for this site.",
			},
			"script": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "The deployment script content.",
			},
			"auto_source_env": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
				MarkdownDescription: "Whether to automatically source environment variables into the deployment script.",
			},
		},
	}
}

func (r *ForgeDeploymentSettingsResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *ForgeDeploymentSettingsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan ForgeDeploymentSettingsResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	serverID := int(plan.ServerID.ValueInt64())
	siteID := int(plan.SiteID.ValueInt64())

	// Handle quick deployment
	if plan.QuickDeploy.ValueBool() {
		if err := r.client.EnableQuickDeployment(ctx, serverID, siteID); err != nil {
			resp.Diagnostics.AddError("Error enabling quick deployment", err.Error())
			return
		}
	} else {
		if err := r.client.DisableQuickDeployment(ctx, serverID, siteID); err != nil {
			resp.Diagnostics.AddError("Error disabling quick deployment", err.Error())
			return
		}
	}

	// Update deployment script if provided
	if !plan.Script.IsNull() {
		scriptReq := forge_client.UpdateDeploymentScriptRequest{
			Content:    plan.Script.ValueString(),
			AutoSource: plan.AutoSourceEnv.ValueBool(),
		}
		if err := r.client.UpdateDeploymentScript(ctx, serverID, siteID, scriptReq); err != nil {
			resp.Diagnostics.AddError("Error updating deployment script", err.Error())
			return
		}
	} else {
		// Get existing script if not provided
		script, err := r.client.GetDeploymentScript(ctx, serverID, siteID)
		if err != nil {
			resp.Diagnostics.AddError("Error getting deployment script", err.Error())
			return
		}
		plan.Script = types.StringValue(script)
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *ForgeDeploymentSettingsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state ForgeDeploymentSettingsResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	serverID := int(state.ServerID.ValueInt64())
	siteID := int(state.SiteID.ValueInt64())

	// Get site details to check quick_deploy status
	site, err := r.client.GetSite(ctx, serverID, siteID)
	if err != nil {
		resp.Diagnostics.AddError("Error reading site", err.Error())
		return
	}
	state.QuickDeploy = types.BoolValue(site.QuickDeploy)

	// Get deployment script
	script, err := r.client.GetDeploymentScript(ctx, serverID, siteID)
	if err != nil {
		resp.Diagnostics.AddError("Error getting deployment script", err.Error())
		return
	}
	state.Script = types.StringValue(script)

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *ForgeDeploymentSettingsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan ForgeDeploymentSettingsResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state ForgeDeploymentSettingsResourceModel
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	serverID := int(plan.ServerID.ValueInt64())
	siteID := int(plan.SiteID.ValueInt64())

	// Update quick deployment setting if it has changed
	if !plan.QuickDeploy.Equal(state.QuickDeploy) {
		if plan.QuickDeploy.ValueBool() {
			if err := r.client.EnableQuickDeployment(ctx, serverID, siteID); err != nil {
				resp.Diagnostics.AddError("Error enabling quick deployment", err.Error())
				return
			}
		} else {
			if err := r.client.DisableQuickDeployment(ctx, serverID, siteID); err != nil {
				resp.Diagnostics.AddError("Error disabling quick deployment", err.Error())
				return
			}
		}
	}

	// Update deployment script if it has changed
	if !plan.Script.Equal(state.Script) || !plan.AutoSourceEnv.Equal(state.AutoSourceEnv) {
		scriptReq := forge_client.UpdateDeploymentScriptRequest{
			Content:    plan.Script.ValueString(),
			AutoSource: plan.AutoSourceEnv.ValueBool(),
		}
		if err := r.client.UpdateDeploymentScript(ctx, serverID, siteID, scriptReq); err != nil {
			resp.Diagnostics.AddError("Error updating deployment script", err.Error())
			return
		}
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *ForgeDeploymentSettingsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state ForgeDeploymentSettingsResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// For deployment settings, we don't actually delete anything from the API
	// We just disable quick deployment as a cleanup step
	if state.QuickDeploy.ValueBool() {
		serverID := int(state.ServerID.ValueInt64())
		siteID := int(state.SiteID.ValueInt64())

		if err := r.client.DisableQuickDeployment(ctx, serverID, siteID); err != nil {
			resp.Diagnostics.AddError("Error disabling quick deployment", err.Error())
			return
		}
	}
}

func (r *ForgeDeploymentSettingsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	idParts := strings.Split(req.ID, ":")
	if len(idParts) != 2 {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import identifier with format: server_id:site_id. Got: %q", req.ID),
		)
		return
	}

	serverID, err := strconv.ParseInt(idParts[0], 10, 64)
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid Server ID",
			fmt.Sprintf("Unable to parse server ID: %s. Error: %s", idParts[0], err),
		)
		return
	}

	siteID, err := strconv.ParseInt(idParts[1], 10, 64)
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid Site ID",
			fmt.Sprintf("Unable to parse site ID: %s. Error: %s", idParts[1], err),
		)
		return
	}

	// Create a new state model with the imported values
	var state ForgeDeploymentSettingsResourceModel
	state.ServerID = types.Int64Value(serverID)
	state.SiteID = types.Int64Value(siteID)

	script, err := r.client.GetDeploymentScript(ctx, int(serverID), int(siteID))
	if err != nil {
		resp.Diagnostics.AddError("Error getting deployment script", err.Error())
		return
	}
	state.Script = types.StringValue(script)
	state.AutoSourceEnv = types.BoolValue(false) // Default to false, as we don't know the current state

	site, err := r.client.GetSite(ctx, int(serverID), int(siteID))
	if err != nil {
		resp.Diagnostics.AddError("Error reading site", err.Error())
		return
	}
	state.QuickDeploy = types.BoolValue(site.QuickDeploy)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
