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

package util

import (
	"fmt"

	"github.com/samber/lo"
)

type JFrogErrors struct {
	Errors []JFrogError `json:"errors"`
}

func (e JFrogErrors) String() string {
	return lo.Reduce(
		e.Errors,
		func(agg string, err JFrogError, _ int) string {
			if agg == "" {
				return err.Message
			}

			return fmt.Sprintf("%s %s.", agg, err.Message)
		},
		"",
	)
}

type JFrogError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}
