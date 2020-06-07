#!/bin/bash

baseDir=$1

request=$(cat)

repository_url=$(echo ${request} | jq -r '.source.repository_url // ""')
username=$(echo ${request} | jq -r '.source.username // ""')
password=$(echo ${request} | jq -r '.source.password // ""')
skip_tls_validation=$(echo ${request} | jq -r '.source.skip_tls_validation // ""')

repository=$(echo ${request} | jq -r '.params.repository // "."')
version_file=$(echo ${request} | jq -r  '.params.version_file // ""')

auth_args=""
[[ -z "${username}" ]] || auth_args="--username '${username}' --password '${password}'"

version_arg=""
[[ -f "${version_file}" ]] && version_arg="--version $(cat ${version_file})"

repo_name="put-${RANDOM}"
helm repo add ${repo_name} ${repository_url} ${auth_args}

helm push ${repository} ${repo_name}
