package check

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"path"
	"sort"
	"strings"

	"github.com/blang/semver"
	resource "github.com/jghiloni/helm-resource"
	"gopkg.in/yaml.v2"
)

type CheckRequest struct {
	Source  resource.Source   `json:"source"`
	Version *resource.Version `json:"version"`
}

type CheckResponse []resource.Version

func DoCheck(client resource.HTTPClient, req CheckRequest) (CheckResponse, error) {

	u, err := url.ParseRequestURI(req.Source.RepositoryURL)
	if err != nil {
		return CheckResponse{}, err
	}
	u.Path = path.Join(u.Path, "index.yaml")

	httpReq, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return CheckResponse{}, err
	}

	if req.Source.Username != "" {
		httpReq.SetBasicAuth(req.Source.Username, req.Source.Password)
	}

	resp, err := client.Do(httpReq)
	if err != nil {
		return CheckResponse{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return CheckResponse{}, fmt.Errorf("Received bad HTTP response: %q", resp.Status)
	}

	repo := resource.HelmChartRepository{}
	err = yaml.NewDecoder(resp.Body).Decode(&repo)
	if err != nil {
		return CheckResponse{}, err
	}

	chartVersions, ok := repo.Entries[req.Source.ChartName]
	if !ok {
		return CheckResponse{}, fmt.Errorf("No chart %q found", req.Source.ChartName)
	}

	sortBy := strings.TrimSpace(req.Source.SortBy)
	if sortBy == "" {
		sortBy = "semver"
	}

	if sortBy != "semver" && sortBy != "created" {
		return CheckResponse{}, fmt.Errorf("Sort criteria is %q, but it must be semver or created", sortBy)
	}

	sort.Slice(chartVersions, func(i, j int) bool {
		switch sortBy {
		case "semver":
			v1, e1 := semver.Parse(chartVersions[i].Version)
			if e1 != nil {
				log.Printf("Error parsing semver %q\n", chartVersions[i].Version)
				return false
			}

			v2, e2 := semver.Parse(chartVersions[j].Version)
			if e2 != nil {
				log.Printf("Error parsing semver %q\n", chartVersions[j].Version)
				return false
			}

			return v1.LT(v2)
		case "created":
			t1, t2 := chartVersions[i].Created, chartVersions[j].Created
			return t1.Before(t2)
		}

		return false
	})

	versions := []resource.Version{}
	if req.Version != nil {
		ourVersion := -1
		for i := range chartVersions {
			if chartVersions[i].Version == req.Version.Version {
				ourVersion = i
				break
			}
		}

		if ourVersion == -1 {
			return CheckResponse{}, fmt.Errorf("Requested version %q not found", req.Version.Version)
		}

		newVersions := chartVersions[ourVersion:]
		for _, v := range newVersions {
			log.Println(v.Version)
			versions = append(versions, resource.Version{
				Version: v.Version,
			})
		}

		return versions, nil
	}

	return []resource.Version{
		{Version: chartVersions[len(chartVersions)-1].Version},
	}, nil
}
