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
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-validators/helpers/validatordiag"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// Ensure our implementation satisfies the validator.String interface.
var _ validator.String = &regexNotMatchesValidator{}

type regexNotMatchesValidator struct {
	regexp  *regexp.Regexp
	message string
}

func (v regexNotMatchesValidator) Description(_ context.Context) string {
	if v.message != "" {
		return v.message
	}
	return fmt.Sprintf("value must not match regular expression '%s'", v.regexp)
}

func (v regexNotMatchesValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v regexNotMatchesValidator) ValidateString(ctx context.Context, request validator.StringRequest, response *validator.StringResponse) {
	if request.ConfigValue.IsNull() || request.ConfigValue.IsUnknown() {
		return
	}

	value := request.ConfigValue.ValueString()

	if v.regexp.MatchString(value) {
		response.Diagnostics.Append(validatordiag.InvalidAttributeValueMatchDiagnostic(
			request.Path,
			v.Description(ctx),
			value,
		))
	}
}

// RegexNotMatches returns an AttributeValidator which ensures that any configured
// attribute value:
//
//   - Is a string.
//   - Not matches the given regular expression https://github.com/google/re2/wiki/Syntax.
//
// Null (unconfigured) and unknown (known after apply) values are skipped.
// Optionally an error message can be provided to return something friendlier
// than "value must not match regular expression 'regexp'".
func RegexNotMatches(regexp *regexp.Regexp, message string) validator.String {
	return regexNotMatchesValidator{
		regexp:  regexp,
		message: message,
	}
}
