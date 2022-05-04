package util

import (
	"context"
	"fmt"
	"reflect"
	"regexp"
	"sort"
	"strings"

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
	if v, ok := d.GetOkExists(key); ok && (!onlyIfChanged || d.HasChange(key)) {
		return BoolPtr(v.(bool))
	}
	return nil
}

func (d *ResourceData) GetBool(key string, onlyIfChanged bool) bool {
	if v, ok := d.GetOkExists(key); ok && (!onlyIfChanged || d.HasChange(key)) {
		return v.(bool)
	}
	return false
}

func (d *ResourceData) GetInt(key string, onlyIfChanged bool) int {
	if v, ok := d.GetOkExists(key); ok && (!onlyIfChanged || d.HasChange(key)) {
		return v.(int)
	}
	return 0
}

func (d *ResourceData) GetSet(key string) []string {
	if v, ok := d.GetOkExists(key); ok {
		arr := CastToStringArr(v.(*schema.Set).List())
		return arr
	}
	return nil
}

func (d *ResourceData) GetList(key string) []string {
	if v, ok := d.GetOkExists(key); ok {
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

type HclPredicate func(hcl string) bool

func SchemaHasKey(skeema map[string]*schema.Schema) HclPredicate {
	return func(key string) bool {
		_, ok := skeema[key]
		return ok
	}
}

func FormatCommaSeparatedString(thing interface{}) string {
	fields := strings.Fields(thing.(string))
	sort.Strings(fields)
	return strings.Join(fields, ",")
}

type PackFunc func(repo interface{}, d *schema.ResourceData) error

// UniversalPack consider making this a function that takes a predicate of what to include and returns
// a function that does the job. This would allow for the legacy code to specify which keys to keep and not
func UniversalPack(predicate HclPredicate) PackFunc {

	return func(payload interface{}, d *schema.ResourceData) error {
		setValue := MkLens(d)

		var errors []error

		values := lookup(payload, predicate)

		for hcl, value := range values {
			if predicate != nil && predicate(hcl) {
				errors = setValue(hcl, value)
			}
		}

		if errors != nil && len(errors) > 0 {
			return fmt.Errorf("failed saving state %q", errors)
		}
		return nil
	}
}

// fieldToHcl this function is meant to use the HCL provided in the tag, or create a snake_case from the field name
// it actually works as expected, but dynamically working with these names was catching edge cases everywhere and
// it was/is a time sink to catch.
func fieldToHcl(field reflect.StructField) string {

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

type AutoMapper func(field reflect.StructField, thing reflect.Value) map[string]interface{}

func findInspector(kind reflect.Kind) AutoMapper {
	switch kind {
	case reflect.Struct:
		return func(f reflect.StructField, t reflect.Value) map[string]interface{} {
			return lookup(t.Interface(), nil)
		}
	case reflect.Ptr:
		return func(field reflect.StructField, thing reflect.Value) map[string]interface{} {
			deref := reflect.Indirect(thing)
			if deref.CanAddr() {
				result := deref.Interface()
				if deref.Kind() == reflect.Struct {
					result = []interface{}{lookup(deref.Interface(), nil)}
				}
				return map[string]interface{}{
					fieldToHcl(field): result,
				}
			}
			return map[string]interface{}{}
		}
	case reflect.Slice:
		return func(field reflect.StructField, thing reflect.Value) map[string]interface{} {
			return map[string]interface{}{
				fieldToHcl(field): CastToInterfaceArr(thing.Interface().([]string)),
			}
		}
	}
	return func(field reflect.StructField, thing reflect.Value) map[string]interface{} {
		return map[string]interface{}{
			fieldToHcl(field): thing.Interface(),
		}
	}
}

var allowAllPredicate = func(hcl string) bool {
	return true
}

func lookup(payload interface{}, predicate HclPredicate) map[string]interface{} {

	if predicate == nil {
		predicate = allowAllPredicate
	}

	values := map[string]interface{}{}
	var t = reflect.TypeOf(payload)
	var v = reflect.ValueOf(payload)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
		v = v.Elem()
	}
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		thing := v.Field(i)

		shouldLookup := true
		if thing.Kind() != reflect.Struct {
			hcl := fieldToHcl(field)
			shouldLookup = predicate(hcl)
		}

		if shouldLookup {
			typeInspector := findInspector(thing.Kind())
			for key, value := range typeInspector(field, thing) {
				if _, ok := values[key]; !ok {
					values[key] = value
				}
			}
		}
	}
	return values
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
		features = append(features, Feature{FeatureId: featureUsage} )
	}

	usage := UsageStruct{productId, features}

	_, err := client.R().
		SetBody(usage).
		Post("artifactory/api/system/usage")

	if err != nil {
		tflog.Info(ctx, fmt.Sprintf("failed to send usage: %v", err))
	}
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

func CheckArtifactoryLicense(client *resty.Client) diag.Diagnostics {

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
		return diag.Errorf("Failed to check for license. %s", err)
	}

	var licenseType string
	if len(licensesWrapper.Licenses) > 0 {
		licenseType = licensesWrapper.Licenses[0].Type
	} else {
		licenseType = licensesWrapper.Type
	}

	if matched, _ := regexp.MatchString(`Enterprise`, licenseType); !matched {
		return diag.Errorf("Artifactory Projects requires Enterprise license to work with Terraform!")
	}

	return nil
}
