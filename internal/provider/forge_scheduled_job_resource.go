package provider

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"terraform-provider-laravel/internal/forge_client"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ resource.Resource = &ForgeScheduledJobResource{}
var _ resource.ResourceWithImportState = &ForgeScheduledJobResource{}

type ForgeScheduledJobResource struct {
	client *forge_client.Client
}

type ForgeScheduledJobResourceModel struct {
	ID        types.Int64  `tfsdk:"id"`
	ServerID  types.Int64  `tfsdk:"server_id"`
	Command   types.String `tfsdk:"command"`
	Frequency types.String `tfsdk:"frequency"`
	User      types.String `tfsdk:"user"`
	Minute    types.String `tfsdk:"minute"`
	Hour      types.String `tfsdk:"hour"`
	Day       types.String `tfsdk:"day"`
	Month     types.String `tfsdk:"month"`
	Weekday   types.String `tfsdk:"weekday"`
	Status    types.String `tfsdk:"status"`
	CreatedAt types.String `tfsdk:"created_at"`
}

func NewForgeScheduledJobResource() resource.Resource {
	return &ForgeScheduledJobResource{}
}

func (r *ForgeScheduledJobResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_forge_scheduled_job"
}

func (r *ForgeScheduledJobResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Forge scheduled job resource. This resource allows you to manage scheduled jobs on Forge servers.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Computed: true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"server_id": schema.Int64Attribute{
				Required:            true,
				MarkdownDescription: "The ID of the server where the scheduled job will be run.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"command": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The command to run on the server.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"user": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("forge"),
				MarkdownDescription: "The user under which the command will be run.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplaceIfConfigured(),
				},
			},
			"frequency": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The frequency in which the job should run. Valid values are `minutely`, `hourly`, `nightly`, `weekly`, `monthly`, `reboot`, and `custom`",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplaceIfConfigured(),
				},
			},
			"minute": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The minute at which the job should run. Required if frequency is `custom`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplaceIfConfigured(),
				},
			},
			"hour": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The hour at which the job should run. Required if frequency is `custom`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplaceIfConfigured(),
				},
			},
			"day": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The day at which the job should run. Required if frequency is `custom`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplaceIfConfigured(),
				},
			},
			"month": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The month at which the job should run. Required if frequency is `custom`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplaceIfConfigured(),
				},
			},
			"weekday": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The weekday at which the job should run. Required if frequency is `custom`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplaceIfConfigured(),
				},
			},
			"status": schema.StringAttribute{
				Computed: true,
			},
			"created_at": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

