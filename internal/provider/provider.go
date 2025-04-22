package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ provider.Provider                       = &LaravelProvider{}
	_ provider.ProviderWithFunctions          = &LaravelProvider{}
	_ provider.ProviderWithEphemeralResources = &LaravelProvider{}
)

type LaravelProvider struct {
	version string
}

type LaravelProviderModel struct {
	// Envoyer Configuration
	EnvoyerAPIToken types.String `tfsdk:"envoyer_api_token"`
	EnvoyerEnvKey   types.String `tfsdk:"envoyer_env_key"`
	EnvoyerBaseURL  types.String `tfsdk:"envoyer_base_url"`

	// Forge Configuration
	ForgeAPIToken types.String `tfsdk:"forge_api_token"`
	ForgeBaseURL  types.String `tfsdk:"forge_base_url"`

	// Advanced Configuration Options
	RequestTimeout types.Int64 `tfsdk:"request_timeout"`
	MaxRetries     types.Int64 `tfsdk:"max_retries"`
	RetryDelay     types.Int64 `tfsdk:"retry_delay"`
	EnableCache    types.Bool  `tfsdk:"enable_cache"`
	CacheTTL       types.Int64 `tfsdk:"cache_ttl"`
}

func (p *LaravelProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "laravel"
	resp.Version = p.version
}

func (p *LaravelProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"envoyer_api_token": schema.StringAttribute{
				MarkdownDescription: "Envoyer API token (Bearer token). Required if using Envoyer resources/data sources.",
				Optional:            true,
				Sensitive:           true,
			},
			"envoyer_env_key": schema.StringAttribute{
				MarkdownDescription: "Optional override of the Envoyer env-lock key.",
				Optional:            true,
				Sensitive:           true,
			},
			"envoyer_base_url": schema.StringAttribute{
				MarkdownDescription: "Optional override of the Envoyer API base URL (defaults to `https://envoyer.io/api`).",
				Optional:            true,
			},
			"forge_api_token": schema.StringAttribute{
				MarkdownDescription: "Forge API token (Bearer token). Required if using Forge resources/data sources.",
				Optional:            true,
				Sensitive:           true,
			},
			"forge_base_url": schema.StringAttribute{
				MarkdownDescription: "Optional override of the Forge API base URL (defaults to `https://forge.laravel.com/api/v1`).",
				Optional:            true,
			},
			"request_timeout": schema.Int64Attribute{
				MarkdownDescription: "Timeout for API requests in seconds. Default is 30 seconds.",
				Optional:            true,
			},
			"max_retries": schema.Int64Attribute{
				MarkdownDescription: "Maximum number of retries for failed API requests. Default is 3.",
				Optional:            true,
			},
			"retry_delay": schema.Int64Attribute{
				MarkdownDescription: "Delay between retries in seconds. Default is 5 seconds.",
				Optional:            true,
			},
			"enable_cache": schema.BoolAttribute{
				MarkdownDescription: "Enable caching of API responses. Default is false.",
				Optional:            true,
			},
			"cache_ttl": schema.Int64Attribute{
				MarkdownDescription: "Time-to-live for cached API responses in seconds. Default is 300 seconds (5 minutes).",
				Optional:            true,
			},
		},
	}
}

func (p *LaravelProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	tflog.Info(ctx, "Configuring Laravel provider")

	var config LaravelProviderModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if (config.ForgeAPIToken.IsNull() || config.ForgeAPIToken.IsUnknown()) &&
		(config.EnvoyerAPIToken.IsNull() || config.EnvoyerAPIToken.IsUnknown()) {
		resp.Diagnostics.AddError(
			"Missing API Tokens",
			"You must provide at least one of forge_api_token or envoyer_api_token depending on which services you need to use.",
		)
		return
	}

	providerConfig, diags := createProviderConfig(ctx, config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.DataSourceData = providerConfig
	resp.ResourceData = providerConfig

	if providerConfig.Forge != nil {
		tflog.Info(ctx, "Laravel Forge client configured successfully")
	}
	if providerConfig.Envoyer != nil {
		tflog.Info(ctx, "Laravel Envoyer client configured successfully")
	}
}

func (p *LaravelProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewEnvoyerProjectResource,
		NewEnvoyerDeploymentResource,
		NewEnvoyerServerResource,
		NewEnvoyerHookResource,
		NewEnvoyerEnvironmentResource,
		NewForgeServerResource,
		NewForgeSiteResource,
		NewForgeWorkerResource,
		NewForgeRecipeResource,
		NewForgeRecipeRunResource,
		// NewForgeDaemonResource,
		// NewForgeFirewallRuleResource,
		NewForgeSSHKeyResource,
		NewForgeCertificateResource,
		NewForgeCertificateSigningRequestResource,
		NewForgeCertificateSigningRequestInstallationResource,
		NewForgeScheduledJobResource,
		// NewForgeDatabaseResource,
		// NewForgeDatabaseUserResource,
		// NewForgeNginxTemplateResource,
		// NewForgeRedirectRuleResource,
		// NewForgeMonitorResource,
	}
}

func (p *LaravelProvider) EphemeralResources(ctx context.Context) []func() ephemeral.EphemeralResource {
	return []func() ephemeral.EphemeralResource{}
}

func (p *LaravelProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewEnvoyerProjectDataSource,
		NewEnvoyerServersDataSource,
		NewEnvoyerActionsDataSource,
		NewForgeCredentialsDataSource,
		// NewForgeServersDataSource,
		// NewForgeSitesDataSource,
		// NewForgePHPVersionsDataSource,
		// NewForgeRegionsDataSource,
		// NewForgeUserDataSource,
		// NewForgeSSHKeysDataSource,
		// NewForgeJobsDataSource,
	}
}

func (p *LaravelProvider) Functions(ctx context.Context) []func() function.Function {
	return []func() function.Function{}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &LaravelProvider{
			version: version,
		}
	}
}
