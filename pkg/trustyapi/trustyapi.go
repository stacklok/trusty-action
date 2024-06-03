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

func GenerateReportContent(dependencies []string, ecosystem string, globalThreshold float64, repoActivityThreshold float64, authorActivityThreshold float64, provenanceThreshold float64, typosquattingThreshold float64,
	failOnMalicious bool, failOnDeprecated bool, failOnArchived bool) (string, bool) {
	var (
		failedReportBuilder  strings.Builder
		successReportBuilder strings.Builder
		failAction           bool // Flag to track if the GitHub Action should fail
	)

	failedReportBuilder.WriteString("## 🔴 Failed Dependency Checks\n\n")
	successReportBuilder.WriteString("## 🟢 Successful Dependency Checks\n\n")

	// The following loop generates the report for each dependency and then adds
	// it to the existing reportBuilder, between the header and footer.
	for _, dep := range dependencies {
		log.Printf("Analyzing dependency: %s\n", dep)
		report, shouldFail := ProcessDependency(dep, ecosystem, globalThreshold, repoActivityThreshold, authorActivityThreshold, provenanceThreshold, typosquattingThreshold,
			failOnMalicious, failOnDeprecated, failOnArchived)

		if shouldFail {
			if strings.TrimSpace(report) != "" {
				failedReportBuilder.WriteString(report)
			}
			failAction = true
		} else {
			if strings.TrimSpace(report) != "" {
				successReportBuilder.WriteString(report)
			}
		}
	}

	finalReportBuilder := strings.Builder{}
	finalReportBuilder.WriteString("## 🐻 Trusty Dependency Analysis Action Report \n\n")
	if failedReportBuilder.Len() > len("## 🔴 Failed Dependency Checks\n\n") {
		finalReportBuilder.WriteString(failedReportBuilder.String())
		finalReportBuilder.WriteString("\n")
	}
	if successReportBuilder.Len() > len("## 🟢 Successful Dependency Checks\n\n") {
		finalReportBuilder.WriteString(successReportBuilder.String())
	}

	finalReportBuilder.WriteString("> 🌟 If you like this action, why not try out [Minder](https://github.com/stacklok/minder), the secure supply chain platform. It has vastly more protections and is also free (as in :beer:) to opensource projects.\n")

	// Build the comment body from the report builder
	commentBody := finalReportBuilder.String()

	return commentBody, failAction

}

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
	repoActivityThreshold float64,
	authorActivityThreshold float64,
	provenanceThreshold float64,
	typosquattingThreshold float64,
	failOnMalicious bool,
	failOnDeprecated bool,
	failOnArchived bool) {

	reportContent, failAction := GenerateReportContent(dependencies, ecosystem, globalThreshold, repoActivityThreshold, authorActivityThreshold, provenanceThreshold, typosquattingThreshold,
		failOnMalicious, failOnDeprecated, failOnArchived)

	if strings.TrimSpace(reportContent) != "## 🐻 Trusty Dependency Analysis Action Report \n\n" {
		_, _, err := ghClient.Issues.CreateComment(ctx, owner, repo, prNumber, &github.IssueComment{Body: &reportContent})
		if err != nil {
			log.Printf("error posting comment to PR: %v\n", err)
		} else {
			log.Printf("posted comment to PR: %s/%s#%d\n", owner, repo, prNumber)
		}
	} else {
		log.Println("No report content to post, skipping comment.")
	}

	if failAction {
		log.Println("Failing the GitHub Action due to dependencies not meeting the required criteria.")
		os.Exit(1)
	}
}

func getScoreIcon(score float64, threshold float64) string {
	scoreIcon := ":white_check_mark:"
	if score < threshold {
		scoreIcon = ":x:"
	}
	return scoreIcon
}

func getBoolIcon(b bool, fail bool) string {
	icon := ":white_check_mark:"
	if b && fail {
		icon = ":x:"
	}
	return icon
}

