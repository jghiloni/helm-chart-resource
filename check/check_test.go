package check_test

import (
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	resource "github.com/jghiloni/helm-resource"
	"github.com/jghiloni/helm-resource/check"
)

func TestFirstCheck(t *testing.T) {
	log.SetFlags(log.Lshortfile | log.Ldate | log.Ltime | log.LUTC)
	client := &fakeClient{}

	source := resource.Source{
		RepositoryURL: "https://example.com/",
		ChartName:     "concourse",
	}

	checkReq := check.CheckRequest{
		Source: source,
	}

	t.Run("It works on a public repo with no cursor", func(t *testing.T) {
		resp, err := check.DoCheck(client, checkReq)
		if err != nil {
			t.Fatalf("Unexpected error %v", err)
		}

		if len(resp) != 1 {
			t.Fatalf("There should be exactly 1 version returned, but there were %d", len(resp))
		}

		if resp[0].Version != "11.1.0" {
			t.Fatalf("Expected version %s to be 11.1.0", resp[0].Version)
		}
	})

	t.Run("It works on a public repo with a cursor", func(t *testing.T) {
		checkReq.Version = &resource.Version{Version: "10.3.0"}

		resp, err := check.DoCheck(client, checkReq)
		if err != nil {
			t.Fatalf("Unexpected error %v", err)
		}

		if len(resp) != 4 {
			t.Fatalf("There should be exactly 4 version returned, but there were %d", len(resp))
		}

		expected := []resource.Version{
			{Version: "10.3.0"},
			{Version: "11.0.0"},
			{Version: "11.0.1"},
			{Version: "11.1.0"},
		}

		for i := range resp {
			if resp[i] != expected[i] {
				t.Fatalf("Expected version %q to be %q", resp[i].Version, expected[i].Version)
			}
		}
	})
}

type fakeClient struct{}

