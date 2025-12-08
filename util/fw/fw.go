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
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ValueToSetType ensures we have a types.Set literal
func ValueToSetType(v attr.Value) types.Set {
	if vb, ok := v.(types.Set); ok {
		return vb
	}
	if vb, ok := v.(*types.Set); ok {
		return *vb
	}
	panic(fmt.Sprintf("cannot pass type %T to conv.ValueToSetType", v))
}

// AttributeValueToString will attempt to execute the appropriate AttributeStringerFunc from the ones registered.
func AttributeValueToString(v attr.Value) string {
	if s, ok := v.(types.String); ok {
		return s.ValueString()
	}
	return v.String()
}

func StringSetToStrings(v attr.Value) []string {
	vt := ValueToSetType(v)
	out := make([]string, len(vt.Elements()))
	for i, ve := range vt.Elements() {
		out[i] = AttributeValueToString(ve)
	}
	return out
}

func UnableToCreateResourceError(resp *resource.CreateResponse, err string) {
	resp.Diagnostics.AddError(
		"Unable to Create Resource",
		"An unexpected error occurred while creating the resource update request. "+
			"Please report this issue to the provider developers.\n\n"+
			"Error: "+err,
	)
}

func UnableToUpdateResourceError(resp *resource.UpdateResponse, err string) {
	resp.Diagnostics.AddError(
		"Unable to Update Resource",
		"An unexpected error occurred while updating the resource update request. "+
			"Please report this issue to the provider developers.\n\n"+
			"Error: "+err,
	)
}

func UnableToRefreshResourceError(resp *resource.ReadResponse, err string) {
	resp.Diagnostics.AddError(
		"Unable to Refresh Resource",
		"An unexpected error occurred while attempting to refresh resource state. "+
			"Please retry the operation or report this issue to the provider developers.\n\n"+
			"Error: "+err,
	)
}

func UnableToDeleteResourceError(resp *resource.DeleteResponse, err string) {
	resp.Diagnostics.AddError(
		"Unable to Delete Resource",
		"An unexpected error occurred while attempting to delete the resource. "+
			"Please retry the operation or report this issue to the provider developers.\n\n"+
			"Error: "+err,
	)
}
