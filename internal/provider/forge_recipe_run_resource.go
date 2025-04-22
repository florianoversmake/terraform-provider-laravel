package provider

import (
	"context"
	"fmt"

	"terraform-provider-laravel/internal/forge_client"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ resource.Resource = &ForgeRecipeRunResource{}

// ForgeRecipeRunResource implements a Terraform resource for a Forge site.
type ForgeRecipeRunResource struct {
	client *forge_client.Client
}

// Resource model.
type ForgeRecipeRunResourceModel struct {
	RecipeID types.Int64   `tfsdk:"recipe_id"`
	Servers  []types.Int64 `tfsdk:"servers"`
	Notify   types.Bool    `tfsdk:"notify"`
}

func NewForgeRecipeRunResource() resource.Resource {
	return &ForgeRecipeRunResource{}
}

func (r *ForgeRecipeRunResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_forge_recipe_run"
}

func (r *ForgeRecipeRunResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Forge recipe run resource. This resource allows you to run a custom recipe on a server in Forge.",
		Attributes: map[string]schema.Attribute{
			"recipe_id": schema.Int64Attribute{
				Required:            true,
				MarkdownDescription: "The ID of the recipe to run.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"servers": schema.ListAttribute{
				ElementType:         types.Int64Type,
				Required:            true,
				MarkdownDescription: "List of server IDs to run the recipe on.",
				PlanModifiers: []planmodifier.List{
					listplanmodifier.RequiresReplace(),
				},
			},
			"notify": schema.BoolAttribute{
				Optional: true,
				Computed: true,
				Default:  booldefault.StaticBool(false),
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.RequiresReplace(),
				},
			},
		},
	}
}

func (r *ForgeRecipeRunResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *ForgeRecipeRunResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan ForgeRecipeRunResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var serverIds []int64
	for _, serverId := range plan.Servers {
		serverIds = append(serverIds, serverId.ValueInt64())
	}

	// Build the RunRecipeRequest payload.
	payload := forge_client.RunRecipeRequest{
		Servers: serverIds,
		Notify:  plan.Notify.ValueBool(),
	}

	// Call CreateRecipe on the client.
	err := r.client.RunRecipe(ctx, int(plan.RecipeID.ValueInt64()), payload)
	if err != nil {
		resp.Diagnostics.AddError("Error creating recipe", err.Error())
		return
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *ForgeRecipeRunResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state ForgeRecipeRunResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *ForgeRecipeRunResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan ForgeRecipeRunResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *ForgeRecipeRunResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state ForgeRecipeRunResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
