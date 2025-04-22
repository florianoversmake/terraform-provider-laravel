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

var _ resource.Resource = &ForgeSSHKeyResource{}
var _ resource.ResourceWithImportState = &ForgeSSHKeyResource{}

// ForgeSSHKeyResource implements a Terraform resource for a Forge site.
type ForgeSSHKeyResource struct {
	client *forge_client.Client
}

// Resource model.
type ForgeSSHKeyResourceModel struct {
	ID        types.Int64  `tfsdk:"id"`
	ServerID  types.Int64  `tfsdk:"server_id"`
	Key       types.String `tfsdk:"key"`
	Username  types.String `tfsdk:"username"`
	Name      types.String `tfsdk:"name"`
	Status    types.String `tfsdk:"status"`
	CreatedAt types.String `tfsdk:"created_at"`
}

func NewForgeSSHKeyResource() resource.Resource {
	return &ForgeSSHKeyResource{}
}

func (r *ForgeSSHKeyResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_forge_ssh_key"
}

func (r *ForgeSSHKeyResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Forge SSH key resource. This resource allows you to manage SSH keys in Forge.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Computed: true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"server_id": schema.Int64Attribute{
				Required:            true,
				MarkdownDescription: "The ID of the server to which the SSH key will be added.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplaceIfConfigured(),
				},
			},
			"key": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The public SSH key to be added.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplaceIfConfigured(),
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The name of the SSH key.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplaceIfConfigured(),
				},
			},
			"username": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("forge"),
				MarkdownDescription: "The username associated with the SSH key. If not provided, defaults to 'forge'.",
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

func (r *ForgeSSHKeyResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *ForgeSSHKeyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan ForgeSSHKeyResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	trimmedKey := strings.Trim(plan.Key.ValueString(), "\n")
	if trimmedKey == "" {
		resp.Diagnostics.AddError("Invalid SSH key", "The SSH key cannot be empty.")
		return
	}

	// Build the CreateSSHKeyRequest payload.
	payload := forge_client.CreateSSHKeyRequest{
		Name:     plan.Name.ValueString(),
		Key:      trimmedKey,
		Username: plan.Username.ValueString(),
	}

	// Call CreateSSHKey on the client using the provided server_id.
	sshKey, err := r.client.CreateSSHKey(ctx, int(plan.ServerID.ValueInt64()), payload)
	if err != nil {
		resp.Diagnostics.AddError("Error creating ssh key", err.Error())
		return
	}

	// Update plan state with response values.
	plan.ID = types.Int64Value(sshKey.ID)
	plan.Name = types.StringValue(sshKey.Name)
	plan.Username = types.StringValue(sshKey.Username)
	plan.Status = types.StringValue(sshKey.Status)
	plan.CreatedAt = types.StringValue(sshKey.CreatedAt)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *ForgeSSHKeyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state ForgeSSHKeyResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	sshKey, err := r.client.GetSSHKey(ctx, int(state.ServerID.ValueInt64()), int(state.ID.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError("Error reading ssh key", err.Error())
		return
	}

	state.Name = types.StringValue(sshKey.Name)
	state.Username = types.StringValue(sshKey.Username)
	state.ServerID = types.Int64Value(state.ServerID.ValueInt64())
	state.ID = types.Int64Value(sshKey.ID)
	state.Status = types.StringValue(sshKey.Status)
	state.CreatedAt = types.StringValue(sshKey.CreatedAt)

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *ForgeSSHKeyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan ForgeSSHKeyResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// We need to delete the existing SSH key and create a new one with the updated values.

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *ForgeSSHKeyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state ForgeSSHKeyResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteSSHKey(ctx, int(state.ServerID.ValueInt64()), int(state.ID.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError("Error deleting ssh key", err.Error())
		return
	}
}

func (r *ForgeSSHKeyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Expect import ID in format "server_id:key_id"
	parts := splitCompositeID(req.ID, 2)
	if parts == nil {
		resp.Diagnostics.AddError("Invalid import ID", "Expected format: server_id:key_id")
		return
	}
	serverID, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		resp.Diagnostics.AddError("Invalid server_id", err.Error())
		return
	}
	keyID, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		resp.Diagnostics.AddError("Invalid key_id", err.Error())
		return
	}

	key, err := r.client.GetSSHKey(ctx, int(serverID), int(keyID))
	if err != nil {
		resp.Diagnostics.AddError("Error reading ssh key", err.Error())
		return
	}

	var stateModel ForgeSSHKeyResourceModel
	stateModel.ID = types.Int64Value(key.ID)
	stateModel.ServerID = types.Int64Value(serverID)
	stateModel.Key = types.StringUnknown()
	stateModel.Name = types.StringValue(key.Name)
	stateModel.Username = types.StringValue(key.Username)
	stateModel.Status = types.StringValue(key.Status)
	stateModel.CreatedAt = types.StringValue(key.CreatedAt)

	diags := resp.State.Set(ctx, stateModel)
	resp.Diagnostics.Append(diags...)
}
