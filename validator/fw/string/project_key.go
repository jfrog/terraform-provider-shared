// Copyright (c) JFrog Ltd. (2025)
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package string

import (
	"context"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-validators/helpers/validatordiag"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// Ensure our implementation satisfies the validator.String interface.
var _ validator.String = &projectKeyValidator{}

type projectKeyValidator struct{}

func (v projectKeyValidator) Description(_ context.Context) string {
	return "value must be a valid email address"
}

func (v projectKeyValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v projectKeyValidator) ValidateString(ctx context.Context, request validator.StringRequest, response *validator.StringResponse) {
	if request.ConfigValue.IsNull() || request.ConfigValue.IsUnknown() {
		return
	}

	value := request.ConfigValue.ValueString()
	re := regexp.MustCompile(`^[a-z][a-z0-9\-]{1,31}$`)
	if !re.MatchString(value) {
		response.Diagnostics.Append(validatordiag.InvalidAttributeValueMatchDiagnostic(
			request.Path,
			"must be 2 - 32 lowercase alphanumeric and hyphen characters",
			value,
		))
	}
}

func ProjectKey() validator.String {
	return projectKeyValidator{}
}
