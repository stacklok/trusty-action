package trustyapi

import (
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
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
	// TODO(puerco): This test is really load testing the trusty
	// API but does not really provide much value. I will remove
	// it in a future iteration.
	t.Parallel()
	const (
		ecosystem               = "go"
		scoreThreshold          = 5.0
		repoActivityThreshold   = 5.0
		authorActivityThreshold = 5.0
		provenanceThreshold     = 5.0
		typosquattingThreshold  = 5.0
	)
	for _, tc := range []struct {
		pkg          string
		expectedFail bool
	}{
		{"github.com/alecthomas/units", true},
		{"github.com/prometheus/client_golang", false},
		{"github.com/prometheus/common", false},
		{"github.com/Tinkoff/libvirt-exporter", true},
		{"github.com/beorn7/perks", true},
		{"golang.org/x/sys", true},
		{"gopkg.in/alecthomas/kingpin.v2", true},
		{"github.com/matttproud/golang_protobuf_extensions", true},
		{"github.com/prometheus/client_model", true},
		{"libvirt.org/go/libvirt", true},
		{"github.com/alecthomas/template", true},
		{"github.com/golang/protobuf", false},
		{"github.com/prometheus/procfs", false},
	} {
		tc := tc
		t.Run(fmt.Sprintf("test-%s", tc.pkg), func(t *testing.T) {
			t.Parallel()
			_, shouldFail := ProcessDependency(tc.pkg, ecosystem, repoActivityThreshold, authorActivityThreshold, provenanceThreshold, typosquattingThreshold, scoreThreshold, true, true, true)
			require.Truef(t, shouldFail == tc.expectedFail, "Dependency %s failed check unexpectedly, expected %v, got %v", tc.pkg, tc.expectedFail, shouldFail)

			// Temporarily disabling this checl as the package is not being
			// ingested by trusty at the moment:
			//nolint:gocritic
			/*
				if tc.pkg == "github.com/Tinkoff/libvirt-exporter" {
					require.Truef(
						t, strings.Contains(report, "Archived"),
						"Expected report to contain 'Archived' for %s", tc.pkg,
					)
				}
			*/
		})
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

	if !strings.Contains(report, "https://www.trustypkg.dev/npm/sigstore") {
		t.Errorf("Link to package page not found")
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
