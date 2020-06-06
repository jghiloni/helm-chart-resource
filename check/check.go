package check

import (
	"fmt"
	"log"
	"sort"
	"strings"

	"github.com/blang/semver"
	resource "github.com/jghiloni/helm-resource"
	"github.com/jghiloni/helm-resource/repository"
)

type Request struct {
	Source  resource.Source   `json:"source"`
	Version *resource.Version `json:"version"`
}

type Response []resource.Version

func RunCommand(client resource.HTTPClient, req Request) (Response, error) {
	repo, err := repository.Fetch(client, req.Source)
	if err != nil {
		return Response{}, err
	}

	chartVersions, ok := repo.Entries[req.Source.ChartName]
	if !ok {
		return Response{}, fmt.Errorf("No chart %q found", req.Source.ChartName)
	}

	sortBy := strings.TrimSpace(req.Source.SortBy)
	if sortBy == "" {
		sortBy = "semver"
	}

	if sortBy != "semver" && sortBy != "created" {
		return Response{}, fmt.Errorf("Sort criteria is %q, but it must be semver or created", sortBy)
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
			return []resource.Version{
				{Version: chartVersions[len(chartVersions)-1].Version},
			}, nil
		}

		newVersions := chartVersions[ourVersion:]
		for _, v := range newVersions {
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
