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

package githubapi

import (
	"context"

	"github.com/google/go-github/v60/github" // Make sure to use the version of go-github you need
	"golang.org/x/oauth2"
)

// GitHubClient wraps the GitHub client and its context
type GitHubClient struct {
	Client *github.Client
	Ctx    context.Context
}

// NewGitHubClient initializes and returns a new GitHubClient
func NewGitHubClient(token string) *GitHubClient {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	return &GitHubClient{
		Client: client,
		Ctx:    ctx,
	}
}

// GetFileContent fetches the content of a file from a specific branch in a repository
func (g *GitHubClient) GetFileContent(owner, repo, path, ref string) (string, error) {
	opts := &github.RepositoryContentGetOptions{Ref: ref}
	fileContent, _, _, err := g.Client.Repositories.GetContents(g.Ctx, owner, repo, path, opts)
	if err != nil {
		return "", err
	}

	content, err := fileContent.GetContent()
	if err != nil {
		return "", err
	}

	return content, nil
}
