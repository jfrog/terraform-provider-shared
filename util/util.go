package util

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"regexp"
	"strings"
	"text/template"

	"github.com/go-resty/resty/v2"
	"github.com/hashicorp/go-version"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type ProviderMetadata struct {
	Client             *resty.Client
	ProductId          string
	ArtifactoryVersion string
	XrayVersion        string
}

func resourceFeatureUsage(resourceName, method string) string {
	return fmt.Sprintf("Resource/%s/%s", resourceName, method)
}

func SendUsageResourceCreate(ctx context.Context, req *resty.Request, productId, resourceName string) {
	SendUsage(ctx, req, productId, resourceFeatureUsage(resourceName, "CREATE"))
}

func SendUsageResourceRead(ctx context.Context, req *resty.Request, productId, resourceName string) {
	SendUsage(ctx, req, productId, resourceFeatureUsage(resourceName, "READ"))
}

func SendUsageResourceUpdate(ctx context.Context, req *resty.Request, productId, resourceName string) {
	SendUsage(ctx, req, productId, resourceFeatureUsage(resourceName, "UPDATE"))
}

func SendUsageResourceDelete(ctx context.Context, req *resty.Request, productId, resourceName string) {
	SendUsage(ctx, req, productId, resourceFeatureUsage(resourceName, "DELETE"))
}

type Feature struct {
	FeatureId string `json:"featureId"`
}

type UsageStruct struct {
	ProductId string    `json:"productId"`
	Features  []Feature `json:"features"`
}

func SendUsage(ctx context.Context, req *resty.Request, productId string, featureUsages ...string) {
	if req == nil {
		tflog.Info(ctx, "SendUsage req is nil. Skipping.")
		return
	}

	features := []Feature{
		{FeatureId: "Partner/ACC-007450"},
	}

	for _, featureUsage := range featureUsages {
		features = append(features, Feature{FeatureId: featureUsage})
	}

	usage := UsageStruct{productId, features}

	resp, err := req.
		SetBody(usage).
		Post("artifactory/api/system/usage")

	if err != nil {
		tflog.Info(ctx, fmt.Sprintf("failed to send usage: %v", err))
	}

	if resp.IsError() {
		tflog.Info(ctx, fmt.Sprintf("failed to send usage: %v", resp.String()))
	}
}

type OIDCAccessTokenRequest struct {
	GrantType        string `json:"grant_type"`
	SubjectTokenType string `json:"subject_token_type"`
	SubjectToken     string `json:"subject_token"`
	ProviderName     string `json:"provider_name"`
}

type OIDCAccessTokenResponse struct {
	AccessToken string `json:"access_token"`
}

// OIDCTokenExchange use TFC_WORKLOAD_IDENTITY_TOKEN env var value to exchange for a access token using
// OIDC provider configured on JFrog platform
func OIDCTokenExchange(ctx context.Context, client *resty.Client, providerName string) (string, error) {
	if client == nil {
		return "", fmt.Errorf("client is nil")
	}

	tfcWorkloadIdentityToken := CheckEnvVars([]string{"TFC_WORKLOAD_IDENTITY_TOKEN"}, "")
	if tfcWorkloadIdentityToken == "" || providerName == "" {
		tflog.Info(ctx, "either TFC_WORKLOAD_IDENTITY_TOKEN or provider name is not set")
		return "", nil
	}

	payload := OIDCAccessTokenRequest{
		GrantType:        "urn:ietf:params:oauth:grant-type:token-exchange",
		SubjectTokenType: "urn:ietf:params:oauth:token-type:id_token",
		SubjectToken:     tfcWorkloadIdentityToken,
		ProviderName:     providerName,
	}

	var result OIDCAccessTokenResponse
	response, err := client.R().
		SetBody(payload).
		SetResult(&result).
		Post("/access/api/v1/oidc/token")

	if err != nil {
		return "", err
	}

	if response.IsError() {
		return "", fmt.Errorf(response.String())
	}

	return result.AccessToken, nil
}

func CheckArtifactoryLicense(client *resty.Client, licenseTypesToCheck ...string) error {
	if len(licenseTypesToCheck) == 0 {
		return fmt.Errorf("licenseTypesToCheck is empty")
	}

	type License struct {
		Type string `json:"type"`
	}

	type LicensesWrapper struct {
		License
		Licenses []License `json:"licenses"` // HA licenses returns as an array instead
	}

	licensesWrapper := LicensesWrapper{}
	resp, err := client.R().
		SetResult(&licensesWrapper).
		Get("/artifactory/api/system/license")

	if err != nil {
		return fmt.Errorf("failed to check for license. If your usage doesn't require admin permission, you can set `check_license` attribute to `false` to skip this check. %s", err.Error())
	}

	if resp.IsError() {
		return fmt.Errorf("failed to check for license. If your usage doesn't require admin permission, you can set `check_license` attribute to `false` to skip this check. %s", resp.String())
	}

	var licenseType string
	if len(licensesWrapper.Licenses) > 0 {
		licenseType = licensesWrapper.Licenses[0].Type
	} else {
		licenseType = licensesWrapper.Type
	}

	licenseTypesToCheckRegex := fmt.Sprintf("(?:%s)", strings.Join(licenseTypesToCheck, "|"))
	if matched, _ := regexp.MatchString(licenseTypesToCheckRegex, licenseType); !matched {
		licenseTypesToCheckMessage := strings.Join(licenseTypesToCheck, " or ")
		return fmt.Errorf("artifactory requires %s license to work with Terraform! If your usage doesn't require a license, you can set `check_license` attribute to `false` to skip this check", licenseTypesToCheckMessage)
	}

	return nil
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
	resp, err := client.R().
		SetResult(&artifactoryVersion).
		Get("/artifactory/api/system/version")

	if err != nil {
		return "", fmt.Errorf("failed to get Artifactory version. %s", err)
	}

	if resp.IsError() {
		return "", fmt.Errorf("failed to get Artifactory version. %s", resp.String())
	}

	return artifactoryVersion.Version, nil
}

func GetXrayVersion(client *resty.Client) (string, error) {
	type XrayVersion struct {
		Version string `json:"xray_version"`
	}

	xrayVersion := XrayVersion{}
	resp, err := client.R().
		SetResult(&xrayVersion).
		Get("/xray/api/v1/system/version")

	if err != nil {
		return "", fmt.Errorf("failed to get Xray version. %s", err)
	}

	if resp.IsError() {
		return "", fmt.Errorf("failed to get Xray version. %s", resp.String())
	}

	return xrayVersion.Version, nil
}

func CheckEnvVars(vars []string, dv string) string {
	for _, k := range vars {
		if v := os.Getenv(k); v != "" {
			return v
		}
	}
	return dv
}

func ExecuteTemplate(name, temp string, fields interface{}) string {
	var tpl bytes.Buffer
	if err := template.Must(template.New(name).Parse(temp)).Execute(&tpl, fields); err != nil {
		panic(err)
	}

	return tpl.String()
}
