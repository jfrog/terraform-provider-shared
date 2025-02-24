package testutil

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"reflect"
	"strings"
	"testing"

	"github.com/go-resty/resty/v2"
	tfjson "github.com/hashicorp/terraform-json"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/samber/lo"
)

func RandomInt() int {
	return rand.Intn(10000000)
}

func RandBool() bool {
	return RandomInt()%2 == 0
}

func RandSelect(items ...interface{}) interface{} {
	return items[RandomInt()%len(items)]
}

func GetEnvVarWithFallback(t *testing.T, envVars ...string) string {
	envVarValue, err := schema.MultiEnvDefaultFunc(envVars, nil)()
	if envVarValue == "" || envVarValue == nil || err != nil {
		t.Fatalf("%s must be set for acceptance tests", strings.Join(envVars, " or "))
		return ""
	}

	return envVarValue.(string)
}

func MapToTestChecks(fqrn string, fields map[string]interface{}) []resource.TestCheckFunc {
	var result []resource.TestCheckFunc
	for key, value := range fields {
		switch reflect.TypeOf(value).Kind() {
		case reflect.Slice:
			for i, lv := range value.([]interface{}) {
				result = append(result, resource.TestCheckResourceAttr(
					fqrn,
					fmt.Sprintf("%s.%d", key, i),
					fmt.Sprintf("%v", lv),
				))
			}
		case reflect.Map:
			// this also gets generated, but it's value is '1', which is also the size. So, I don't know
			// what it means
			// content_synchronisation.0.%
			resource.TestCheckResourceAttr(
				fqrn,
				fmt.Sprintf("%s.#", key),
				fmt.Sprintf("%d", len(value.(map[string]interface{}))),
			)
		default:
			result = append(result, resource.TestCheckResourceAttr(fqrn, key, fmt.Sprintf(`%v`, value)))
		}
	}
	return result
}

type CheckFun func(id string, request *resty.Request) (*resty.Response, error)

func MkNames(name, resource string) (int, string, string) {
	id := RandomInt()
	n := fmt.Sprintf("%s%d", name, id)
	return id, fmt.Sprintf("%s.%s", resource, n), n
}

var ConfigPlanChecks = func(resourceName string) resource.ConfigPlanChecks {
	return resource.ConfigPlanChecks{
		PreApply: []plancheck.PlanCheck{
			DebugPlan(resourceName, "PreApply"),
		},
		PostApplyPreRefresh: []plancheck.PlanCheck{
			DebugPlan(resourceName, "PostApplyPreRefresh"),
		},
		PostApplyPostRefresh: []plancheck.PlanCheck{
			DebugPlan(resourceName, "PostApplyPostRefresh"),
		},
	}
}

var _ plancheck.PlanCheck = PlanCheck{}

type PlanCheck struct {
	Stage        string
	ResourceName string
}

func (p PlanCheck) CheckPlan(ctx context.Context, req plancheck.CheckPlanRequest, resp *plancheck.CheckPlanResponse) {
	var err error

	rc, err := json.Marshal(req.Plan.ResourceChanges[0])
	if err != nil {
		resp.Error = err
		return
	}

	pv, err := json.Marshal(req.Plan.PlannedValues)
	if err != nil {
		resp.Error = err
		return
	}

	ps, err := json.Marshal(req.Plan.PriorState)
	if err != nil {
		resp.Error = err
		return
	}

	rd, err := json.Marshal(req.Plan.ResourceDrift)
	if err != nil {
		resp.Error = err
		return
	}

	tflog.Debug(ctx, "CheckPlan", map[string]interface{}{
		"stage":                                  p.Stage,
		"req.Plan.ResourceChanges.ResourceDrift": string(rd),
		"req.Plan.ResourceChanges":               string(rc),
		"req.Plan.PlannedValues":                 string(pv),
		"req.Plan.PriorState":                    string(ps),
	})

	if len(req.Plan.ResourceDrift) > 0 {
		drifts := lo.FilterMap(req.Plan.ResourceDrift, func(c *tfjson.ResourceChange, index int) (string, bool) {
			tflog.Debug(ctx, "CheckPlan", map[string]interface{}{
				"p.ResourceName":                         p.ResourceName,
				"fmt.Sprintf(\"%s.%s\", c.Type, c.Name)": fmt.Sprintf("%s.%s", c.Type, c.Name),
			})
			driftsMessage := fmt.Sprintf("Name: %s.%s\n\nBefore: %v\n\nAfter: %v\n", c.Type, c.Name, c.Change.Before, c.Change.After)
			shouldInclude := p.ResourceName == "" || p.ResourceName == fmt.Sprintf("%s.%s", c.Type, c.Name)
			return driftsMessage, shouldInclude
		})
		if len(drifts) > 0 {
			resp.Error = fmt.Errorf("expected empty plan, but has resource drift(s):\n\n%v", strings.Join(drifts, "\n\n"))
		}
		return
	}

	var errStrings []string
	for _, change := range req.Plan.OutputChanges {
		if !change.Actions.NoOp() {
			errStrings = append(errStrings, fmt.Sprintf("expected empty plan, but %s has output change(s):\n\nbefore: %v\n\nafter: %v\n\nunknown: %v", change.Actions, change.Before, change.After, change.AfterUnknown))
		}
	}

	for _, rc := range req.Plan.ResourceChanges {
		if !rc.Change.Actions.NoOp() {
			errStrings = append(errStrings, fmt.Sprintf("expected empty plan, but %s has planned action(s): %v\n\nbefore: %v\n\nafter: %v\n\nunknown: %v", rc.Address, rc.Change.Actions, rc.Change.Before, rc.Change.After, rc.Change.AfterUnknown))
		}
	}

	if len(errStrings) > 0 {
		resp.Error = fmt.Errorf("%s", strings.Join(errStrings, "\n"))
		return
	}
}

func DebugPlan(resourceName, stage string) plancheck.PlanCheck {
	return PlanCheck{
		ResourceName: resourceName,
		Stage:        stage,
	}
}
