#!/bin/bash

# export MONGODB_PASSWORD=$(kubectl get secret --namespace blogsvc mongodb -o jsonpath="{.data.mongodb-password}" | base64 --decode)
#  blogserver -grpc-port=9090 -http-port=8080 -db-host=localhost:27017 -log-level=-1 -db-user=ampc -db-password=changeme
# blogserver -grpc-port=9090 -http-port=8080 -db-host=mongodb.blogsvc:27017 -log-level=-1 -db-user=ampc -db-password=changeme
 blogserver -grpc-port=${GRPC_PORT} -http-port=${REST_PORT} -db-host=${MGODB_HOST} -log-level=-1 -db-user=${MGODB_USER} -db-password=${MGODB_PWD} -db-schema=${MGODB_SCHEMA} -db-collection=${MGODB_COLLECTION}