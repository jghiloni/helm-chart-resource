#!/bin/bash

set -e 

exec 3>&1 # make stdout available as fd 3 for the result
exec 1>&2 # redirect all output to stderr for logging

baseDir=$1

request=$(cat)
sourceDir=$1

cd ${sourceDir}

tmpfile=$(mktemp)
echo "Adding debug information to $tmpfile"
echo "${request}" > $tmpfile

repository_url=$(echo ${request} | jq -r '.source.repository_url // ""')
chart=$(echo ${request} | jq -r '.source.chart // ""')
username=$(echo ${request} | jq -r '.source.username // ""')
password=$(echo ${request} | jq -r '.source.password // ""')
skip_tls_validation=$(echo ${request} | jq -r '.source.skip_tls_validation // ""')

repository=$(echo ${request} | jq -r '.params.repository // "."')
version_file=$(echo ${request} | jq -r  '.params.version_file // ""')

auth_args=""
[[ -z "${username}" ]] || auth_args="--username ${username} --password ${password}"

version_arg=""
[[ -f "${version_file}" ]] && version_arg="--version $(cat ${version_file})"

repo_name="put-${RANDOM}"
helm repo add ${repo_name} ${repository_url} ${auth_args}

helm push ${repository} ${repo_name}

[ -f "$(dirname ${repository})/metadata.json" ] && \
    metadata=$(cat "${repository}/metadata.json") || \
    metadata="[ {\"name\": \"repository\", \"value\": \"${repository_url}\"}, {\"name\": \"chart\", \"value\": \"${chart}\"} ]"

version=$(helm search repo "${repo_name}/${chart}" -o json | jq -r --arg chart "${repo_name}/${chart}" '.[] | select(.name==$chart) | .version')

jq -n --arg version "${version}" --argjson metadata "${metadata}" '{"version": {"version": $version}, "metadata": $metadata }' >&3