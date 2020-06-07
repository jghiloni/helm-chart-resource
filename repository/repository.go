package repository

import (
	"fmt"
	"net/http"
	"net/url"
	"path"

	resource "github.com/jghiloni/helm-resource"
	"gopkg.in/yaml.v2"
)

func Fetch(client resource.HTTPClient, source resource.Source) (resource.HelmChartRepository, error) {
	u, err := url.ParseRequestURI(source.RepositoryURL)
	if err != nil {
		return resource.HelmChartRepository{}, err
	}
	u.Path = path.Join(u.Path, "index.yaml")

	httpReq, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return resource.HelmChartRepository{}, err
	}

	if source.Username != "" {
		httpReq.SetBasicAuth(source.Username, source.Password)
	}

	resp, err := client.Do(httpReq)
	if err != nil {
		return resource.HelmChartRepository{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return resource.HelmChartRepository{}, fmt.Errorf("Received bad HTTP response: %q", resp.Status)
	}

	repo := resource.HelmChartRepository{}
	err = yaml.NewDecoder(resp.Body).Decode(&repo)
	if err != nil {
		return resource.HelmChartRepository{}, err
	}

	return repo, nil
}
