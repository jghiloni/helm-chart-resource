package in

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"time"

	resource "github.com/jghiloni/helm-resource"
	"github.com/jghiloni/helm-resource/repository"
	"gopkg.in/yaml.v2"
)

type Params struct {
	SkipDownload bool     `json:"skip_download"`
	Globs        []string `json:"globs"`
}

type Request struct {
	Source  resource.Source  `json:"source"`
	Version resource.Version `json:"version"`
	Params  Params           `json:"params"`
}

type Response struct {
	Version  resource.Version         `json:"version"`
	Metadata []resource.MetadataField `json:"metadata"`
}

func RunCommand(baseDir string, client resource.HTTPClient, req Request) (Response, error) {
	repo, err := repository.Fetch(client, req.Source)
	if err != nil {
		return Response{}, err
	}

	chartVersions, ok := repo.Entries[req.Source.ChartName]
	if !ok {
		return Response{}, fmt.Errorf("No chart %q found", req.Source.ChartName)
	}

	chartInfo := resource.HelmChartInfo{}
	for _, info := range chartVersions {
		if info.Version == req.Version.Version {
			chartInfo = info
			break
		}
	}

	if chartInfo.Version != req.Version.Version {
		return Response{}, fmt.Errorf("No chart with version %q found", req.Version.Version)
	}

	if !req.Params.SkipDownload {
		for _, chartURL := range chartInfo.URLs {
			doDownload := false
			if len(req.Params.Globs) == 0 {
				doDownload = true
			} else {
				for _, glob := range req.Params.Globs {
					match, err := filepath.Match(glob, chartURL)
					if err != nil {
						return Response{}, err
					}

					if match {
						doDownload = match
						break
					}
				}
			}

			if doDownload {
				u, err := url.Parse(req.Source.RepositoryURL)
				if err != nil {
					return Response{}, err
				}

				u.Path = path.Join(u.Path, chartURL)
				target := filepath.Join(baseDir, chartURL)

				if err = os.MkdirAll(filepath.Dir(target), 0755); err != nil {
					return Response{}, err
				}

				targetFile, err := os.Create(target)
				if err != nil {
					return Response{}, err
				}
				defer targetFile.Close()

				httpReq, err := http.NewRequest(http.MethodGet, u.String(), nil)
				if err != nil {
					return Response{}, err
				}

				httpResp, err := client.Do(httpReq)
				if err != nil {
					return Response{}, err
				}
				defer httpResp.Body.Close()

				_, err = io.Copy(targetFile, httpResp.Body)
				if err != nil {
					return Response{}, err
				}
			}
		}
	}

	versionFile, err := os.Create(filepath.Join(baseDir, "version"))
	if err != nil {
		return Response{}, err
	}
	defer versionFile.Close()
	versionFile.WriteString(req.Version.Version)

	metadataFile, err := os.Create(filepath.Join(baseDir, "metadata.yml"))
	if err != nil {
		return Response{}, err
	}
	defer metadataFile.Close()

	metadata := []resource.MetadataField{
		{Name: "digest", Value: chartInfo.Digest},
		{Name: "app_version", Value: chartInfo.AppVersion},
		{Name: "created", Value: chartInfo.Created.Format(time.RFC3339)},
	}

	err = yaml.NewEncoder(metadataFile).Encode(metadata)
	if err != nil {
		return Response{}, err
	}

	response := Response{
		Version:  req.Version,
		Metadata: metadata,
	}

	return response, nil
}
