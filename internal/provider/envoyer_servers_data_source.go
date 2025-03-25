// Copyright (c) HashiCorp, Inc.

package provider

import (
	"context"
	"fmt"
	"terraform-provider-laravel/internal/envoyer_client"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure ServersDataSource satisfies the DataSource interface.
var _ datasource.DataSource = &EnvoyerServersDataSource{}

func NewEnvoyerServersDataSource() datasource.DataSource {
	return &EnvoyerServersDataSource{}
}

type EnvoyerServersDataSource struct {
	client *envoyer_client.Client
}

type EnvoyerServersDataSourceModel struct {
	ProjectID types.Int64          `tfsdk:"project_id"`
	Servers   []EnvoyerServerModel `tfsdk:"servers"`
}

type EnvoyerServerModel struct {
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

func (d *EnvoyerServersDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_envoyer_servers"
}

func (d *EnvoyerServersDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"project_id": schema.Int64Attribute{
				Description: "The ID of the project to list servers for",
				Required:    true,
			},
			"servers": schema.ListNestedAttribute{
				Description: "List of servers",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id":                        schema.Int64Attribute{Computed: true},
						"project_id":                schema.Int64Attribute{Computed: true},
						"name":                      schema.StringAttribute{Computed: true},
						"connect_as":                schema.StringAttribute{Computed: true},
						"ip_address":                schema.StringAttribute{Computed: true},
						"port":                      schema.StringAttribute{Computed: true},
						"php_version":               schema.StringAttribute{Computed: true},
						"receives_code_deployments": schema.BoolAttribute{Computed: true},
						"should_restart_fpm":        schema.BoolAttribute{Computed: true},
						"deployment_path":           schema.StringAttribute{Computed: true},
						"php_path":                  schema.StringAttribute{Computed: true},
						"composer_path":             schema.StringAttribute{Computed: true},
						"public_key":                schema.StringAttribute{Computed: true},
					},
				},
			},
		},
	}
}

func (d *EnvoyerServersDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	d.client = providerConfig.Envoyer
}

func (d *EnvoyerServersDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state EnvoyerServersDataSourceModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Call the client's ListServers method using the given project ID.
	servers, err := d.client.ListServers(ctx, int(state.ProjectID.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError("Error reading servers", err.Error())
		return
	}

	var serverModels []EnvoyerServerModel
	for _, s := range servers {
		serverModels = append(serverModels, EnvoyerServerModel{
			ID:                      types.Int64Value(s.ID),
			ProjectID:               types.Int64Value(s.ProjectID),
			Name:                    types.StringValue(s.Name),
			ConnectAs:               types.StringValue(s.ConnectAs),
			IPAddress:               types.StringValue(s.IPAddress),
			Port:                    types.StringValue(s.Port),
			PHPVersion:              types.StringValue(s.PHPVersion),
			ReceivesCodeDeployments: types.BoolValue(s.ReceivesCodeDeploys),
			ShouldRestartFPM:        types.BoolValue(s.ShouldRestartFPM),
			DeploymentPath:          types.StringValue(s.DeploymentPath),
			PHPPath:                 types.StringValue(s.PHPPath),
			ComposerPath:            types.StringValue(s.ComposerPath),
			PublicKey:               types.StringValue(s.PublicKey),
		})
	}
	state.Servers = serverModels

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}