// processDependency analyzes a dependency by making an API request to TrustyPkg.dev and returns a formatted report.
// It takes the dependency name, ecosystem, and score threshold as input parameters.
// The function constructs the query URL, makes the API request, and processes the response.
// If the Trusty score of the dependency is above the score threshold, it skips the dependency.
// Otherwise, it formats the report using Markdown and includes information about the dependency's Trusty score,
// whether it is malicious, deprecated or archived, and recommended alternative packages if available.
// The function returns the formatted report as a string.
func ProcessDependency(dep string, ecosystem string, globalThreshold float64, repoActivityThreshold float64, authorActivityThreshold float64, provenanceThreshold float64, typosquattingThreshold float64,
	failOnMalicious bool, failOnDeprecated bool, failOnArchived bool) (string, bool) {
	var reportBuilder strings.Builder
	shouldFail := false

	// Construct the query URL for the API request
	baseURL := "https://gh.trustypkg.dev/v1/report"
	queryParams := url.Values{}
	queryParams.Add("package_name", dep)
	queryParams.Add("package_type", ecosystem)
	requestURL := fmt.Sprintf("%s?%s", baseURL, queryParams.Encode())

	// Create a channel for receiving the result
	resultChan := make(chan Package)

	// Make the API request to gh.trustypkg.dev
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
	reportBuilder.WriteString(fmt.Sprintf("### :package: [%s](https://www.trustypkg.dev/%s/%s) - %.2f\n\n", dep, ecosystem, dep, result.Summary.Score))

	// Highlight if the package is malicious, deprecated or archived
	if result.PackageData.Origin == "malicious" {
		reportBuilder.WriteString(fmt.Sprintf("⚠ **Malicious** (This package is marked as Malicious. Proceed with extreme caution!) %s\n", getBoolIcon(result.PackageData.Origin == "malicious", failOnMalicious)))
	}
	if result.PackageData.IsDeprecated {
		reportBuilder.WriteString(fmt.Sprintf("⚠ **Deprecated** (This package is marked as Deprecated. Proceed with caution!) %s\n", getBoolIcon(result.PackageData.IsDeprecated, failOnDeprecated)))
	}

	if result.PackageData.Archived {
		reportBuilder.WriteString(fmt.Sprintf("⚠ **Archived** (This package is marked as Archived. Proceed with caution!) %s\n", getBoolIcon(result.PackageData.Archived, failOnArchived)))
	}

	// scores
	reportBuilder.WriteString("<details>\n")
	reportBuilder.WriteString(fmt.Sprintf("<summary>📉 <b>Trusty Score: %.2f %s</b></summary>\n\n", result.Summary.Score, getScoreIcon(result.Summary.Score, globalThreshold)))
	reportBuilder.WriteString("| Category | Score | Passed |\n")
	reportBuilder.WriteString("| --- | --- | --- |\n")
	reportBuilder.WriteString(fmt.Sprintf("| Repo activity   | `%.2f` | %s |\n", result.Summary.Description.ActivityRepo, getScoreIcon(result.Summary.Description.ActivityRepo, repoActivityThreshold)))
	reportBuilder.WriteString(fmt.Sprintf("| Author activity | `%.2f` | %s |\n", result.Summary.Description.ActivityUser, getScoreIcon(result.Summary.Description.ActivityUser, authorActivityThreshold)))
	reportBuilder.WriteString(fmt.Sprintf("| Provenance      | `%.2f` | %s |\n", result.Summary.Description.Provenance, getScoreIcon(result.Summary.Description.Provenance, provenanceThreshold)))
	reportBuilder.WriteString(fmt.Sprintf("| Typosquatting   | `%.2f` | %s |\n", result.Summary.Description.Typosquatting, getScoreIcon(result.Summary.Description.Typosquatting, typosquattingThreshold)))
	reportBuilder.WriteString("</details>\n")

	// write provenance information
	reportBuilder.WriteString("<details>\n")
	if result.Provenance.Description.Provenance.Issuer != "" {
		reportBuilder.WriteString("<summary><strong>Proof of origin (Provenance)</strong>&nbsp;&nbsp;\n")
		reportBuilder.WriteString("<img src='https://cd.foundation/wp-content/uploads/sites/78/2023/05/sigstore_stacked-color-1024x698.png' alt='Sigstore' height='35'></summary>\n\n")
		reportBuilder.WriteString("Built and signed with sigstore using GitHub Actions.\n")
		reportBuilder.WriteString("| | |\n")
		reportBuilder.WriteString("| --- | --- |\n")
		reportBuilder.WriteString(fmt.Sprintf("| Source repo | %s |\n", result.Provenance.Description.Provenance.SourceRepo))
		reportBuilder.WriteString(fmt.Sprintf("| Github Action Workflow | %s |\n", result.Provenance.Description.Provenance.Workflow))
		reportBuilder.WriteString(fmt.Sprintf("| Issuer | %s |\n", result.Provenance.Description.Provenance.Issuer))
		reportBuilder.WriteString(fmt.Sprintf("| Rekor Public Ledger | %s |\n", result.Provenance.Description.Provenance.Transparency))
	} else {
		// need to write regular provenance info
		reportBuilder.WriteString("<summary><strong>Proof of origin (Provenance)</strong>&nbsp;&nbsp;\n")
		reportBuilder.WriteString("<img src='https://www.trustypkg.dev/stacklok-horizontal.svg' alt='Stacklok' height='15'></summary>\n\n")
		reportBuilder.WriteString("| | |\n")
		reportBuilder.WriteString("| --- | --- |\n")
		reportBuilder.WriteString(fmt.Sprintf("| Number of versions | %.0f |\n", result.Provenance.Description.Hp.Versions))
		reportBuilder.WriteString(fmt.Sprintf("| Number of Git Tags/Releases | %.0f |\n", result.Provenance.Description.Hp.Tags))
		reportBuilder.WriteString(fmt.Sprintf("| Number of versions matched to Git Tags/Releases | %.0f |\n", result.Provenance.Description.Hp.Common))
	}
	reportBuilder.WriteString("\n[Learn more about source of origin provenance](https://docs.stacklok.com/trusty/understand/provenance)\n")
	reportBuilder.WriteString("</details>\n")

	// Include alternative packages in a Markdown table if available and if the package is deprecated, archived or malicious
	if result.Alternatives.Packages != nil && len(result.Alternatives.Packages) > 0 {
		reportBuilder.WriteString("<details>\n")
		reportBuilder.WriteString("<summary><strong>Alternative Packages</strong> 💡</summary>\n\n")
		reportBuilder.WriteString("| Package | Score | Trusty Link |\n")
		reportBuilder.WriteString("| ------- | ----- | ---------- |\n")
		for _, alt := range result.Alternatives.Packages {
			altURL := fmt.Sprintf("https://www.trustypkg.dev/%s/%s", ecosystem, url.QueryEscape(alt.PackageName))
			reportBuilder.WriteString(fmt.Sprintf("| `%s` | `%.2f` | [`%s`](%s) |\n", alt.PackageName, float64(alt.Score), alt.PackageName, altURL))
		}
		reportBuilder.WriteString("</details>\n")
	} else {
		reportBuilder.WriteString("No alternative packages found.\n")
	}

	reportBuilder.WriteString("\n---\n\n")

	// Check if the Trusty score is below the scoreThreshold, if IsDeprecated, isMalicious, Archived, if so shouldFail is set to true
	if (failOnDeprecated && result.PackageData.IsDeprecated) ||
		(failOnMalicious && result.PackageData.Origin == "malicious") ||
		(failOnArchived && result.PackageData.Archived) ||
		result.Summary.Score < globalThreshold || result.Summary.Description.ActivityRepo < repoActivityThreshold ||
		result.Summary.Description.ActivityUser < authorActivityThreshold || result.Summary.Description.Provenance < provenanceThreshold ||
		result.Summary.Description.Typosquatting < typosquattingThreshold {
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
