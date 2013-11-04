#!/bin/bash

export GOPATH=`pwd`

echo installing...
go install mudlib

if [ $? == 0 ]; then
  echo testing...
  go test mudlib
fi

if [ $? == 0 ]; then
  echo building...
  go build -ldflags "-s" gomud
fi

if [ $? == 0 ]; then
  echo SUCCESS
fi
