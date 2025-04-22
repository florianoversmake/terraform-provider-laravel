package provider

import (
	"context"
	"fmt"
	"strconv"

	"terraform-provider-laravel/internal/forge_client"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ resource.Resource = &ForgeRecipeResource{}
var _ resource.ResourceWithImportState = &ForgeRecipeResource{}

// ForgeRecipeResource implements a Terraform resource for a Forge site.
type ForgeRecipeResource struct {
	client *forge_client.Client
}

// Resource model.
type ForgeRecipeResourceModel struct {
	ID        types.Int64  `tfsdk:"id"`
	Name      types.String `tfsdk:"name"`
	User      types.String `tfsdk:"user"`
	Script    types.String `tfsdk:"script"`
	CreatedAt types.String `tfsdk:"created_at"`
}

func NewForgeRecipeResource() resource.Resource {
	return &ForgeRecipeResource{}
}

func (r *ForgeRecipeResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_forge_recipe"
}

func (r *ForgeRecipeResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Forge recipe resource. This resource allows you to manage custom recipes in Forge.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Computed: true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The name of the recipe.",
			},
			"user": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The username to run the recipe as.",
			},
			"script": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The script content of the recipe.",
			},
			"created_at": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

func (r *ForgeRecipeResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *ForgeRecipeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan ForgeRecipeResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Build the CreateRecipeRequest payload.
	payload := forge_client.CreateRecipeRequest{
		Name:   plan.Name.ValueString(),
		User:   plan.User.ValueString(),
		Script: plan.Script.ValueString(),
	}

	// Call CreateRecipe on the client.
	recipe, err := r.client.CreateRecipe(ctx, payload)
	if err != nil {
		resp.Diagnostics.AddError("Error creating recipe", err.Error())
		return
	}

	// Update plan state with response values.
	plan.ID = types.Int64Value(recipe.ID)
	plan.CreatedAt = types.StringValue(recipe.CreatedAt)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *ForgeRecipeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state ForgeRecipeResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	recipe, err := r.client.GetRecipe(ctx, int(state.ID.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError("Error reading recipe", err.Error())
		return
	}

	state.Name = types.StringValue(recipe.Name)
	state.User = types.StringValue(recipe.User)
	state.Script = types.StringValue(recipe.Script)
	state.CreatedAt = types.StringValue(recipe.CreatedAt)
	state.ID = types.Int64Value(recipe.ID)

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *ForgeRecipeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan ForgeRecipeResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Build the CreateRecipeRequest payload.
	payload := forge_client.CreateRecipeRequest{
		Name:   plan.Name.ValueString(),
		User:   plan.User.ValueString(),
		Script: plan.Script.ValueString(),
	}

	// Call UpdateRecipe on the client.
	recipe, err := r.client.UpdateRecipe(ctx, int(plan.ID.ValueInt64()), payload)
	if err != nil {
		resp.Diagnostics.AddError("Error updating recipe", err.Error())
		return
	}
	// Update plan state with response values.
	plan.ID = types.Int64Value(recipe.ID)
	plan.CreatedAt = types.StringValue(recipe.CreatedAt)
	plan.Name = types.StringValue(recipe.Name)
	plan.User = types.StringValue(recipe.User)
	plan.Script = types.StringValue(recipe.Script)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *ForgeRecipeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state ForgeRecipeResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteRecipe(ctx, int(state.ID.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError("Error deleting recipe", err.Error())
		return
	}
}

func (r *ForgeRecipeResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	id, err := strconv.ParseInt(req.ID, 10, 64)
	if err != nil {
		resp.Diagnostics.AddError("Error parsing ID", err.Error())
		return
	}

	recipe, err := r.client.GetRecipe(ctx, int(id))
	if err != nil {
		resp.Diagnostics.AddError("Error reading recipe", err.Error())
		return
	}

	var stateModel ForgeRecipeResourceModel
	stateModel.ID = types.Int64Value(recipe.ID)
	stateModel.Name = types.StringValue(recipe.Name)
	stateModel.User = types.StringValue(recipe.User)
	stateModel.Script = types.StringValue(recipe.Script)
	stateModel.CreatedAt = types.StringValue(recipe.CreatedAt)

	diags := resp.State.Set(ctx, stateModel)
	resp.Diagnostics.Append(diags...)
}
