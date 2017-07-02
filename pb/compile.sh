#!/bin/bash

cd `dirname $0`
protoc --go_out=plugins=grpc:. god.proto
