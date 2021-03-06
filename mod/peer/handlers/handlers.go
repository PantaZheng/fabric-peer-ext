/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package handlers

import (
	"github.com/hyperledger/fabric/core/handlers/auth"
	"github.com/hyperledger/fabric/core/handlers/decoration"

	"github.com/trustbloc/fabric-peer-ext/pkg/handler/authfilter"
)

// Get allows extensions to load an Auth filter.
// Returns the Auth filter if found; nil otherwise
func GetAuthFilter(name string) auth.Filter {
	return authfilter.Get(name)
}

// GetDecorator allows extensions to load a Decorator.
// Returns the Decorator if found; nil otherwise
func GetDecorator(name string) decoration.Decorator {
	// Not implemented
	return nil
}
