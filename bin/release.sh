#!/bin/sh

# keep track of script's total run time
SECONDS=0

HERE=$(cd $(dirname "$0"); pwd)
ROOT=${HERE}/../
LIB="${HERE}/lib"
. ${LIB}/props.functions

## if any simple command fails, abort script.
set -e

## lookup the local application version from project.props file
propsGetCurrentVersion APP_VERS
echo "Will build and deploy version ==> ${APP_VERS}"
echo ""

## make sure that we want this release
read -p "Build and release 'sideload-${APP_VERS}'? (y/N) " -r REPLY;
if [ "xy" != "x$REPLY" ]; then
    echo "whew! that was close!"
    exit 1;
fi

## watch commands as they execute. set +xv negates this
set -xv

## build and release project
cd $ROOT
goreleaser release --rm-dist

ELAPSED_TIME=$(date -u -d @"$SECONDS" +'%-Mm %-Ss')
echo "elapsed time: ${ELAPSED_TIME}"
