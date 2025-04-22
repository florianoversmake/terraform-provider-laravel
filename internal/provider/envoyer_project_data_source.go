package provider

import (
	"context"
	"fmt"
	"terraform-provider-laravel/internal/envoyer_client"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var _ datasource.DataSource = &EnvoyerProjectDataSource{}

func NewEnvoyerProjectDataSource() datasource.DataSource {
	return &EnvoyerProjectDataSource{}
}

// EnvoyerProjectDataSource allows reading an existing Envoyer project by ID.
type EnvoyerProjectDataSource struct {
	client *envoyer_client.Client
}

// EnvoyerProjectDataSourceModel describes the data source schema.
type EnvoyerProjectDataSourceModel struct {
	ID                   types.Int64  `tfsdk:"id"`
	UserID               types.Int64  `tfsdk:"user_id"`
	Version              types.Int64  `tfsdk:"version"`
	Name                 types.String `tfsdk:"name"`
	RepoProvider         types.String `tfsdk:"repo_provider"` // Renamed from Provider to avoid conflict.
	Repository           types.String `tfsdk:"repository"`
	Type                 types.String `tfsdk:"type"`
	Branch               types.String `tfsdk:"branch"`
	PushToDeploy         types.Bool   `tfsdk:"push_to_deploy"`
	WebhookID            types.String `tfsdk:"webhook_id"`
	Status               types.String `tfsdk:"status"`
	ShouldDeployAgain    types.Int64  `tfsdk:"should_deploy_again"`
	DeploymentStartedAt  types.String `tfsdk:"deployment_started_at"`
	DeploymentFinishedAt types.String `tfsdk:"deployment_finished_at"`
	LastDeploymentStatus types.String `tfsdk:"last_deployment_status"`
	DailyDeploys         types.Int64  `tfsdk:"daily_deploys"`
	WeeklyDeploys        types.Int64  `tfsdk:"weekly_deploys"`
	LastDeploymentTook   types.Int64  `tfsdk:"last_deployment_took"`
	RetainDeployments    types.Int64  `tfsdk:"retain_deployments"`

	EnvironmentServers types.List `tfsdk:"environment_servers"`

	// Folders                 types.String `tfsdk:"folders"`
	Monitor                types.String `tfsdk:"monitor"`
	NewYorkStatus          types.String `tfsdk:"new_york_status"`
	LondonStatus           types.String `tfsdk:"london_status"`
	SingaporeStatus        types.String `tfsdk:"singapore_status"`
	Token                  types.String `tfsdk:"token"`
	CreatedAt              types.String `tfsdk:"created_at"`
	UpdatedAt              types.String `tfsdk:"updated_at"`
	InstallDevDependencies types.Bool   `tfsdk:"install_dev_dependencies"`
	InstallDependencies    types.Bool   `tfsdk:"install_dependencies"`
	QuietComposer          types.Bool   `tfsdk:"quiet_composer"`
	//Servers                 types.List   `tfsdk:"servers"`
	HasEnvironment          types.Bool   `tfsdk:"has_environment"`
	HasMonitoringError      types.Bool   `tfsdk:"has_monitoring_error"`
	HasMissingHeartbeats    types.Bool   `tfsdk:"has_missing_heartbeats"`
	LastDeployedBranch      types.String `tfsdk:"last_deployed_branch"`
	LastDeploymentID        types.Int64  `tfsdk:"last_deployment_id"`
	LastDeploymentAuthor    types.String `tfsdk:"last_deployment_author"`
	LastDeploymentAvatar    types.String `tfsdk:"last_deployment_avatar"`
	LastDeploymentHash      types.String `tfsdk:"last_deployment_hash"`
	LastDeploymentTimestamp types.String `tfsdk:"last_deployment_timestamp"`
}

func (d *EnvoyerProjectDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_envoyer_project"
}

func (d *EnvoyerProjectDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Data source for retrieving an existing Envoyer project by ID.",

		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				MarkdownDescription: "Project ID.",
				Required:            true,
			},
			"environment_servers": schema.ListAttribute{
				MarkdownDescription: "Environment servers.",
				ElementType:         types.Int64Type,
				Computed:            true,
			},
			"user_id": schema.Int64Attribute{
				MarkdownDescription: "User ID.",
				Computed:            true,
			},
			"version": schema.Int64Attribute{
				MarkdownDescription: "Version.",
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Project name.",
				Computed:            true,
			},
			"repo_provider": schema.StringAttribute{
				MarkdownDescription: "Repository provider.",
				Computed:            true,
			},
			"repository": schema.StringAttribute{
				MarkdownDescription: "Repository URL.",
				Computed:            true,
			},
			"type": schema.StringAttribute{
				MarkdownDescription: "Project type.",
				Computed:            true,
			},
			"branch": schema.StringAttribute{
				MarkdownDescription: "Branch.",
				Computed:            true,
			},
			"push_to_deploy": schema.BoolAttribute{
				MarkdownDescription: "Push to deploy.",
				Computed:            true,
			},
			"webhook_id": schema.StringAttribute{
				MarkdownDescription: "Webhook ID.",
				Computed:            true,
			},
			"status": schema.StringAttribute{
				MarkdownDescription: "Status.",
				Computed:            true,
			},
			"should_deploy_again": schema.Int64Attribute{
				MarkdownDescription: "Should deploy again.",
				Computed:            true,
			},
			"deployment_started_at": schema.StringAttribute{
				MarkdownDescription: "Deployment started at.",
				Computed:            true,
			},
			"deployment_finished_at": schema.StringAttribute{
				MarkdownDescription: "Deployment finished at.",
				Computed:            true,
			},
			"last_deployment_status": schema.StringAttribute{
				MarkdownDescription: "Last deployment status.",
				Computed:            true,
			},
			"daily_deploys": schema.Int64Attribute{
				MarkdownDescription: "Daily deploys.",
				Computed:            true,
			},
			"weekly_deploys": schema.Int64Attribute{
				MarkdownDescription: "Weekly deploys.",
				Computed:            true,
			},
			"last_deployment_took": schema.Int64Attribute{
				MarkdownDescription: "Last deployment took.",
				Computed:            true,
			},
			"retain_deployments": schema.Int64Attribute{
				MarkdownDescription: "Retain deployments.",
				Computed:            true,
			},
			"monitor": schema.StringAttribute{
				MarkdownDescription: "Monitor.",
				Computed:            true,
			},
			"new_york_status": schema.StringAttribute{
				MarkdownDescription: "New York status.",
				Computed:            true,
			},
			"london_status": schema.StringAttribute{
				MarkdownDescription: "London status.",
				Computed:            true,
			},
			"singapore_status": schema.StringAttribute{
				MarkdownDescription: "Singapore status.",
				Computed:            true,
			},
			"token": schema.StringAttribute{
				MarkdownDescription: "Token.",
				Computed:            true,
			},
			"created_at": schema.StringAttribute{
				MarkdownDescription: "Created at.",
				Computed:            true,
			},
			"updated_at": schema.StringAttribute{
				MarkdownDescription: "Updated at.",
				Computed:            true,
			},
			"install_dev_dependencies": schema.BoolAttribute{
				MarkdownDescription: "Install dev dependencies.",
				Computed:            true,
			},
			"install_dependencies": schema.BoolAttribute{
				MarkdownDescription: "Install dependencies.",
				Computed:            true,
			},
			"quiet_composer": schema.BoolAttribute{
				MarkdownDescription: "Quiet Composer.",
				Computed:            true,
			},
			"has_environment": schema.BoolAttribute{
				MarkdownDescription: "Has environment.",
				Computed:            true,
			},
			"has_monitoring_error": schema.BoolAttribute{
				MarkdownDescription: "Has monitoring error.",
				Computed:            true,
			},
			"has_missing_heartbeats": schema.BoolAttribute{
				MarkdownDescription: "Has missing heartbeats.",
				Computed:            true,
			},
			"last_deployed_branch": schema.StringAttribute{
				MarkdownDescription: "Last deployed branch.",
				Computed:            true,
			},
			"last_deployment_id": schema.Int64Attribute{
				MarkdownDescription: "Last deployment ID.",
				Computed:            true,
			},
			"last_deployment_author": schema.StringAttribute{
				MarkdownDescription: "Last deployment author.",
				Computed:            true,
			},
			"last_deployment_avatar": schema.StringAttribute{
				MarkdownDescription: "Last deployment avatar.",
				Computed:            true,
			},
			"last_deployment_hash": schema.StringAttribute{
				MarkdownDescription: "Last deployment hash.",
				Computed:            true,
			},
			"last_deployment_timestamp": schema.StringAttribute{
				MarkdownDescription: "Last deployment timestamp.",
				Computed:            true,
			},
		},
	}
}

