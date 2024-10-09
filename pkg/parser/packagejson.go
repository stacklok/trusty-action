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
	"encoding/json"

	"github.com/stacklok/trusty-action/pkg/types"
)

// ParsePackageJSON parses package.json content and extracts dependencies
func ParsePackageJSON(content string) ([]types.Dependency, error) {
	var parsedContent struct {
		Dependencies map[string]string `json:"dependencies"`
	}
	if err := json.Unmarshal([]byte(content), &parsedContent); err != nil {
		return nil, err
	}

	deps := make([]types.Dependency, len(parsedContent.Dependencies))
	i := 0
	for name, version := range parsedContent.Dependencies {
		deps[i] = types.Dependency{Name: name, Version: version}
		i++
	}
	return deps, nil
}
