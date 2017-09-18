#!/bin/bash

function run-gate {
    pushd catasb/local/gate/
    ./run_gate.sh -k || CLUSTER_SETUP_ERROR=true
    popd
}

function get-port {
    PORT="6443"
}

function cluster-login {
    # Do nothing
    :
}
