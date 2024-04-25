package trustyapi

import (
	"fmt"
	"log"
	"strings"
	"testing"
)

func TestProcessDependencies(t *testing.T) {
	fmt.Println("in test processing")
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
