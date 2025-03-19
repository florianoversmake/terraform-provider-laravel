// Copyright (c) HashiCorp, Inc.

package provider

import (
	"context"
	"fmt"
	"strconv"

	"terraform-provider-laravel/internal/envoyer_client"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure EnvironmentResource implements required interfaces.
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
	Servers   []types.Int64 `tfsdk:"servers"` // Optional list of server IDs
}

func (r *EnvoyerEnvironmentResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_envoyer_environment"
}

func (r *EnvoyerEnvironmentResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"project_id": schema.Int64Attribute{
				Required: true,
			},
			"contents": schema.StringAttribute{
				Optional: true,
				Computed: true,
			},
			"servers": schema.ListAttribute{
				ElementType: types.Int64Type,
				Optional:    true,
			},
		},
	}
}

func (r *EnvoyerEnvironmentResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	providerConfig, ok := req.ProviderData.(*providerConfig)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Provider Configure Type",
			fmt.Sprintf("Expected *providerConfig, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	r.client = providerConfig.Envoyer
}

func (r *EnvoyerEnvironmentResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan EnvoyerEnvironmentResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var serverIds []int64

	for _, serverId := range plan.Servers {
		serverIds = append(serverIds, serverId.ValueInt64())
	}

	payload := envoyer_client.UpdateEnvironmentRequest{
		Contents: plan.Contents.ValueString(),
		Servers:  serverIds,
	}

	contents, err := r.client.UpdateEnvironment(ctx, int(plan.ProjectID.ValueInt64()), payload)
	if err != nil {
		resp.Diagnostics.AddError("Error creating/updating environment", err.Error())
		return
	}

	plan.Contents = types.StringValue(contents)
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

	contents, err := r.client.GetEnvironment(ctx, int(state.ProjectID.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError("Error reading environment", err.Error())
		return
	}

	servers, err := r.client.GetEnvironmentServers(ctx, int(state.ProjectID.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError("Error reading environment servers", err.Error())
		return
	}

	state.Contents = types.StringValue(contents)

	serverIds := make([]types.Int64, 0, len(servers))
	for _, server := range servers {
		serverIds = append(serverIds, types.Int64Value(server))
	}
	state.Servers = serverIds

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *EnvoyerEnvironmentResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan EnvoyerEnvironmentResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var serverIds []int64

	for _, serverId := range plan.Servers {
		serverIds = append(serverIds, serverId.ValueInt64())
	}

	payload := envoyer_client.UpdateEnvironmentRequest{
		Contents: plan.Contents.ValueString(),
		Servers:  serverIds,
	}

	contents, err := r.client.UpdateEnvironment(ctx, int(plan.ProjectID.ValueInt64()), payload)
	if err != nil {
		resp.Diagnostics.AddError("Error updating environment", err.Error())
		return
	}
	plan.Contents = types.StringValue(contents)
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

	resp.State.RemoveResource(ctx)
}

func (r *EnvoyerEnvironmentResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	projectId, err := strconv.ParseInt(req.ID, 10, 64)
	if err != nil {
		resp.Diagnostics.AddError("Invalid project ID", err.Error())
		return
	}

	environment, err := r.client.GetEnvironment(ctx, int(projectId))

	if err != nil {
		resp.Diagnostics.AddError("Error reading environment", err.Error())
		return
	}

	serverIds, err := r.client.GetEnvironmentServers(ctx, int(projectId))
	if err != nil {
		resp.Diagnostics.AddError("Error reading environment servers", err.Error())
		return
	}

	state := EnvoyerEnvironmentResourceModel{
		ProjectID: types.Int64Value(projectId),
		Contents:  types.StringValue(environment),
	}

	serverIdsList := make([]types.Int64, 0, len(serverIds))
	for _, serverId := range serverIds {
		serverIdsList = append(serverIdsList, types.Int64Value(serverId))
	}
	state.Servers = serverIdsList

	diags := resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}
