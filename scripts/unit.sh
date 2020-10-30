#!/bin/bash
#
# Copyright SecureKey Technologies Inc. All Rights Reserved.
#
# SPDX-License-Identifier: Apache-2.0
#
set -e



# Packages to exclude
PKGS=`go list github.com/trustbloc/fabric-peer-ext/pkg/... 2> /dev/null | \
                                                 grep -v /mocks | \
                                                 grep -v /api | \
                                                 grep -v /protos | \
                                                 grep -v /github.com`
echo "Running pkg unit tests..."
go test $PKGS -count=1 -tags "testing" -race -coverprofile=coverage.txt -covermode=atomic  -p 1 -timeout=10m