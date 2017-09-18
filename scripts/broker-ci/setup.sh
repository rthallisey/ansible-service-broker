#!/bin/bash

BROKER_DIR="$(dirname "${BASH_SOURCE}")/../.."
source "${BROKER_DIR}/scripts/broker-ci/error.sh"
source "${BROKER_DIR}/scripts/broker-ci/logs.sh"

BUILD_ERROR=false
MAKE_DEPLOY_ERROR=false
CLUSTER_SETUP_ERROR=false
RESOURCE_ERROR=false

set -ex

function cluster-gather {
    if [ -n "${KUBERNETES}" ]; then
	source "${BROKER_DIR}/scripts/broker-ci/kubernetes/gate.sh"
    else
	source "${BROKER_DIR}/scripts/broker-ci/openshift/gate.sh"
    fi
}

function cluster-setup () {
    #git clone https://github.com/fusor/catasb
    git clone https://github.com/rthallisey/catasb
    git fetch
    git checkout kubernetes-playbook

    cat <<EOF > "catasb/config/my_vars.yml"
---
dockerhub_user_name: brokerciuser
dockerhub_org_name: ansibleplaybookbundle
dockerhub_user_password: brokerciuser
EOF

    run-gate
    env-error-check "cluster-setup"

    # Kubernetes port: 6443
    # OpenShift port: 8443
    get-port

    cat <<EOF > "scripts/my_local_dev_vars"
CLUSTER_HOST=172.17.0.1
CLUSTER_PORT=${PORT}

# BROKER_IP_ADDR must be the IP address of where to reach broker
#   it should not be 127.0.0.1, needs to be an address the pods will be able to reach
BROKER_IP_ADDR=${CLUSTER_HOST}
DOCKERHUB_USERNAME="brokerciuser"
DOCKERHUB_PASSWORD="brokerciuser"
DOCKERHUB_ORG="ansibleplaybookbundle"

BOOTSTRAP_ON_STARTUP="true"
BEARER_TOKEN_FILE=""
BROKER_INSECURE="false"
CA_FILE=""
REFRESH_INTERVAL="600s"

# Always, IfNotPresent, Never
IMAGE_PULL_POLICY="Always"
EOF
}

function make-build-image {
    set +x
    RETRIES=3
    for x in $(seq $RETRIES); do
	make build-image
	if [ $? -eq 0 ]; then
	    print-with-green "Broker container completed building."
	    break
	else
	    print-with-yellow "Broker container failed to build."
	    print-with-yellow "Retrying..."
	fi
    done
    if [ "${x}" -eq "${RETRIES}" ]; then
	print-with-red "Broker container failed to build."
	BUILD_ERROR=true
    fi
    env-error-check "make-build-image"
    set -x
}

function make-deploy {
    make deploy
    NAMESPACE="ansible-service-broker" ./scripts/broker-ci/wait-for-resource.sh create pod asb >> /tmp/wait-for-pods-log 2>&1
    env-error-check "make-deploy"
}

function local-env() {
    cluster-login
    make-build-image
    make-deploy
}

echo "========== Broker CI ==========="
echo "Gathering cluster"
cluster-gather

echo "Setting up cluster"
cluster-setup

echo "Setting up local environment"
local-env

set +e
