#!/bin/bash

version="0.0.1"
GOOS="linux"

buildDate=`date "+%Y-%m-%d"`

buildOS=$(go version)

buildTime=`date +%Y%m%d%H%M`

versionTime=${version}'_'${buildTime}

GOOS=${GOOS} go build -o gdav_${versionTime}  
