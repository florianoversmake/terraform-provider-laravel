package provider

import (
	"context"
	"fmt"
	"strconv"
	"terraform-provider-laravel/internal/forge_client"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure ServerResource satisfies the expected interfaces.
var _ resource.Resource = &ForgeServerResource{}
var _ resource.ResourceWithImportState = &ForgeServerResource{}

// Parameters
// Key	Description
// ubuntu_version	The version of Ubuntu to create the server with. Valid values are "20.04", "22.04", and "24.04". "24.04" is used by default if no value is defined. It is recommended to always specify a version as the default may change at any time.
// type	The type of server to create. Valid values are app, web, loadbalancer, cache, database, worker, meilisearch. app is used by default if no value is defined.
// provider	The server provider. Valid values are ocean2 for Digital Ocean, akamai (Linode), vultr2, aws, hetzner and custom.
// size	The instance type (aws)
// disk_size	The size of the disk in GB. Valid when the provider is aws. Minimum of 8GB. Example: 20.
// circle	The ID of a circle to create the server within.
// credential_id	This is only required when the provider is not custom.
// region	The name of the region where the server will be created. This value is not required you are building a Custom VPS server. Valid region identifiers.
// ip_address	The IP Address of the server. Only required when the provider is custom.
// private_ip_address	The Private IP Address of the server. Only required when the provider is custom.
// php_version	Valid values are php84, php83, php82, php81, php80, php74, php73,php72,php82, php70, and php56.
// database	The name of the database Forge should create when building the server. If omitted, forge will be used.
// database_type	Valid values are mysql8, mariadb106, mariadb1011, mariadb114, postgres, postgres13, postgres14, postgres15, postgres16 or postgres17.
// network	An array of server IDs that the server should be able to connect to.
// recipe_id	An optional ID of a recipe to run after provisioning.
// aws_vpc_id	ID of the existing VPC
// aws_subnet_id	ID of the existing subnet
// aws_vpc_name	When creating a new one
// hetzner_network_id	ID of the existing VPC
// ocean2_vpc_uuid	UUID of the existing VPC
// ocean2_vpc_name	When creating a new one
// vultr2_network_id	ID of the existing private network
// vultr2_network_name	When creating a new one

// ForgeServerResourceModel defines the schema data model for the server.
type ForgeServerResourceModel struct {
	ID             types.Int64  `tfsdk:"id"`
	ServerProvider types.String `tfsdk:"server_provider"`

	UbuntuVersion types.String `tfsdk:"ubuntu_version"` // readonly after creation
	Name          types.String `tfsdk:"name"`
	CredentialID  types.Int64  `tfsdk:"credential_id"` // readonly after creation
	Type          types.String `tfsdk:"type"`          // readonly after creation
	Circle        types.Int64  `tfsdk:"circle"`        // readonly after creation
	PhpVersion    types.String `tfsdk:"php_version"`   // readonly after creation
	DatabaseType  types.String `tfsdk:"database_type"` // readonly after creation
	Database      types.String `tfsdk:"database"`      // readonly after creation
	Network       types.List   `tfsdk:"network"`
	RecipeID      types.Int64  `tfsdk:"recipe_id"`

	// required for create only when provider is custom
	IpAddress        types.String `tfsdk:"ip_address"`
	PrivateIpAddress types.String `tfsdk:"private_ip_address"`
	SshPort          types.Int32  `tfsdk:"ssh_port"`

	// Generic for all providers except custom
	Region types.String `tfsdk:"region"` // readonly after creation
	Size   types.String `tfsdk:"size"`   // readonly after creation

	// AWS
	DiskSize    types.Int32  `tfsdk:"disk_size"`     // readonly after creation
	Identifier  types.String `tfsdk:"identifier"`    // known after creation
	AwsVpcID    types.String `tfsdk:"aws_vpc_id"`    // readonly after creation
	AwsSubnetID types.String `tfsdk:"aws_subnet_id"` // readonly after creation
	AwsVpcName  types.String `tfsdk:"aws_vpc_name"`  // readonly after creation

	// computed
	LocalPublicKey      types.String `tfsdk:"local_public_key"`     // readonly, known after creation
	Revoked             types.Bool   `tfsdk:"revoked"`              // readonly, known after creation
	IsReady             types.Bool   `tfsdk:"is_ready"`             // readonly, known after creation
	SudoPassword        types.String `tfsdk:"sudo_password"`        // readonly, known after creation
	DatabasePassword    types.String `tfsdk:"database_password"`    // readonly, known after creation
	MeilisearchPassword types.String `tfsdk:"meilisearch_password"` // readonly, known after creation
	ProvisionCommand    types.String `tfsdk:"provision_command"`    // readonly, known after creation
}

type Tag struct {
	ID        types.Int64  `tfsdk:"id"`
	Name      types.String `tfsdk:"name"`
	CreatedAt types.String `tfsdk:"created_at"`
}

// ForgeServerResource implements the resource.Resource interface.
type ForgeServerResource struct {
	client *forge_client.Client
}

// NewForgeServerResource is a helper function to instantiate the resource.
func NewForgeServerResource() resource.Resource {
	return &ForgeServerResource{}
}

func (r *ForgeServerResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_forge_server"
}

func (r *ForgeServerResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Computed: true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"server_provider": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"ubuntu_version": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString("24.04"),
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplaceIfConfigured(),
				},
			},
			"name": schema.StringAttribute{
				Required: true,
			},
			"credential_id": schema.Int64Attribute{
				Optional: true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplaceIfConfigured(),
				},
			},
			"type": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString("app"),
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplaceIfConfigured(),
				},
			},
			"circle": schema.Int64Attribute{
				Optional: true,
				Computed: true,
				Default:  int64default.StaticInt64(0),
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplaceIfConfigured(),
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"php_version": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString("php82"),
			},
			"database": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString(""),
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplaceIfConfigured(),
				},
				MarkdownDescription: "The name of the database Forge should create when building the server. If omitted, forge will be used.",
			},
			"database_type": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString(""),
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplaceIfConfigured(),
				},
				MarkdownDescription: "Valid values are mysql8, mariadb106, mariadb1011, mariadb114, postgres, postgres13, postgres14, postgres15, postgres16 or postgres17.",
			},
			"network": schema.ListAttribute{
				MarkdownDescription: "An array of server IDs that the server should be able to connect to.",
				ElementType:         types.Int64Type,
				Computed:            true,
				Optional:            true,
				Default:             listdefault.StaticValue(types.ListValueMust(types.Int64Type, []attr.Value{})),
			},
			"recipe_id": schema.Int64Attribute{
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplaceIfConfigured(),
					int64planmodifier.UseStateForUnknown(),
				},
				Optional:            true,
				MarkdownDescription: "An optional ID of a recipe to run after provisioning.",
			},
			"ip_address": schema.StringAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					stringplanmodifier.RequiresReplaceIfConfigured(),
				},
			},
			"private_ip_address": schema.StringAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					stringplanmodifier.RequiresReplaceIfConfigured(),
				},
			},
			"ssh_port": schema.Int32Attribute{
				Optional: true,
				Computed: true,
				Default:  int32default.StaticInt32(22),
				PlanModifiers: []planmodifier.Int32{
					int32planmodifier.RequiresReplaceIfConfigured(),
				},
			},
			"identifier": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"size": schema.StringAttribute{
				Optional: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"region": schema.StringAttribute{
				Optional: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"disk_size": schema.Int32Attribute{
				Optional: true,
				PlanModifiers: []planmodifier.Int32{
					int32planmodifier.RequiresReplaceIfConfigured(),
				},
			},
			"aws_vpc_id": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The ID of the VPC to launch the server in.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplaceIfConfigured(),
				},
			},
			"aws_subnet_id": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The ID of the subnet to launch the server in.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplaceIfConfigured(),
				},
			},
			"aws_vpc_name": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "When creating a new one",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplaceIfConfigured(),
				},
			},
			"local_public_key": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"revoked": schema.BoolAttribute{
				Computed: true,
				Optional: true,
				Default:  booldefault.StaticBool(false),
			},
			"is_ready": schema.BoolAttribute{
				Computed: true,
			},
			"sudo_password": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"database_password": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"meilisearch_password": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"provision_command": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *ForgeServerResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r ForgeServerResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var data ForgeServerResourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *ForgeServerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan ForgeServerResourceModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if plan.ServerProvider.ValueString() == "custom" {
		if plan.IpAddress.IsNull() {
			resp.Diagnostics.AddAttributeError(path.Root("ip_address"), "missing ip_address", "ip_address is required when provider is custom")
		}
		if plan.PrivateIpAddress.IsNull() {
			resp.Diagnostics.AddAttributeError(path.Root("private_ip_address"), "missing private_ip_address", "private_ip_address is required when provider is custom")
		}
	} else {
		if plan.CredentialID.IsNull() {
			resp.Diagnostics.AddAttributeError(path.Root("credential_id"), "missing credential_id", "credential_id is required when provider is not custom")
		}
		if !plan.IpAddress.IsNull() {
			resp.Diagnostics.AddAttributeError(path.Root("ip_address"), "invalid ip_address", "ip_address is not allowed when provider is not custom")
		}
		if !plan.PrivateIpAddress.IsNull() {
			resp.Diagnostics.AddAttributeError(path.Root("private_ip_address"), "invalid private_ip_address", "private_ip_address is not allowed when provider is not custom")
		}

		if plan.ServerProvider.ValueString() == "aws" {
			if plan.AwsVpcID.IsNull() {
				resp.Diagnostics.AddAttributeError(path.Root("aws_vpc_id"), "missing aws_vpc_id", "aws_vpc_id is required when provider is aws")
			}
			if plan.AwsSubnetID.IsNull() {
				resp.Diagnostics.AddAttributeError(path.Root("aws_subnet_id"), "missing aws_subnet_id", "aws_subnet_id is required when provider is aws")
			}
		} else {
			if !plan.AwsVpcID.IsNull() {
				resp.Diagnostics.AddAttributeError(path.Root("aws_vpc_id"), "invalid aws_vpc_id", "aws_vpc_id is not allowed when provider is not aws")
			}
			if !plan.AwsSubnetID.IsNull() {
				resp.Diagnostics.AddAttributeError(path.Root("aws_subnet_id"), "invalid aws_subnet_id", "aws_subnet_id is not allowed when provider is not aws")
			}
		}
	}

	networkElements := make([]int64, 0)
	for _, v := range plan.Network.Elements() {
		networkElements = append(networkElements, (v.(types.Int64)).ValueInt64())
	}

	payload := forge_client.CreateServerRequest{
		UbuntuVersion:    plan.UbuntuVersion.ValueString(),
		Name:             plan.Name.ValueString(),
		Type:             plan.Type.ValueString(),
		Provider:         plan.ServerProvider.ValueString(),
		CredentialID:     plan.CredentialID.ValueInt64Pointer(),
		Circle:           plan.Circle.ValueInt64Pointer(),
		PHPVersion:       plan.PhpVersion.ValueString(),
		DatabaseType:     plan.DatabaseType.ValueStringPointer(),
		Database:         plan.Database.ValueStringPointer(),
		Network:          networkElements,
		RecipeID:         plan.RecipeID.ValueInt64Pointer(),
		IPAddress:        plan.IpAddress.ValueStringPointer(),
		PrivateIPAddress: plan.PrivateIpAddress.ValueStringPointer(),
		SSHPort:          plan.SshPort.ValueInt32Pointer(),
		Region:           plan.Region.ValueStringPointer(),
		Size:             plan.Size.ValueStringPointer(),
		DiskSize:         plan.DiskSize.ValueInt32Pointer(),
		AWSVPCID:         plan.AwsVpcID.ValueStringPointer(),
		AWSSubnetID:      plan.AwsSubnetID.ValueStringPointer(),
		AWSVPCName:       plan.AwsVpcName.ValueStringPointer(),
	}

	response, err := r.client.CreateServer(ctx, payload)
	if err != nil {
		resp.Diagnostics.AddError("Error creating server", err.Error())
		return
	}

	// wait for server to be ready
	err = r.client.WaitForServerToBeReady(ctx, int(response.Server.ID))
	if err != nil {
		resp.Diagnostics.AddError("Error waiting for server to be ready", err.Error())
		return
	}

	plan.ID = types.Int64Value(int64(response.Server.ID))
	plan.IpAddress = types.StringPointerValue(response.Server.IPAddress)
	plan.PrivateIpAddress = types.StringPointerValue(response.Server.PrivateIPAddress)
	plan.SshPort = types.Int32Value(int32(response.Server.SSHPort))
	plan.Identifier = types.StringValue(response.Server.Identifier)
	plan.LocalPublicKey = types.StringValue(response.Server.LocalPublicKey)
	plan.Revoked = types.BoolValue(response.Server.Revoked)
	plan.IsReady = types.BoolValue(response.Server.IsReady)
	plan.SudoPassword = types.StringValue(response.SudoPassword)
	plan.DatabasePassword = types.StringPointerValue(response.DatabasePassword)
	plan.MeilisearchPassword = types.StringPointerValue(response.MeilisearchPassword)
	plan.ProvisionCommand = types.StringPointerValue(response.ProvisionCommand)

	regionId, err := r.client.GetRegionIDByName(ctx, plan.ServerProvider.ValueString(), response.Server.Region)
	if err != nil {
		resp.Diagnostics.AddError("Error getting region ID", err.Error())
		return
	}

	plan.Region = types.StringValue(regionId)

	sizeSize, err := r.client.GetRegionSizeSizeByID(ctx, plan.ServerProvider.ValueString(), regionId, response.Server.Size)
	if err != nil {
		resp.Diagnostics.AddError("Error getting size size", err.Error())
		return
	}

	plan.Size = types.StringValue(sizeSize)

	listValue, diags := types.ListValueFrom(ctx, types.Int64Type, response.Server.Network)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}
	plan.Network = listValue

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *ForgeServerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state ForgeServerResourceModel
	diags := req.State.Get(ctx, &state)
	//resp.Diagnostics.AddWarning("Dumping plan1", fmt.Sprintf("%+v", state))
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	server, err := r.client.GetServer(ctx, int(state.ID.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError("Error reading server", err.Error())
		return
	}

	state.ServerProvider = types.StringValue(server.Provider)
	state.UbuntuVersion = types.StringValue(server.UbuntuVersion)
	state.Name = types.StringValue(server.Name)
	state.CredentialID = types.Int64Value(int64(server.CredentialID))
	state.Type = types.StringValue(server.Type)
	state.PhpVersion = types.StringValue(server.PHPVersion)
	state.DatabaseType = types.StringValue(server.DatabaseType)
	state.IpAddress = types.StringPointerValue(server.IPAddress)
	state.PrivateIpAddress = types.StringPointerValue(server.PrivateIPAddress)
	state.SshPort = types.Int32Value(int32(server.SSHPort))
	state.LocalPublicKey = types.StringValue(server.LocalPublicKey)
	state.Revoked = types.BoolValue(server.Revoked)
	state.IsReady = types.BoolValue(server.IsReady)
	state.SudoPassword = types.StringValue(state.SudoPassword.ValueString())
	state.Identifier = types.StringValue(server.Identifier)

	regionId, err := r.client.GetRegionIDByName(ctx, state.ServerProvider.ValueString(), server.Region)
	if err != nil {
		resp.Diagnostics.AddError("Error getting region ID", err.Error())
		return
	}

	state.Region = types.StringValue(regionId)

	sizeSize, err := r.client.GetRegionSizeSizeByID(ctx, state.ServerProvider.ValueString(), regionId, server.Size)
	if err != nil {
		resp.Diagnostics.AddError("Error getting size size", err.Error())
		return
	}

	state.Size = types.StringValue(sizeSize)

	state.DatabasePassword = types.StringPointerValue(state.DatabasePassword.ValueStringPointer())
	state.Circle = types.Int64PointerValue(state.Circle.ValueInt64Pointer())
	state.RecipeID = types.Int64PointerValue(state.RecipeID.ValueInt64Pointer())

	listValue, diags := types.ListValueFrom(ctx, types.Int64Type, server.Network)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}
	state.Network = listValue

	// state.DiskSize = types.Int32PointerValue(state.DiskSize)
	// state.AwsVpcID = types.StringPointerValue(state.AwsVpcID.ValueStringPointer())
	// state.AwsSubnetID = types.StringPointerValue(state.AwsSubnetID.ValueStringPointer())
	// state.AwsVpcName = types.StringPointerValue(state.AwsVpcName.ValueStringPointer())

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)

	// dump plan to console
	//resp.Diagnostics.AddWarning("Dumping plan2", fmt.Sprintf("%+v", state))
	//resp.Diagnostics.AddWarning("Dumping state", fmt.Sprintf("%+v", resp.State))
	//resp.Diagnostics.AddWarning("Dumping server", fmt.Sprintf("%+v", server))
}

