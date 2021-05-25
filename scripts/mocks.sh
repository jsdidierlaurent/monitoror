#!/usr/bin/env bash
# Do not use this script manually, Use makefile

set -e

#######################################################################
# This script is to generate mocks for monitorables and api directory #
#######################################################################

find . -name "mocks" -type d -print0 | xargs -r0 -- rm -r
go generate ./...

# Generating mocks for external interfaces
mockery --name Client --output pkg/goazuredevops/build/mocks --dir "$(go list -m -f '{{ .Dir }}' github.com/jsdidierlaurent/azure-devops-go-api/azuredevops)/build"
mockery --name Client --output pkg/goazuredevops/release/mocks --dir "$(go list -m -f '{{ .Dir }}' github.com/jsdidierlaurent/azure-devops-go-api/azuredevops)/release"
