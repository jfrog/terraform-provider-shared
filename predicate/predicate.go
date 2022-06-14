package predicate

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

type HclPredicate func(hcl string) bool

func All(predicates ...HclPredicate) HclPredicate {
	return func(hcl string) bool {
		for _, predicate := range predicates {
			if !predicate(hcl) {
				return false
			}
		}
		return true
	}
}
func Any(predicates ...HclPredicate) HclPredicate {
	return func(hcl string) bool {
		for _, predicate := range predicates {
			if predicate(hcl) {
				return true
			}
		}
		return false
	}
}
var True = func(hcl string) bool {
	return true
}
var NoClass = Ignore("class", "rclass")
var NoPassword = Ignore("password")


func Ignore(names ...string) HclPredicate {
	set := map[string]interface{}{}
	for _, name := range names {
		set[name] = nil
	}
	return func(hcl string) bool {
		_, found := set[hcl]
		return !found
	}
}
func SchemaHasKey(skeema map[string]*schema.Schema) HclPredicate {
	return func(key string) bool {
		_, ok := skeema[key]
		return ok
	}
}
