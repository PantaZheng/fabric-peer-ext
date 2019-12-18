// Copyright SecureKey Technologies Inc. All Rights Reserved.
//
// SPDX-License-Identifier: Apache-2.0

module github.com/trustbloc/fabric-peer-ext/test/bddtests

require (
	github.com/DATA-DOG/godog v0.7.13
	github.com/btcsuite/btcutil v0.0.0-20190425235716-9e5f4b9a998d
	github.com/hyperledger/fabric-protos-go v0.0.0
	github.com/hyperledger/fabric-sdk-go v1.0.0-beta1.0.20191219180315-e1055f391525
	github.com/pkg/errors v0.8.1
	github.com/spf13/viper v1.0.2
	github.com/trustbloc/fabric-peer-test-common v0.1.1-0.20191128150623-95d6e4c6a2ac
)

replace github.com/hyperledger/fabric-protos-go => github.com/trustbloc/fabric-protos-go-ext v0.1.1-0.20191126151100-5a61374c2e1b

go 1.13
