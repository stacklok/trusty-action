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
	"github.com/BurntSushi/toml"
	"github.com/stacklok/trusty-action/pkg/types"
)

// ParseCargoToml parses the content of a Cargo.toml file and returns a slice of dependencies.
// It takes a string parameter `content` which represents the content of the Cargo.toml file.
// The function returns a slice of `types.Dependency` and an error if any occurred during parsing.
func ParseCargoToml(content string) ([]types.Dependency, error) {
	var conf struct {
		Dependencies map[string]string `toml:"dependencies"`
	}
	if _, err := toml.Decode(content, &conf); err != nil {
		return nil, err
	}

	deps := make([]types.Dependency, len(conf.Dependencies))
	i := 0
	for name, version := range conf.Dependencies {
		deps[i] = types.Dependency{Name: name, Version: version}
		i++
	}
	return deps, nil
}
