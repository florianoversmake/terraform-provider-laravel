package provider

import (
	"context"
	"fmt"
	"strconv"

	"terraform-provider-laravel/internal/forge_client"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ resource.Resource = &ForgeCertificateResource{}
var _ resource.ResourceWithImportState = &ForgeCertificateResource{}

// ForgeCertificateResource implements a Terraform resource for a Forge site.
type ForgeCertificateResource struct {
	client *forge_client.Client
}

// Resource model.
type ForgeCertificateResourceModel struct {
	ID            types.Int64  `tfsdk:"id"`
	ServerID      types.Int64  `tfsdk:"server_id"`
	SiteID        types.Int64  `tfsdk:"site_id"`
	Domain        types.String `tfsdk:"domain"`
	Key           types.String `tfsdk:"key"`
	Certificate   types.String `tfsdk:"certificate"`
	RequestStatus types.String `tfsdk:"request_status"`
	Existing      types.Bool   `tfsdk:"existing"`
	Active        types.Bool   `tfsdk:"active"`
	CreatedAt     types.Int64  `tfsdk:"created_at"`
}

func NewForgeCertificateResource() resource.Resource {
	return &ForgeCertificateResource{}
}

func (r *ForgeCertificateResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_forge_certificate"
}

func (r *ForgeCertificateResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Forge certificate resource. This resource allows you to manage SSL certificates in Forge.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Computed: true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"server_id": schema.Int64Attribute{
				Required:            true,
				MarkdownDescription: "The ID of the server the certificate is associated with.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"site_id": schema.Int64Attribute{
				Required:            true,
				MarkdownDescription: "The ID of the site the certificate is associated with.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"key": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The key of the certificate.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"certificate": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The certificate.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"domain": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The domain of the certificate.",
			},
			"request_status": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The request status of the certificate.",
			},
			"existing": schema.BoolAttribute{
				Computed:            true,
				MarkdownDescription: "Whether the certificate already exists.",
			},
			"active": schema.BoolAttribute{
				Computed:            true,
				Optional:            true,
				MarkdownDescription: "Whether the certificate is active.",
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.RequiresReplaceIfConfigured(),
				},
			},
			"created_at": schema.Int64Attribute{
				Computed: true,
			},
		},
	}
}

func (r *ForgeCertificateResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *ForgeCertificateResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan ForgeCertificateResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Build the CreateCertificateRequest payload.
	payload := forge_client.CreateCertificateRequest{
		Type:        "existing",
		Key:         plan.Key.ValueString(),
		Certificate: plan.Certificate.ValueString(),
	}

	// Call CreateRecipe on the client.
	certificate, err := r.client.CreateCertificate(ctx, int(plan.ServerID.ValueInt64()), int(plan.SiteID.ValueInt64()), payload)
	if err != nil {
		resp.Diagnostics.AddError("Error creating certificate", err.Error())
		return
	}

	if plan.Active.ValueBool() && !certificate.Active {
		err = r.client.ActivateCertificate(ctx, int(plan.ServerID.ValueInt64()), int(plan.SiteID.ValueInt64()), int(certificate.ID))
		if err != nil {
			resp.Diagnostics.AddError("Error activating certificate", err.Error())
			return
		}
		certificate.Active = true
	}

	// Update plan state with response values.
	plan.ID = types.Int64Value(certificate.ID)
	plan.CreatedAt = types.Int64Value(certificate.CreatedAt)
	plan.Domain = types.StringValue(certificate.Domain)
	plan.RequestStatus = types.StringValue(certificate.RequestStatus)
	plan.Existing = types.BoolValue(certificate.Existing)
	plan.Active = types.BoolValue(certificate.Active)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *ForgeCertificateResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state ForgeCertificateResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	certificate, err := r.client.GetCertificate(ctx, int(state.ServerID.ValueInt64()), int(state.SiteID.ValueInt64()), int(state.ID.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError("Error reading certificate", err.Error())
		return
	}

	state.ID = types.Int64Value(certificate.ID)
	state.CreatedAt = types.Int64Value(certificate.CreatedAt)
	state.Domain = types.StringValue(certificate.Domain)
	state.RequestStatus = types.StringValue(certificate.RequestStatus)
	state.Existing = types.BoolValue(certificate.Existing)
	state.Active = types.BoolValue(certificate.Active)

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *ForgeCertificateResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan ForgeCertificateResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *ForgeCertificateResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state ForgeCertificateResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteCertificate(ctx, int(state.ServerID.ValueInt64()), int(state.SiteID.ValueInt64()), int(state.ID.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError("Error deleting certificate", err.Error())
		return
	}
}

func (r *ForgeCertificateResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	parts := splitCompositeID(req.ID, 3)
	if parts == nil {
		resp.Diagnostics.AddError("Invalid import ID", "Expected format: server_id:site_id:certificate_id")
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
	certificateID, err := strconv.ParseInt(parts[2], 10, 64)
	if err != nil {
		resp.Diagnostics.AddError("Invalid certificate_id", err.Error())
		return
	}

	certificate, err := r.client.GetCertificate(ctx, int(serverID), int(siteID), int(certificateID))
	if err != nil {
		resp.Diagnostics.AddError("Error reading certificate", err.Error())
		return
	}
	var stateModel ForgeCertificateResourceModel
	stateModel.ID = types.Int64Value(certificateID)
	stateModel.ServerID = types.Int64Value(serverID)
	stateModel.SiteID = types.Int64Value(siteID)
	stateModel.Domain = types.StringValue(certificate.Domain)
	stateModel.RequestStatus = types.StringValue(certificate.RequestStatus)
	stateModel.Existing = types.BoolValue(certificate.Existing)
	stateModel.Active = types.BoolValue(certificate.Active)
	stateModel.CreatedAt = types.Int64Value(certificate.CreatedAt)
	stateModel.Key = types.StringValue("")
	stateModel.Certificate = types.StringValue("")

	diags := resp.State.Set(ctx, stateModel)
	resp.Diagnostics.Append(diags...)
}
