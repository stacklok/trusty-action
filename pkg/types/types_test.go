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

package types

import (
	"reflect"
	"testing"
)

// TestConvertDepsToMap tests the ConvertDepsToMap function for converting a slice of Dependency structs to a map
func TestConvertDepsToMap(t *testing.T) {
	deps := []Dependency{
		{Name: "dep1", Version: "1.0"},
		{Name: "dep2", Version: "2.0"},
	}

	expected := map[string]string{
		"dep1": "1.0",
		"dep2": "2.0",
	}

	result := ConvertDepsToMap(deps)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

// TestDiffDependencies tests the DiffDependencies function for finding added dependencies
func TestDiffDependencies(t *testing.T) {
	oldDeps := map[string]string{
		"dep1": "1.0",
	}

	newDeps := map[string]string{
		"dep1": "1.0",
		"dep2": "2.0",
	}

	expectedAdded := map[string]string{
		"dep2": "2.0",
	}

	addedDeps := DiffDependencies(oldDeps, newDeps)
	if !reflect.DeepEqual(addedDeps, expectedAdded) {
		t.Errorf("Expected added dependencies %v, got %v", expectedAdded, addedDeps)
	}

	// Test case where there are no new dependencies added
	newDepsSame := map[string]string{
		"dep1": "1.0",
	}

	expectedEmpty := map[string]string{}

	addedDepsSame := DiffDependencies(oldDeps, newDepsSame)
	if !reflect.DeepEqual(addedDepsSame, expectedEmpty) {
		t.Errorf("Expected no added dependencies, got %v", addedDepsSame)
	}
}
