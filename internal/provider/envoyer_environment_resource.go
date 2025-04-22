package provider

import (
	"context"
	"fmt"
	"strconv"
	"terraform-provider-laravel/internal/envoyer_client"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var _ resource.Resource = &EnvoyerEnvironmentResource{}
var _ resource.ResourceWithImportState = &EnvoyerEnvironmentResource{}

func NewEnvoyerEnvironmentResource() resource.Resource {
	return &EnvoyerEnvironmentResource{}
}

type EnvoyerEnvironmentResource struct {
	client *envoyer_client.Client
}

type EnvoyerEnvironmentResourceModel struct {
	ProjectID types.Int64   `tfsdk:"project_id"`
	Contents  types.String  `tfsdk:"contents"`
	Servers   []types.Int64 `tfsdk:"servers"`
}

func (r *EnvoyerEnvironmentResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_envoyer_environment"
}

func (r *EnvoyerEnvironmentResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Envoyer environment resource. This resource allows you to manage environment variables in Envoyer.",
		Attributes: map[string]schema.Attribute{
			"project_id": schema.Int64Attribute{
				Required:    true,
				Description: "The ID of the Envoyer project.",
			},
			"contents": schema.StringAttribute{
				Required:    true,
				Description: "The raw environment variable contents.",
			},
			"servers": schema.ListAttribute{
				ElementType: types.Int64Type,
				Optional:    true,
				Description: "List of server IDs that should receive these environment variables.",
			},
		},
	}
}

func (r *EnvoyerEnvironmentResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	if providerConfig.Envoyer == nil {
		resp.Diagnostics.AddError(
			"Envoyer Client Not Configured",
			"This resource requires the Envoyer API token to be configured in the provider. "+
				"Please set the 'envoyer_api_token' attribute in the provider configuration.",
		)
		return
	}

	r.client = providerConfig.Envoyer
}

// validateConfig ensures that only one of contents or env_var blocks is used.
func (r *EnvoyerEnvironmentResource) validateConfig(plan EnvoyerEnvironmentResourceModel) diag.Diagnostics {
	var diags diag.Diagnostics

	hasContents := !plan.Contents.IsNull() && plan.Contents.ValueString() != ""

	if !hasContents {
		diags.AddError(
			"Missing Configuration",
			"'contents' must be provided.",
		)
	}

	return diags
}

func (r *EnvoyerEnvironmentResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan EnvoyerEnvironmentResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Validate configuration
	resp.Diagnostics.Append(r.validateConfig(plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var serverIds []int64
	for _, serverId := range plan.Servers {
		serverIds = append(serverIds, serverId.ValueInt64())
	}

	// Update environment
	payload := envoyer_client.UpdateEnvironmentRequest{
		Contents: plan.Contents.ValueString(),
		Servers:  serverIds,
	}

	tflog.Debug(ctx, "Creating/updating Envoyer environment", map[string]interface{}{
		"project_id": plan.ProjectID.ValueInt64(),
		"contents":   plan.Contents.ValueString(),
	})

	_, err := r.client.UpdateEnvironment(ctx, int(plan.ProjectID.ValueInt64()), payload)
	if err != nil {
		resp.Diagnostics.AddError("Error creating/updating environment", err.Error())
		return
	}

	tflog.Debug(ctx, "Created/updated Envoyer environment", map[string]interface{}{
		"project_id": plan.ProjectID.ValueInt64(),
		"contents":   plan.Contents.ValueString(),
	})
	// Update plan with the returned content
	plan.Contents = types.StringValue(plan.Contents.ValueString())

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *EnvoyerEnvironmentResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state EnvoyerEnvironmentResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get environment content from the API
	env, err := r.client.GetEnvironment(ctx, int(state.ProjectID.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError("Error reading environment", err.Error())
		return
	}

	// Get server IDs
	servers, err := r.client.GetEnvironmentServers(ctx, int(state.ProjectID.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError("Error reading environment servers", err.Error())
		return
	}

	// Update server IDs in state
	serverIds := make([]types.Int64, 0, len(servers))
	for _, server := range servers {
		serverIds = append(serverIds, types.Int64Value(server))
	}
	state.Servers = serverIds
	// Update contents in state
	state.Contents = types.StringValue(env)

	// Set the state
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *EnvoyerEnvironmentResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan EnvoyerEnvironmentResourceModel
	// var state EnvoyerEnvironmentResourceModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	// diags = req.State.Get(ctx, &state)
	// resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Validate configuration
	resp.Diagnostics.Append(r.validateConfig(plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Extract server IDs
	var serverIds []int64
	for _, serverId := range plan.Servers {
		serverIds = append(serverIds, serverId.ValueInt64())
	}

	// Update environment
	payload := envoyer_client.UpdateEnvironmentRequest{
		Contents: plan.Contents.ValueString(),
		Servers:  serverIds,
	}

	tflog.Debug(ctx, "Updating Envoyer environment", map[string]interface{}{
		"project_id": plan.ProjectID.ValueInt64(),
		"contents":   plan.Contents.ValueString(),
	})
	_, err := r.client.UpdateEnvironment(ctx, int(plan.ProjectID.ValueInt64()), payload)
	if err != nil {
		resp.Diagnostics.AddError("Error updating environment", err.Error())
		return
	}
	tflog.Debug(ctx, "Updated Envoyer environment", map[string]interface{}{
		"project_id": plan.ProjectID.ValueInt64(),
		"contents":   plan.Contents.ValueString(),
	})
	// Update plan with the returned content
	plan.Contents = types.StringValue(plan.Contents.ValueString())
	// Update the state
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *EnvoyerEnvironmentResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state EnvoyerEnvironmentResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// For environment resources, we don't actually delete anything from the API
	// We just remove the resource from Terraform state
	resp.State.RemoveResource(ctx)
}

func (r *EnvoyerEnvironmentResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	projectId, err := strconv.ParseInt(req.ID, 10, 64)
	if err != nil {
		resp.Diagnostics.AddError("Invalid project ID", err.Error())
		return
	}

	// Set ID and let Read method handle the rest
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("project_id"), projectId)...)

	// Read the state to populate the rest of the attributes
	var state EnvoyerEnvironmentResourceModel
	diags := resp.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	// Read the environment from the API
	env, err := r.client.GetEnvironment(ctx, int(state.ProjectID.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError("Error reading environment", err.Error())
		return
	}
	// Get server IDs
	servers, err := r.client.GetEnvironmentServers(ctx, int(state.ProjectID.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError("Error reading environment servers", err.Error())
		return
	}
	// Update server IDs in state
	serverIds := make([]types.Int64, 0, len(servers))
	for _, server := range servers {
		serverIds = append(serverIds, types.Int64Value(server))
	}
	state.Servers = serverIds
	// Update contents in state
	state.Contents = types.StringValue(env)
	// Set the state
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}
