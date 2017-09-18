#!/bin/bash

function run-gate {
    pushd catasb/local/gate/
    ./run_gate.sh || CLUSTER_SETUP_ERROR=true
    popd
}

function get-port {
    PORT="8443"
}


function cluster-login {
    oc login --insecure-skip-tls-verify 172.17.0.1:8443 -u admin -p admin
    oc project default
}