// Configure obtains the *http.Client from the provider.
func (d *EnvoyerProjectDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *EnvoyerProjectDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data EnvoyerProjectDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	projID := data.ID.ValueInt64()
	if projID <= 0 {
		resp.Diagnostics.AddError(
			"Invalid Project ID",
			fmt.Sprintf("Project ID must be positive, got: %d", projID),
		)
		return
	}

	project, err := d.client.GetProject(ctx, int(projID))
	if err != nil {
		resp.Diagnostics.AddError("Failed to read project", err.Error())
		return
	}

	tflog.Debug(ctx, "Reading Envoyer project from data source", map[string]any{
		"project_id": projID,
	})

	data.ID = types.Int64Value(project.ID)
	data.EnvironmentServers.ElementsAs(ctx, &project.EnvironmentServers, false)
	data.UserID = types.Int64Value(project.UserID)
	data.Version = types.Int64Value(project.Version)
	data.Name = types.StringValue(project.Name)
	data.RepoProvider = types.StringValue(project.Provider)
	data.Repository = types.StringValue(project.Repository)
	data.Type = types.StringValue(project.Type)
	data.Branch = types.StringValue(project.Branch)
	data.PushToDeploy = types.BoolValue(project.PushToDeploy)
	data.WebhookID = types.StringPointerValue(project.WebhookID)
	data.Status = types.StringPointerValue(project.Status)
	data.ShouldDeployAgain = types.Int64Value(project.ShouldDeployAgain)
	if project.DeploymentStartedAt != nil {
		data.DeploymentStartedAt = types.StringValue(project.DeploymentStartedAt.String())
	} else {
		data.DeploymentStartedAt = types.StringNull()
	}
	data.DeploymentFinishedAt = types.StringValue(project.DeploymentFinishedAt.String())
	data.LastDeploymentStatus = types.StringValue(project.LastDeploymentStatus)
	data.DailyDeploys = types.Int64Value(project.DailyDeploys)
	data.WeeklyDeploys = types.Int64Value(project.WeeklyDeploys)
	data.LastDeploymentTook = types.Int64Value(project.LastDeploymentTook)
	data.RetainDeployments = types.Int64Value(project.RetainDeployments)
	data.Monitor = types.StringValue(project.Monitor)
	data.NewYorkStatus = types.StringValue(project.NewYorkStatus)
	data.LondonStatus = types.StringValue(project.LondonStatus)
	data.SingaporeStatus = types.StringValue(project.SingaporeStatus)
	data.Token = types.StringValue(project.Token)
	data.CreatedAt = types.StringValue(project.CreatedAt.String())
	data.UpdatedAt = types.StringValue(project.UpdatedAt.String())
	data.InstallDevDependencies = types.BoolValue(project.InstallDevDependencies)
	data.InstallDependencies = types.BoolValue(project.InstallDependencies)
	data.QuietComposer = types.BoolValue(project.QuietComposer)
	data.HasEnvironment = types.BoolValue(project.HasEnvironment)
	data.HasMonitoringError = types.BoolValue(project.HasMonitoringError)
	data.HasMissingHeartbeats = types.BoolValue(project.HasMissingHeartbeats)
	data.LastDeployedBranch = types.StringValue(project.LastDeployedBranch)
	data.LastDeploymentID = types.Int64Value(project.LastDeploymentID)
	data.LastDeploymentAuthor = types.StringValue(project.LastDeploymentAuthor)
	data.LastDeploymentAvatar = types.StringValue(project.LastDeploymentAvatar)
	data.LastDeploymentHash = types.StringValue(project.LastDeploymentHash)
	data.LastDeploymentTimestamp = types.StringValue(project.LastDeploymentTimestamp)

	resp.State.Set(ctx, &data)
}
