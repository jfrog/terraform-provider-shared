package packer

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/jfrog/terraform-provider-shared/predicate"
	utilsdk "github.com/jfrog/terraform-provider-shared/util/sdk"

	"fmt"
	"reflect"
)

type PackFunc func(repo interface{}, d *schema.ResourceData) error

// Universal consider making this a function that takes a predicate of what to include and returns
// a function that does the job. This would allow for the legacy code to specify which keys to keep and not
func Universal(predicate predicate.HclPredicate) PackFunc {

	return func(payload interface{}, d *schema.ResourceData) error {
		setValue := utilsdk.MkLens(d)

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
func Compose(packers ...PackFunc) PackFunc {
	return func(repo interface{}, d *schema.ResourceData) error {
		var errors []error

		for _, packer := range packers {
			err := packer(repo, d)
			if err != nil {
				errors = append(errors, err)
			}
		}
		if errors != nil && len(errors) > 0 {
			return fmt.Errorf("failed saving state %q", errors)
		}
		return nil
	}
}

func Default(skeema map[string]*schema.Schema) PackFunc {
	return Universal(
		predicate.All(
			predicate.SchemaHasKey(skeema),
			predicate.NoPassword,
			predicate.NoClass,
		),
	)
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
					utilsdk.FieldToHcl(field): result,
				}
			}
			return map[string]interface{}{}
		}
	case reflect.Slice:
		return func(field reflect.StructField, thing reflect.Value) map[string]interface{} {
			return map[string]interface{}{
				utilsdk.FieldToHcl(field): utilsdk.CastToInterfaceArr(thing.Interface().([]string)),
			}
		}
	}
	return func(field reflect.StructField, thing reflect.Value) map[string]interface{} {
		return map[string]interface{}{
			utilsdk.FieldToHcl(field): thing.Interface(),
		}
	}
}

func lookup(payload interface{}, pred predicate.HclPredicate) map[string]interface{} {

	if pred == nil {
		pred = predicate.True
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
			hcl := utilsdk.FieldToHcl(field)
			shouldLookup = pred(hcl)
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
