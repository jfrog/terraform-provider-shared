// Copyright (c) JFrog Ltd. (2025)
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package util

import (
	"context"
	"fmt"

	"github.com/go-resty/resty/v2"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/jfrog/terraform-provider-shared/client"
	validator_string "github.com/jfrog/terraform-provider-shared/validator/fw/string"
)

type JFrogProvider struct {
	TypeName  string
	Meta      ProviderMetadata
	ProductID string
	Version   string
}

type ProviderMetadata struct {
	Client             *resty.Client
	ProductId          string
	ArtifactoryVersion string
	AccessVersion      string
	XrayVersion        string
}

type JFrogProviderModel struct {
	Url                  types.String `tfsdk:"url"`
	AccessToken          types.String `tfsdk:"access_token"`
	OIDCProviderName     types.String `tfsdk:"oidc_provider_name"`
	TFCCredentialTagName types.String `tfsdk:"tfc_credential_tag_name"`
}

func (p *JFrogProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	// Check environment variables, first available OS variable will be assigned to the var
	url := CheckEnvVars([]string{"JFROG_URL"}, "")
	accessToken := CheckEnvVars([]string{"JFROG_ACCESS_TOKEN"}, "")

	var config JFrogProviderModel

	// Read configuration data into model
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if config.Url.ValueString() != "" {
		url = config.Url.ValueString()
	}

	if url == "" {
		resp.Diagnostics.AddError(
			"Missing URL Configuration",
			"While configuring the provider, the url was not found in the JFROG_URL environment variable or provider configuration block url attribute.",
		)
		return
	}

	restyClient, err := client.Build(url, p.ProductID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating Resty client",
			err.Error(),
		)
		return
	}

	oidcProviderName := config.OIDCProviderName.ValueString()
	if oidcProviderName != "" {
		oidcAccessToken, err := OIDCTokenExchange(ctx, restyClient, oidcProviderName, config.TFCCredentialTagName.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				"Failed OIDC ID token exchange",
				err.Error(),
			)
			return
		}

		// use token from OIDC provider, which should take precedence over
		// environment variable data, if found.
		if oidcAccessToken != "" {
			accessToken = oidcAccessToken
		}
	}

	// use token from configuration, which should take precedence over
	// environment variable data or OIDC provider, if found.
	if config.AccessToken.ValueString() != "" {
		accessToken = config.AccessToken.ValueString()
	}

	if accessToken == "" {
		resp.Diagnostics.AddWarning(
			"Missing JFrog Access Token",
			"Access Token was not found in the JFROG_ACCESS_TOKEN environment variable, provider configuration block access_token attribute, or Terraform Cloud TFC_WORKLOAD_IDENTITY_TOKEN environment variable. Platform functionality will be affected.",
		)
	}

	artifactoryVersion := ""
	if len(accessToken) > 0 {
		_, err = client.AddAuth(restyClient, "", accessToken)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error adding Auth to Resty client",
				err.Error(),
			)
			return
		}

		version, err := GetArtifactoryVersion(restyClient)
		if err != nil {
			resp.Diagnostics.AddWarning(
				"Error getting Artifactory version",
				fmt.Sprintf("Provider functionality might be affected by the absence of Artifactory version. %v", err),
			)
		}

		artifactoryVersion = version

		featureUsage := fmt.Sprintf("Terraform/%s", req.TerraformVersion)
		go SendUsage(ctx, restyClient.R(), p.ProductID, featureUsage)
	}

	accessVersion := ""
	if len(accessToken) > 0 {
		_, err = client.AddAuth(restyClient, "", accessToken)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error adding Auth to Resty client",
				err.Error(),
			)
			return
		}

		version, err := GetAccessVersion(restyClient)
		if err != nil {
			resp.Diagnostics.AddWarning(
				"Error getting Access version",
				fmt.Sprintf("Provider functionality might be affected by the absence of Access version. %v", err),
			)
		}

		accessVersion = version

		featureUsage := fmt.Sprintf("Terraform/%s", req.TerraformVersion)
		go SendUsage(ctx, restyClient.R(), p.ProductID, featureUsage)
	}

	featureUsage := fmt.Sprintf("Terraform/%s", req.TerraformVersion)
	go SendUsage(ctx, restyClient.R(), p.ProductID, featureUsage)

	meta := ProviderMetadata{
		Client:             restyClient,
		ArtifactoryVersion: artifactoryVersion,
		AccessVersion:      accessVersion,
		ProductId:          p.ProductID,
	}

	p.Meta = meta

	resp.DataSourceData = meta
	resp.ResourceData = meta
}

func (p *JFrogProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = p.TypeName
	resp.Version = p.Version
}

func (p *JFrogProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"url": schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					validator_string.IsURLHttpOrHttps(),
				},
				MarkdownDescription: "JFrog Platform URL. This can also be sourced from the `JFROG_URL` environment variable.",
			},
			"access_token": schema.StringAttribute{
				Optional:  true,
				Sensitive: true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
				MarkdownDescription: "This is a access token that can be given to you by your admin under `Platform Configuration -> User Management -> Access Tokens`. This can also be sourced from the `JFROG_ACCESS_TOKEN` environment variable.",
			},
			"oidc_provider_name": schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
				MarkdownDescription: "OIDC provider name. See [Configure an OIDC Integration](https://jfrog.com/help/r/jfrog-platform-administration-documentation/configure-an-oidc-integration) for more details.",
			},
			"tfc_credential_tag_name": schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
				Description: "Terraform Cloud Workload Identity Token tag name. Use for generating multiple TFC workload identity tokens. When set, the provider will attempt to use env var with this tag name as suffix. **Note:** this is case sensitive, so if set to `JFROG`, then env var `TFC_WORKLOAD_IDENTITY_TOKEN_JFROG` is used instead of `TFC_WORKLOAD_IDENTITY_TOKEN`. See [Generating Multiple Tokens](https://developer.hashicorp.com/terraform/cloud-docs/workspaces/dynamic-provider-credentials/manual-generation#generating-multiple-tokens) on HCP Terraform for more details.",
			},
		},
	}
}
