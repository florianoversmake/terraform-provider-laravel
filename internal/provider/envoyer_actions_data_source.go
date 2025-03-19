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

var _ datasource.DataSource = &EnvoyerActionsDataSource{}

func NewEnvoyerActionsDataSource() datasource.DataSource {
	return &EnvoyerActionsDataSource{}
}

type EnvoyerActionsDataSource struct {
	client *envoyer_client.Client
}

type EnvoyerActionsDataSourceModel struct {
	Actions []EnvoyerActionModel `tfsdk:"actions"`
}

type EnvoyerActionModel struct {
	ID        types.Int64  `tfsdk:"id"`
	Version   types.Int64  `tfsdk:"version"`
	Name      types.String `tfsdk:"name"`
	View      types.String `tfsdk:"view"`
	Sequence  types.Int64  `tfsdk:"sequence"`
	CreatedAt types.String `tfsdk:"created_at"`
	UpdatedAt types.String `tfsdk:"updated_at"`
}

func (d *EnvoyerActionsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_envoyer_actions"
}

func (d *EnvoyerActionsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"actions": schema.ListNestedAttribute{
				Computed:    true,
				Description: "List of actions available in Envoyer",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id":         schema.Int64Attribute{Computed: true},
						"version":    schema.Int64Attribute{Computed: true},
						"name":       schema.StringAttribute{Computed: true},
						"view":       schema.StringAttribute{Computed: true},
						"sequence":   schema.Int64Attribute{Computed: true},
						"created_at": schema.StringAttribute{Computed: true},
						"updated_at": schema.StringAttribute{Computed: true},
					},
				},
			},
		},
	}
}

func (d *EnvoyerActionsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *EnvoyerActionsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state EnvoyerActionsDataSourceModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Assume a client.ListActions method is implemented.
	actions, err := d.client.ListActions(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Error reading actions", err.Error())
		return
	}

	var actionModels []EnvoyerActionModel
	for _, a := range actions {
		actionModels = append(actionModels, EnvoyerActionModel{
			ID:        types.Int64Value(a.ID),
			Version:   types.Int64Value(a.Version),
			Name:      types.StringValue(a.Name),
			View:      types.StringValue(a.View),
			Sequence:  types.Int64Value(a.Sequence),
			CreatedAt: types.StringValue(a.CreatedAt),
			UpdatedAt: types.StringValue(a.UpdatedAt),
		})
	}
	state.Actions = actionModels

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}
