package provider

import (
	"context"
	"fmt"
	"regexp"
	"strconv"

	"terraform-provider-laravel/internal/forge_client"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ resource.Resource = &ForgeWorkerResource{}
var _ resource.ResourceWithImportState = &ForgeWorkerResource{}

type ForgeWorkerResource struct {
	client *forge_client.Client
}

// Note: we add a "directory" field because CreateWorkerRequest requires it.
type ForgeWorkerResourceModel struct {
	ID               types.Int64  `tfsdk:"id"`
	ServerID         types.Int64  `tfsdk:"server_id"`
	SiteID           types.Int64  `tfsdk:"site_id"`
	WorkerConnection types.String `tfsdk:"worker_connection"`
	Timeout          types.Int64  `tfsdk:"timeout"`
	Sleep            types.Int64  `tfsdk:"sleep"`
	Tries            types.Int64  `tfsdk:"tries"`
	Processes        types.Int64  `tfsdk:"processes"`
	StopWaitSecs     types.Int64  `tfsdk:"stop_wait_secs"`
	Delay            types.Int64  `tfsdk:"delay"`
	Daemon           types.Bool   `tfsdk:"daemon"`
	Force            types.Bool   `tfsdk:"force"`
	PHPVersion       types.String `tfsdk:"php_version"`
	Queue            types.String `tfsdk:"queue"`
	Memory           types.Int64  `tfsdk:"memory"`
	Directory        types.String `tfsdk:"directory"`
	Command          types.String `tfsdk:"command"`
	Status           types.String `tfsdk:"status"`
	CreatedAt        types.String `tfsdk:"created_at"`
}

func NewForgeWorkerResource() resource.Resource {
	return &ForgeWorkerResource{}
}

func (r *ForgeWorkerResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_forge_worker"
}

func (r *ForgeWorkerResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Forge worker resource. This resource allows you to manage workers in Forge.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Computed: true,
			},
			"server_id": schema.Int64Attribute{
				Required: true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
				MarkdownDescription: "The ID of the server where the worker is created.",
			},
			"site_id": schema.Int64Attribute{
				Required: true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
				MarkdownDescription: "The ID of the site where the worker is created.",
			},
			"worker_connection": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				MarkdownDescription: "The connection string for the worker. Like `sync`, `database`, `beanstalkd`, `sqs`, `redis`...",
			},
			"timeout": schema.Int64Attribute{
				Optional: true,
				Computed: true,
				Default:  int64default.StaticInt64(60),
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
				MarkdownDescription: "The timeout for the worker in seconds. Default is 60 seconds.",
			},
			"sleep": schema.Int64Attribute{
				Optional: true,
				Computed: true,
				Default:  int64default.StaticInt64(3),
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
				MarkdownDescription: "The sleep time for the worker in seconds. Default is 3 seconds.",
			},
			"tries": schema.Int64Attribute{
				Optional: true,
				Computed: true,
				Default:  int64default.StaticInt64(0),
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
				MarkdownDescription: "The number of tries for the worker. Default is 0 (unlimited).",
			},
			"processes": schema.Int64Attribute{
				Optional: true,
				Computed: true,
				Default:  int64default.StaticInt64(1),
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
				MarkdownDescription: "The number of processes for the worker. Default is 1.",
			},
			"stop_wait_secs": schema.Int64Attribute{
				Optional: true,
				Computed: true,
				Default:  int64default.StaticInt64(10),
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
				MarkdownDescription: "The number of seconds to wait for the worker to stop. Default is 10 seconds. You should ensure that the value of stopwaitsecs is greater than the number of seconds consumed by your longest running job. Otherwise, Supervisor may kill the job before it is finished processing.",
			},
			"delay": schema.Int64Attribute{
				Optional: true,
				Computed: true,
				Default:  int64default.StaticInt64(0),
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
				MarkdownDescription: "The delay time for the worker in seconds. Default is 0 seconds.",
			},
			"daemon": schema.BoolAttribute{
				Optional: true,
				Computed: true,
				Default:  booldefault.StaticBool(true),
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.RequiresReplace(),
				},
				MarkdownDescription: "Whether the worker should run as a daemon. Default is true.",
			},
			"force": schema.BoolAttribute{
				Optional: true,
				Computed: true,
				Default:  booldefault.StaticBool(false),
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.RequiresReplace(),
				},
				MarkdownDescription: "To force your queue workers to process jobs even if maintenance mode is enabled, you may use force option.",
			},
			"php_version": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString("php"),
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				MarkdownDescription: "The PHP version to use for the worker. Default is 'php' (System default).",
			},
			"queue": schema.StringAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplaceIfConfigured(),
				},
				MarkdownDescription: "The queue name for the worker. Default is empty string (no specific queue).",
			},
			"memory": schema.Int64Attribute{
				Optional: true,
				Computed: true,
				Default:  int64default.StaticInt64(128),
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
				MarkdownDescription: "The memory limit for the worker in megabytes. Default is 128 MB.",
			},
			"directory": schema.StringAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				MarkdownDescription: "The directory where the worker is located. Default is empty string (current directory).",
			},
			"command": schema.StringAttribute{
				Computed: true,
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

func (r *ForgeWorkerResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *ForgeWorkerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan ForgeWorkerResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Build the artisan queue:work command from individual parameters.
	cmd := buildQueueWorkerCommand(plan)

	var dir *string
	if !plan.Directory.IsNull() && plan.Directory.ValueString() != "" {
		d := plan.Directory.ValueString()
		dir = &d
	}

	stopWait := int(plan.StopWaitSecs.ValueInt64())
	createReq := forge_client.CreateWorkerRequest{
		Name:         fmt.Sprintf("queue-worker-%s", plan.WorkerConnection.ValueString()),
		Command:      cmd,
		User:         "forge",
		Directory:    dir,
		Processes:    int(plan.Processes.ValueInt64()),
		StopWaitSecs: &stopWait,
	}

	worker, err := r.client.CreateWorker(ctx, int(plan.ServerID.ValueInt64()), int(plan.SiteID.ValueInt64()), createReq)
	if err != nil {
		resp.Diagnostics.AddError("Error creating worker", err.Error())
		return
	}

	plan.ID = types.Int64Value(worker.ID)
	plan.Processes = types.Int64Value(int64(worker.Processes))
	plan.Command = types.StringValue(worker.Command)
	plan.Status = types.StringValue(worker.Status)
	plan.CreatedAt = types.StringValue(worker.CreatedAt)
	if worker.Directory != nil {
		plan.Directory = types.StringValue(*worker.Directory)
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *ForgeWorkerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state ForgeWorkerResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	worker, err := r.client.GetWorker(ctx, int(state.ServerID.ValueInt64()), int(state.SiteID.ValueInt64()), int(state.ID.ValueInt64()))
	if err != nil {
		_, ok := err.(*forge_client.ErrorWorkerNotFound)
		if ok {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading worker", err.Error())
		return
	}

	// The new API returns a simplified background process model.
	// We update computed fields from the API response and preserve
	// the config fields from state (they all have RequiresReplace).
	state.Processes = types.Int64Value(int64(worker.Processes))
	state.Command = types.StringValue(worker.Command)
	state.Status = types.StringValue(worker.Status)
	state.CreatedAt = types.StringValue(worker.CreatedAt)
	if worker.Directory != nil {
		state.Directory = types.StringValue(*worker.Directory)
	}
	// Parse memory from the command if present.
	mem, err := parseMemoryFromCommand(worker.Command)
	if err != nil {
		mem = int(state.Memory.ValueInt64())
	}
	state.Memory = types.Int64Value(int64(mem))

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *ForgeWorkerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// No update API exists for workers so we simply pass through.
	var plan ForgeWorkerResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *ForgeWorkerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state ForgeWorkerResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteWorker(ctx, int(state.ServerID.ValueInt64()), int(state.SiteID.ValueInt64()), int(state.ID.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError("Error deleting worker", err.Error())
		return
	}
}

func (r *ForgeWorkerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Expect import ID in the format "server_id:site_id:worker_id"
	parts := splitCompositeID(req.ID, 3)
	if parts == nil {
		resp.Diagnostics.AddError("Invalid import ID", "Expected format: server_id:site_id:worker_id")
		return
	}
	serverID, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		resp.Diagnostics.AddError("Invalid server_id", err.Error())
		return
	}
	siteID, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		resp.Diagnostics.AddError("Invalid site_id", err.Error())
		return
	}
	workerID, err := strconv.ParseInt(parts[2], 10, 64)
	if err != nil {
		resp.Diagnostics.AddError("Invalid worker_id", err.Error())
		return
	}

	worker, err := r.client.GetWorker(ctx, int(serverID), int(siteID), int(workerID))
	if err != nil {
		resp.Diagnostics.AddError("Error reading worker", err.Error())
		return
	}

	var stateModel ForgeWorkerResourceModel
	stateModel.ID = types.Int64Value(worker.ID)
	stateModel.ServerID = types.Int64Value(serverID)
	stateModel.SiteID = types.Int64Value(siteID)
	stateModel.Processes = types.Int64Value(int64(worker.Processes))
	stateModel.Command = types.StringValue(worker.Command)
	stateModel.Status = types.StringValue(worker.Status)
	stateModel.CreatedAt = types.StringValue(worker.CreatedAt)
	if worker.Directory != nil {
		stateModel.Directory = types.StringValue(*worker.Directory)
	} else {
		stateModel.Directory = types.StringValue("")
	}

	// Set defaults for fields that the new API no longer returns.
	// These are preserved for backward compatibility with existing configs.
	stateModel.WorkerConnection = types.StringValue("")
	stateModel.Timeout = types.Int64Value(60)
	stateModel.Delay = types.Int64Value(0)
	stateModel.Sleep = types.Int64Value(3)
	stateModel.Tries = types.Int64Value(0)
	stateModel.StopWaitSecs = types.Int64Value(10)
	stateModel.Daemon = types.BoolValue(true)
	stateModel.Force = types.BoolValue(false)
	stateModel.PHPVersion = types.StringValue("php")
	stateModel.Queue = types.StringValue("")
	mem, err := parseMemoryFromCommand(worker.Command)
	if err != nil {
		mem = 128
	}
	stateModel.Memory = types.Int64Value(int64(mem))

	diags := resp.State.Set(ctx, &stateModel)
	resp.Diagnostics.Append(diags...)
}

// buildQueueWorkerCommand constructs a php artisan queue:work command from resource model fields.
func buildQueueWorkerCommand(plan ForgeWorkerResourceModel) string {
	phpVersion := plan.PHPVersion.ValueString()
	if phpVersion == "" {
		phpVersion = "php"
	}
	connection := plan.WorkerConnection.ValueString()
	cmd := fmt.Sprintf("%s artisan queue:work %s", phpVersion, connection)

	if !plan.Queue.IsNull() && plan.Queue.ValueString() != "" {
		cmd += fmt.Sprintf(" --queue=%s", plan.Queue.ValueString())
	}
	if !plan.Delay.IsNull() && plan.Delay.ValueInt64() > 0 {
		cmd += fmt.Sprintf(" --delay=%d", plan.Delay.ValueInt64())
	}
	if !plan.Memory.IsNull() {
		cmd += fmt.Sprintf(" --memory=%d", plan.Memory.ValueInt64())
	}
	if !plan.Sleep.IsNull() {
		cmd += fmt.Sprintf(" --sleep=%d", plan.Sleep.ValueInt64())
	}
	if !plan.Timeout.IsNull() {
		cmd += fmt.Sprintf(" --timeout=%d", plan.Timeout.ValueInt64())
	}
	if !plan.Tries.IsNull() && plan.Tries.ValueInt64() > 0 {
		cmd += fmt.Sprintf(" --tries=%d", plan.Tries.ValueInt64())
	}
	if !plan.Force.IsNull() && plan.Force.ValueBool() {
		cmd += " --force"
	}
	if !plan.Daemon.IsNull() && plan.Daemon.ValueBool() {
		cmd += " --daemon"
	}
	return cmd
}

// parseMemoryFromCommand parses a memory value from the worker command string.
// It expects the flag format "--memory 128" or "--memory=128".
func parseMemoryFromCommand(command string) (int, error) {
	re := regexp.MustCompile(`--memory(?:=|\s+)(\d+)`)
	matches := re.FindStringSubmatch(command)
	if len(matches) >= 2 {
		return strconv.Atoi(matches[1])
	}
	return 0, fmt.Errorf("memory flag not found in command")
}
