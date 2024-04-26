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

package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/stacklok/trusty-action/pkg/githubapi"
	"github.com/stacklok/trusty-action/pkg/parser"
	"github.com/stacklok/trusty-action/pkg/trustyapi"
	"github.com/stacklok/trusty-action/pkg/types"

	"github.com/google/go-github/v60/github"
	"golang.org/x/oauth2"
)

func parseScore(scoreStr string, defaultScore string) float64 {
	if scoreStr == "" {
		scoreStr = defaultScore
	}
	score, err := strconv.ParseFloat(scoreStr, 64)
	if err != nil {
		log.Printf("Invalid score threshold value: %s\n", scoreStr)
		return 0
	}
	return score
}

func main() {
	ctx := context.Background()

	globalThreshold := parseScore(os.Getenv("INPUT_THRESHOLDS_GLOBAL"), "5")
	activityThreshold := parseScore(os.Getenv("INPUT_THRESHOLDS_ACTIVITY"), "0")
	provenanceThreshold := parseScore(os.Getenv("INPUT_THRESHOLDS_PROVENANCE"), "0")
	typosquattingThreshold := parseScore(os.Getenv("INPUT_THRESHOLDS_TYPOSQUATTING"), "0")

	// Split the GITHUB_REPOSITORY environment variable to get owner and repo
	repoFullName := os.Getenv("GITHUB_REPOSITORY")
	if repoFullName == "" {
		log.Println("GITHUB_REPOSITORY environment variable is not set.")
		os.Exit(1)
	}
	repoParts := strings.Split(repoFullName, "/")
	if len(repoParts) != 2 {
		log.Println("Invalid GITHUB_REPOSITORY format. Expected format is 'owner/repo'.")
		os.Exit(1)
	}
	owner, repo := repoParts[0], repoParts[1]

	// Read the event file to get the pull request number
	eventPath := os.Getenv("GITHUB_EVENT_PATH")
	if eventPath == "" {
		log.Println("GITHUB_EVENT_PATH environment variable is not set.")
		os.Exit(1)
	}

	eventData, err := os.ReadFile(eventPath)
	if err != nil {
		log.Printf("Error reading event payload file: %v\n", err)
		os.Exit(1)
	}

	var eventPayload struct {
		PullRequest struct {
			Number int `json:"number"`
		} `json:"pull_request"`
	}
	if err := json.Unmarshal(eventData, &eventPayload); err != nil {
		log.Printf("Error parsing event payload JSON: %v\n", err)
		os.Exit(1)
	}
	prNumber := eventPayload.PullRequest.Number
	if prNumber == 0 {
		log.Println("Pull request number not found in event payload.")
		os.Exit(1)
	}

	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		log.Println("GITHUB_TOKEN environment variable is not set.")
		os.Exit(1)
	}

	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	tc := oauth2.NewClient(ctx, ts)
	ghClient := github.NewClient(tc)

	githubClient := githubapi.NewGitHubClient(token)

	// To debug overide the owner, repo and prNumber below
	// owner := "lukehinds" // Replace with the owner of the repository
	// repo := "bad-npm"    // Replace with the repository name
	// prNumber := 54       // Replace with the actual PR number you want to analyze

	// Fetch PR to get base and head refs (using the original GitHub client)
	pr, _, err := ghClient.PullRequests.Get(ctx, owner, repo, prNumber)
	if err != nil {
		log.Printf("Error fetching PR: %v\n", err)
		return
	}
	baseRef, headRef := *pr.Base.Ref, *pr.Head.Ref

	// Get the files that were changed in the PR
	files, _, err := ghClient.PullRequests.ListFiles(ctx, owner, repo, prNumber, nil)
	if err != nil {
		log.Printf("Error fetching changed files: %v\n", err)
		return
	}

	for _, file := range files {

		baseContent, err := githubClient.GetFileContent(owner, repo, *file.Filename, baseRef)
		if err != nil {
			log.Printf("Error fetching base content for %s: %v\n", *file.Filename, err)
			continue
		}
		headContent, err := githubClient.GetFileContent(owner, repo, *file.Filename, headRef)
		if err != nil {
			log.Printf("Error fetching head content for %s: %v\n", *file.Filename, err)
			continue
		}

		// Parse the contents to get slices of Dependency structs and the ecosystem from the base content
		baseDeps, ecosystem, err := parser.Parse(*file.Filename, baseContent) // Use ecosystem from base content parsing
		if err != nil {
			log.Printf("Error parsing base dependencies for %s: %v\n", *file.Filename, err)
			continue
		}

		// Ignore the ecosystem from head content parsing
		headDeps, _, err := parser.Parse(*file.Filename, headContent) // Ignore ecosystem from head content parsing
		if err != nil {
			log.Printf("Error parsing head dependencies for %s: %v\n", *file.Filename, err)
			continue
		}

		// Convert slices to maps
		baseDepsMap := make(map[string]string)
		for _, dep := range baseDeps {
			baseDepsMap[dep.Name] = dep.Version
		}
		headDepsMap := make(map[string]string)
		for _, dep := range headDeps {
			headDepsMap[dep.Name] = dep.Version
		}

		// Find added dependencies
		addedDepsMap := types.DiffDependencies(baseDepsMap, headDepsMap)

		// Extract dependency names from the addedDepsMap
		var addedDepNames []string
		for depName := range addedDepsMap {
			addedDepNames = append(addedDepNames, depName)
		}
		// Debug print
		log.Printf("Added dependencies: %v\n", addedDepNames)

		// In your main application where you call ProcessDependencies
		trustyapi.BuildReport(ctx, ghClient, owner, repo, prNumber, addedDepNames, ecosystem, globalThreshold, activityThreshold, provenanceThreshold, typosquattingThreshold)

	}
}
