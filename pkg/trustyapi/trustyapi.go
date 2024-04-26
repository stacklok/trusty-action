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

package trustyapi

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"net/http"
	"net/url"
	"strings"

	"github.com/google/go-github/v60/github"
)

// BuildReport analyzes the dependencies of a PR and generates a report based on their Trusty scores.
// It takes the following parameters:
// - ctx: The context.Context for the function.
// - ghClient: A pointer to a github.Client for interacting with the GitHub API.
// - owner: The owner of the repository.
// - repo: The name of the repository.
// - prNumber: The number of the pull request.
// - dependencies: A slice of strings representing the dependencies to be analyzed.
// - ecosystem: The ecosystem of the dependencies (e.g., "npm", "pip", "maven").
// - scoreThreshold: The threshold for Trusty scores below which a warning will be generated.
//
// The function generates a report and posts it as a comment on the pull request.
func BuildReport(ctx context.Context,
	ghClient *github.Client,
	owner,
	repo string,
	prNumber int,
	dependencies []string,
	ecosystem string,
	globalThreshold float64,
	activithTreshold float64,
	provenanceThreshold float64,
	typosquattingThreshold float64) {

	var (
		reportBuilder strings.Builder
		failAction    bool // Flag to track if the GitHub Action should fail
	)

	reportHeader := "## 🐻 Trusty Dependency Analysis Action Report \n\n"

	reportBuilder.WriteString(reportHeader)

	warningMessage := fmt.Sprintf("#### The following dependencies have Trusty scores below the set threshold of `%.2f`:\n\n", scoreThreshold)
	reportBuilder.WriteString(warningMessage)

	// The following loop generates the report for each dependency and then adds
	// it to the existing reportBuilder, between the header and footer.
	for _, dep := range dependencies {
		log.Printf("Analyzing dependency: %s\n", dep)
		report, shouldFail := ProcessDependency(dep, ecosystem, scoreThreshold)
		// Check if the report is not just whitespace
		if strings.TrimSpace(report) != "" {
			reportBuilder.WriteString(report)
		}
		if shouldFail {
			failAction = true // Update this line
		}
	}

	reportFooter := "> 🌟 If you like this action, why not try out [Minder](https://github.com/stacklok/minder), the secure supply chain platform. It has vastly more protections and is also free (as in :beer:) to opensource projects.\n"
	reportBuilder.WriteString(reportFooter)

	// Build the comment body from the report builder
	commentBody := reportBuilder.String()

	// Trim whitespace for accurate comparison
	trimmedCommentBody := strings.TrimSpace(commentBody)
	trimmedHeaderAndFooter := strings.TrimSpace(reportHeader + warningMessage + reportFooter)

	// Check if the comment body has more content than just the header and footer combined
	if len(trimmedCommentBody) > len(trimmedHeaderAndFooter) {
		_, _, err := ghClient.Issues.CreateComment(ctx, owner, repo, prNumber, &github.IssueComment{Body: &commentBody})
		log.Printf("posting comment to PR: %s/%s#%d\n", owner, repo, prNumber)
		if err != nil {
			log.Printf("error posting comment to PR: %v\n", err)
		}
	} else {
		log.Println("No report content to post, skipping comment.")
	}

	if failAction {
		log.Println("Failing the GitHub Action due to dependencies not meeting the required criteria.")
		os.Exit(1)
	}

}

