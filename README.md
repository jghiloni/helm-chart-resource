# helm-chart-resource

This resource allows you to monitor a [Helm](https://helm.sh) repository for
new versions to a chart.

To use:

```
resource_types:
- name: helm-chart
  type: registry-image
  source:
    repository: jghiloni/helm-chart-resource
    tag: v0.1.1

resources:
- name: concourse-helm
  type: helm-chart
  source:
    repository_url: https://concourse-charts.storage.googleapis.com
    chart: concourse
```

## Source Configuration
* `repository_url`: *Required*. The Base URL of the Helm Repository (The URL you would use in `helm repo add`).
* `chart`: *Required*. The name of the helm chart.
* `username`: *Optional*. If HTTP Basic Authorization is required, the username to authenticate.
* `password`: *Optional*. If HTTP Basic Authorization is required, the password to authenticate.
* `skip_tls_validation`: *Optional*. Defaults to `false`. Please don't.
* `sort_by`: *Optional*. Defaults to `semver`. If versions are not semantically versioned or want to version by date
  created, use `created` instead.

## Behavior

### `check`: Discover new chart versions
Reports the latest version for the specified chart in the repository.

### `in`: Fetches the chart files from the repository
Fetches all files specified in the chart's `urls` section In addition, the following files
are created, regardless of whether or not `skip_download` is true:
* `version`: The version number of the fetched chart
* `metadata.json`: A json file with the following contents:
  * chart digest
  * application version
  * chart created date

#### Parameters
* `skip_download`: Default `false`. If `true`, no files will be downloaded.

### `out`: Pushes a new version of a chart

#### Parameters
* `repository`: *Required*. Can be a directory or tarball.
* `version_file`: *Optional*. If set, the pushed chart will have the version specified in `version_file`,
  overriding the version specified in Chart.yaml.
