package provider

import (
	"context"
	"fmt"
	"strings"
	"terraform-provider-laravel/internal/envoyer_client"
	"terraform-provider-laravel/internal/forge_client"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// providerConfig contains the configured client connections.
type providerConfig struct {
	Forge   *forge_client.Client
	Envoyer *envoyer_client.Client
}

// createProviderConfig creates and configures the API clients based on provider configuration.
func createProviderConfig(ctx context.Context, config LaravelProviderModel) (*providerConfig, diag.Diagnostics) {
	var diags diag.Diagnostics
	provConfig := &providerConfig{}

	// Configure Forge client if token is provided
	if !config.ForgeAPIToken.IsNull() && !config.ForgeAPIToken.IsUnknown() {
		forgeClient, forgeDiags := configureForgeClient(ctx, config)
		diags.Append(forgeDiags...)
		if !diags.HasError() {
			provConfig.Forge = forgeClient
			tflog.Info(ctx, "Forge client configured successfully")
		}
	} else {
		tflog.Info(ctx, "Forge API token not provided, Forge client will not be configured")
	}

	// Configure Envoyer client if token is provided
	if !config.EnvoyerAPIToken.IsNull() && !config.EnvoyerAPIToken.IsUnknown() {
		envoyerClient, envoyerDiags := configureEnvoyerClient(ctx, config)
		diags.Append(envoyerDiags...)
		if !diags.HasError() {
			provConfig.Envoyer = envoyerClient
			tflog.Info(ctx, "Envoyer client configured successfully")
		}
	} else {
		tflog.Info(ctx, "Envoyer API token not provided, Envoyer client will not be configured")
	}

	return provConfig, diags
}

func configureForgeClient(ctx context.Context, config LaravelProviderModel) (*forge_client.Client, diag.Diagnostics) {
	var diags diag.Diagnostics

	forgeAPIToken := config.ForgeAPIToken.ValueString()
	forgeBaseURL := forge_client.DefaultBaseURL
	if !config.ForgeBaseURL.IsNull() && config.ForgeBaseURL.ValueString() != "" {
		forgeBaseURL = strings.TrimSuffix(config.ForgeBaseURL.ValueString(), "/")
	}

	// Create the client
	client := forge_client.NewClient(forgeAPIToken)
	client.WithBaseURL(forgeBaseURL)

	// Configure advanced options
	if !config.RequestTimeout.IsNull() {
		timeout := time.Duration(config.RequestTimeout.ValueInt64()) * time.Second
		httpClient := client.HTTPClient()
		httpClient.Timeout = timeout
		client.WithHTTPClient(httpClient)
	}

	if !config.MaxRetries.IsNull() {
		retries := int(config.MaxRetries.ValueInt64())
		retryDelay := 5 * time.Second
		if !config.RetryDelay.IsNull() {
			retryDelay = time.Duration(config.RetryDelay.ValueInt64()) * time.Second
		}
		client.WithRetryConfig(retries, retryDelay)
	}

	if !config.EnableCache.IsNull() && config.EnableCache.ValueBool() {
		cache := forge_client.NewMemoryCache()
		cacheConfig := forge_client.CacheConfig{
			Enabled:             true,
			TTL:                 5 * time.Minute,
			CleanupInterval:     10 * time.Minute,
			CacheErrorResponses: false,
		}

		if !config.CacheTTL.IsNull() {
			cacheConfig.TTL = time.Duration(config.CacheTTL.ValueInt64()) * time.Second
		}

		client.WithCache(cache)
		client.WithCacheConfig(cacheConfig)

		tflog.Info(ctx, "Enabled caching for Forge API client", map[string]interface{}{
			"ttl": cacheConfig.TTL.String(),
		})
	}

	// Test the client with a basic request
	_, err := client.GetUser(ctx)
	if err != nil {
		diags.AddError(
			"Failed to connect to Forge API",
			fmt.Sprintf("Error connecting to Forge API: %s", err),
		)
		return nil, diags
	}

	tflog.Info(ctx, "Successfully connected to Forge API", map[string]interface{}{
		"base_url": forgeBaseURL,
	})

	return client, diags
}

func configureEnvoyerClient(ctx context.Context, config LaravelProviderModel) (*envoyer_client.Client, diag.Diagnostics) {
	var diags diag.Diagnostics

	envoyerAPIToken := config.EnvoyerAPIToken.ValueString()

	envoyerEnvKey := ""
	if !config.EnvoyerEnvKey.IsNull() && config.EnvoyerEnvKey.ValueString() != "" {
		envoyerEnvKey = config.EnvoyerEnvKey.ValueString()
	}

	envoyerBaseURL := envoyer_client.DefaultBaseURL
	if !config.EnvoyerBaseURL.IsNull() && config.EnvoyerBaseURL.ValueString() != "" {
		envoyerBaseURL = strings.TrimSuffix(config.EnvoyerBaseURL.ValueString(), "/")
	}

	// Create the client
	client := envoyer_client.NewClient(envoyerAPIToken, envoyerEnvKey)
	client.WithBaseURL(envoyerBaseURL)
	// client.WithDebug(true)

	tflog.Debug(ctx, "Configuring Envoyer client", map[string]interface{}{
		"base_url": envoyerBaseURL,
		"env_key":  envoyerEnvKey,
	})

	// Test the client with a basic request
	_, err := client.ListProjects(ctx)
	if err != nil {
		diags.AddError(
			"Failed to connect to Envoyer API",
			fmt.Sprintf("Error connecting to Envoyer API: %s", err),
		)
		return nil, diags
	}

	tflog.Info(ctx, "Successfully connected to Envoyer API", map[string]interface{}{
		"base_url": envoyerBaseURL,
	})

	return client, diags
}
