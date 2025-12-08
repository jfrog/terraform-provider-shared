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

package fw

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

type requireIfDefinedValidator struct {
	expressions path.Expressions
}

// Description returns a plaintext string describing the validator.
func (v requireIfDefinedValidator) Description(_ context.Context) string {
	return "required if parent attribute is defined/not null"
}

// MarkdownDescription returns a Markdown formatted string describing the validator.
func (v requireIfDefinedValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

// ValidateBool performs the validation logic for the validator.
func (v requireIfDefinedValidator) ValidateObject(ctx context.Context, req validator.ObjectRequest, resp *validator.ObjectResponse) {
	// If the current attribute configuration is null or unknown, there
	// cannot be any value comparisons, so exit early without error.
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}
	// Combine the given path expressions with the current attribute path
	// expression. This call automatically handles relative and absolute
	// expressions.
	expressions := req.PathExpression.MergeExpressions(v.expressions...)

	// For each expression, find matching paths.
	for _, expression := range expressions {
		// Find paths matching the expression in the configuration data.
		matchedPaths, diags := req.Config.PathMatches(ctx, expression)
		// Collect all errors
		resp.Diagnostics.Append(diags...)
		if diags.HasError() {
			continue
		}

		// For each matched path, get the data and compare.
		for _, matchedPath := range matchedPaths {
			// Fetch the generic attr.Value at the given path. This ensures any
			// potential parent value of a different type, which can be a null
			// or unknown value, can be safely checked without raising a type
			// conversion error.
			var matchedPathValue attr.Value

			diags = req.Config.GetAttribute(ctx, matchedPath, &matchedPathValue)
			// Collect all errors
			resp.Diagnostics.Append(diags...)
			if diags.HasError() {
				continue
			}

			// If the matched path value is null or unknown, we cannot compare
			// values, so continue to other matched paths.
			if matchedPathValue.IsNull() || matchedPathValue.IsUnknown() {
				resp.Diagnostics.AddAttributeError(
					matchedPath,
					"Invalid Attribute Value",
					fmt.Sprintf("Attribute %s must be set when %s is defined.", matchedPath, req.Path),
				)
			}
		}
	}
}

// SetInSlice tests if the provided value matches the value of an element
// in the valid slice
func RequireIfDefined(expressions ...path.Expression) validator.Object {
	return &requireIfDefinedValidator{
		expressions: expressions,
	}
}
