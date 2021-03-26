#!/bin/bash

echo "Enter pgo username:"
read username

echo "Enter pgo password:"
read password

export PGO_USERNAME=$username
export PGO_PASSWORD=$password

envsubst < scripts/yaml/InstallPostgresOperator.yaml > InstallPostgresOperatorFilled.yaml

kubectl create namespace pgo
kubectl apply -f InstallPostgresOperatorFilled.yaml
rm InstallPostgresOperatorFilled.yaml

echo "Operator will take several minutes until it comes up"
