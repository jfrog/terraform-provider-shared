package fw

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure our implementation satisfies the validator.Bool interface.
var _ validator.Bool = &booleanValidator{}

type booleanValidator struct {
	conflictingBools bool
	expressions      path.Expressions
}

// Description returns a plaintext string describing the validator.
func (v booleanValidator) Description(_ context.Context) string {
	return fmt.Sprintf("The boolean attributes can not be set to %t for all the following attributes: %s", v.conflictingBools, v.expressions)
}

// MarkdownDescription returns a Markdown formatted string describing the validator.
func (v booleanValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

// ValidateBool performs the validation logic for the validator.
func (v booleanValidator) ValidateBool(ctx context.Context, req validator.BoolRequest, resp *validator.BoolResponse) {
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
		resp.Diagnostics.Append(diags...)
		// Collect all errors
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

			resp.Diagnostics.Append(diags...)

			// Collect all errors
			if diags.HasError() {
				continue
			}

			// If the matched path value is null or unknown, we cannot compare
			// values, so continue to other matched paths.
			if matchedPathValue.IsNull() || matchedPathValue.IsUnknown() {
				continue
			}

			// Now that we know the matched path value is not null or unknown,
			// it is safe to attempt converting it to the intended attr.Value
			// implementation, in this case a types.Bool value.
			var matchedPathConfig types.Bool

			diags = tfsdk.ValueAs(ctx, matchedPathValue, &matchedPathConfig)

			resp.Diagnostics.Append(diags...)

			// If the matched path value was not able to be converted from
			// attr.Value to the intended types.Bool implementation, it most
			// likely means that the path expression was not pointing at a
			// types.BoolType attribute. Collect the error and continue to
			// other matched paths.
			if diags.HasError() {
				continue
			}

			if matchedPathConfig.ValueBool() == v.conflictingBools && req.ConfigValue.ValueBool() == v.conflictingBools {
				resp.Diagnostics.AddAttributeError(
					matchedPath,
					"Invalid Attribute Value",
					fmt.Sprintf("Attribute %s can not be set to %t, when %s is %t.", req.Path, req.ConfigValue.ValueBool(), matchedPath.Steps(), matchedPathConfig.ValueBool()),
				)
			}
		}
	}
}

// BoolConflict checks that any Bool values in the paths described by the
// path.Expression are less than the current attribute value.
func BoolConflict(conflictingBools bool, expressions ...path.Expression) validator.Bool {
	return &booleanValidator{
		conflictingBools: conflictingBools,
		expressions:      expressions,
	}
}
