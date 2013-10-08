#!/bin/bash

export GOPATH=`pwd`

go install mudlib && \
  go test mudlib && \
  go build -ldflags "-s" gomud

if [ $? == 0 ]
then
  echo SUCCESS
fi
