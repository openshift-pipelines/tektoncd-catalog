#!/usr/bin/env bash

set -ex

echo "Run e2e tests :)"
./automation/e2e-deploy-resources.sh $1

# TODO: handle platform arg
# TODO: create an actual tool for this
./tasks/buildah/0.6/tests/run.sh $1
