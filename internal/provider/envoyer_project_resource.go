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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure ProjectResource implements interfaces.
var (
	_ resource.Resource                = &EnvoyerProjectResource{}
	_ resource.ResourceWithImportState = &EnvoyerProjectResource{}
)

// NewEnvoyerProjectResource is a helper for Terraform to instantiate our resource.
func NewEnvoyerProjectResource() resource.Resource {
	return &EnvoyerProjectResource{}
}

// EnvoyerProjectResource implements the resource logic for an Envoyer project.
type EnvoyerProjectResource struct {
	client *envoyer_client.Client
}

// EnvoyerProjectResourceModel is the Terraform schema for envoyer_project.
type EnvoyerProjectResourceModel struct {
	ID                types.Int64  `tfsdk:"id"`
	Name              types.String `tfsdk:"name"`
	RepoProvider      types.String `tfsdk:"repo_provider"` // Renamed from Provider to avoid conflict.
	Repository        types.String `tfsdk:"repository"`
	Type              types.String `tfsdk:"type"`
	Branch            types.String `tfsdk:"branch"`
	RetainDeployments types.Int64  `tfsdk:"retain_deployments"`
	Monitor           types.String `tfsdk:"monitor"`
	ComposerDev       types.Bool   `tfsdk:"composer_dev"`
	Composer          types.Bool   `tfsdk:"composer"`
	ComposerQuiet     types.Bool   `tfsdk:"composer_quiet"`

	// This is a special field that is not part of the actual envoyer api.
	DeleteProtection types.Bool `tfsdk:"delete_protection"`
}

func (r *EnvoyerProjectResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_envoyer_project"
}

func (r *EnvoyerProjectResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a project in Envoyer.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "The unique ID of the Envoyer project.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The project's display name.",
				Required:            true,
			},
			"repo_provider": schema.StringAttribute{
				MarkdownDescription: "Repository provider: `github`, `bitbucket`, `gitlab`, or `gitlab-self`.",
				Required:            true,
			},
			"repository": schema.StringAttribute{
				MarkdownDescription: "Repository slug or URL for self-hosted.",
				Required:            true,
			},
			"branch": schema.StringAttribute{
				MarkdownDescription: "Default git branch to deploy from.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("main"),
			},
			"type": schema.StringAttribute{
				MarkdownDescription: "Project type: `laravel-5`, `laravel-4`, or `other`.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("laravel-5"),
			},
			"retain_deployments": schema.Int64Attribute{
				MarkdownDescription: "Number of deployments to retain.",
				Optional:            true,
				Computed:            true,
				Default:             int64default.StaticInt64(5),
			},
			"monitor": schema.StringAttribute{
				MarkdownDescription: "Uptime monitoring URL (optional).",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"composer_dev": schema.BoolAttribute{
				MarkdownDescription: "Installation of dev dependencies.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
			},
			"composer": schema.BoolAttribute{
				MarkdownDescription: "Installation of dependencies.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
			},
			"composer_quiet": schema.BoolAttribute{
				MarkdownDescription: "Whether composer should run quietly.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"delete_protection": schema.BoolAttribute{
				MarkdownDescription: "Prevent this project from being deleted.",
				Optional:            true,
				Default:             booldefault.StaticBool(false),
				Computed:            true,
			},
		},
	}
}

