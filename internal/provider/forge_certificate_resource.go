package provider

import (
	"context"
	"fmt"
	"strconv"

	"terraform-provider-laravel/internal/forge_client"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ resource.Resource = &ForgeCertificateResource{}
var _ resource.ResourceWithImportState = &ForgeCertificateResource{}

// ForgeCertificateResource implements a Terraform resource for a Forge certificate.
type ForgeCertificateResource struct {
	client *forge_client.Client
}

// Resource model.
type ForgeCertificateResourceModel struct {
	ID            types.Int64  `tfsdk:"id"`
	ServerID      types.Int64  `tfsdk:"server_id"`
	SiteID        types.Int64  `tfsdk:"site_id"`
	DomainID      types.Int64  `tfsdk:"domain_id"`
	Key           types.String `tfsdk:"key"`
	Certificate   types.String `tfsdk:"certificate"`
	RequestStatus types.String `tfsdk:"request_status"`
	Status        types.String `tfsdk:"status"`
	CreatedAt     types.String `tfsdk:"created_at"`
}

func NewForgeCertificateResource() resource.Resource {
	return &ForgeCertificateResource{}
}

func (r *ForgeCertificateResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_forge_certificate"
}

func (r *ForgeCertificateResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Forge certificate resource. This resource allows you to manage SSL certificates in Forge. In the new API, certificates are managed per domain record.",
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
			"domain_id": schema.Int64Attribute{
				Required:            true,
				MarkdownDescription: "The ID of the domain record the certificate is associated with.",
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
			"request_status": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The request status of the certificate.",
			},
			"status": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The status of the certificate.",
			},
			"created_at": schema.StringAttribute{
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

	// Call CreateCertificate on the client with domainID.
	certificate, err := r.client.CreateCertificate(ctx, int(plan.ServerID.ValueInt64()), int(plan.SiteID.ValueInt64()), int(plan.DomainID.ValueInt64()), payload)
	if err != nil {
		resp.Diagnostics.AddError("Error creating certificate", err.Error())
		return
	}

	// Update plan state with response values.
	plan.ID = types.Int64Value(certificate.ID)
	plan.CreatedAt = types.StringValue(certificate.CreatedAt)
	plan.RequestStatus = types.StringValue(certificate.RequestStatus)
	plan.Status = types.StringValue(certificate.Status)

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

	certificate, err := r.client.GetCertificate(ctx, int(state.ServerID.ValueInt64()), int(state.SiteID.ValueInt64()), int(state.DomainID.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError("Error reading certificate", err.Error())
		return
	}

	state.ID = types.Int64Value(certificate.ID)
	state.CreatedAt = types.StringValue(certificate.CreatedAt)
	state.RequestStatus = types.StringValue(certificate.RequestStatus)
	state.Status = types.StringValue(certificate.Status)

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

	err := r.client.DeleteCertificate(ctx, int(state.ServerID.ValueInt64()), int(state.SiteID.ValueInt64()), int(state.DomainID.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError("Error deleting certificate", err.Error())
		return
	}
}

func (r *ForgeCertificateResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	parts := splitCompositeID(req.ID, 3)
	if parts == nil {
		resp.Diagnostics.AddError("Invalid import ID", "Expected format: server_id:site_id:domain_id")
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
	domainID, err := strconv.ParseInt(parts[2], 10, 64)
	if err != nil {
		resp.Diagnostics.AddError("Invalid domain_id", err.Error())
		return
	}

	certificate, err := r.client.GetCertificate(ctx, int(serverID), int(siteID), int(domainID))
	if err != nil {
		resp.Diagnostics.AddError("Error reading certificate", err.Error())
		return
	}
	var stateModel ForgeCertificateResourceModel
	stateModel.ID = types.Int64Value(certificate.ID)
	stateModel.ServerID = types.Int64Value(serverID)
	stateModel.SiteID = types.Int64Value(siteID)
	stateModel.DomainID = types.Int64Value(domainID)
	stateModel.RequestStatus = types.StringValue(certificate.RequestStatus)
	stateModel.Status = types.StringValue(certificate.Status)
	stateModel.CreatedAt = types.StringValue(certificate.CreatedAt)
	stateModel.Key = types.StringValue("")
	stateModel.Certificate = types.StringValue("")

	diags := resp.State.Set(ctx, stateModel)
	resp.Diagnostics.Append(diags...)
}
