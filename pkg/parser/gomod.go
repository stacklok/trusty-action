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

// ParseGoMod parses the content of a go.mod file and returns a slice of dependencies.
// Each dependency is represented by a types.Dependency struct, containing the name and version.
func ParseGoMod(content string) ([]types.Dependency, error) {
	var deps []types.Dependency
	lines := strings.Split(content, "\n")
	inRequireBlock := false // Flag to track whether we are inside a require block

	for _, line := range lines {
		trimmedLine := strings.TrimSpace(line)
		if trimmedLine == "" || strings.HasPrefix(trimmedLine, "//") {
			continue // Skip empty lines and comments
		}

		if strings.HasPrefix(trimmedLine, "require (") {
			inRequireBlock = true // We're entering a require block
			continue
		}

		if trimmedLine == ")" && inRequireBlock {
			inRequireBlock = false // We're leaving a require block
			continue
		}

		if inRequireBlock || strings.HasPrefix(trimmedLine, "require") {
			parts := strings.Fields(trimmedLine)
			if len(parts) >= 2 {
				var depName, depVersion string
				if inRequireBlock {
					depName = parts[0]
					depVersion = parts[1]
				} else { // Inline require
					depName = parts[1]
					depVersion = parts[2]
				}
				deps = append(deps, types.Dependency{Name: depName, Version: depVersion})
			}
		}
	}

	return deps, nil
}
