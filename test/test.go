package test

import (
	"bytes"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"math/rand"
	"reflect"
	"strings"
	"testing"
	"text/template"
	"time"
)

func RandomInt() int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(10000000)
}

func RandBool() bool {
	return RandomInt()%2 == 0
}

func RandSelect(items ...interface{}) interface{} {
	return items[RandomInt()%len(items)]
}

func ExecuteTemplate(name, temp string, fields interface{}) string {
	var tpl bytes.Buffer
	if err := template.Must(template.New(name).Parse(temp)).Execute(&tpl, fields); err != nil {
		panic(err)
	}

	return tpl.String()
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
