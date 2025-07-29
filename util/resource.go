package util

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
)

type JFrogResource struct {
	ProviderData            *ProviderMetadata
	TypeName                string
	ValidArtifactoryVersion string
	ValidXrayVersion        string
	DocumentEndpoint        string
	CollectionEndpoint      string
}

func (r *JFrogResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = r.TypeName
}

func (r *JFrogResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}
	m := req.ProviderData.(ProviderMetadata)
	r.ProviderData = &m
}

func (r JFrogResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	if r.ProviderData == nil || r.ValidArtifactoryVersion == "" {
		return
	}

	valid, err := CheckVersion(r.ProviderData.ArtifactoryVersion, r.ValidArtifactoryVersion)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to verify Artifactory version",
			err.Error(),
		)
	}

	if !valid {
		resp.Diagnostics.AddError(
			"Incompatible Artifactory version",
			fmt.Sprintf("This resource is only supported by Artifactory version %s or later.", r.ValidArtifactoryVersion),
		)
	}
}

func (r JFrogResource) ValidateXrayConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	if r.ProviderData == nil || r.ValidXrayVersion == "" {
		return
	}

	valid, err := CheckVersion(r.ProviderData.XrayVersion, r.ValidXrayVersion)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to verify Xray version",
			err.Error(),
		)
	}

	if !valid {
		resp.Diagnostics.AddError(
			"Incompatible Xray version",
			fmt.Sprintf("This resource is only supported by Xray version %s or later.", r.ValidXrayVersion),
		)
	}
}