// processDependency analyzes a dependency by making an API request to TrustyPkg.dev and returns a formatted report.
// It takes the dependency name, ecosystem, and score threshold as input parameters.
// The function constructs the query URL, makes the API request, and processes the response.
// If the Trusty score of the dependency is above the score threshold, it skips the dependency.
// Otherwise, it formats the report using Markdown and includes information about the dependency's Trusty score,
// whether it is malicious, deprecated or archived, and recommended alternative packages if available.
// The function returns the formatted report as a string.
func ProcessDependency(dep string, ecosystem string, scoreThreshold float64) (string, bool) {
	var reportBuilder strings.Builder
	shouldFail := false

	// Construct the query URL for the API request
	baseURL := "https://api.trustypkg.dev/v1/report"
	queryParams := url.Values{}
	queryParams.Add("package_name", dep)
	queryParams.Add("package_type", ecosystem)
	requestURL := fmt.Sprintf("%s?%s", baseURL, queryParams.Encode())

	// Create a channel for receiving the result
	resultChan := make(chan Package)

	// Make the API request to api.trustypkg.dev
	fetchPackageData(requestURL, dep, ecosystem, resultChan)

	// Wait for the result
	result, ok := <-resultChan
	if !ok {
		log.Printf("Error processing dependency: %s\n", dep)
	} else {
		// Process the result
		log.Printf("Processing result for dependency: %s\n", dep)
	}

	// Format the report using Markdown
	reportBuilder.WriteString(fmt.Sprintf("### :package: Dependency: [`%s`](https://www.trustypkg.dev/%s/%s)\n", dep, ecosystem, dep))
	// Highlight if the package is malicious, deprecated or archived
	if result.PackageData.Origin == "malicious" {
		reportBuilder.WriteString("### **⚠️ Malicious** (This package is marked as Malicious. Proceed with extreme caution!)\n\n")
	}
	if result.PackageData.IsDeprecated {
		reportBuilder.WriteString("### **⚠️ Deprecated** (This package is marked as Deprecated. Proceed with caution!)\n\n")
	}

	if result.PackageData.Archived {
		reportBuilder.WriteString("### **⚠️ Archived** (This package is marked as Archived. Proceed with caution!)\n\n")
	}

	reportBuilder.WriteString(fmt.Sprintf("### 📉 Trusty Score: `%.2f`\n", result.Summary.Score))

	// write provenance information
	if result.Provenance.Description.Provenance.Issuer != "" {
		reportBuilder.WriteString("### ![Sigstore](https://www.trustypkg.dev/icons/sigstore-horizontal.svg) Proof of origin (Provenance):\n")
		reportBuilder.WriteString("Built and signed with sigstore using GitHub Actions.\n")
		reportBuilder.WriteString(fmt.Sprintf("· Source repo: `%s`\n", result.Provenance.Description.Provenance.SourceRepo))
		reportBuilder.WriteString(fmt.Sprintf("· Github Action Workflow: `%s`\n", result.Provenance.Description.Provenance.Workflow))
		reportBuilder.WriteString(fmt.Sprintf("· Issuer: `%s`\n", result.Provenance.Description.Provenance.Issuer))
		reportBuilder.WriteString(fmt.Sprintf("· Rekor Public Ledger: `%s`\n", result.Provenance.Description.Provenance.Transparency))
	} else {
		// need to write regular provenance info
		reportBuilder.WriteString("### :key: Proof of origin (Provenance):\n")
		reportBuilder.WriteString(fmt.Sprintf("· Number of versions: %.0f\n", result.Provenance.Description.Hp.Versions))
		reportBuilder.WriteString(fmt.Sprintf("· Number of Git Tags/Releases: %.0f\n", result.Provenance.Description.Hp.Tags))
		reportBuilder.WriteString(fmt.Sprintf("· Number of versions matched to Git Tags/Releases: %.0f\n", result.Provenance.Description.Hp.Common))
	}
	reportBuilder.WriteString("[Learn more about source of origin provenance](https://docs.stacklok.com/trusty/understand/provenance)\n")

	// Include alternative packages in a Markdown table if available and if the package is deprecated, archived or malicious
	if result.Alternatives.Packages != nil && len(result.Alternatives.Packages) > 0 {
		reportBuilder.WriteString("### :bulb: Recommended Alternative Packages\n")
		reportBuilder.WriteString("| Package | Score | Trusty Link |\n")
		reportBuilder.WriteString("| ------- | ----- | ---------- |\n")
		for _, alt := range result.Alternatives.Packages {
			altURL := fmt.Sprintf("https://www.trustypkg.dev/%s/%s", ecosystem, url.QueryEscape(alt.PackageName))
			reportBuilder.WriteString(fmt.Sprintf("| `%s` | `%.2f` | [`%s`](%s) |\n", alt.PackageName, float64(alt.Score), alt.PackageName, altURL))
		}
	} else {
		reportBuilder.WriteString("#### No alternative packages found.\n")
	}

	reportBuilder.WriteString("\n---\n\n")

	// Check if the Trusty score is below the scoreThreshold, if IsDeprecated, isMalicious, Archived, if so shouldFail is set to true
	if result.PackageData.IsDeprecated ||
		result.PackageData.Origin == "malicious" ||
		result.PackageData.Archived ||
		result.Summary.Score < scoreThreshold {
		shouldFail = true
	}

	return reportBuilder.String(), shouldFail
}

// fetchPackageData fetches package data from the specified request URL for a given dependency and ecosystem.
// It sends the result to the provided resultChan channel.
// The function runs asynchronously in a goroutine.
// If an error occurs during the API request, reading the response, or unmarshaling the response,
// the function logs the error and closes the resultChan channel.
// It handles different statuses returned by the API response and takes appropriate actions for each status.
// For "complete" status, it logs a success message and sends the data to the resultChan channel before closing it.
// For "failed" status, it logs a failure message and closes the resultChan channel.
// For "pending" and "scoring" statuses, it waits for 5 seconds before retrying.
// For any other unexpected status, it logs an error message and closes the resultChan channel.
func fetchPackageData(requestURL, dep, ecosystem string, resultChan chan<- Package) {
	go func() {
		var data Package
		for {
			resp, err := http.Get(requestURL)
			if err != nil {
				log.Printf("Error making API request for %s in %s ecosystem: %v\n", dep, ecosystem, err)
				close(resultChan)
				return
			}

			body, err := io.ReadAll(resp.Body)
			resp.Body.Close() // Ensure the body is closed after reading
			if err != nil {
				log.Printf("Error reading API response for %s in %s ecosystem: %v\n", dep, ecosystem, err)
				close(resultChan)
				return
			}

			if err := json.Unmarshal(body, &data); err != nil {
				log.Printf("Error unmarshaling API response for %s in %s ecosystem: %v\n", dep, ecosystem, err)
				close(resultChan)
				return
			}

			// Handle different statuses
			switch data.PackageData.Status {
			case "complete":
				log.Printf("API request for %s in %s ecosystem complete\n", dep, ecosystem)
				resultChan <- data
				close(resultChan)
				return
			case "failed":
				// Handle failure, log error, and close channel
				log.Printf("API request for %s in %s ecosystem failed\n", dep, ecosystem)
				close(resultChan)
				return
			case "pending", "scoring":
				// Wait before retrying for these statuses
				time.Sleep(5 * time.Second)
			default:
				// Handle unexpected status
				log.Printf("Unexpected status for %s in %s ecosystem: %s\n", dep, ecosystem, data.PackageData.Status)
				close(resultChan)
				return
			}
		}
	}()
}
