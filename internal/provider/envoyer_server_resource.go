// Copyright (c) HashiCorp, Inc.

package provider

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"terraform-provider-laravel/internal/envoyer_client"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure ServerResource satisfies the expected interfaces.
var _ resource.Resource = &EnvoyerServerResource{}
var _ resource.ResourceWithImportState = &EnvoyerServerResource{}

// EnvoyerServerResourceModel defines the schema data model for the server.
type EnvoyerServerResourceModel struct {
	ID                      types.Int64  `tfsdk:"id"`
	ProjectID               types.Int64  `tfsdk:"project_id"`
	Name                    types.String `tfsdk:"name"`
	ConnectAs               types.String `tfsdk:"connect_as"`
	IPAddress               types.String `tfsdk:"ip_address"`
	Port                    types.String `tfsdk:"port"`
	PHPVersion              types.String `tfsdk:"php_version"`
	ReceivesCodeDeployments types.Bool   `tfsdk:"receives_code_deployments"`
	ShouldRestartFPM        types.Bool   `tfsdk:"should_restart_fpm"`
	DeploymentPath          types.String `tfsdk:"deployment_path"`
	PHPPath                 types.String `tfsdk:"php_path"`
	ComposerPath            types.String `tfsdk:"composer_path"`
	PublicKey               types.String `tfsdk:"public_key"`
}

// EnvoyerServerResource implements the resource.Resource interface.
type EnvoyerServerResource struct {
	client *envoyer_client.Client
}

// NewEnvoyerServerResource is a helper function to instantiate the resource.
func NewEnvoyerServerResource() resource.Resource {
	return &EnvoyerServerResource{}
}

func (r *EnvoyerServerResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_envoyer_server"
}

func (r *EnvoyerServerResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Computed: true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"project_id": schema.Int64Attribute{
				Required: true,
			},
			"name": schema.StringAttribute{
				Required: true,
			},
			"connect_as": schema.StringAttribute{
				Required: true,
			},
			"ip_address": schema.StringAttribute{
				Required: true,
			},
			"port": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString("22"),
			},
			"php_version": schema.StringAttribute{
				Required: true,
			},
			"receives_code_deployments": schema.BoolAttribute{
				Optional: true,
				Computed: true,
				Default:  booldefault.StaticBool(true),
			},
			"should_restart_fpm": schema.BoolAttribute{
				Optional: true,
				Computed: true,
				Default:  booldefault.StaticBool(true),
			},
			"deployment_path": schema.StringAttribute{
				Required: true,
			},
			"php_path": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString("php"),
			},
			"composer_path": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString("composer"),
			},
			"public_key": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *EnvoyerServerResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *EnvoyerServerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan EnvoyerServerResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	port, err := strconv.Atoi(plan.Port.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid port", "Port must be a valid integer")
		return
	}

	// Build the client request payload. You may need to adapt this
	// depending on how you want to derive the host or other properties.
	reqPayload := envoyer_client.CreateServerRequest{
		Name:                plan.Name.ValueString(),
		ConnectAs:           plan.ConnectAs.ValueString(),
		Host:                plan.IPAddress.ValueString(), // Derive from your configuration if needed.
		Port:                port,                         // Use default port or plan value.
		PHPVersion:          plan.PHPVersion.ValueString(),
		ReceivesCodeDeploys: plan.ReceivesCodeDeployments.ValueBool(),
		DeploymentPath:      plan.DeploymentPath.ValueString(),
		RestartFpm:          plan.ShouldRestartFPM.ValueBool(),
		ComposerPath:        plan.ComposerPath.ValueString(),
		PHPPath:             plan.PHPPath.ValueString(),
	}

	created, err := r.client.CreateServer(ctx, int(plan.ProjectID.ValueInt64()), reqPayload)
	if err != nil {
		resp.Diagnostics.AddError("Error creating server", fmt.Sprintf("Could not create server: %s", err.Error()))
		return
	}

	// Map the response back to the resource state.
	plan.ID = types.Int64Value(created.ID)
	plan.IPAddress = types.StringValue(created.IPAddress)
	plan.Port = types.StringValue(created.Port)
	plan.PublicKey = types.StringValue(created.PublicKey)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *EnvoyerServerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state EnvoyerServerResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	server, err := r.client.GetServer(ctx, int(state.ProjectID.ValueInt64()), int(state.ID.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError("Error reading server", fmt.Sprintf("Could not read server: %s", err.Error()))
		return
	}

	state.Name = types.StringValue(server.Name)
	state.ConnectAs = types.StringValue(server.ConnectAs)
	state.IPAddress = types.StringValue(server.IPAddress)
	state.Port = types.StringValue(server.Port)
	state.PHPVersion = types.StringValue(server.PHPVersion)
	state.ReceivesCodeDeployments = types.BoolValue(server.ReceivesCodeDeploys)
	state.ShouldRestartFPM = types.BoolValue(server.ShouldRestartFPM)
	state.DeploymentPath = types.StringValue(server.DeploymentPath)
	state.PHPPath = types.StringValue(server.PHPPath)
	state.ComposerPath = types.StringValue(server.ComposerPath)
	state.PublicKey = types.StringValue(server.PublicKey)

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *EnvoyerServerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan EnvoyerServerResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	port, err := strconv.Atoi(plan.Port.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid port", "Port must be a valid integer")
		return
	}

	reqPayload := envoyer_client.CreateServerRequest{
		Name:                plan.Name.ValueString(),
		ConnectAs:           plan.ConnectAs.ValueString(),
		Host:                plan.IPAddress.ValueString(),
		Port:                port,
		PHPVersion:          plan.PHPVersion.ValueString(),
		ReceivesCodeDeploys: plan.ReceivesCodeDeployments.ValueBool(),
		DeploymentPath:      plan.DeploymentPath.ValueString(),
		RestartFpm:          plan.ShouldRestartFPM.ValueBool(),
		ComposerPath:        plan.ComposerPath.ValueString(),
		PHPPath:             plan.PHPPath.ValueString(),
	}

	updated, err := r.client.UpdateServer(ctx, int(plan.ProjectID.ValueInt64()), int(plan.ID.ValueInt64()), reqPayload)
	if err != nil {
		resp.Diagnostics.AddError("Error updating server", fmt.Sprintf("Could not update server: %s", err.Error()))
		return
	}

	plan.Name = types.StringValue(updated.Name)
	// (Map the remaining updated fields in a similar fashion.)
	plan.IPAddress = types.StringValue(updated.IPAddress)
	plan.Port = types.StringValue(updated.Port)
	plan.PublicKey = types.StringValue(updated.PublicKey)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *EnvoyerServerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state EnvoyerServerResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.DeleteServer(ctx, int(state.ProjectID.ValueInt64()), int(state.ID.ValueInt64())); err != nil {
		resp.Diagnostics.AddError("Error deleting server", fmt.Sprintf("Could not delete server: %s", err.Error()))
		return
	}
}

