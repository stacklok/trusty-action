package parser

import (
	"reflect"
	"sort"
	"testing"

	"github.com/stacklok/trusty-action/pkg/types"
)

func TestParse(t *testing.T) {
	tests := []struct {
		filename  string
		content   string
		expected  []types.Dependency
		ecosystem string
		err       error
	}{
		{
			filename:  "meh.txt",
			content:   "some content",
			expected:  []types.Dependency{},
			ecosystem: "none",
			err:       nil,
		},
		{
			filename: "package.json",
			content:  "{\"dependencies\": {\"express\": \"^4.17.1\", \"lodash\": \"^4.17.21\"}}",
			expected: []types.Dependency{
				{Name: "express", Version: "^4.17.1"},
				{Name: "lodash", Version: "^4.17.21"},
			},
			ecosystem: "npm",
			err:       nil,
		},
		{
			filename: "go.mod",
			content:  "module example.com\n\ngo 1.16\n\nrequire (\n\tgithub.com/google/go-github/v60 v60.0.0\n)",
			expected: []types.Dependency{
				{Name: "github.com/google/go-github/v60", Version: "v60.0.0"},
			},
			ecosystem: "go",
			err:       nil,
		},
		{
			filename: "Cargo.toml",
			content:  "[dependencies]\nrand = \"0.8.4\"\n",
			expected: []types.Dependency{
				{Name: "rand", Version: "0.8.4"},
			},
			ecosystem: "crates",
			err:       nil,
		},
		{
			filename: "requirements.txt",
			content:  "requests==2.25.1\n",
			expected: []types.Dependency{
				{Name: "requests", Version: "2.25.1"},
			},
			ecosystem: "pypi",
			err:       nil,
		},
		{
			filename: "pom.xml",
			content:  "<project>\n\t<dependencies>\n\t\t<dependency>\n\t\t\t<groupId>org.apache.maven</groupId>\n\t\t\t<artifactId>maven-core</artifactId>\n\t\t\t<version>3.8.1</version>\n\t\t</dependency>\n\t</dependencies>\n</project>",
			expected: []types.Dependency{
				{Name: "org.apache.maven:maven-core", Version: "3.8.1"},
			},
			ecosystem: "maven",
			err:       nil,
		},
	}

	for _, test := range tests {
		deps, ecosystem, err := Parse(test.filename, test.content)

		sort.Slice(deps, func(i, j int) bool { return deps[i].Name < deps[j].Name })
		sort.Slice(test.expected, func(i, j int) bool { return test.expected[i].Name < test.expected[j].Name })

		if !reflect.DeepEqual(deps, test.expected) {
			t.Errorf("Expected dependencies %v, but got %v", test.expected, deps)
		}
		if ecosystem != test.ecosystem {
			t.Errorf("Expected ecosystem %s, but got %s", test.ecosystem, ecosystem)
		}
		if err != test.err {
			t.Errorf("Expected error %v, but got %v", test.err, err)
		}
	}
}
