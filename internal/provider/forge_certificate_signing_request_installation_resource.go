package provider

import (
	"context"
	"fmt"

	"terraform-provider-laravel/internal/forge_client"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ resource.Resource = &ForgeCertificateSigningRequestInstallationResource{}

// ForgeCertificateSigningRequestInstallationResource implements a Terraform resource for a Forge site.
type ForgeCertificateSigningRequestInstallationResource struct {
	client *forge_client.Client
}

// Resource model.
type ForgeCertificateSigningRequestInstallationResourceModel struct {
	CertificateSigningRequestID types.Int64  `tfsdk:"certificate_signing_request_id"`
	ServerID                    types.Int64  `tfsdk:"server_id"`
	SiteID                      types.Int64  `tfsdk:"site_id"`
	Certificate                 types.String `tfsdk:"certificate"`
	AddIntermediates            types.Bool   `tfsdk:"add_intermediates"`
}

func NewForgeCertificateSigningRequestInstallationResource() resource.Resource {
	return &ForgeCertificateSigningRequestInstallationResource{}
}

func (r *ForgeCertificateSigningRequestInstallationResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_forge_certificate_signing_request_installation"
}

func (r *ForgeCertificateSigningRequestInstallationResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Forge certificate signing request installation resource. This resource allows you to install a certificate signing request on a server in Forge.",
		Attributes: map[string]schema.Attribute{
			"certificate_signing_request_id": schema.Int64Attribute{
				Required:            true,
				MarkdownDescription: "The ID of the certificate signing request.",
			},
			"server_id": schema.Int64Attribute{
				Required:            true,
				MarkdownDescription: "The ID of the server the certificate signing request is associated with.",
			},
			"site_id": schema.Int64Attribute{
				Required:            true,
				MarkdownDescription: "The ID of the site the certificate signing request is associated with.",
			},
			"certificate": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The certificate to install.",
			},
			"add_intermediates": schema.BoolAttribute{
				Computed:            true,
				Optional:            true,
				MarkdownDescription: "Automatically Add Intermediate Certificates.",
				Default:             booldefault.StaticBool(true),
			},
		},
	}
}

func (r *ForgeCertificateSigningRequestInstallationResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *ForgeCertificateSigningRequestInstallationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan ForgeCertificateSigningRequestInstallationResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Build the InstallCertificateRequest payload.
	payload := forge_client.InstallCertificateRequest{
		Certificate:      plan.Certificate.ValueString(),
		AddIntermediates: plan.AddIntermediates.ValueBool(),
	}

	// Call InstallCertificate on the client.
	err := r.client.InstallCertificate(ctx, int(plan.ServerID.ValueInt64()), int(plan.SiteID.ValueInt64()), int(plan.CertificateSigningRequestID.ValueInt64()), payload)
	if err != nil {
		resp.Diagnostics.AddError("Error installing certificate signing request", err.Error())
		return
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *ForgeCertificateSigningRequestInstallationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state ForgeCertificateSigningRequestInstallationResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *ForgeCertificateSigningRequestInstallationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan ForgeCertificateSigningRequestInstallationResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *ForgeCertificateSigningRequestInstallationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state ForgeCertificateSigningRequestInstallationResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
