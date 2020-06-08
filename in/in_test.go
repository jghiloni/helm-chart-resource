package in_test

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	resource "github.com/jghiloni/helm-resource"
	"github.com/jghiloni/helm-resource/in"
)

func TestInCommand(t *testing.T) {
	client := &fakeClient{}

	req := in.Request{
		Source: resource.Source{
			RepositoryURL: "http://localhost:8080",
			ChartName:     "concourse",
		},
		Version: resource.Version{Version: "11.1.0"},
	}

	expected := []resource.MetadataField{
		{Name: "repository", Value: "http://localhost:8080"},
		{Name: "chart", Value: "concourse"},
		{Name: "digest", Value: "86f5f3bd5380eaf6331b6413b5628ceed7116f316ab83c302191c319d168a2d7"},
		{Name: "app_version", Value: "6.2.0"},
		{Name: "created", Value: "2020-06-05T14:01:19Z"},
	}

	t.Run("Downloading Everything", func(t *testing.T) {
		baseDir, err := ioutil.TempDir(os.TempDir(), "helm-test-")
		if err != nil {
			t.Fatal(err)
		}
		defer func() {
			os.RemoveAll(baseDir)
		}()

		resp, err := in.RunCommand(baseDir, client, req)
		if err != nil {
			t.Fatal(err)
		}

		if req.Version != resp.Version {
			t.Fatalf("Requested version %q does not match emitted version %q", req.Version, resp.Version)
		}

		for len(expected) != len(resp.Metadata) {
			t.Fatalf("Emitted metadata does not match expected data")
		}

		for i := range resp.Metadata {
			if resp.Metadata[i] != expected[i] {
				t.Fatalf("%v does not match %v", resp.Metadata[i], expected[i])
			}
		}

		checkFiles(t, baseDir, true)
	})

	t.Run("Downloading Nothing", func(t *testing.T) {
		req.Params = in.Params{SkipDownload: true}
		baseDir, err := ioutil.TempDir(os.TempDir(), "helm-test-")
		if err != nil {
			t.Fatal(err)
		}
		defer func() {
			os.RemoveAll(baseDir)
		}()

		resp, err := in.RunCommand(baseDir, client, req)
		if err != nil {
			t.Fatal(err)
		}

		if req.Version != resp.Version {
			t.Fatalf("Requested version %q does not match emitted version %q", req.Version, resp.Version)
		}

		for len(expected) != len(resp.Metadata) {
			t.Fatalf("Emitted metadata does not match expected data")
		}

		for i := range resp.Metadata {
			if resp.Metadata[i] != expected[i] {
				t.Fatalf("%v does not match %v", resp.Metadata[i], expected[i])
			}
		}

		checkFiles(t, baseDir, false)
	})
}

func checkFiles(t *testing.T, baseDir string, checkForTarball bool) {
	_, err := os.Stat(filepath.Join(baseDir, "version"))
	if err != nil {
		t.Fatal(err)
	}

	_, err = os.Stat(filepath.Join(baseDir, "metadata.json"))
	if err != nil {
		t.Fatal(err)
	}

	if checkForTarball {
		_, err := os.Stat(filepath.Join(baseDir, "concourse-11.1.0.tgz"))
		if err != nil {
			t.Fatal(err)
		}

		_, err = os.Stat(filepath.Join(baseDir, "external-file.tgz"))
		if err != nil {
			t.Fatal(err)
		}
	}
}

type fakeClient struct{}

func (f *fakeClient) Do(req *http.Request) (*http.Response, error) {
	w := httptest.NewRecorder()
	if strings.HasSuffix(req.URL.Path, "/index.yaml") {
		w.WriteString(chartYAML)
	} else if strings.HasSuffix(req.URL.Path, ".tgz") {
		w.WriteString("12345")
	} else {
		w.WriteHeader(http.StatusNotFound)
	}

	return w.Result(), nil
}

var chartYAML = `apiVersion: v1
entries:
  concourse:
  - apiVersion: v1
    appVersion: 6.2.0
    created: "2020-06-05T14:01:19.680138326Z"
    description: Concourse is a simple and scalable CI system.
    digest: 86f5f3bd5380eaf6331b6413b5628ceed7116f316ab83c302191c319d168a2d7
    engine: gotpl
    home: https://concourse-ci.org/
    icon: https://avatars1.githubusercontent.com/u/7809479
    keywords:
    - ci
    - concourse
    - concourse.ci
    maintainers:
    - email: cscosta@pivotal.io
      name: cirocosta
    - email: will@autonomic.ai
      name: william-tran
    - email: byoussef@pivotal.io
      name: YoussB
    - email: tsilva@pivotal.io
      name: taylorsilva
    name: concourse
    sources:
    - https://github.com/concourse/concourse
    - https://github.com/helm/charts
    urls:
    - concourse-11.1.0.tgz
    - https://some-external-site/external-file.tgz
    version: 11.1.0`
