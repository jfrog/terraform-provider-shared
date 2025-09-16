package util

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
	"text/template"

	"github.com/go-resty/resty/v2"
	"github.com/hashicorp/go-version"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

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
func OIDCTokenExchange(ctx context.Context, client *resty.Client, providerName, credentialTag string) (string, error) {
	if client == nil {
		return "", fmt.Errorf("client is nil")
	}

	if providerName == "" {
		return "", fmt.Errorf("provider name is not set")
	}

	tfcWorkloadIdentityTokenEnvVars := []string{"TFC_WORKLOAD_IDENTITY_TOKEN"}
	if credentialTag != "" {
		tfcWorkloadIdentityTokenEnvVars = append(
			tfcWorkloadIdentityTokenEnvVars,
			fmt.Sprintf("TFC_WORKLOAD_IDENTITY_TOKEN_%s", credentialTag),
		)
	}

	tfcWorkloadIdentityToken := CheckEnvVars(tfcWorkloadIdentityTokenEnvVars, "")
	if tfcWorkloadIdentityToken == "" {
		return "", fmt.Errorf("env var %s is not set", strings.Join(tfcWorkloadIdentityTokenEnvVars, " or "))
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
		return "", fmt.Errorf("OIDC token exchange failed: %s", response.String())
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

func CheckXrayVersion(client *resty.Client, minVersion string, customMessage string) (string, error) {
	// Skip version check if disabled via environment variable
	if GetBoolEnvVar([]string{"SKIP_XRAY_VERSION_CHECK"}, false) {
		return "", nil
	}

	type versionCheckError struct {
		currentVersion string
		minVersion     string
		customMessage  string
	}

	vErr := func(current, min, msg string) error {
		err := &versionCheckError{
			currentVersion: current,
			minVersion:     min,
			customMessage:  msg,
		}
		if err.customMessage != "" {
			return fmt.Errorf(err.customMessage, err.minVersion, err.currentVersion)
		}
		return fmt.Errorf("xray version %s is not supported - minimum required version is %s", err.currentVersion, err.minVersion)
	}
	version, err := GetXrayVersion(client)
	if err != nil {
		return "", fmt.Errorf("failed to get Xray version: %v", err)
	}

	supported, err := CheckVersion(version, minVersion)
	if err != nil {
		return version, fmt.Errorf("failed to check version compatibility: %v", err)
	}

	if !supported {
		return version, vErr(version, minVersion, customMessage)
	}

	return version, nil
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

func GetAccessVersion(client *resty.Client) (string, error) {
	type AccessVersion struct {
		Version string `json:"name"`
	}

	accessVersion := AccessVersion{}
	resp, err := client.R().
		SetResult(&accessVersion).
		Get("/access/api/v1/system/version")

	if err != nil {
		return "", fmt.Errorf("failed to get Access version. %s", err)
	}

	if resp.IsError() {
		return "", fmt.Errorf("failed to get Access version. %s", resp.String())
	}

	return accessVersion.Version, nil
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

func CheckCatalogHealth(client *resty.Client) error {
	type CatalogEntitlements struct {
		EntitledForCatalog bool `json:"entitled_for_catalog"`
		HasCentralToken    bool `json:"has_central_token"`
		TokenExpired       bool `json:"token_expired"`
	}

	type CatalogCentral struct {
		CentralConnectionWorking bool `json:"central_connection_working"`
	}

	type CatalogHealthResponse struct {
		Entitlements        CatalogEntitlements `json:"entitlements"`
		Central             CatalogCentral      `json:"central"`
		DbConnectionWorking bool                `json:"db_connection_working"`
		OneModelAvailable   bool                `json:"one_model_available"`
		Code                string              `json:"code"`
	}

	catalogHealth := CatalogHealthResponse{}
	resp, err := client.R().
		SetResult(&catalogHealth).
		Get("/catalog/api/v1/system/app_health")

	if err != nil {
		log.Printf("[DEBUG] Catalog health check failed with error: %s", err.Error())
		return fmt.Errorf("failed to validate catalog health. %s", err)
	}

	if resp.IsError() {
		log.Printf("[ERROR] Catalog health check returned error response: %s", resp.String())
		return fmt.Errorf("failed to validate catalog health. %s", resp.String())
	}

	// Check if catalog is healthy
	if catalogHealth.Code != "OK" {
		log.Printf("[ERROR] Catalog health check failed with code: %s", catalogHealth.Code)
		return fmt.Errorf("catalog health check failed with code: %s", catalogHealth.Code)
	}

	// Check entitlements
	if !catalogHealth.Entitlements.EntitledForCatalog {
		log.Printf("[ERROR] Catalog is not entitled for use")
		return fmt.Errorf("catalog is not entitled for use")
	}

	if catalogHealth.Entitlements.TokenExpired {
		log.Printf("[ERROR] Catalog token has expired")
		return fmt.Errorf("catalog token has expired")
	}

	// Check connections
	if !catalogHealth.DbConnectionWorking {
		log.Printf("[ERROR] Catalog database connection is not working")
		return fmt.Errorf("catalog database connection is not working")
	}

	if !catalogHealth.Central.CentralConnectionWorking {
		log.Printf("[ERROR] Catalog central connection is not working")
		return fmt.Errorf("catalog central connection is not working")
	}

	if !catalogHealth.OneModelAvailable {
		log.Printf("[ERROR] Catalog model is not available")
		return fmt.Errorf("catalog model is not available")
	}

	return nil
}

func CheckEnvVars(vars []string, defaultValue string) string {
	for _, k := range vars {
		if v := os.Getenv(k); v != "" {
			return v
		}
	}
	return defaultValue
}

func GetBoolEnvVar(vars []string, defaultValue bool) bool {
	for _, k := range vars {
		if v := os.Getenv(k); v != "" {
			return v == "true"
		}
	}
	return defaultValue
}

func ExecuteTemplate(name, temp string, fields interface{}) string {
	var tpl bytes.Buffer
	if err := template.Must(template.New(name).Parse(temp)).Execute(&tpl, fields); err != nil {
		panic(err)
	}

	return tpl.String()
}

type Identifiable interface {
	Id() string
}