func (f *fakeClient) Do(req *http.Request) (*http.Response, error) {
	w := httptest.NewRecorder()
	if strings.HasSuffix(req.URL.Path, "/index.yaml") {
		w.WriteString(chartYAML)
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
    version: 11.1.0
  - apiVersion: v1
    appVersion: 6.1.0
    created: "2020-06-01T18:45:39.44313152Z"
    description: Concourse is a simple and scalable CI system.
    digest: e27764741330f034461b60d94b37737409beffb6873e342cfb6669fcbb5fcd49
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
    - concourse-11.0.1.tgz
    version: 11.0.1
  - apiVersion: v1
    appVersion: 6.1.0
    created: "2020-06-01T13:36:47.463924583Z"
    description: Concourse is a simple and scalable CI system.
    digest: c4764a303a0b06bc3649ec9e8e121b92a2b0fec12260b9ef3c6a68b3aab8b51d
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
    - concourse-11.0.0.tgz
    version: 11.0.0
  - apiVersion: v1
    appVersion: 6.1.0
    created: "2020-05-25T22:00:16.253889228Z"
    description: Concourse is a simple and scalable CI system.
    digest: 604f7ed41d1846dfb84e59a6f1f75c200113ecec22dd0a99a86d6e3a47432d7a
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
    - concourse-10.3.0.tgz
    version: 10.3.0
  - apiVersion: v1
    appVersion: 6.1.0
    created: "2020-05-19T18:44:24.925651524Z"
    description: Concourse is a simple and scalable CI system.
    digest: 16662cf627371b61ec51104886d6282836660913d5aeeaff08acdc88d2aee764
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
    - concourse-10.2.3.tgz
    version: 10.2.3
  - apiVersion: v1
    appVersion: 6.1.0
    created: "2020-05-19T14:38:51.06430037Z"
    description: Concourse is a simple and scalable CI system.
    digest: eef7992d2f5c86976e1660b74a45a4d7832a3962b74a221a972070a609bf9278
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
    - concourse-10.2.2.tgz
    version: 10.2.2
  - apiVersion: v1
    appVersion: 6.1.0
    created: "2020-05-15T19:53:07.563447085Z"
    description: Concourse is a simple and scalable CI system.
    digest: 7c00ebd0dbf259b15230240fdff7e411195085fc2f0a8bce7ff77cbec55975cf
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
    - concourse-10.2.1.tgz
    version: 10.2.1
  - apiVersion: v1
    appVersion: 6.1.0
    created: "2020-05-13T16:50:00.899359116Z"
    description: Concourse is a simple and scalable CI system.
    digest: 68d4be70e62d0ddaf4580da97495f584c5855c67e9da64dd61b6222fc6e89999
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
    - concourse-10.2.0.tgz
    version: 10.2.0
  - apiVersion: v1
    appVersion: 6.1.0
    created: "2020-05-12T15:01:06.078595753Z"
    description: Concourse is a simple and scalable CI system.
    digest: 0fc630f824363227f3781f76b7c5523e5750a2d6fbf5d0f9e5838cd5a1963da3
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
    - concourse-10.1.0.tgz
    version: 10.1.0
  - apiVersion: v1
    appVersion: 6.0.0
    created: "2020-05-08T14:45:58.415978279Z"
    description: Concourse is a simple and scalable CI system.
    digest: b3a08f1dd211f09b2e5b901dcb5c1d954a086333cf305d39d845f1a1abce3508
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
    - concourse-10.0.7.tgz
    version: 10.0.7
  - apiVersion: v1
    appVersion: 6.0.0
    created: "2020-05-04T20:53:06.855591033Z"
    description: Concourse is a simple and scalable CI system.
    digest: 0e2cff3cd9b0f273691c4facfe10a5ebd130ebc88ed9ba27aad9233b6441feb9
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
    - concourse-10.0.6.tgz
    version: 10.0.6
  - apiVersion: v1
    appVersion: 6.0.0
    created: "2020-04-24T21:41:21.036858696Z"
    description: Concourse is a simple and scalable CI system.
    digest: 1d5d96a2c134f70d740c7f085bb67a8e902329d24007bd5a814034154c3f9ec2
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
    - concourse-10.0.5.tgz
    version: 10.0.5
  - apiVersion: v1
    appVersion: 6.0.0
    created: "2020-04-16T14:08:25.136351865Z"
    description: Concourse is a simple and scalable CI system.
    digest: b653ada0494bc7ad8a6df50faf3e3a682760242072911d3b5abc1a17032679a5
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
    - concourse-10.0.4.tgz
    version: 10.0.4
  - apiVersion: v1
    appVersion: 6.0.0
    created: "2020-04-15T17:59:24.26523054Z"
    description: Concourse is a simple and scalable CI system.
    digest: 3334712bbc750273f0e080126b39de4c0853fae8dc84881ac8f871137933892e
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
    - concourse-10.0.3.tgz
    version: 10.0.3
  - apiVersion: v1
    appVersion: 6.0.0
    created: "2020-04-14T15:12:44.864500067Z"
    description: Concourse is a simple and scalable CI system.
    digest: 3343619dd0437a8b5018d2ba33d1b3d3ec35a34581a63db8973d669c01ffd4f0
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
    - concourse-10.0.2.tgz
    version: 10.0.2
  - apiVersion: v1
    appVersion: 6.0.0
    created: "2020-03-31T21:35:52.672318133Z"
    description: Concourse is a simple and scalable CI system.
    digest: f68d1486f169f84348d7775d93ee7d8daf8d01e98346839eacf19c8a3089dfc5
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
    - concourse-10.0.1.tgz
    version: 10.0.1
  - apiVersion: v1
    appVersion: 6.0.0
    created: "2020-03-25T19:15:47.811103577Z"
    description: Concourse is a simple and scalable CI system.
    digest: ef326816d6f1fdc5566a93705d73d5230a6a7d31f6f72c34942600865e61a02c
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
    - concourse-10.0.0.tgz
    version: 10.0.0
  - apiVersion: v1
    appVersion: 5.8.1
    created: "2020-03-26T21:05:09.250165154Z"
    description: Concourse is a simple and scalable CI system.
    digest: 8ca470bf8991675545f2cf1d7a2dbb2786372e9dee73a2113d320a643eda0141
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
    - concourse-9.1.3.tgz
    version: 9.1.3
  - apiVersion: v1
    appVersion: 5.8.0
    created: "2020-03-09T14:21:08.370335773Z"
    description: Concourse is a simple and scalable CI system.
    digest: 109ca3388301e609bc76d7c640e0c85e90000ea5bcd33ef1d7e84938e0426b29
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
    - concourse-9.1.1.tgz
    version: 9.1.1
  - apiVersion: v1
    appVersion: 5.8.0
    created: "2020-02-07T20:00:16.577542819Z"
    description: Concourse is a simple and scalable CI system.
    digest: cd7e61bd35a33a0868abeef26a55387f4e2ee742b47ba0068391c3c2abac4372
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
    name: concourse
    sources:
    - https://github.com/concourse/concourse
    - https://github.com/helm/charts
    urls:
    - concourse-9.1.0.tgz
    version: 9.1.0
  - apiVersion: v1
    appVersion: 5.8.0
    created: "2020-01-08T17:31:01.637418216Z"
    description: Concourse is a simple and scalable CI system.
    digest: 4bc12ea6517e412d1a3551fd8327a27f0376da187a8e8f6fd8e48b16eaf23ed4
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
    name: concourse
    sources:
    - https://github.com/concourse/concourse
    - https://github.com/helm/charts
    urls:
    - concourse-9.0.0.tgz
    version: 9.0.0
  - apiVersion: v1
    appVersion: 5.7.2
    created: "2019-12-13T17:04:17.506247-05:00"
    description: Concourse is a simple and scalable CI system.
    digest: 712cda8b48083e902bca94749176fcdf59abbfb5cef797ca33f114bc8fc5b982
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
    name: concourse
    sources:
    - https://github.com/concourse/concourse
    - https://github.com/helm/charts
    urls:
    - concourse-8.4.1.tgz
    version: 8.4.1
  - apiVersion: v1
    appVersion: 5.5.11
    created: "2020-04-27T14:53:21.827617001Z"
    description: Concourse is a simple and scalable CI system.
    digest: eef147175570a49945228cfa1e79881a4db72ffab73c4a5115723ccd0d05a68a
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
    name: concourse
    sources:
    - https://github.com/concourse/concourse
    - https://github.com/helm/charts
    urls:
    - concourse-8.2.13.tgz
    version: 8.2.13
  - apiVersion: v1
    appVersion: 5.5.11
    created: "2020-04-24T22:01:27.355505653Z"
    description: Concourse is a simple and scalable CI system.
    digest: c5e0064387a7a4c526c3ba1019899d1c0641adaa7b33136c52fe48447902a9c5
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
    name: concourse
    sources:
    - https://github.com/concourse/concourse
    - https://github.com/helm/charts
    urls:
    - concourse-8.2.12.tgz
    version: 8.2.12
  - apiVersion: v1
    appVersion: 5.5.9
    created: "2020-03-24T20:42:05.158128244Z"
    description: Concourse is a simple and scalable CI system.
    digest: 153bfd17157bb367cc83cb6b61d84dec9efea933d53d65cb8657807b742d738f
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
    name: concourse
    sources:
    - https://github.com/concourse/concourse
    - https://github.com/helm/charts
    urls:
    - concourse-8.2.10.tgz
    version: 8.2.10
  - apiVersion: v1
    appVersion: 5.5.7
    created: "2020-02-26T21:11:47.413907973Z"
    description: Concourse is a simple and scalable CI system.
    digest: 6a74c56d63418063d5f6d83737726133b3b32c1ffc8ed279f44ff38a943da1ba
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
    name: concourse
    sources:
    - https://github.com/concourse/concourse
    - https://github.com/helm/charts
    urls:
    - concourse-8.2.7.tgz
    version: 8.2.7
  - apiVersion: v1
    appVersion: 5.6.0
    created: "2019-11-07T22:20:52.233240077Z"
    description: Concourse is a simple and scalable CI system.
    digest: b3c152fe69b2131e71009e61d73ef135c005461307a948129badccd4fb6f9da4
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
    name: concourse
    sources:
    - https://github.com/concourse/concourse
    - https://github.com/helm/charts
    urls:
    - concourse-8.2.6.tgz
    version: 8.2.6
generated: "2020-06-05T14:01:19.673532266Z"
`