func (r *EnvoyerServerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	projectAndServerID := req.ID
	parts := strings.Split(projectAndServerID, "/")
	if len(parts) != 2 {
		resp.Diagnostics.AddError("Invalid import ID", "Expected format <project_id>/<server_id>")
		return
	}

	projectID, err := strconv.Atoi(parts[0])
	if err != nil {
		resp.Diagnostics.AddError("Invalid project ID", "Could not parse project ID")
		return
	}

	serverID, err := strconv.Atoi(parts[1])
	if err != nil {
		resp.Diagnostics.AddError("Invalid server ID", "Could not parse server ID")
		return
	}

	server, err := r.client.GetServer(ctx, projectID, serverID)
	if err != nil {
		resp.Diagnostics.AddError("Error importing server", fmt.Sprintf("Could not import server: %s", err.Error()))
		return
	}

	state := EnvoyerServerResourceModel{
		ID:                      types.Int64Value(server.ID),
		ProjectID:               types.Int64Value(server.ProjectID),
		Name:                    types.StringValue(server.Name),
		ConnectAs:               types.StringValue(server.ConnectAs),
		IPAddress:               types.StringValue(server.IPAddress),
		Port:                    types.StringValue(server.Port),
		PHPVersion:              types.StringValue(server.PHPVersion),
		ReceivesCodeDeployments: types.BoolValue(server.ReceivesCodeDeploys),
		ShouldRestartFPM:        types.BoolValue(server.ShouldRestartFPM),
		DeploymentPath:          types.StringValue(server.DeploymentPath),
		PHPPath:                 types.StringValue(server.PHPPath),
		ComposerPath:            types.StringValue(server.ComposerPath),
		PublicKey:               types.StringValue(server.PublicKey),
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
