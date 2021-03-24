#!/bin/bash

echo "Enter pgo hostname:"
read hostname

export PGO_HOST=$hostname

envsubst < scripts/yaml/PostgresOperatorIngress.yaml > PostgresOperatorIngressFilled.yaml

kubectl apply -f PostgresOperatorIngressFilled.yaml
rm PostgresOperatorIngressFilled.yaml

