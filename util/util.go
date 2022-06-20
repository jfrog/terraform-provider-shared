package util

import (
	"bytes"
	"context"
	"fmt"
	"math"
	"reflect"
	"regexp"
	"sort"
	"strings"
	"text/template"

	"github.com/go-resty/resty/v2"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type ResourceData struct{ *schema.ResourceData }

func (d *ResourceData) GetString(key string, onlyIfChanged bool) string {
	if v, ok := d.GetOk(key); ok && (!onlyIfChanged || d.HasChange(key)) {
		return v.(string)
	}
	return ""
}

func BoolPtr(v bool) *bool { return &v }

func (d *ResourceData) GetBoolRef(key string, onlyIfChanged bool) *bool {
	if v, ok := d.GetOk(key); ok && (!onlyIfChanged || d.HasChange(key)) {
		return BoolPtr(v.(bool))
	}
	return nil
}

func (d *ResourceData) GetBool(key string, onlyIfChanged bool) bool {
	if v, ok := d.GetOk(key); ok && (!onlyIfChanged || d.HasChange(key)) {
		return v.(bool)
	}
	return false
}

func (d *ResourceData) GetInt(key string, onlyIfChanged bool) int {
	if v, ok := d.GetOk(key); ok && (!onlyIfChanged || d.HasChange(key)) {
		return v.(int)
	}
	return 0
}

func (d *ResourceData) GetSet(key string) []string {
	if v, ok := d.GetOk(key); ok {
		arr := CastToStringArr(v.(*schema.Set).List())
		return arr
	}
	return nil
}

func (d *ResourceData) GetList(key string) []string {
	if v, ok := d.GetOk(key); ok {
		arr := CastToStringArr(v.([]interface{}))
		return arr
	}
	return []string{}
}

func CastToStringArr(arr []interface{}) []string {
	cpy := make([]string, 0, len(arr))
	for _, r := range arr {
		cpy = append(cpy, r.(string))
	}

	return cpy
}

func CastToInterfaceArr(arr []string) []interface{} {
	cpy := make([]interface{}, 0, len(arr))
	for _, r := range arr {
		cpy = append(cpy, r)
	}

	return cpy
}

func MergeSchema(schemata ...map[string]*schema.Schema) map[string]*schema.Schema {
	result := map[string]*schema.Schema{}
	for _, schma := range schemata {
		for k, v := range schma {
			result[k] = v
		}
	}
	return result
}

type Lens func(key string, value interface{}) []error

func MkLens(d *schema.ResourceData) Lens {
	var errors []error
	return func(key string, value interface{}) []error {
		if err := d.Set(key, value); err != nil {
			errors = append(errors, err)
		}
		return errors
	}
}

type Schema map[string]*schema.Schema

func FormatCommaSeparatedString(thing interface{}) string {
	fields := strings.Fields(thing.(string))
	sort.Strings(fields)
	return strings.Join(fields, ",")
}

func toHclFormat(thing interface{}) string {
	switch thing.(type) {
	case string:
		return fmt.Sprintf(`"%s"`, thing.(string))
	case []interface{}:
		var result []string
		for _, e := range thing.([]interface{}) {
			result = append(result, toHclFormat(e))
		}
		return fmt.Sprintf("[%s]", strings.Join(result, ","))
	case map[string]interface{}:
		return fmt.Sprintf("\n\t%s\n\t\t\t\t", FmtMapToHcl(thing.(map[string]interface{})))
	default:
		return fmt.Sprintf("%v", thing)
	}
}

func FmtMapToHcl(fields map[string]interface{}) string {
	var allPairs []string
	max := float64(0)
	for key := range fields {
		max = math.Max(max, float64(len(key)))
	}
	for key, value := range fields {
		hcl := toHclFormat(value)
		format := toHclFormatString(3, int(max), value)
		allPairs = append(allPairs, fmt.Sprintf(format, key, hcl))
	}

	return strings.Join(allPairs, "\n")
}

func toHclFormatString(tabs, max int, value interface{}) string {
	prefix := ""
	suffix := ""
	delimiter := "="
	if reflect.TypeOf(value).Kind() == reflect.Map {
		delimiter = ""
		prefix = "{"
		suffix = "}"
	}
	return fmt.Sprintf("%s%%-%ds %s %s%s%s", strings.Repeat("\t", tabs), max, delimiter, prefix, "%s", suffix)
}

// FieldToHcl this function is meant to use the HCL provided in the tag, or create a snake_case from the field name
// it actually works as expected, but dynamically working with these names was catching edge cases everywhere and
// it was/is a time sink to catch.
func FieldToHcl(field reflect.StructField) string {

	if field.Tag.Get("hcl") != "" {
		return field.Tag.Get("hcl")
	}
	var lowerFields []string
	rgx := regexp.MustCompile("([A-Z][a-z]+)")
	fields := rgx.FindAllStringSubmatch(field.Name, -1)
	for _, matches := range fields {
		for _, match := range matches[1:] {
			lowerFields = append(lowerFields, strings.ToLower(match))
		}
	}
	result := strings.Join(lowerFields, "_")
	return result
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
func ExecuteTemplate(name, temp string, fields interface{}) string {
	var tpl bytes.Buffer
	if err := template.Must(template.New(name).Parse(temp)).Execute(&tpl, fields); err != nil {
		panic(err)
	}

	return tpl.String()
}
func MergeMaps(schemata ...map[string]interface{}) map[string]interface{} {
	result := map[string]interface{}{}
	for _, schma := range schemata {
		for k, v := range schma {
			result[k] = v
		}
	}
	return result
}
func applyTelemetry(productId, resource, verb string, f func(context.Context, *schema.ResourceData, interface{}) diag.Diagnostics) func(context.Context, *schema.ResourceData, interface{}) diag.Diagnostics {
	if f == nil {
		panic("attempt to apply telemetry to a nil function")
	}
	return func(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
		// best effort. Go routine it
		featureUsage := fmt.Sprintf("Resource/%s/%s", resource, verb)
		go SendUsage(ctx, meta.(*resty.Client), productId, featureUsage)
		return f(ctx, data, meta)
	}
}

func AddTelemetry(productId string, resourceMap map[string]*schema.Resource) map[string]*schema.Resource {
	for name, skeema := range resourceMap {
		if skeema.Create != nil {
			panic(fmt.Sprintf("[%s] deprecated Create function in use", name))
		}
		if skeema.Read != nil {
			panic(fmt.Sprintf("[%s] deprecated Read function in use", name))
		}
		if skeema.Update != nil {
			panic(fmt.Sprintf("[%s] deprecated Update function in use", name))
		}
		if skeema.Delete != nil {
			panic(fmt.Sprintf("[%s] deprecated Delete function in use", name))
		}
	}

	for name, skeema := range resourceMap {
		if skeema.CreateContext != nil {
			skeema.CreateContext = applyTelemetry(productId, name, "CREATE", skeema.CreateContext)
		}
		if skeema.ReadContext != nil {
			skeema.ReadContext = applyTelemetry(productId, name, "READ", skeema.ReadContext)
		}
		if skeema.UpdateContext != nil {
			skeema.UpdateContext = applyTelemetry(productId, name, "UPDATE", skeema.UpdateContext)
		}
		if skeema.DeleteContext != nil {
			skeema.DeleteContext = applyTelemetry(productId, name, "DELETE", skeema.DeleteContext)
		}
	}
	return resourceMap
}

func CheckArtifactoryLicense(client *resty.Client, licenseTypesToCheck ...string) diag.Diagnostics {
	if len(licenseTypesToCheck) == 0 {
		return diag.Errorf("failed: licenseTypesToCheck is empty")
	}

	type License struct {
		Type string `json:"type"`
	}

	type LicensesWrapper struct {
		License
		Licenses []License `json:"licenses"` // HA licenses returns as an array instead
	}

	licensesWrapper := LicensesWrapper{}
	_, err := client.R().
		SetResult(&licensesWrapper).
		Get("/artifactory/api/system/license")

	if err != nil {
		return diag.Errorf("Failed to check for license. If your usage doesn't require admin permission, you can set `check_license` attribute to `false` to skip this check. %s", err)
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
		return diag.Errorf("Artifactory requires %s license to work with Terraform! If your usage doesn't require a license, you can set `check_license` attribute to `false` to skip this check.", licenseTypesToCheckMessage)
	}

	return nil
}

type Identifiable interface {
	Id() string
}
