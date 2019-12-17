/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package scc

import (
	"testing"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	pb "github.com/hyperledger/fabric-protos-go/peer"
	"github.com/hyperledger/fabric/core/ledger"
	"github.com/stretchr/testify/require"
	"github.com/trustbloc/fabric-peer-ext/pkg/chaincode/builder"
	"github.com/trustbloc/fabric-peer-ext/pkg/mocks"
)

func TestCreateSCC(t *testing.T) {
	t.Run("No dependencies -> Success", func(t *testing.T) {
		sccBuilder = builder.New()
		Register(newSCCWithNoDependencies)

		descs := Create()
		require.NotNil(t, descs)
	})

	t.Run("Dependencies not satisfied -> panic", func(t *testing.T) {
		sccBuilder = builder.New()
		Register(newSCCWithDependencies)

		require.Panics(t, func() {
			Create()
		})
	})
	t.Run("Dependencies satisfied -> Success", func(t *testing.T) {
		sccBuilder = builder.New()
		Register(newSCCWithDependencies)

		descs := Create(mocks.NewQueryExecutorProvider())
		require.NotNil(t, descs)
	})
}

type testSCC struct {
}

func newSCCWithNoDependencies() *testSCC {
	return &testSCC{}
}

type queryExecutorProvider interface {
	GetQueryExecutorForLedger(cid string) (ledger.QueryExecutor, error)
}

func newSCCWithDependencies(qeProvider queryExecutorProvider) *testSCC {
	return &testSCC{}
}

func (scc *testSCC) Name() string              { return "testscc" }
func (scc *testSCC) Chaincode() shim.Chaincode { return scc }

func (scc *testSCC) Init(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success(nil)
}

func (scc *testSCC) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success(nil)
}