// Configure is called when the resource is instantiated. We get the provider data (the Envoyer client).
func (r *EnvoyerProjectResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

// Create is called when Terraform needs to create a new project in Envoyer.
func (r *EnvoyerProjectResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan EnvoyerProjectResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// fix repo
	// git@github.com:laravel/laravel.git -> laravel/laravel
	splitRepo := strings.Split(plan.Repository.ValueString(), ":")
	if len(splitRepo) != 2 {
		resp.Diagnostics.AddError("Invalid repository format", "Expected format: git@github.com:laravel/laravel.git")
		return
	}

	fixRepo := strings.TrimSuffix(splitRepo[1], ".git")

	// >>> CALL Envoyer's CreateProject endpoint here:
	// e.g.,
	envoyerReq := envoyer_client.CreateProjectRequest{
		Name:              plan.Name.ValueString(),
		Provider:          plan.RepoProvider.ValueString(),
		Repository:        fixRepo,
		Branch:            plan.Branch.ValueString(),
		Type:              plan.Type.ValueString(),
		RetainDeployments: int(plan.RetainDeployments.ValueInt64()),
		Monitor:           plan.Monitor.ValueString(),
		Composer:          plan.Composer.ValueBool(),
		ComposerDev:       plan.ComposerDev.ValueBool(),
		ComposerQuiet:     plan.ComposerQuiet.ValueBool(),
	}
	createdProject, err := r.client.CreateProject(ctx, envoyerReq)

	if err != nil {
		resp.Diagnostics.AddError("Error creating Envoyer project", err.Error())
		return
	}

	tflog.Debug(ctx, "Created Envoyer project", map[string]any{
		"project_id": createdProject.ID,
	})

	// Set the ID in the state so Terraform knows it's created.
	plan.ID = types.Int64Value(createdProject.ID)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Read is called to refresh Terraform state from the Envoyer API.
func (r *EnvoyerProjectResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state EnvoyerProjectResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	projID := state.ID.ValueInt64()

	if projID <= 0 {
		resp.Diagnostics.AddError(
			"Invalid Project ID",
			fmt.Sprintf("Project ID must be positive, got: %d", projID),
		)
		return
	}

	project, err := r.client.GetProject(ctx, int(projID))
	if err != nil {
		resp.Diagnostics.AddError("Failed to read project", err.Error())
		return
	}

	tflog.Debug(ctx, "Read Envoyer project", map[string]any{
		"project_id": projID,
		"repository": project.Repository,
		"branch":     project.Branch,
	})

	if !state.Repository.IsNull() && state.Repository.ValueString() != project.Repository {
		splitRepo := strings.Split(state.Repository.ValueString(), ":")
		if len(splitRepo) != 2 {
			resp.Diagnostics.AddError("Invalid repository format", "Expected format: git@github.com:laravel/laravel.git")
			return
		}
		fixRepo := strings.TrimSuffix(splitRepo[1], ".git")

		if fixRepo != project.Repository {
			state.Repository = types.StringValue(project.Repository)
		}
	}

	// Set the state to the latest data from Envoyer
	state.Name = types.StringValue(project.Name)
	state.RepoProvider = types.StringValue(project.Provider)

	state.Branch = types.StringValue(project.Branch)
	state.Type = types.StringValue(project.Type)
	state.RetainDeployments = types.Int64Value(project.RetainDeployments)
	state.Monitor = types.StringValue(project.Monitor)
	state.ComposerDev = types.BoolValue(project.InstallDevDependencies)
	state.Composer = types.BoolValue(project.InstallDependencies)
	state.ComposerQuiet = types.BoolValue(project.QuietComposer)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Update is called when Terraform detects configuration changes that require the Envoyer project to be updated.
func (r *EnvoyerProjectResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan EnvoyerProjectResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	projID := plan.ID.ValueInt64()

	if projID <= 0 {
		resp.Diagnostics.AddError(
			"Invalid Project ID",
			fmt.Sprintf("Project ID must be positive, got: %d", projID),
		)
		return
	}

	envoyerReq := envoyer_client.UpdateProjectRequest{
		Name:              plan.Name.ValueString(),
		Monitor:           plan.Monitor.ValueString(),
		RetainDeployments: int(plan.RetainDeployments.ValueInt64()),
		Composer:          plan.Composer.ValueBool(),
		ComposerDev:       plan.ComposerDev.ValueBool(),
		ComposerQuiet:     plan.ComposerQuiet.ValueBool(),
	}

	err := r.client.UpdateProject(ctx, int(projID), envoyerReq)
	if err != nil {
		resp.Diagnostics.AddError("Error updating Envoyer project", err.Error())
		return
	}

	tflog.Debug(ctx, "Updated Envoyer project", map[string]any{
		"project_id": projID,
	})

	// fix repo
	// git@github.com:laravel/laravel.git -> laravel/laravel
	splitRepo := strings.Split(plan.Repository.ValueString(), ":")
	if len(splitRepo) != 2 {
		resp.Diagnostics.AddError("Invalid repository format", "Expected format: git@github.com:laravel/laravel.git")
		return
	}

	fixRepo := strings.TrimSuffix(splitRepo[1], ".git")

	envoyerSourceReq := envoyer_client.UpdateProjectSourceRequest{
		Provider:     plan.RepoProvider.ValueString(),
		Repository:   fixRepo,
		Branch:       plan.Branch.ValueString(),
		PushToDeploy: false,
	}

	err = r.client.UpdateProjectSource(ctx, int(projID), envoyerSourceReq)
	if err != nil {
		resp.Diagnostics.AddError("Error updating Envoyer project source", err.Error())
		return
	}

	// redundancy
	plan.ID = types.Int64Value(projID)
	plan.Name = types.StringValue(plan.Name.ValueString())
	plan.Monitor = types.StringValue(plan.Monitor.ValueString())
	plan.RetainDeployments = types.Int64Value(plan.RetainDeployments.ValueInt64())
	plan.Composer = types.BoolValue(plan.Composer.ValueBool())
	plan.ComposerDev = types.BoolValue(plan.ComposerDev.ValueBool())
	plan.ComposerQuiet = types.BoolValue(plan.ComposerQuiet.ValueBool())
	plan.DeleteProtection = types.BoolValue(plan.DeleteProtection.ValueBool())
	plan.Repository = types.StringValue(plan.Repository.ValueString())

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Delete is called when Terraform wants to remove this project from Envoyer.
func (r *EnvoyerProjectResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state EnvoyerProjectResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	projID := state.ID.ValueInt64()

	if projID <= 0 {
		resp.Diagnostics.AddError(
			"Invalid Project ID",
			fmt.Sprintf("Project ID must be positive, got: %d", projID),
		)
		return
	}

	if state.DeleteProtection.ValueBool() {
		resp.Diagnostics.AddError("Project is protected from deletion", "Set `delete_protection` to false to delete this project.")
		return
	}

	err := r.client.DeleteProject(ctx, int(projID))
	if err != nil {
		resp.Diagnostics.AddError("Error deleting Envoyer project", err.Error())
		return
	}

	tflog.Debug(ctx, "Deleted Envoyer project", map[string]any{
		"project_id": projID,
	})
}

// ImportState allows `terraform import envoyer_project.<name> <projectID>`.
func (r *EnvoyerProjectResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	projectId, err := strconv.ParseInt(req.ID, 10, 64)

	if err != nil {
		resp.Diagnostics.AddError("Invalid project ID", err.Error())
		return
	}

	project, err := r.client.GetProject(ctx, int(projectId))

	if err != nil {
		resp.Diagnostics.AddError("Failed to import project", err.Error())
		return
	}

	state := EnvoyerProjectResourceModel{
		ID:                types.Int64Value(project.ID),
		Name:              types.StringValue(project.Name),
		RepoProvider:      types.StringValue(project.Provider),
		Repository:        types.StringValue(project.Repository),
		Branch:            types.StringValue(project.Branch),
		Type:              types.StringValue(project.Type),
		RetainDeployments: types.Int64Value(project.RetainDeployments),
		Monitor:           types.StringValue(project.Monitor),
		ComposerDev:       types.BoolValue(project.InstallDevDependencies),
		Composer:          types.BoolValue(project.InstallDependencies),
		ComposerQuiet:     types.BoolValue(project.QuietComposer),
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
