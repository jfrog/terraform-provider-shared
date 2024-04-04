package fw

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/go-resty/resty/v2"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ValueToSetType ensures we have a types.Set literal
func ValueToSetType(v attr.Value) types.Set {
	if vb, ok := v.(types.Set); ok {
		return vb
	}
	if vb, ok := v.(*types.Set); ok {
		return *vb
	}
	panic(fmt.Sprintf("cannot pass type %T to conv.ValueToSetType", v))
}

// AttributeValueToString will attempt to execute the appropriate AttributeStringerFunc from the ones registered.
func AttributeValueToString(v attr.Value) string {
	if s, ok := v.(types.String); ok {
		return s.ValueString()
	}
	return v.String()
}

func StringSetToStrings(v attr.Value) []string {
	vt := ValueToSetType(v)
	out := make([]string, len(vt.Elements()))
	for i, ve := range vt.Elements() {
		out[i] = AttributeValueToString(ve)
	}
	return out
}

func UnableToCreateResourceError(resp *resource.CreateResponse, err string) {
	resp.Diagnostics.AddError(
		"Unable to Create Resource",
		"An unexpected error occurred while creating the resource update request. "+
			"Please report this issue to the provider developers.\n\n"+
			"Error: "+err,
	)
}

func UnableToUpdateResourceError(resp *resource.UpdateResponse, err string) {
	resp.Diagnostics.AddError(
		"Unable to Update Resource",
		"An unexpected error occurred while updating the resource update request. "+
			"Please report this issue to the provider developers.\n\n"+
			"Error: "+err,
	)
}

func UnableToRefreshResourceError(resp *resource.ReadResponse, err string) {
	resp.Diagnostics.AddError(
		"Unable to Refresh Resource",
		"An unexpected error occurred while attempting to refresh resource state. "+
			"Please retry the operation or report this issue to the provider developers.\n\n"+
			"Error: "+err,
	)
}

func UnableToDeleteResourceError(resp *resource.DeleteResponse, err string) {
	resp.Diagnostics.AddError(
		"Unable to Delete Resource",
		"An unexpected error occurred while attempting to delete the resource. "+
			"Please retry the operation or report this issue to the provider developers.\n\n"+
			"Error: "+err,
	)
}

func CheckArtifactoryLicense(client *resty.Client, licenseTypesToCheck ...string) (ds diag.Diagnostics) {
	if len(licenseTypesToCheck) == 0 {
		ds.AddError("licenseTypesToCheck is empty", "")
		return
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
		ds.AddError("Failed to check for license. If your usage doesn't require admin permission, you can set `check_license` attribute to `false` to skip this check.", err.Error())
		return
	}

	if resp.IsError() {
		ds.AddError("Failed to check for license. If your usage doesn't require admin permission, you can set `check_license` attribute to `false` to skip this check.", resp.String())
		return
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
		ds.AddError(fmt.Sprintf("Artifactory requires %s license to work with Terraform! If your usage doesn't require a license, you can set `check_license` attribute to `false` to skip this check.", licenseTypesToCheckMessage), "")
		return
	}

	return
}