func (r *ForgeServerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan ForgeServerResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)

	// dump plan to console
	//resp.Diagnostics.AddWarning("Dumping plan", fmt.Sprintf("%+v", plan))
	//resp.Diagnostics.AddWarning("Dumping state", fmt.Sprintf("%+v", resp.State))
}

func (r *ForgeServerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state ForgeServerResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// err := r.client.DeleteServer(ctx, int(state.ID.ValueInt64()))
	// if err != nil {
	// 	resp.Diagnostics.AddError("Error deleting server", err.Error())
	// 	return
	// }
}

func (r *ForgeServerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	id, err := strconv.ParseInt(req.ID, 10, 64)
	if err != nil {
		resp.Diagnostics.AddError("Error parsing ID", err.Error())
		return
	}

	server, err := r.client.GetServer(ctx, int(id))
	if err != nil {
		resp.Diagnostics.AddError("Error reading server", err.Error())
		return
	}

	networkListValue, diags := types.ListValueFrom(ctx, types.Int64Type, server.Network)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	state := ForgeServerResourceModel{
		ID:                  types.Int64Value(int64(server.ID)),
		ServerProvider:      types.StringValue(server.Provider),
		UbuntuVersion:       types.StringValue(server.UbuntuVersion),
		Name:                types.StringValue(server.Name),
		CredentialID:        types.Int64Value(int64(server.CredentialID)),
		Type:                types.StringValue(server.Type),
		Circle:              types.Int64Value(0),
		PhpVersion:          types.StringValue(server.PHPVersion),
		DatabaseType:        types.StringValue(server.DatabaseType),
		Database:            types.StringValue(""),
		Network:             networkListValue,
		RecipeID:            types.Int64Null(),
		IpAddress:           types.StringPointerValue(server.IPAddress),
		PrivateIpAddress:    types.StringPointerValue(server.PrivateIPAddress),
		SshPort:             types.Int32Value(int32(server.SSHPort)),
		Region:              types.StringValue(server.Region),
		Size:                types.StringValue(server.Size),
		DiskSize:            types.Int32Null(),
		AwsVpcID:            types.StringNull(),
		AwsSubnetID:         types.StringNull(),
		AwsVpcName:          types.StringNull(),
		LocalPublicKey:      types.StringValue(server.LocalPublicKey),
		Revoked:             types.BoolValue(server.Revoked),
		IsReady:             types.BoolValue(server.IsReady),
		SudoPassword:        types.StringValue(""),
		DatabasePassword:    types.StringValue(""),
		MeilisearchPassword: types.StringValue(""),
		ProvisionCommand:    types.StringValue(""),
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
