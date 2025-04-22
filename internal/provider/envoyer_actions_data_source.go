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
	Filters []Filter             `tfsdk:"filter"`
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
		MarkdownDescription: "Data source for listing Envoyer actions. Use the `filter` block to filter the actions by specific fields.",
		Blocks: map[string]schema.Block{
			"filter": schema.ListNestedBlock{
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Required:    true,
							Description: "The field name to filter by (e.g., 'name' or 'type')",
						},
						"values": schema.ListAttribute{
							ElementType: types.StringType,
							Required:    true,
							Description: "The list of values to match for the specified field",
						},
					},
				},
				Description: "Filter block for selecting specific credentials.",
			},
		},
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

	d.client = providerConfig.Envoyer
}

func (d *EnvoyerActionsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state EnvoyerActionsDataSourceModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	actions, err := d.client.ListActions(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Error reading actions", err.Error())
		return
	}

	filteredActions := filterEnvoyerActions(actions, state.Filters)

	var actionModels []EnvoyerActionModel
	for _, a := range filteredActions {
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

func filterEnvoyerActions(actions []envoyer_client.Action, filters []Filter) []envoyer_client.Action {
	if len(filters) == 0 {
		return actions
	}

	var filtered []envoyer_client.Action

	for _, c := range actions {
		match := true
		for _, f := range filters {
			switch f.Name.ValueString() {
			case "name":
				if !matchesFilter(c.Name, f.Values) {
					match = false
				}
			case "view":
				if !matchesFilter(c.View, f.Values) {
					match = false
				}
			case "sequence":
				if !matchesFilter(fmt.Sprintf("%d", c.Sequence), f.Values) {
					match = false
				}
			case "version":
				if !matchesFilter(fmt.Sprintf("%d", c.Version), f.Values) {
					match = false
				}
			default:
				// Ignore unknown filters
				match = false
			}
		}

		if match {
			filtered = append(filtered, c)
		}
	}

	return filtered
}
