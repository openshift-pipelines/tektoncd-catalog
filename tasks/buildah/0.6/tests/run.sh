#!/usr/bin/env bash

set -e
    
KUBECTL=kubectl
if [[ "$1" == "openshift" ]]; then
    KUBECTL=oc    
fi

tns=buildah-test

# Run a secure registry as a sidecar to allow the tasks to push to this registry using the certs.
# It will create a configmap `sslcert` with certificate available at key `ca.crt`
function add_sidecar_secure_registry() {
    TMD=$(mktemp -d)

    # Generate SSL Certificate
    openssl req -newkey rsa:4096 -nodes -sha256 -keyout "${TMD}"/ca.key -x509 -days 365 \
            -addext "subjectAltName = DNS:registry" \
            -out "${TMD}"/ca.crt -subj "/C=FR/ST=IDF/L=Paris/O=Tekton/OU=Catalog/CN=registry"

    # Create a configmap from these certs
    ${KUBECTL} create -n "${tns}" configmap sslcert \
            --from-file=ca.crt="${TMD}"/ca.crt --from-file=ca.key="${TMD}"/ca.key

    # Add a secure internal registry as sidecar
    ${KUBECTL} create -n "${tns}" -f ./internal-registry/internal-registry.yaml
}

cd "$(dirname "$0")"
${KUBECTL} create namespace ${tns}

add_sidecar_secure_registry

# Add git-clone
${KUBECTL} -n ${tns} apply -f https://raw.githubusercontent.com/tektoncd/catalog/main/task/git-clone/0.7/git-clone.yaml
${KUBECTL} -n ${tns} apply -f ../buildah.yaml

${KUBECTL} -n ${tns} create -n buildah-test -f ./run.yaml

tkn_pr_status() {
  namespace=$1
  name=$2
  ${KUBECTL} -n $namespace get $name -o jsonpath='{.status.conditions[?(@.type == "Succeeded")].status}'
}

tkn_pr_done () {
  status=$(tkn_pr_status "$@")
  [ "$status" == "True" ] || [ "$status" == "False" ]
}

function tkn_pr_wait {
  # usage: wait_for_pr <timeout_secs> <namespace> <name>
  timeout=$(($1 + $(date +%s)))
  shift

  while :; do
    if [ $(date +%s) -gt $timeout ]; then
        echo "Timeout exceeded waiting for pipeline run to complete"
        return 1
    fi

    if tkn_pr_done "$@"; then
        echo "Pipelinerun has finished"
        return 0
    fi

    echo "Waiting..."
    sleep 30
  done
}

for pipeline in $(${KUBECTL} -n ${tns} get pipelinerun --output=name); do
    tkn_pr_wait 600 ${tns} ${pipeline}
done

set -x

fail=""
for pipeline in $(${KUBECTL} -n ${tns} get pipelinerun --output=name); do
    status=$(${KUBECTL} -n ${tns} get ${pipeline} --output=jsonpath='{.status.conditions[*].status}')
    reason=$(${KUBECTL} -n ${tns} get ${pipeline} --output=jsonpath='{.status.conditions[*].reason}')
    if [[ "${status}" != "True" ]]; then
        echo "Pipeline ${pipeline} failed with reason : ${reason}"
        fail="true"
    fi
done

if [[ -n ${fail} ]]; then
    ${KUBECTL} -n ${tns} get pipelineruns -o yaml
    ${KUBECTL} -n ${tns} get taskruns -o yaml
    for pod in $(${KUBECTL} -n ${tns} get pods --output=name); do
        ${KUBECTL} -n ${tns} logs ${pod} --all-containers
    done
    exit 1
fi
