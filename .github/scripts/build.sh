#!/bin/bash

OUTPUT_DIR=$PWD/dist
mkdir -p ${OUTPUT_DIR}

CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ${OUTPUT_DIR}/scheduler-service-broker -ldflags "-X github.com/rabobank/scheduler-service-broker/conf.VERSION=${VERSION} -X github.com/rabobank/scheduler-service-broker/conf.COMMIT=${COMMIT}" .
