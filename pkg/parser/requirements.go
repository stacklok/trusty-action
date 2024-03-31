//
// Copyright 2024 Stacklok, Inc.
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

package parser

import (
	"strings"

	"github.com/stacklok/trusty-action/pkg/types"
)

// ParseRequirementsTxt parses requirements.txt content and extracts dependencies
func ParseRequirementsTxt(content string) ([]types.Dependency, error) {
	var deps []types.Dependency
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" && !strings.HasPrefix(line, "#") {
			parts := strings.SplitN(line, "==", 2)
			if len(parts) == 2 {
				// Convert package name to lowercase
				packageName := strings.ToLower(parts[0])
				deps = append(deps, types.Dependency{Name: packageName, Version: parts[1]})
			}
		}
	}
	return deps, nil
}
