package trustyapi

import (
	"log"
	"strings"
	"testing"
)

func TestProcessGoDependencies(t *testing.T) {
	ecosystem := "go"
	scoreThreshold := 5.0

	dependencies := []string{"github.com/alecthomas/units", "github.com/prometheus/client_golang", "github.com/prometheus/common", "github.com/Tinkoff/libvirt-exporter",
		"github.com/beorn7/perks", "golang.org/x/sys", "gopkg.in/alecthomas/kingpin.v2", "github.com/matttproud/golang_protobuf_extensions", "github.com/prometheus/client_model",
		"libvirt.org/go/libvirt", "github.com/alecthomas/template", "github.com/golang/protobuf", "github.com/prometheus/procfs"}
	expectedFail := []bool{false, false, false, true, true, true, true, true, false, true, true, false, false, true}

	for i, dep := range dependencies {
		log.Printf("Analyzing dependency: %s\n", dep)
		report, shouldFail := ProcessDependency(dep, ecosystem, scoreThreshold)
		if shouldFail != expectedFail[i] {
			t.Errorf("Dependency %s failed check unexpectedly, expected %v, got %v", dep, expectedFail[i], shouldFail)
		}
		if dep == "github.com/Tinkoff/libvirt-exporter" {
			if !strings.Contains(report, "Archived") {
				t.Errorf("Expected report to contain 'Archived' for %s", dep)
			}
		}
	}
}

func TestProcessDeprecatedDependencies(t *testing.T) {
	ecosystem := "npm"
	scoreThreshold := 10.0

	dependencies := []string{"@types/google-cloud__storage", "cutjs", "scriptoni", "stryker-mocha-framework", "grunt-html-smoosher", "moesif-express", "swagger-methods",
		"@syncfusion/ej2-heatmap", "@cnbritain/wc-buttons", "gulp-google-cdn"}

	for _, dep := range dependencies {
		log.Printf("Analyzing dependency: %s\n", dep)
		report, _ := ProcessDependency(dep, ecosystem, scoreThreshold)
		if !strings.Contains(report, "Deprecated") {
			t.Errorf("Expected report to contain 'Deprecated' for %s", dep)
		}
	}

}

func TestProcessMaliciousDependencies(t *testing.T) {
	ecosystem := "pypi"
	scoreThreshold := 10.0

	dependencies := []string{"lyft-service", "types-for-adobe", "booto3", "google-requests", "reqargs"}

	for _, dep := range dependencies {
		log.Printf("Analyzing dependency: %s\n", dep)
		report, _ := ProcessDependency(dep, ecosystem, scoreThreshold)
		if !strings.Contains(report, "Malicious") {
			t.Errorf("Expected report to contain 'Malicious' for %s", dep)
		}
	}

}
