#!/bin/bash

export GOPATH=`pwd`

go install mudlib && go build gomud

if [ $? == 0 ]
then
  echo SUCCESS
fi
