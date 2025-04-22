package provider

import (
	"context"
	"fmt"
	"terraform-provider-laravel/internal/envoyer_client"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure DeploymentResource satisfies the expected interfaces.
var _ resource.Resource = &EnvoyerDeploymentResource{}

// EnvoyerDeploymentResourceModel defines the schema data model for the server.
type EnvoyerDeploymentResourceModel struct {
	ID        types.Int64  `tfsdk:"id"`
	ProjectID types.Int64  `tfsdk:"project_id"`
	UserID    types.Int64  `tfsdk:"user_id"`
	From      types.String `tfsdk:"from"`   // branch or tag
	Branch    types.String `tfsdk:"branch"` // for branch
	Tag       types.String `tfsdk:"tag"`    // for tag
	CreatedAt types.String `tfsdk:"created_at"`
	UpdatedAt types.String `tfsdk:"updated_at"`
	Status    types.String `tfsdk:"status"`
}

// EnvoyerDeploymentResource implements the resource.Resource interface.
type EnvoyerDeploymentResource struct {
	client *envoyer_client.Client
}

// NewEnvoyerDeploymentResource is a helper function to instantiate the resource.
func NewEnvoyerDeploymentResource() resource.Resource {
	return &EnvoyerDeploymentResource{}
}

func (r *EnvoyerDeploymentResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_envoyer_deployment"
}

func (r *EnvoyerDeploymentResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Envoyer deployment resource. This resource allows you to create deployments in Envoyer.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Computed: true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
				MarkdownDescription: "The ID of the deployment.",
			},
			"user_id": schema.Int64Attribute{
				Computed: true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
				MarkdownDescription: "The ID of the user who initiated the deployment.",
			},
			"project_id": schema.Int64Attribute{
				Required:            true,
				MarkdownDescription: "The ID of the Envoyer project.",
			},
			"from": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The source of the deployment. Can be `branch` or `tag`. This determines how the deployment is created.",
			},
			"branch": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The branch to deploy. This is only used if `from` is set to `branch`.",
			},
			"tag": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The tag to deploy. This is only used if `from` is set to `tag`.",
			},
			"created_at": schema.StringAttribute{
				MarkdownDescription: "When the deployment was created.",
				Computed:            true,
			},
			"updated_at": schema.StringAttribute{
				MarkdownDescription: "When the deployment was last updated.",
				Computed:            true,
			},
			"status": schema.StringAttribute{
				MarkdownDescription: "The status of the deployment.",
				Computed:            true,
			},
		},
	}
}

func (r *EnvoyerDeploymentResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *EnvoyerDeploymentResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan EnvoyerDeploymentResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if plan.From.ValueString() == "branch" && plan.Branch.IsNull() {
		resp.Diagnostics.AddError("Invalid branch", "Branch must be specified when from is set to 'branch'")
		return
	} else if plan.From.ValueString() == "tag" && plan.Tag.IsNull() {
		resp.Diagnostics.AddError("Invalid tag", "Tag must be specified when from is set to 'tag'")
		return
	}

	reqPayload := envoyer_client.CreateProjectDeploymentRequest{
		From:   plan.From.ValueString(),
		Branch: plan.Branch.ValueStringPointer(),
		Tag:    plan.Tag.ValueStringPointer(),
	}

	err := r.client.CreateProjectDeployment(ctx, int(plan.ProjectID.ValueInt64()), reqPayload)
	if err != nil {
		resp.Diagnostics.AddError("Error creating deployment", fmt.Sprintf("Could not create deployment: %s", err.Error()))
		return
	}

	deployments, err := r.client.ListProjectDeployments(ctx, int(plan.ProjectID.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError("Error listing deployments", fmt.Sprintf("Could not list deployments: %s", err.Error()))
		return
	}
	if len(deployments) == 0 {
		resp.Diagnostics.AddError("Error listing deployments", "No deployments found")
		return
	}
	created := deployments[0]
	// Wait for the deployment to finish.

	finished, err := r.client.WaitFoirDeploymentToFinish(ctx, int(plan.ProjectID.ValueInt64()), int(created.ID))
	if err != nil {
		resp.Diagnostics.AddError("Error waiting for deployment", fmt.Sprintf("Could not wait for deployment to finish: %s", err.Error()))
		return
	}
	if finished == nil {
		resp.Diagnostics.AddError("Error waiting for deployment", "Deployment finished with no response")
		return
	}

	// Map the response back to the resource state.
	plan.ID = types.Int64Value(finished.ID)
	plan.UserID = types.Int64Value(finished.UserID)
	plan.CreatedAt = types.StringValue(finished.CreatedAt.Format(time.RFC3339))
	plan.UpdatedAt = types.StringValue(finished.UpdatedAt.Format(time.RFC3339))
	plan.Status = types.StringValue(finished.Status)
	plan.ProjectID = types.Int64Value(finished.ProjectID)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *EnvoyerDeploymentResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state EnvoyerDeploymentResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Fetch the deployment from the API.
	deployment, err := r.client.GetProjectDeployment(ctx, int(state.ProjectID.ValueInt64()), int(state.ID.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError("Error reading deployment", fmt.Sprintf("Could not read deployment: %s", err.Error()))
		return
	}

	finished, err := r.client.WaitFoirDeploymentToFinish(ctx, int(state.ProjectID.ValueInt64()), int(state.ID.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError("Error waiting for deployment", fmt.Sprintf("Could not wait for deployment to finish: %s", err.Error()))
		return
	}
	if finished == nil {
		resp.Diagnostics.AddError("Error waiting for deployment", "Deployment finished with no response")
		return
	}

	// Map the response back to the resource state.
	state.ID = types.Int64Value(deployment.ID)
	state.UserID = types.Int64Value(finished.UserID)
	state.CreatedAt = types.StringValue(finished.CreatedAt.Format(time.RFC3339))
	state.UpdatedAt = types.StringValue(finished.UpdatedAt.Format(time.RFC3339))
	state.Status = types.StringValue(finished.Status)
	state.ProjectID = types.Int64Value(finished.ProjectID)
	if state.From.ValueString() == "branch" {
		state.Branch = types.StringValue(finished.CommitBranch)
	}

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *EnvoyerDeploymentResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan EnvoyerDeploymentResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *EnvoyerDeploymentResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state EnvoyerDeploymentResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}
