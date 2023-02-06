#!/usr/bin/env bash

set -ex

if oc get namespace tekton-pipelines > /dev/null 2>&1; then
  exit 0
fi

# Deploy Openshift Pipelines
# TODO: add support for installing nightly long term
cat <<EOF | oc apply -f-
apiVersion: operators.coreos.com/v1alpha1
kind: Subscription
metadata:
  name: openshift-pipeline-operator
  namespace: openshift-operators
spec:
  channel: latest
  name: openshift-pipelines-operator-rh
  source: redhat-operators
  sourceNamespace: openshift-marketplace
EOF

# wait for tekton pipelines
kubectl rollout status -n openshift-operators deployment/openshift-pipelines-operator --timeout 10m

# wait until clustertasks tekton CRD is properly deployed
timeout 10m bash <<- EOF
  until kubectl get crd tasks.tekton.dev; do
    sleep 5
  done
EOF

# wait until tekton pipelines webhook is created
timeout 10m bash <<- EOF
  until kubectl get deployment tekton-pipelines-webhook -n openshift-pipelines; do
    sleep 5
  done
EOF

# wait until tekton pipelines webhook is online
kubectl wait -n openshift-pipelines deployment tekton-pipelines-webhook --for condition=Available --timeout 10m

