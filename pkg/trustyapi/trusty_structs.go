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

type Package struct {
	PackageName string `json:"package_name"`
	PackageType string `json:"package_type"`
	PackageData struct {
		ID                 string `json:"id"`
		Status             string `json:"status"`
		StatusCode         any    `json:"status_code"`
		Name               string `json:"name"`
		Version            string `json:"version"`
		Author             string `json:"author"`
		AuthorEmail        string `json:"author_email"`
		PackageDescription string `json:"package_description"`
		RepoDescription    string `json:"repo_description"`
		Origin             string `json:"origin"`
		StargazersCount    int    `json:"stargazers_count"`
		WatchersCount      int    `json:"watchers_count"`
		HomePage           string `json:"home_page"`
		HasIssues          bool   `json:"has_issues"`
		HasProjects        bool   `json:"has_projects"`
		HasDownloads       bool   `json:"has_downloads"`
		ForksCount         int    `json:"forks_count"`
		Archived           bool   `json:"archived"`
		IsDeprecated       bool   `json:"is_deprecated"`
		Disabled           bool   `json:"disabled"`
		OpenIssuesCount    int    `json:"open_issues_count"`
		Visibility         string `json:"visibility"`
		Forks              int    `json:"forks"`
		DefaultBranch      string `json:"default_branch"`
		NetworkCount       int    `json:"network_count"`
		SubscribersCount   int    `json:"subscribers_count"`
		RepositoryName     string `json:"repository_name"`
		ContributorCount   int    `json:"contributor_count"`
		PublicRepos        int    `json:"public_repos"`
		PublicGists        int    `json:"public_gists"`
		Followers          int    `json:"followers"`
		Following          int    `json:"following"`
		Owner              struct {
			Author          string `json:"author"`
			AuthorEmail     string `json:"author_email"`
			Login           string `json:"login"`
			AvatarURL       string `json:"avatar_url"`
			GravatarID      string `json:"gravatar_id"`
			URL             string `json:"url"`
			HTMLURL         string `json:"html_url"`
			Company         string `json:"company"`
			Blog            string `json:"blog"`
			Location        string `json:"location"`
			Email           string `json:"email"`
			Hireable        bool   `json:"hireable"`
			TwitterUsername string `json:"twitter_username"`
			PublicRepos     int    `json:"public_repos"`
			PublicGists     any    `json:"public_gists"`
			Followers       int    `json:"followers"`
			Following       int    `json:"following"`
		} `json:"owner"`
		Contributors []struct {
			Author          string `json:"author"`
			AuthorEmail     string `json:"author_email"`
			Login           string `json:"login"`
			AvatarURL       string `json:"avatar_url"`
			GravatarID      string `json:"gravatar_id"`
			URL             string `json:"url"`
			HTMLURL         string `json:"html_url"`
			Company         any    `json:"company"`
			Blog            any    `json:"blog"`
			Location        string `json:"location"`
			Email           string `json:"email"`
			Hireable        bool   `json:"hireable"`
			TwitterUsername any    `json:"twitter_username"`
			PublicRepos     int    `json:"public_repos"`
			PublicGists     any    `json:"public_gists"`
			Followers       int    `json:"followers"`
			Following       int    `json:"following"`
		} `json:"contributors"`
		LastUpdate string `json:"last_update"`
	} `json:"package_data"`
	Summary struct {
		Score       float64 `json:"score"`
		Description struct {
			Activity      float64 `json:"activity"`
			Provenance    float64 `json:"provenance"`
			Typosquatting float64 `json:"typosquatting"`
			ActivityUser  float64 `json:"activity_user"`
			ActivityRepo  float64 `json:"activity_repo"`
		} `json:"description"`
		UpdatedAt string `json:"updated_at"`
	} `json:"summary"`
	Provenance struct {
		Score       float64 `json:"score"`
		Description struct {
			Hp struct {
				Tags     float64 `json:"tags"`
				Common   float64 `json:"common"`
				Overlap  float64 `json:"overlap"`
				Versions float64 `json:"versions"`
				OverTime struct {
				} `json:"over_time"`
			} `json:"hp"`
			Score      float64 `json:"score"`
			Status     string  `json:"status"`
			Provenance struct {
				Issuer       string `json:"issuer"`
				Workflow     string `json:"workflow"`
				SourceRepo   string `json:"source_repo"`
				TokenIssuer  string `json:"token_issuer"`
				Transparency string `json:"transparency"`
			} `json:"provenance"`
		} `json:"description"`
		UpdatedAt string `json:"updated_at"`
	} `json:"provenance"`
	Activity struct {
		Score       float64 `json:"score"`
		Description struct {
			Repo float64 `json:"repo"`
			User float64 `json:"user"`
		} `json:"description"`
		UpdatedAt string `json:"updated_at"`
	} `json:"activity"`
	Typosquatting struct {
		Score       float64 `json:"score"`
		Description struct {
			TotalSimilarNames int `json:"total_similar_names"`
		} `json:"description"`
		UpdatedAt string `json:"updated_at"`
	} `json:"typosquatting"`
	Alternatives struct {
		Status   string `json:"status"`
		Packages []struct {
			ID              string  `json:"id"`
			PackageName     string  `json:"package_name"`
			PackageType     string  `json:"package_type"`
			RepoDescription string  `json:"repo_description"`
			Score           float64 `json:"score"`
			Provenance      struct {
				Score       float64 `json:"score"`
				Description struct {
					Hp struct {
						Tags     float64 `json:"tags"`
						Common   float64 `json:"common"`
						Overlap  float64 `json:"overlap"`
						Versions float64 `json:"versions"`
						OverTime struct {
						} `json:"over_time"`
					} `json:"hp"`
					Score      float64 `json:"score"`
					Status     string  `json:"status"`
					Provenance struct {
						Issuer       string `json:"issuer"`
						Workflow     string `json:"workflow"`
						SourceRepo   string `json:"source_repo"`
						TokenIssuer  string `json:"token_issuer"`
						Transparency string `json:"transparency"`
					} `json:"provenance"`
				} `json:"description"`
				UpdatedAt string `json:"updated_at"`
			} `json:"provenance"`
		} `json:"packages"`
	} `json:"alternatives"`
	SimilarPackageNames []struct {
		ID              string  `json:"id"`
		PackageName     string  `json:"package_name"`
		PackageType     string  `json:"package_type"`
		RepoDescription string  `json:"repo_description"`
		Score           float64 `json:"score"`
		Provenance      struct {
			Score       float64 `json:"score"`
			Description struct {
				Hp struct {
					Tags     float64 `json:"tags"`
					Common   float64 `json:"common"`
					Overlap  float64 `json:"overlap"`
					Versions float64 `json:"versions"`
					OverTime struct {
					} `json:"over_time"`
				} `json:"hp"`
				Score      float64 `json:"score"`
				Status     string  `json:"status"`
				Provenance struct {
					Issuer       string `json:"issuer"`
					Workflow     string `json:"workflow"`
					SourceRepo   string `json:"source_repo"`
					TokenIssuer  string `json:"token_issuer"`
					Transparency string `json:"transparency"`
				} `json:"provenance"`
			} `json:"description"`
			UpdatedAt string `json:"updated_at"`
		} `json:"provenance"`
	} `json:"similar_package_names"`
	SameOriginPackagesCount int `json:"same_origin_packages_count"`
}
