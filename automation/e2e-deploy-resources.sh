#!/usr/bin/env bash

set -ex

# If openshift-pipelines already exists, probably no need to deploy it.
if oc get namespace openshift-pipelines > /dev/null 2>&1; then
  exit 0
fi

# Deploy Openshift Pipelines
# TODO: add support for installing multiple version long term
# Make sure openshift allows custom catalog sources (right ?)
oc patch operatorhub.config.openshift.io/cluster -p='{"spec":{"disableAllDefaultSources":true}}' --type=merge
sleep 2
# Add a custom catalog-source
# FIXME: use the real openshift-pipeline nightly, not mine (vdemeest)
cat <<EOF | oc apply -f-
apiVersion: operators.coreos.com/v1alpha1
kind: CatalogSource
metadata:                      
  name: custom-osp-nightly
  namespace: openshift-marketplace         
spec:                                                                                                                                                                                                                                                
  sourceType: grpc                     
  image: quay.io/openshift-pipeline/openshift-pipelines-operator-index:1.10
  displayName: "Custom OSP Nightly"
  updateStrategy:
    registryPoll:
      interval: 30m                                                                                                                                                                                                                                  
EOF
sleep 10
# Create the "correct" subscription
oc delete subscription pipelines -n openshift-operators || true
cat <<EOF | oc apply -f-
apiVersion: operators.coreos.com/v1alpha1
kind: Subscription
metadata:
  name: openshift-pipelines-operator
  namespace: openshift-operators
spec:
  channel: latest
  name: openshift-pipelines-operator-rh
  source: custom-osp-nightly
  sourceNamespace: openshift-marketplace
EOF

# This deploys a released version
# cat <<EOF | oc apply -f-
# apiVersion: operators.coreos.com/v1alpha1
# kind: Subscription
# metadata:
#   name: openshift-pipeline-operator
#   namespace: openshift-operators
# spec:
#   channel: latest
#   name: openshift-pipelines-operator-rh
#   source: redhat-operators
#   sourceNamespace: openshift-marketplace
# EOF

# wait until tekton pipelines operator is created
timeout 2m bash <<- EOF
  until oc get deployment openshift-pipelines-operator -n openshift-operators; do
    sleep 5
  done
EOF

# wait for tekton pipelines
oc rollout status -n openshift-operators deployment/openshift-pipelines-operator --timeout 10m

# wait until clustertasks tekton CRD is properly deployed
timeout 10m bash <<- EOF
  until oc get crd tasks.tekton.dev; do
    sleep 5
  done
EOF

# wait until tekton pipelines webhook is created
timeout 2m bash <<- EOF
  until oc get deployment tekton-pipelines-webhook -n openshift-pipelines; do
    sleep 5
  done
EOF

# wait until tekton pipelines webhook is online
oc wait -n openshift-pipelines deployment tekton-pipelines-webhook --for condition=Available --timeout 10m

tkn version
