package util

import (
	"context"
	"fmt"

	"github.com/go-resty/resty/v2"
	"github.com/hashicorp/go-version"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type ProvderMetadata struct {
	Client             *resty.Client
	ArtifactoryVersion string
	XrayVersion        string
}

func SendUsage(ctx context.Context, client *resty.Client, productId string, featureUsages ...string) {
	type Feature struct {
		FeatureId string `json:"featureId"`
	}
	type UsageStruct struct {
		ProductId string    `json:"productId"`
		Features  []Feature `json:"features"`
	}

	features := []Feature{
		{FeatureId: "Partner/ACC-007450"},
	}

	for _, featureUsage := range featureUsages {
		features = append(features, Feature{FeatureId: featureUsage})
	}

	usage := UsageStruct{productId, features}

	_, err := client.R().
		SetBody(usage).
		Post("artifactory/api/system/usage")

	if err != nil {
		tflog.Info(ctx, fmt.Sprintf("failed to send usage: %v", err))
	}
}

func CheckVersion(versionToCheck string, supportedVersion string) (bool, error) {
	v1, err := version.NewVersion(versionToCheck)
	if err != nil {
		return false, fmt.Errorf("could not parse version: %s", versionToCheck)
	}

	v2, err := version.NewVersion(supportedVersion)
	if err != nil {
		return false, fmt.Errorf("could not parse version: %s", supportedVersion)
	}

	return v1.GreaterThanOrEqual(v2), nil
}

func GetArtifactoryVersion(client *resty.Client) (string, error) {
	type ArtifactoryVersion struct {
		Version string `json:"version"`
	}

	artifactoryVersion := ArtifactoryVersion{}
	_, err := client.R().
		SetResult(&artifactoryVersion).
		Get("/artifactory/api/system/version")

	if err != nil {
		return "", fmt.Errorf("failed to get Artifactory version. %s", err)
	}

	return artifactoryVersion.Version, nil
}

func GetXrayVersion(client *resty.Client) (string, error) {
	type XrayVersion struct {
		Version string `json:"xray_version"`
	}

	xrayVersion := XrayVersion{}
	_, err := client.R().
		SetResult(&xrayVersion).
		Get("/xray/api/v1/system/version")

	if err != nil {
		return "", fmt.Errorf("failed to get Xray version. %s", err)
	}

	return xrayVersion.Version, nil
}
