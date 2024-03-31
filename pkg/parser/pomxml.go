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
	"encoding/xml"

	"github.com/stacklok/trusty-action/pkg/types"
)

// ParsePomXml parses the content of a POM XML file and returns a slice of dependencies
// along with their group ID, artifact ID, and version.
// The content parameter is the string representation of the POM XML file.
// The function returns a slice of types.Dependency and an error.
// If the parsing is successful, the slice of dependencies is populated and the error is nil.
// If an error occurs during parsing, the function returns nil for the slice of dependencies
// and the error that occurred.
func ParsePomXml(content string) ([]types.Dependency, error) {
	var project struct {
		Dependencies struct {
			Dependency []struct {
				GroupId    string `xml:"groupId"`
				ArtifactId string `xml:"artifactId"`
				Version    string `xml:"version"`
			} `xml:"dependency"`
		} `xml:"dependencies"`
	}
	if err := xml.Unmarshal([]byte(content), &project); err != nil {
		return nil, err
	}

	var deps []types.Dependency
	for _, d := range project.Dependencies.Dependency {
		deps = append(deps, types.Dependency{Name: d.GroupId + ":" + d.ArtifactId, Version: d.Version})
	}
	return deps, nil
}
