package trustyapi

import (
	"log"
	"strings"
	"testing"
	"time"
)

func TestReportBuilder(t *testing.T) {
	dependencies := []string{"next", "react", "bugsnagmw", "scriptoni", "notifyjs"}

	result, failAction := GenerateReportContent(dependencies, "npm", 5.0, 5.0, 5.0, 5.0, 5.0, true, true, true)
	// fmt.Println(result) // this is normally used to display and validate the report output, uncomment for debugging
	if result == "" {
		t.Errorf("Report is empty")
	}
	if !failAction {
		t.Errorf("Fail action is false")
	}
}

func TestProcessGoDependencies(t *testing.T) {
	ecosystem := "go"
	scoreThreshold := 5.0
	repoActivityThreshold := 5.0
	authorActivityThreshold := 5.0
	provenanceThreshold := 5.0
	typosquattingThreshold := 5.0

	dependencies := []string{"github.com/alecthomas/units", "github.com/prometheus/client_golang", "github.com/prometheus/common", "github.com/Tinkoff/libvirt-exporter",
		"github.com/beorn7/perks", "golang.org/x/sys", "gopkg.in/alecthomas/kingpin.v2", "github.com/matttproud/golang_protobuf_extensions", "github.com/prometheus/client_model",
		"libvirt.org/go/libvirt", "github.com/alecthomas/template", "github.com/golang/protobuf", "github.com/prometheus/procfs"}
	expectedFail := []bool{true, false, false, true, true, true, true, true, true, true, true, false, false, true}

	for i, dep := range dependencies {
		log.Printf("Analyzing dependency: %s\n", dep)
		report, shouldFail := ProcessDependency(dep, ecosystem, repoActivityThreshold, authorActivityThreshold, provenanceThreshold, typosquattingThreshold, scoreThreshold, true, true, true)
		if shouldFail != expectedFail[i] {
			t.Errorf("Dependency %s failed check unexpectedly, expected %v, got %v", dep, expectedFail[i], shouldFail)
		}
		if dep == "github.com/Tinkoff/libvirt-exporter" {
			if !strings.Contains(report, "Archived") {
				t.Errorf("Expected report to contain 'Archived' for %s", dep)
			}
		}
		time.Sleep(1 * time.Second)
	}
}

func TestProcessDeprecatedDependencies(t *testing.T) {
	ecosystem := "npm"
	scoreThreshold := 5.0

	dependencies := []string{"@types/google-cloud__storage", "cutjs", "scriptoni", "stryker-mocha-framework", "grunt-html-smoosher", "moesif-express", "swagger-methods",
		"@syncfusion/ej2-heatmap", "@cnbritain/wc-buttons", "gulp-google-cdn"}

	for _, dep := range dependencies {
		log.Printf("Analyzing dependency: %s\n", dep)
		report, _ := ProcessDependency(dep, ecosystem, scoreThreshold, 0.0, 0.0, 0.0, 0.0, true, true, true)
		if !strings.Contains(report, "Deprecated") {
			t.Errorf("Expected report to contain 'Deprecated' for %s", dep)
		}
	}

}

func TestProcessMaliciousDependencies(t *testing.T) {
	ecosystem := "pypi"
	scoreThreshold := 5.0

	dependencies := []string{"lyft-service", "types-for-adobe", "reqargs"}

	for _, dep := range dependencies {
		log.Printf("Analyzing dependency: %s\n", dep)
		report, _ := ProcessDependency(dep, ecosystem, scoreThreshold, 0.0, 0.0, 0.0, 0.0, true, true, true)
		if !strings.Contains(report, "Malicious") {
			t.Errorf("Expected report to contain 'Malicious' for %s", dep)
		}
	}

}

func TestProcessSigstoreProvenance(t *testing.T) {
	ecosystem := "npm"
	scoreThreshold := 5.0

	report, _ := ProcessDependency("sigstore", ecosystem, scoreThreshold, 0.0, 0.0, 0.0, 0.0, true, true, true)
	if !strings.Contains(report, "sigstore") {
		t.Errorf("Expected report to contain 'sigstore'")
	}
	if !strings.Contains(report, "Source repo: `https://github.com/sigstore/sigstore-js`") {
		t.Errorf("Source repo not matching")
	}
	if !strings.Contains(report, "Github Action Workflow: `.github/workflows/release.yml`") {
		t.Errorf("Github workflow not matching")
	}
	if !strings.Contains(report, "Issuer: `CN=sigstore-intermediate,O=sigstore.dev`") {
		t.Errorf("Issuer not matching")
	}
}

func TestProcessHistoricalProvenance(t *testing.T) {
	ecosystem := "npm"
	scoreThreshold := 5.0

	report, _ := ProcessDependency("openpgp", ecosystem, scoreThreshold, 0.0, 0.0, 0.0, 0.0, true, true, true)
	if !strings.Contains(report, "Number of versions") {
		t.Errorf("Versions for historical provenance not populated")
	}
	if !strings.Contains(report, "Number of Git Tags/Releases") {
		t.Errorf("Tags for historical provenance not populated")
	}
	if !strings.Contains(report, "Number of versions matched to Git Tags/Releases") {
		t.Errorf("Matched for historical provenance not populated")
	}

}
