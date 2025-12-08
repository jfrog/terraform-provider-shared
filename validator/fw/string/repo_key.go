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
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/helpers/validatordiag"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// Ensure our implementation satisfies the validator.String interface.
var _ validator.String = &repoKeyValidator{}

type repoKeyValidator struct{}

func (v repoKeyValidator) Description(_ context.Context) string {
	return "value must be a valid email address"
}

func (v repoKeyValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v repoKeyValidator) ValidateString(ctx context.Context, request validator.StringRequest, response *validator.StringResponse) {
	if request.ConfigValue.IsNull() || request.ConfigValue.IsUnknown() {
		return
	}

	value := request.ConfigValue.ValueString()

	if len(value) == 0 || len(value) > 64 {
		response.Diagnostics.Append(validatordiag.InvalidAttributeValueLengthDiagnostic(
			request.Path,
			"must be 1 - 64 alphanumeric and hyphen characters",
			value,
		))
	}

	if strings.ContainsAny(value, " !@#$%^&*()+={}[]:;<>,/?~`|\\") {
		response.Diagnostics.Append(validatordiag.InvalidAttributeValueMatchDiagnostic(
			request.Path,
			"cannot contain spaces or special characters: !@#$%^&*()+={}[]:;<>,/?~`|\\",
			value,
		))
	}
}

func RepoKey() validator.String {
	return repoKeyValidator{}
}
