#!/bin/bash

echo "Enter username:"
read username 

echo "Enter nginx-Namespace:"
read nginxNamespace

export KUB_TEMP_PREFIX=$username
export NGINX_NAMESPACE=$nginxNamespace

envsubst < .github/create_service_account.yaml > new_user.yaml

kubectl apply -f new_user.yaml
rm new_user.yaml

echo "Namespace is:"
echo $KUB_TEMP_PREFIX"-namespace"

echo "Token is:"
kubectl -n $KUB_TEMP_PREFIX-namespace describe secret $(kubectl -n $KUB_TEMP_PREFIX-namespace get secret | (grep $KUB_TEMP_PREFIX-user || echo "$_") | awk '{print $1}') | grep token: | awk '{print $2}'



