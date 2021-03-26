#!/bin/bash

echo "Enter username:"
read username

echo "Enter nginx-Namespace:"
read nginxNamespace

export KUB_TEMP_PREFIX=$username
export NGINX_NAMESPACE=$nginxNamespace

kubectl delete namespace $KUB_TEMP_PREFIX-namespace
kubectl delete rolebinding $KUB_TEMP_PREFIX-nginx-rolebinding -n $NGINX_NAMESPACE
kubectl delete clusterrole $KUB_TEMP_PREFIX-role
kubectl delete clusterrole $KUB_TEMP_PREFIX-nginx-role



