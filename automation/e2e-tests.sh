#!/usr/bin/env bash

set -ex

echo "Run e2e tests :)"
./automation/e2e-deploy-resources.sh

echo make cluster-test
