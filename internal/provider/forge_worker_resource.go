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

// Ensure ForgeWorkerResource satisfies required interfaces.
var _ resource.Resource = &ForgeWorkerResource{}
var _ resource.ResourceWithImportState = &ForgeWorkerResource{}

// ForgeWorkerResource implements a Terraform resource for a Forge worker.
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

// NewForgeWorkerResource returns a new instance.
func NewForgeWorkerResource() resource.Resource {
	return &ForgeWorkerResource{}
}

func (r *ForgeWorkerResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_forge_worker"
}

func (r *ForgeWorkerResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Computed: true,
			},
			"server_id": schema.Int64Attribute{
				Required: true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"site_id": schema.Int64Attribute{
				Required: true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"worker_connection": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"timeout": schema.Int64Attribute{
				Optional: true,
				Computed: true,
				Default:  int64default.StaticInt64(60),
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"sleep": schema.Int64Attribute{
				Optional: true,
				Computed: true,
				Default:  int64default.StaticInt64(3),
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"tries": schema.Int64Attribute{
				Optional: true,
				Computed: true,
				Default:  int64default.StaticInt64(0),
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"processes": schema.Int64Attribute{
				Optional: true,
				Computed: true,
				Default:  int64default.StaticInt64(1),
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"stop_wait_secs": schema.Int64Attribute{
				Optional: true,
				Computed: true,
				Default:  int64default.StaticInt64(10),
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"delay": schema.Int64Attribute{
				Optional: true,
				Computed: true,
				Default:  int64default.StaticInt64(0),
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"daemon": schema.BoolAttribute{
				Optional: true,
				Computed: true,
				Default:  booldefault.StaticBool(true),
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.RequiresReplace(),
				},
			},
			"force": schema.BoolAttribute{
				Optional: true,
				Computed: true,
				Default:  booldefault.StaticBool(false),
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.RequiresReplace(),
				},
			},
			"php_version": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString("php"),
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"queue": schema.StringAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplaceIfConfigured(),
				},
			},
			"memory": schema.Int64Attribute{
				Optional: true,
				Computed: true,
				Default:  int64default.StaticInt64(128),
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"directory": schema.StringAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
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
	providerConfig, ok := req.ProviderData.(*providerConfig)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Provider Configure Type",
			fmt.Sprintf("Expected *providerConfig, got: %T. Please report this issue to the provider developers.", req.ProviderData),
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

	createReq := forge_client.CreateWorkerRequest{
		Connection: plan.WorkerConnection.ValueString(),
		TimeOut:    int(plan.Timeout.ValueInt64()),
		Sleep:      int(plan.Sleep.ValueInt64()),
		Delay:      int(plan.Delay.ValueInt64()),
		Processes:  int(plan.Processes.ValueInt64()),
		Daemon:     plan.Daemon.ValueBool(),
		Force:      plan.Force.ValueBool(),
		PHPVersion: plan.PHPVersion.ValueString(),
		Memory:     int(plan.Memory.ValueInt64()),
		Directory:  plan.Directory.ValueString(),
	}

	// Optional fields.
	if !plan.Tries.IsNull() && !plan.StopWaitSecs.IsUnknown() {
		t := int(plan.Tries.ValueInt64())
		createReq.Tries = &t
	}
	if !plan.StopWaitSecs.IsNull() && !plan.StopWaitSecs.IsUnknown() {
		s := int(plan.StopWaitSecs.ValueInt64())
		createReq.StopWaitSecs = &s
	}
	if !plan.Queue.IsNull() && plan.Queue.ValueString() != "" {
		q := plan.Queue.ValueString()
		createReq.Queue = &q
	}

	worker, err := r.client.CreateWorker(ctx, int(plan.ServerID.ValueInt64()), int(plan.SiteID.ValueInt64()), createReq)
	if err != nil {
		resp.Diagnostics.AddError("Error creating worker", err.Error())
		return
	}

	// Parse memory from worker.Command.
	mem, err := parseMemoryFromCommand(worker.Command)
	if err != nil {
		// Fallback: use the provided memory if parsing fails.
		mem = int(plan.Memory.ValueInt64())
	}

	plan.ID = types.Int64Value(worker.ID)
	plan.WorkerConnection = types.StringValue(worker.Connection)
	plan.Timeout = types.Int64Value(int64(worker.Timeout))
	plan.Sleep = types.Int64Value(int64(worker.Sleep))
	if worker.Tries != nil {
		plan.Tries = types.Int64Value(int64(*worker.Tries))
	} else {
		plan.Tries = types.Int64Value(0)
	}
	plan.Processes = types.Int64Value(int64(worker.Processes))
	if worker.StopWaitSecs != nil {
		plan.StopWaitSecs = types.Int64Value(int64(*worker.StopWaitSecs))
	} else {
		plan.StopWaitSecs = types.Int64Null()
	}
	plan.Daemon = types.BoolValue(worker.Daemon)
	plan.Force = types.BoolValue(worker.Force)
	plan.Delay = types.Int64Value(int64(worker.Delay))

	phpVersion, err := r.client.GetPHPVersionFromDisplayableVersion(ctx, int(plan.ServerID.ValueInt64()), worker.DisplayablePHPVersion)
	if err != nil {
		resp.Diagnostics.AddError("Error getting PHP version", err.Error())
		return
	}

	plan.PHPVersion = types.StringValue(phpVersion.Version)
	if worker.Queue != nil {
		plan.Queue = types.StringValue(*worker.Queue)
	} else {
		plan.Queue = types.StringValue("")
	}
	plan.Memory = types.Int64Value(int64(mem))
	plan.Command = types.StringValue(worker.Command)
	plan.Status = types.StringValue(worker.Status)
	plan.CreatedAt = types.StringValue(worker.CreatedAt)

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

	state.WorkerConnection = types.StringValue(worker.Connection)
	state.Timeout = types.Int64Value(int64(worker.Timeout))
	state.Sleep = types.Int64Value(int64(worker.Sleep))
	if worker.Tries != nil {
		state.Tries = types.Int64Value(int64(*worker.Tries))
	} else {
		state.Tries = types.Int64Value(0)
	}
	state.Processes = types.Int64Value(int64(worker.Processes))
	if worker.StopWaitSecs != nil {
		state.StopWaitSecs = types.Int64Value(int64(*worker.StopWaitSecs))
	} else {
		state.StopWaitSecs = types.Int64Value(0)
	}
	state.Daemon = types.BoolValue(worker.Daemon)
	state.Force = types.BoolValue(worker.Force)
	state.Delay = types.Int64Value(int64(worker.Delay))

	phpVersion, err := r.client.GetPHPVersionFromDisplayableVersion(ctx, int(state.ServerID.ValueInt64()), worker.DisplayablePHPVersion)
	if err != nil {
		resp.Diagnostics.AddError("Error getting PHP version", err.Error())
		return
	}

	state.PHPVersion = types.StringValue(phpVersion.Version)
	if worker.Queue != nil {
		state.Queue = types.StringValue(*worker.Queue)
	} else {
		state.Queue = types.StringValue("")
	}
	// Parse memory from the command.
	mem, err := parseMemoryFromCommand(worker.Command)
	if err != nil {
		mem = int(state.Memory.ValueInt64())
	}
	state.Memory = types.Int64Value(int64(mem))
	state.Command = types.StringValue(worker.Command)
	state.Status = types.StringValue(worker.Status)
	state.CreatedAt = types.StringValue(worker.CreatedAt)

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
	stateModel.WorkerConnection = types.StringValue(worker.Connection)
	stateModel.Timeout = types.Int64Value(int64(worker.Timeout))
	stateModel.Delay = types.Int64Value(int64(worker.Delay))
	stateModel.Sleep = types.Int64Value(int64(worker.Sleep))
	if worker.Tries != nil {
		stateModel.Tries = types.Int64Value(int64(*worker.Tries))
	} else {
		stateModel.Tries = types.Int64Value(0)
	}
	stateModel.Processes = types.Int64Value(int64(worker.Processes))
	if worker.StopWaitSecs != nil {
		stateModel.StopWaitSecs = types.Int64Value(int64(*worker.StopWaitSecs))
	} else {
		stateModel.StopWaitSecs = types.Int64Value(0)
	}
	stateModel.Daemon = types.BoolValue(worker.Daemon)
	stateModel.Force = types.BoolValue(worker.Force)

	phpVersion, err := r.client.GetPHPVersionFromDisplayableVersion(ctx, int(serverID), worker.DisplayablePHPVersion)
	if err != nil {
		resp.Diagnostics.AddError("Error getting PHP version", err.Error())
		return
	}

	stateModel.PHPVersion = types.StringValue(phpVersion.Version)
	if worker.Queue != nil {
		stateModel.Queue = types.StringValue(*worker.Queue)
	} else {
		stateModel.Queue = types.StringValue("")
	}
	mem, err := parseMemoryFromCommand(worker.Command)
	if err != nil {
		mem = 0
	}
	stateModel.Memory = types.Int64Value(int64(mem))
	stateModel.Command = types.StringValue(worker.Command)
	stateModel.Status = types.StringValue(worker.Status)
	stateModel.CreatedAt = types.StringValue(worker.CreatedAt)

	diags := resp.State.Set(ctx, &stateModel)
	resp.Diagnostics.Append(diags...)
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
