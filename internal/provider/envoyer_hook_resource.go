package provider

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"terraform-provider-laravel/internal/envoyer_client"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure HookResource implements required interfaces.
var _ resource.Resource = &EnvoyerHookResource{}
var _ resource.ResourceWithImportState = &EnvoyerHookResource{}

func NewEnvoyerHookResource() resource.Resource {
	return &EnvoyerHookResource{}
}

type EnvoyerHookResource struct {
	client *envoyer_client.Client
}

type EnvoyerHookResourceModel struct {
	ID        types.Int64  `tfsdk:"id"`
	ProjectID types.Int64  `tfsdk:"project_id"`
	ActionID  types.Int64  `tfsdk:"action_id"`
	Timing    types.String `tfsdk:"timing"`
	Name      types.String `tfsdk:"name"`
	RunAs     types.String `tfsdk:"run_as"`
	Script    types.String `tfsdk:"script"`
	Sequence  types.Int64  `tfsdk:"sequence"`
	Servers   types.List   `tfsdk:"servers"` // List of server IDs (int64)
	CreatedAt types.String `tfsdk:"created_at"`
	UpdatedAt types.String `tfsdk:"updated_at"`
}

func (r *EnvoyerHookResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_envoyer_hook"
}

func (r *EnvoyerHookResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Computed: true,
			},
			"project_id": schema.Int64Attribute{
				Required: true,
			},
			"action_id": schema.Int64Attribute{
				Required: true,
			},
			"timing": schema.StringAttribute{
				Required: true,
			},
			"name": schema.StringAttribute{
				Required: true,
			},
			"run_as": schema.StringAttribute{
				Required: true,
			},
			"script": schema.StringAttribute{
				Required: true,
			},
			"sequence": schema.Int64Attribute{
				Computed: true,
			},
			"servers": schema.ListAttribute{
				ElementType: types.Int64Type,
				Optional:    true,
			},
			"created_at": schema.StringAttribute{
				Computed: true,
			},
			"updated_at": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

func (r *EnvoyerHookResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	providerConfig, ok := req.ProviderData.(*providerConfig)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Provider Configure Type",
			fmt.Sprintf("Expected *providerConfig, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	r.client = providerConfig.Envoyer
}

func (r *EnvoyerHookResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan EnvoyerHookResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	payload := envoyer_client.CreateHookRequest{
		ProjectID: plan.ProjectID.ValueInt64(),
		ActionID:  plan.ActionID.ValueInt64(),
		Timing:    plan.Timing.ValueString(),
		Name:      plan.Name.ValueString(),
		RunAs:     plan.RunAs.ValueString(),
		Script:    plan.Script.ValueString(),
		Servers:   extractIntSliceFromList(plan.Servers),
	}

	hook, err := r.client.CreateHook(ctx, payload)
	if err != nil {
		resp.Diagnostics.AddError("Error creating hook", err.Error())
		return
	}

	plan.ID = types.Int64Value(hook.ID)
	plan.Sequence = types.Int64Value(hook.Sequence)
	plan.CreatedAt = types.StringValue(hook.CreatedAt)
	plan.UpdatedAt = types.StringValue(hook.UpdatedAt)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *EnvoyerHookResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state EnvoyerHookResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	hook, err := r.client.GetHook(ctx, int(state.ProjectID.ValueInt64()), int(state.ID.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError("Error reading hook", err.Error())
		return
	}

	state.ActionID = types.Int64Value(hook.ActionID)
	state.Timing = types.StringValue(hook.Timing)
	state.Name = types.StringValue(hook.Name)
	state.RunAs = types.StringValue(hook.RunAs)
	state.Script = types.StringValue(hook.Script)
	state.Sequence = types.Int64Value(hook.Sequence)
	state.CreatedAt = types.StringValue(hook.CreatedAt)
	state.UpdatedAt = types.StringValue(hook.UpdatedAt)
	state.Servers = convertIntSliceToList(hook.Servers)

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *EnvoyerHookResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan EnvoyerHookResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	payload := envoyer_client.UpdateHookRequest{
		Servers: extractIntSliceFromList(plan.Servers),
	}

	if err := r.client.UpdateHook(ctx, int(plan.ProjectID.ValueInt64()), int(plan.ID.ValueInt64()), payload); err != nil {
		resp.Diagnostics.AddError("Error updating hook", err.Error())
		return
	}

	// Refresh state
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *EnvoyerHookResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state EnvoyerHookResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.DeleteHook(ctx, int(state.ProjectID.ValueInt64()), int(state.ID.ValueInt64())); err != nil {
		resp.Diagnostics.AddError("Error deleting hook", err.Error())
		return
	}
	resp.State.RemoveResource(ctx)
}

func (r *EnvoyerHookResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Expected import id format: "<project_id>,<hook_id>"
	parts := strings.Split(req.ID, ",")
	if len(parts) != 2 {
		resp.Diagnostics.AddError("Unexpected Import Identifier",
			fmt.Sprintf("Expected format <project_id>,<hook_id>. Got: %s", req.ID))
		return
	}
	projectID, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		resp.Diagnostics.AddError("Error parsing project_id", err.Error())
		return
	}
	hookID, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		resp.Diagnostics.AddError("Error parsing hook_id", err.Error())
		return
	}

	resp.State.SetAttribute(ctx, path.Root("project_id"), projectID)
	resp.State.SetAttribute(ctx, path.Root("id"), hookID)
}

// Helper functions for converting Terraform list values.
func extractIntSliceFromList(list types.List) []int64 {
	var ints []int64
	var values []types.Int64
	diags := list.ElementsAs(context.Background(), &values, false)
	if diags.HasError() {
		return ints
	}
	for _, v := range values {
		ints = append(ints, v.ValueInt64())
	}
	return ints
}

func convertIntSliceToList(ints []int64) types.List {
	var values []int64
	for _, i := range ints {
		values = append(values, i)
	}
	list, _ := types.ListValueFrom(context.Background(), types.Int64Type, values)
	return list
}
