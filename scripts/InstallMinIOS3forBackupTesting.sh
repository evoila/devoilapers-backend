#!/bin/bash

echo "Enter kubernetes hostname:"
read hostname

export MINIO_HOSTNAME=$hostname

echo "Enter minio username:"
read username

export MINIO_USERNAME=$username

echo "Enter minio password:"
read password

export MINIO_PASSWORD=$password

envsubst < scripts/yaml/InstallMinIOS3forBackupTesting.yaml > InstallMinIOS3forBackupTesting.yaml

kubectl apply -f InstallMinIOS3forBackupTesting.yaml
rm InstallMinIOS3forBackupTesting.yaml