func (r *ForgeScheduledJobResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *ForgeScheduledJobResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan ForgeScheduledJobResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	payload := forge_client.CreateJobRequest{
		Command:   plan.Command.ValueString(),
		Frequency: plan.Frequency.ValueString(),
		User:      plan.User.ValueString(),
		Minute:    plan.Minute.ValueString(),
		Hour:      plan.Hour.ValueString(),
		Day:       plan.Day.ValueString(),
		Month:     plan.Month.ValueString(),
		Weekday:   plan.Weekday.ValueString(),
	}

	scheduledJob, err := r.client.CreateJob(ctx, int(plan.ServerID.ValueInt64()), payload)
	if err != nil {
		resp.Diagnostics.AddError("Error creating scheduled job", err.Error())
		return
	}

	// Update plan state with response values.
	plan.ID = types.Int64Value(scheduledJob.ID)
	plan.ServerID = types.Int64Value(plan.ServerID.ValueInt64())
	plan.Command = types.StringValue(scheduledJob.Command)
	plan.Frequency = types.StringValue(strings.ToLower(scheduledJob.Frequency))
	plan.User = types.StringValue(scheduledJob.User)

	// Only parse cron expression for custom frequency
	if strings.EqualFold(strings.ToLower(scheduledJob.Frequency), "custom") {
		err := parseCronExpressionIntoModel(scheduledJob.Cron, &plan)
		if err != nil {
			resp.Diagnostics.AddError("Error parsing cron expression", err.Error())
			return
		}
	}

	plan.Status = types.StringValue(scheduledJob.Status)
	plan.CreatedAt = types.StringValue(scheduledJob.CreatedAt)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *ForgeScheduledJobResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state ForgeScheduledJobResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	scheduledJob, err := r.client.GetJob(ctx, int(state.ServerID.ValueInt64()), int(state.ID.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError("Error reading scheduled job", err.Error())
		return
	}
	state.Command = types.StringValue(scheduledJob.Command)
	state.Frequency = types.StringValue(strings.ToLower(scheduledJob.Frequency))
	state.User = types.StringValue(scheduledJob.User)

	if strings.EqualFold(strings.ToLower(scheduledJob.Frequency), "custom") {
		err := parseCronExpressionIntoModel(scheduledJob.Cron, &state)
		if err != nil {
			resp.Diagnostics.AddError("Error parsing cron expression", err.Error())
			return
		}
	}

	state.Status = types.StringValue(scheduledJob.Status)
	state.CreatedAt = types.StringValue(scheduledJob.CreatedAt)

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *ForgeScheduledJobResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan ForgeScheduledJobResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *ForgeScheduledJobResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state ForgeScheduledJobResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteJob(ctx, int(state.ServerID.ValueInt64()), int(state.ID.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError("Error deleting scheduled job", err.Error())
		return
	}
}

func (r *ForgeScheduledJobResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Expect import ID in format "server_id:scheduled_job_id"
	parts := splitCompositeID(req.ID, 2)
	if parts == nil {
		resp.Diagnostics.AddError("Invalid import ID", "Expected format: server_id:scheduled_job_id")
		return
	}
	serverID, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		resp.Diagnostics.AddError("Invalid server_id", err.Error())
		return
	}
	scheduledJobID, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		resp.Diagnostics.AddError("Invalid scheduled_job_id", err.Error())
		return
	}

	scheduledJob, err := r.client.GetJob(ctx, int(serverID), int(scheduledJobID))
	if err != nil {
		resp.Diagnostics.AddError("Error reading scheduled job", err.Error())
		return
	}

	var stateModel ForgeScheduledJobResourceModel
	stateModel.ID = types.Int64Value(scheduledJob.ID)
	stateModel.ServerID = types.Int64Value(serverID)
	stateModel.Command = types.StringValue(scheduledJob.Command)
	stateModel.Frequency = types.StringValue(strings.ToLower(scheduledJob.Frequency))
	stateModel.User = types.StringValue(scheduledJob.User)

	// Only parse cron expression for custom frequency
	if strings.EqualFold(strings.ToLower(scheduledJob.Frequency), "custom") {
		err := parseCronExpressionIntoModel(scheduledJob.Cron, &stateModel)
		if err != nil {
			resp.Diagnostics.AddError("Error parsing cron expression", err.Error())
			return
		}
	}

	stateModel.Status = types.StringValue(scheduledJob.Status)
	stateModel.CreatedAt = types.StringValue(scheduledJob.CreatedAt)

	diags := resp.State.Set(ctx, stateModel)
	resp.Diagnostics.Append(diags...)
}

func parseCronExpressionIntoModel(cronExpr string, model *ForgeScheduledJobResourceModel) error {
	// Clean the cron expression by removing any extra spaces or characters
	cronExpr = strings.TrimSpace(cronExpr)

	// Split the cron expression by spaces
	cronParts := strings.Fields(cronExpr)

	// Standard cron has 5 parts: minute, hour, day, month, weekday
	if len(cronParts) < 5 {
		return fmt.Errorf("invalid cron expression: %s (expected at least 5 parts, got %d)", cronExpr, len(cronParts))
	}

	// Assign the first 5 parts to the model fields
	model.Minute = types.StringValue(cronParts[0])
	model.Hour = types.StringValue(cronParts[1])
	model.Day = types.StringValue(cronParts[2])
	model.Month = types.StringValue(cronParts[3])
	model.Weekday = types.StringValue(cronParts[4])

	return nil
}
