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
	"log"
	"strings"

	"github.com/stacklok/trusty-action/pkg/types"
)

// ParsingFunction is a function type that takes a string as input and returns
// a slice of dependencies and an error.
type ParsingFunction func(string) ([]types.Dependency, error)

var parsingFunctions = map[string]ParsingFunction{
	"go.mod":           ParseGoMod,
	"Cargo.toml":       ParseCargoToml,
	"requirements.txt": ParseRequirementsTxt,
	"pom.xml":          ParsePomXml,
	"package.json":     ParsePackageJSON,
}

// Parse parses the given filename and content to extract dependencies and
// determine the ecosystem. It iterates through the available parsing function
// based on the file suffix and calls the appropriate function.
// If a matching parsing function is found, it returns the extracted dependencies,
// the determined ecosystem, and any error encountered.
// If no matching parsing function is found, it returns an empty slice of
// dependencies, "none" as the ecosystem, and no error.
func Parse(filename string, content string) ([]types.Dependency, string, error) {
	log.Printf("Parsing file: %s\n", filename)
	for suffix, function := range parsingFunctions {
		if strings.HasSuffix(filename, suffix) {
			deps, err := function(content)
			ecosystem := determineEcosystem(suffix)
			return deps, ecosystem, err
		}
	}
	return []types.Dependency{}, "none", nil
}

// determineEcosystem maps a file suffix to its ecosystem
func determineEcosystem(suffix string) string {
	switch suffix {
	case "go.mod":
		return "go"
	case "Cargo.toml":
		return "crates"
	case "requirements.txt":
		return "pypi"
	case "pom.xml":
		return "maven"
	case "package.json":
		return "npm"
	default:
		return "unknown"
	}
}
