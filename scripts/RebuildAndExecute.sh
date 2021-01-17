#!/bin/sh

cd ../cmd/service
go build
cd ../..
./cmd/service/service start -c configs/appconfig.json 