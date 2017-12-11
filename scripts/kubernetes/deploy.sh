#!/bin/bash

source "$(dirname "${BASH_SOURCE}")/../lib/init.sh"

PROJECT=${ASB_PROJECT}

kubectl get secret -n ${ASB_PROJECT} asb-tls -o jsonpath='{ .data.tls\.key }' > /tmp/key.pem
kubectl get secret -n ${ASB_PROJECT} asb-tls -o jsonpath='{ .data.tls\.crt }' > /tmp/cert.pem

kubectl delete ns ${PROJECT}

retries=25
for r in $(seq $retries); do
    kubectl get ns ansible-service-broker | grep ansible-service-broker
    if [ "$?" -eq 1 ]; then
	break
    fi
    sleep 4
done

kubectl delete clusterrolebindings --ignore-not-found=true asb
kubectl delete pv --ignore-not-found=true etcd

# Render the Kubernetes template
"${TEMPLATE_DIR}/k8s-template.py"

kubectl create ns ${PROJECT}
kubectl create secret tls asb-tls --cert="/tmp/cert.pem" --key="/tmp/key.pem"

kubectl create -f "${TEMPLATE_DIR}/k8s-ansible-service-broker.yaml"
