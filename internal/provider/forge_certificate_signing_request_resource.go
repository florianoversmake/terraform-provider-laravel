package provider

import (
	"context"
	"crypto/x509"
	"encoding/pem"
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

var _ resource.Resource = &ForgeCertificateSigningRequestResource{}
var _ resource.ResourceWithImportState = &ForgeCertificateSigningRequestResource{}

// ForgeCertificateSigningRequestResource implements a Terraform resource for a Forge site.
type ForgeCertificateSigningRequestResource struct {
	client *forge_client.Client
}

// Resource model.
type ForgeCertificateSigningRequestResourceModel struct {
	ID                        types.Int64  `tfsdk:"id"`
	ServerID                  types.Int64  `tfsdk:"server_id"`
	SiteID                    types.Int64  `tfsdk:"site_id"`
	Domain                    types.String `tfsdk:"domain"`
	Country                   types.String `tfsdk:"country"`
	State                     types.String `tfsdk:"state"`
	City                      types.String `tfsdk:"city"`
	Organization              types.String `tfsdk:"organization"`
	Department                types.String `tfsdk:"department"`
	CertificateSigningRequest types.String `tfsdk:"certificate_signing_request"`
	RequestStatus             types.String `tfsdk:"request_status"`
	Existing                  types.Bool   `tfsdk:"existing"`
	Active                    types.Bool   `tfsdk:"active"`
	CreatedAt                 types.Int64  `tfsdk:"created_at"`
}

func NewForgeCertificateSigningRequestResource() resource.Resource {
	return &ForgeCertificateSigningRequestResource{}
}

func (r *ForgeCertificateSigningRequestResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_forge_certificate_signing_request"
}

func (r *ForgeCertificateSigningRequestResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Forge certificate signing request resource. This resource allows you to manage SSL certificate signing requests in Forge.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Computed: true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"server_id": schema.Int64Attribute{
				Required:            true,
				MarkdownDescription: "The ID of the server the certificate signing request is associated with.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"site_id": schema.Int64Attribute{
				Required:            true,
				MarkdownDescription: "The ID of the site the certificate signing request is associated with.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"country": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The country of the certificate signing request.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"state": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The state of the certificate signing request.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"city": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The city of the certificate signing request.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"organization": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The organization of the certificate signing request.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"department": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The department of the certificate signing request.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"domain": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The domain of the certificate.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"certificate_signing_request": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The certificate signing request.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
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
				MarkdownDescription: "Whether the certificate is active.",
			},
			"created_at": schema.Int64Attribute{
				Computed: true,
			},
		},
	}
}

func (r *ForgeCertificateSigningRequestResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *ForgeCertificateSigningRequestResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan ForgeCertificateSigningRequestResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Build the CreateCertificateRequest payload.
	payload := forge_client.CreateCertificateRequest{
		Type:         "new",
		Domain:       plan.Domain.ValueString(),
		Country:      plan.Country.ValueString(),
		State:        plan.State.ValueString(),
		City:         plan.City.ValueString(),
		Organization: plan.Organization.ValueString(),
		Department:   plan.Department.ValueString(),
	}

	// Call CreateRecipe on the client.
	certificate, err := r.client.CreateCertificate(ctx, int(plan.ServerID.ValueInt64()), int(plan.SiteID.ValueInt64()), payload)
	if err != nil {
		resp.Diagnostics.AddError("Error creating certificate signing request", err.Error())
		return
	}

	certificateSigningRequest, err := r.client.GetCertificateCSR(ctx, int(plan.ServerID.ValueInt64()), int(plan.SiteID.ValueInt64()), int(certificate.ID))
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving certificate signing request", err.Error())
		return
	}

	// Update plan state with response values.
	plan.ID = types.Int64Value(certificate.ID)
	plan.CreatedAt = types.Int64Value(certificate.CreatedAt)
	plan.Domain = types.StringValue(certificate.Domain)
	plan.RequestStatus = types.StringValue(certificate.RequestStatus)
	plan.Existing = types.BoolValue(certificate.Existing)
	plan.Active = types.BoolValue(certificate.Active)
	plan.CertificateSigningRequest = types.StringValue(certificateSigningRequest)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *ForgeCertificateSigningRequestResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state ForgeCertificateSigningRequestResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	certificate, err := r.client.GetCertificate(ctx, int(state.ServerID.ValueInt64()), int(state.SiteID.ValueInt64()), int(state.ID.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError("Error reading certificate signing request", err.Error())
		return
	}

	certificateSigningRequest, err := r.client.GetCertificateCSR(ctx, int(state.ServerID.ValueInt64()), int(state.SiteID.ValueInt64()), int(state.ID.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving certificate signing request", err.Error())
		return
	}

	// Decode the CSR and extract subject information
	country, stateValue, city, organization, department, err := extractCSRInfo(certificateSigningRequest)
	if err != nil {
		resp.Diagnostics.AddWarning(
			"Error parsing CSR",
			fmt.Sprintf("Could not extract information from CSR: %s", err.Error()),
		)
	} else {
		// Only update these values if we successfully parsed the CSR
		state.Country = types.StringValue(country)
		state.State = types.StringValue(stateValue)
		state.City = types.StringValue(city)
		state.Organization = types.StringValue(organization)
		state.Department = types.StringValue(department)
	}

	state.ID = types.Int64Value(certificate.ID)
	state.CreatedAt = types.Int64Value(certificate.CreatedAt)
	state.Domain = types.StringValue(certificate.Domain)
	state.RequestStatus = types.StringValue(certificate.RequestStatus)
	state.Existing = types.BoolValue(certificate.Existing)
	state.Active = types.BoolValue(certificate.Active)
	state.CertificateSigningRequest = types.StringValue(certificateSigningRequest)

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *ForgeCertificateSigningRequestResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan ForgeCertificateSigningRequestResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *ForgeCertificateSigningRequestResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state ForgeCertificateSigningRequestResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteCertificate(ctx, int(state.ServerID.ValueInt64()), int(state.SiteID.ValueInt64()), int(state.ID.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError("Error deleting certificate signing request", err.Error())
		return
	}
}

func (r *ForgeCertificateSigningRequestResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
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

	certificateSigningRequest, err := r.client.GetCertificateCSR(ctx, int(serverID), int(siteID), int(certificateID))
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving certificate signing request", err.Error())
		return
	}

	var stateModel ForgeCertificateSigningRequestResourceModel
	stateModel.ID = types.Int64Value(certificateID)
	stateModel.ServerID = types.Int64Value(serverID)
	stateModel.SiteID = types.Int64Value(siteID)
	stateModel.Domain = types.StringValue(certificate.Domain)
	stateModel.RequestStatus = types.StringValue(certificate.RequestStatus)
	stateModel.Existing = types.BoolValue(certificate.Existing)
	stateModel.Active = types.BoolValue(certificate.Active)
	stateModel.CreatedAt = types.Int64Value(certificate.CreatedAt)
	stateModel.CertificateSigningRequest = types.StringValue(certificateSigningRequest)

	// Decode the CSR and extract subject information
	country, stateValue, city, organization, department, err := extractCSRInfo(certificateSigningRequest)
	if err != nil {
		resp.Diagnostics.AddWarning(
			"Error parsing CSR",
			fmt.Sprintf("Could not extract information from CSR: %s", err.Error()),
		)
	} else {
		// Only set these values if we successfully parsed the CSR
		stateModel.Country = types.StringValue(country)
		stateModel.State = types.StringValue(stateValue)
		stateModel.City = types.StringValue(city)
		stateModel.Organization = types.StringValue(organization)
		stateModel.Department = types.StringValue(department)
	}

	diags := resp.State.Set(ctx, stateModel)
	resp.Diagnostics.Append(diags...)
}

func extractCSRInfo(csrPEM string) (country, state, city, organization, department string, err error) {
	block, _ := pem.Decode([]byte(csrPEM))
	if block == nil || block.Type != "CERTIFICATE REQUEST" {
		return "", "", "", "", "", fmt.Errorf("failed to decode PEM block containing CSR")
	}

	csr, err := x509.ParseCertificateRequest(block.Bytes)
	if err != nil {
		return "", "", "", "", "", fmt.Errorf("failed to parse certificate request: %w", err)
	}

	subject := csr.Subject

	// Extract fields from the subject, handling potential empty slices
	if len(subject.Country) > 0 {
		country = subject.Country[0]
	}
	if len(subject.Province) > 0 {
		state = subject.Province[0]
	}
	if len(subject.Locality) > 0 {
		city = subject.Locality[0]
	}
	if len(subject.Organization) > 0 {
		organization = subject.Organization[0]
	}
	if len(subject.OrganizationalUnit) > 0 {
		department = subject.OrganizationalUnit[0]
	}

	return country, state, city, organization, department, nil
}
