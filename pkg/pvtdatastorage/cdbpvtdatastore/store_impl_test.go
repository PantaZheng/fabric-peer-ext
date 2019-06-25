/*
Copyright IBM Corp, SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package pvtdatastorage

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/hyperledger/fabric/common/ledger/util/leveldbhelper"
	"github.com/hyperledger/fabric/core/ledger"
	"github.com/hyperledger/fabric/core/ledger/kvledger/txmgmt/rwsetutil"
	btltestutil "github.com/hyperledger/fabric/core/ledger/pvtdatapolicy/testutil"
	"github.com/hyperledger/fabric/core/ledger/pvtdatastorage"
	"github.com/hyperledger/fabric/core/ledger/util/couchdb"
	"github.com/hyperledger/fabric/extensions/testutil"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
	"github.com/trustbloc/fabric-peer-ext/pkg/pvtdatastorage/common"
	"github.com/trustbloc/fabric-peer-ext/pkg/roles"
	xtestutil "github.com/trustbloc/fabric-peer-ext/pkg/testutil"
)

// This unit tests are copied from fabric, original file from fabric is found in fabric/core/ledger/pvtdatastorage/store_impl_test.go
// modification are made
// 1- setup couchdb
// 2- add TestLookupLastBlock unit test

var couchDBConfig *couchdb.Config

func TestMain(m *testing.M) {
	//setup extension test environment
	_, _, destroy := xtestutil.SetupExtTestEnv()

	viper.Set("peer.fileSystemPath", "/tmp/fabric/core/ledger/pvtdatastorage")
	// Create CouchDB definition from config parameters
	couchDBConfig = xtestutil.TestLedgerConf().StateDB.CouchDB

	code := m.Run()
	//stop couchdb
	destroy()
	os.Exit(code)
}

func TestHasPendingCommit(t *testing.T) {
	t.Run("test error Unmarshal", func(t *testing.T) {
		ledgerId := "ledger"
		env := NewTestStoreEnv(t, ledgerId, nil, couchDBConfig)
		defer env.Cleanup(ledgerId)
		s := env.TestStore.(*store)
		s.missingKeysIndexDB = mockDBHandler{getFunc: func(key []byte) (bytes []byte, e error) {
			return []byte("wrongData"), nil
		}}
		_, err := s.hasPendingCommit()
		require.Error(t, err)
		require.Contains(t, err.Error(), "Unmarshal failed pendingPvtData")
	})
}

func TestGetExpiryDataOfExpiryKey(t *testing.T) {
	t.Run("test error from getExpiryEntriesDB", func(t *testing.T) {
		ledgerId := "ledger"
		env := NewTestStoreEnv(t, ledgerId, nil, couchDBConfig)
		defer env.Cleanup(ledgerId)
		s := env.TestStore.(*store)
		s.db = mockCouchDB{readDocErr: fmt.Errorf("readDoc error")}
		_, err := s.getExpiryDataOfExpiryKey(&common.ExpiryKey{CommittingBlk: 1})
		require.Error(t, err)
		require.Contains(t, err.Error(), "getExpiryEntriesDB failed")
	})

	t.Run("test error from getExpiryEntriesDB", func(t *testing.T) {
		ledgerId := "ledger"
		env := NewTestStoreEnv(t, ledgerId, nil, couchDBConfig)
		defer env.Cleanup(ledgerId)
		s := env.TestStore.(*store)
		var blockPvtData blockPvtDataResponse
		blockPvtData.Expiry = make(map[string][]byte)
		b, err := json.Marshal(blockPvtData)
		require.NoError(t, err)
		s.db = mockCouchDB{readDocValue: &couchdb.CouchDoc{JSONValue: b}}
		expData, err := s.getExpiryDataOfExpiryKey(&common.ExpiryKey{CommittingBlk: 1})
		require.NoError(t, err)
		require.Nil(t, expData)
	})
}

func TestRetrieveBlockPvtEntries(t *testing.T) {
	t.Run("test error NotFoundInIndexErr", func(t *testing.T) {
		ledgerId := "ledger"
		env := NewTestStoreEnv(t, ledgerId, nil, couchDBConfig)
		defer env.Cleanup(ledgerId)
		s := env.TestStore.(*store)
		s.db = mockCouchDB{}
		_, _, _, err := s.retrieveBlockPvtEntries(0)
		require.NoError(t, err)
	})
}

func TestPreparePvtDataDoc(t *testing.T) {
	t.Run("test error from retrieveBlockPvtEntries", func(t *testing.T) {
		ledgerId := "ledger"
		env := NewTestStoreEnv(t, ledgerId, nil, couchDBConfig)
		defer env.Cleanup(ledgerId)
		s := env.TestStore.(*store)
		s.db = mockCouchDB{readDocErr: fmt.Errorf("readDoc error")}
		_, err := s.preparePvtDataDoc(0, nil)
		require.Error(t, err)
		require.Contains(t, err.Error(), "retrieveBlockPvtEntries failed")
	})
}

func TestLastCommittedBlockHeight(t *testing.T) {
	t.Run("test store is empty", func(t *testing.T) {
		ledgerId := "ledger"
		env := NewTestStoreEnv(t, ledgerId, nil, couchDBConfig)
		defer env.Cleanup(ledgerId)
		s := env.TestStore.(*store)
		s.isEmpty = true
		blockNum, err := s.LastCommittedBlockHeight()
		require.NoError(t, err)
		require.Equal(t, blockNum, uint64(0))

	})
}

func TestShutdown(t *testing.T) {
	ledgerId := "ledger"
	env := NewTestStoreEnv(t, ledgerId, nil, couchDBConfig)
	defer env.Cleanup(ledgerId)
	s := env.TestStore.(*store)
	s.Shutdown()
}

func TestGetLastUpdatedOldBlocksPvtData(t *testing.T) {
	t.Run("test error from GetLastUpdatedOldBlocksList", func(t *testing.T) {
		ledgerId := "ledger"
		env := NewTestStoreEnv(t, ledgerId, nil, couchDBConfig)
		defer env.Cleanup(ledgerId)
		s := env.TestStore.(*store)
		s.isLastUpdatedOldBlocksSet = true
		s.missingKeysIndexDB = mockDBHandler{getFunc: func(key []byte) (bytes []byte, e error) {
			return nil, fmt.Errorf("get error")
		}}
		_, err := s.GetLastUpdatedOldBlocksPvtData()
		require.Error(t, err)
		require.Contains(t, err.Error(), "GetLastUpdatedOldBlocksList failed")

	})

}

func TestResetLastUpdatedOldBlocksList(t *testing.T) {
	t.Run("test error from ResetLastUpdatedOldBlocksList", func(t *testing.T) {
		ledgerId := "ledger"
		env := NewTestStoreEnv(t, ledgerId, nil, couchDBConfig)
		defer env.Cleanup(ledgerId)
		s := env.TestStore.(*store)
		s.isLastUpdatedOldBlocksSet = true
		s.missingKeysIndexDB = mockDBHandler{writeBatchErr: fmt.Errorf("writeBatch error")}
		err := s.ResetLastUpdatedOldBlocksList()
		require.Error(t, err)
		require.Contains(t, err.Error(), "ResetLastUpdatedOldBlocksList failed")

	})

}

func TestCheckLastCommittedBlock(t *testing.T) {
	t.Run("test committer logic store is empty", func(t *testing.T) {
		ledgerId := "ledger"
		env := NewTestStoreEnv(t, ledgerId, nil, couchDBConfig)
		defer env.Cleanup(ledgerId)
		s := env.TestStore.(*store)
		s.isEmpty = true
		err := s.checkLastCommittedBlock(0)
		require.Error(t, err)
		require.Contains(t, err.Error(), "The store is empty")

	})

	t.Run("test committer logic blockNum is bigger from lastCommittedBlock", func(t *testing.T) {
		ledgerId := "ledger"
		env := NewTestStoreEnv(t, ledgerId, nil, couchDBConfig)
		defer env.Cleanup(ledgerId)
		s := env.TestStore.(*store)
		s.isEmpty = false
		s.lastCommittedBlock = 1
		err := s.checkLastCommittedBlock(2)
		require.Error(t, err)
		_, ok := err.(*pvtdatastorage.ErrOutOfRange)
		require.True(t, ok)
	})

	t.Run("test endorser logic error from lookupLastBlock", func(t *testing.T) {
		ledgerId := "ledger"
		env := NewTestStoreEnv(t, ledgerId, nil, couchDBConfig)
		defer env.Cleanup(ledgerId)
		s := env.TestStore.(*store)
		s.db = mockCouchDB{readDocErr: fmt.Errorf("readDoc error")}
		rolesValue := make(map[roles.Role]struct{})
		rolesValue[roles.EndorserRole] = struct{}{}
		roles.SetRoles(rolesValue)
		defer func() { roles.SetRoles(nil) }()
		err := s.checkLastCommittedBlock(0)
		require.Error(t, err)
		require.Contains(t, err.Error(), "lookupLastBlock failed")
	})

	t.Run("test endorser logic store is empty", func(t *testing.T) {
		ledgerId := "ledger"
		env := NewTestStoreEnv(t, ledgerId, nil, couchDBConfig)
		defer env.Cleanup(ledgerId)
		s := env.TestStore.(*store)
		s.db = mockCouchDB{}
		rolesValue := make(map[roles.Role]struct{})
		rolesValue[roles.EndorserRole] = struct{}{}
		roles.SetRoles(rolesValue)
		defer func() { roles.SetRoles(nil) }()
		err := s.checkLastCommittedBlock(0)
		require.Error(t, err)
		require.Contains(t, err.Error(), "The store is empty")
	})

	t.Run("test endorser logic blockNum is bigger from lastCommittedBlock", func(t *testing.T) {
		ledgerId := "ledger"
		env := NewTestStoreEnv(t, ledgerId, nil, couchDBConfig)
		defer env.Cleanup(ledgerId)
		s := env.TestStore.(*store)
		jsonBinary, err := json.Marshal(lastCommittedBlockResponse{Data: "1"})
		require.NoError(t, err)
		s.db = mockCouchDB{readDocValue: &couchdb.CouchDoc{JSONValue: jsonBinary}}
		rolesValue := make(map[roles.Role]struct{})
		rolesValue[roles.EndorserRole] = struct{}{}
		roles.SetRoles(rolesValue)
		defer func() { roles.SetRoles(nil) }()
		err = s.checkLastCommittedBlock(2)
		require.Error(t, err)
		_, ok := err.(*pvtdatastorage.ErrOutOfRange)
		require.True(t, ok)
	})

}

func TestRollback(t *testing.T) {
	t.Run("test error from calling rollback on a peer that is not a committer", func(t *testing.T) {
		ledgerId := "ledger"
		env := NewTestStoreEnv(t, ledgerId, nil, couchDBConfig)
		defer env.Cleanup(ledgerId)
		s := env.TestStore.(*store)
		rolesValue := make(map[roles.Role]struct{})
		rolesValue[roles.EndorserRole] = struct{}{}
		roles.SetRoles(rolesValue)
		defer func() { roles.SetRoles(nil) }()
		require.Panics(t, func() {
			s.Rollback()
		})
	})

	t.Run("test error from no pending batch to rollback", func(t *testing.T) {
		ledgerId := "ledger"
		env := NewTestStoreEnv(t, ledgerId, nil, couchDBConfig)
		defer env.Cleanup(ledgerId)
		s := env.TestStore.(*store)
		s.pendingPvtData = &pendingPvtData{BatchPending: false}
		err := s.Rollback()
		require.Error(t, err)
		require.Contains(t, err.Error(), "No pending batch to rollback")
	})

	t.Run("test error from delete", func(t *testing.T) {
		ledgerId := "ledger"
		env := NewTestStoreEnv(t, ledgerId, nil, couchDBConfig)
		defer env.Cleanup(ledgerId)
		s := env.TestStore.(*store)
		s.pendingPvtData = &pendingPvtData{BatchPending: true}
		s.missingKeysIndexDB = mockDBHandler{deleteErr: fmt.Errorf("delete error")}
		err := s.Rollback()
		require.Error(t, err)
		require.Contains(t, err.Error(), "delete PendingCommitKey failed")

	})
}

func TestCommit(t *testing.T) {
	t.Run("test error from calling commit on a peer that is not a committer", func(t *testing.T) {
		ledgerId := "ledger"
		env := NewTestStoreEnv(t, ledgerId, nil, couchDBConfig)
		defer env.Cleanup(ledgerId)
		s := env.TestStore.(*store)
		rolesValue := make(map[roles.Role]struct{})
		rolesValue[roles.EndorserRole] = struct{}{}
		roles.SetRoles(rolesValue)
		defer func() { roles.SetRoles(nil) }()
		require.Panics(t, func() {
			s.Commit()
		})
	})

	t.Run("test error from no pending batch to commit", func(t *testing.T) {
		ledgerId := "ledger"
		env := NewTestStoreEnv(t, ledgerId, nil, couchDBConfig)
		defer env.Cleanup(ledgerId)
		s := env.TestStore.(*store)
		s.pendingPvtData = &pendingPvtData{BatchPending: false}
		err := s.Commit()
		require.Error(t, err)
		require.Contains(t, err.Error(), "No pending batch to commit")
	})

	t.Run("test error from prepareLastCommittedBlockDoc", func(t *testing.T) {
		ledgerId := "ledger"
		env := NewTestStoreEnv(t, ledgerId, nil, couchDBConfig)
		defer env.Cleanup(ledgerId)
		s := env.TestStore.(*store)
		s.pendingPvtData = &pendingPvtData{BatchPending: true}
		s.db = mockCouchDB{readDocErr: fmt.Errorf("readDoc error")}
		err := s.Commit()
		require.Error(t, err)
		require.Contains(t, err.Error(), "prepareLastCommittedBlockDoc failed")

	})

	t.Run("test error from BatchUpdateDocuments", func(t *testing.T) {
		ledgerId := "ledger"
		env := NewTestStoreEnv(t, ledgerId, nil, couchDBConfig)
		defer env.Cleanup(ledgerId)
		s := env.TestStore.(*store)
		s.pendingPvtData = &pendingPvtData{BatchPending: true}
		s.db = mockCouchDB{batchUpdateDocumentsErr: fmt.Errorf("batchUpdateDocuments error")}
		err := s.Commit()
		require.Error(t, err)
		require.Contains(t, err.Error(), "writing private data to CouchDB failed")

	})

	t.Run("test error from missingKeysIndexDB WriteBatch", func(t *testing.T) {
		ledgerId := "ledger"
		env := NewTestStoreEnv(t, ledgerId, nil, couchDBConfig)
		defer env.Cleanup(ledgerId)
		s := env.TestStore.(*store)
		s.pendingPvtData = &pendingPvtData{BatchPending: true}
		s.db = mockCouchDB{}
		s.missingKeysIndexDB = mockDBHandler{writeBatchErr: fmt.Errorf("writeBatch error")}
		err := s.Commit()
		require.Error(t, err)
		require.Contains(t, err.Error(), "WriteBatch failed")

	})
}

func TestRetrieveBlockExpiryData(t *testing.T) {
	t.Run("test error from QueryDocuments", func(t *testing.T) {
		_, err := retrieveBlockExpiryData(mockCouchDB{queryDocumentsErr: fmt.Errorf("QueryDocuments error")}, "")
		require.Error(t, err)
		require.Contains(t, err.Error(), "QueryDocuments error")
	})

	t.Run("test error from Unmarshal", func(t *testing.T) {
		_, err := retrieveBlockExpiryData(mockCouchDB{queryDocumentsValue: []*couchdb.QueryResult{{Value: []byte("wrongData")}}}, "")
		require.Error(t, err)
		require.Contains(t, err.Error(), "result from DB is not JSON encoded")
	})

}

func TestNewErrNotFoundInIndex(t *testing.T) {
	require.Equal(t, NewErrNotFoundInIndex().Error(), "Entry not found in index")

}

func TestCreatePvtDataCouchDB(t *testing.T) {
	err := createPvtDataCouchDB(&mockCouchDB{createNewIndexWithRetryErr: fmt.Errorf("createNewIndexWithRetry error")})
	require.Error(t, err)
	require.Contains(t, err.Error(), "createNewIndexWithRetry error")
}

func TestGetPvtDataCouchInstance(t *testing.T) {
	t.Run("test error from ExistsWithRetry", func(t *testing.T) {
		err := getPvtDataCouchInstance(mockCouchDB{existsWithRetryErr: fmt.Errorf("ExistsWithRetry error")}, "")
		require.Error(t, err)
		require.Contains(t, err.Error(), "ExistsWithRetry error")
	})

	t.Run("test db not exists", func(t *testing.T) {
		err := getPvtDataCouchInstance(mockCouchDB{existsWithRetryValue: false}, "")
		require.Error(t, err)
		require.Contains(t, err.Error(), "DB not found")
	})

	t.Run("test error from IndexDesignDocExistsWithRetry", func(t *testing.T) {
		err := getPvtDataCouchInstance(mockCouchDB{existsWithRetryValue: true, indexDesignDocExistsWithRetryErr: fmt.Errorf("IndexDesignDocExistsWithRetry error")}, "")
		require.Error(t, err)
		require.Contains(t, err.Error(), "IndexDesignDocExistsWithRetry error")
	})

	t.Run("test index not exists", func(t *testing.T) {
		err := getPvtDataCouchInstance(mockCouchDB{existsWithRetryValue: true, indexDesignDocExistsWithRetryValue: false}, "")
		require.Error(t, err)
		require.Contains(t, err.Error(), "DB index not found")
	})
}

func TestRetrieveBlockPvtData(t *testing.T) {
	t.Run("test error from ReadDoc", func(t *testing.T) {
		_, err := retrieveBlockPvtData(mockCouchDB{readDocErr: fmt.Errorf("ReadDoc error")}, "")
		require.Error(t, err)
		require.Contains(t, err.Error(), "ReadDoc error")
	})

	t.Run("test error from Unmarshal", func(t *testing.T) {
		_, err := retrieveBlockPvtData(mockCouchDB{readDocValue: &couchdb.CouchDoc{JSONValue: []byte("wrongData")}}, "")
		require.Error(t, err)
		require.Contains(t, err.Error(), "result from DB is not JSON encoded")
	})

}

func TestNewProviderWithDBDef(t *testing.T) {
	_, err := newProviderWithDBDef(&couchdb.Config{Address: "123"}, nil, nil)
	require.Error(t, err)
	require.Contains(t, err.Error(), "obtaining CouchDB instance failed")
}

func TestOpenStore(t *testing.T) {
	t.Run("test error from createCouchDatabase", func(t *testing.T) {
		removeStorePath()
		conf := testutil.TestLedgerConf().PrivateData
		testStoreProvider, err := NewProvider(conf, testutil.TestLedgerConf())
		require.NoError(t, err)
		_, err = testStoreProvider.OpenStore("_")
		require.Error(t, err)
		require.Contains(t, err.Error(), "createCouchDatabase failed")

	})

	t.Run("test error from newCouchDatabase", func(t *testing.T) {
		removeStorePath()
		roles.IsCommitter()
		rolesValue := make(map[roles.Role]struct{})
		rolesValue[roles.EndorserRole] = struct{}{}
		roles.SetRoles(rolesValue)
		defer func() { roles.SetRoles(nil) }()
		conf := testutil.TestLedgerConf().PrivateData
		testStoreProvider, err := NewProvider(conf, testutil.TestLedgerConf())
		require.NoError(t, err)
		_, err = testStoreProvider.OpenStore("_")
		require.Error(t, err)
		require.Contains(t, err.Error(), "newCouchDatabase failed")

	})
}

func TestInitState(t *testing.T) {
	t.Run("test error from lookupLastBlock", func(t *testing.T) {
		ledgerId := "ledger"
		env := NewTestStoreEnv(t, ledgerId, nil, couchDBConfig)
		defer env.Cleanup(ledgerId)
		s := env.TestStore.(*store)
		s.db = mockCouchDB{readDocErr: fmt.Errorf("readDoc error")}
		err := s.initState()
		require.Error(t, err)
		require.Contains(t, err.Error(), "lookupLastBlock failed")
	})

	t.Run("test error from hasPendingCommit", func(t *testing.T) {
		ledgerId := "ledger"
		env := NewTestStoreEnv(t, ledgerId, nil, couchDBConfig)
		defer env.Cleanup(ledgerId)
		s := env.TestStore.(*store)
		jsonBinary, err := json.Marshal(lastCommittedBlockResponse{Data: "1"})
		require.NoError(t, err)
		s.db = mockCouchDB{readDocValue: &couchdb.CouchDoc{JSONValue: jsonBinary}}
		s.missingKeysIndexDB = mockDBHandler{getFunc: func(key []byte) (bytes []byte, e error) {
			return nil, fmt.Errorf("get error")
		}}
		err = s.initState()
		require.Error(t, err)
		require.Contains(t, err.Error(), "hasPendingCommit failed")
	})

	t.Run("test error from GetLastUpdatedOldBlocksList", func(t *testing.T) {
		ledgerId := "ledger"
		env := NewTestStoreEnv(t, ledgerId, nil, couchDBConfig)
		defer env.Cleanup(ledgerId)
		s := env.TestStore.(*store)
		jsonBinary, err := json.Marshal(lastCommittedBlockResponse{Data: "1"})
		require.NoError(t, err)
		s.db = mockCouchDB{readDocValue: &couchdb.CouchDoc{JSONValue: jsonBinary}}
		s.missingKeysIndexDB = mockDBHandler{getFunc: func(key []byte) ([]byte, error) {
			if bytes.Equal(key, common.LastUpdatedOldBlocksKey) {
				return nil, fmt.Errorf("get error")
			}
			return nil, nil
		}}
		err = s.initState()
		require.Error(t, err)
		require.Contains(t, err.Error(), "getLastUpdatedOldBlocksList failed")
	})

	t.Run("test success", func(t *testing.T) {
		ledgerId := "ledger"
		env := NewTestStoreEnv(t, ledgerId, nil, couchDBConfig)
		defer env.Cleanup(ledgerId)
		s := env.TestStore.(*store)
		jsonBinary, err := json.Marshal(lastCommittedBlockResponse{Data: "1"})
		require.NoError(t, err)
		s.db = mockCouchDB{readDocValue: &couchdb.CouchDoc{JSONValue: jsonBinary}}
		s.missingKeysIndexDB = mockDBHandler{getFunc: func(key []byte) ([]byte, error) {
			if bytes.Equal(key, common.LastUpdatedOldBlocksKey) {
				updatedBlksList := []uint64{1}
				buf := proto.NewBuffer(nil)
				require.NoError(t, buf.EncodeVarint(uint64(len(updatedBlksList))))
				for _, blkNum := range updatedBlksList {
					require.NoError(t, buf.EncodeVarint(blkNum))
				}
				return buf.Bytes(), nil
			}
			if bytes.Equal(key, common.PendingCommitKey) {
				return json.Marshal(pendingPvtData{})
			}
			return nil, nil
		}}
		err = s.initState()
		require.NoError(t, err)
	})
}

func TestStorePurge(t *testing.T) {
	ledgerid := "teststorepurge"
	viper.Set("ledger.pvtdataStore.purgeInterval", 2)
	btlPolicy := btltestutil.SampleBTLPolicy(
		map[[2]string]uint64{
			{"ns-1", "coll-1"}: 1,
			{"ns-1", "coll-2"}: 0,
			{"ns-2", "coll-1"}: 0,
			{"ns-2", "coll-2"}: 4,
			{"ns-3", "coll-1"}: 1,
			{"ns-3", "coll-2"}: 0,
		},
	)
	env := NewTestStoreEnv(t, ledgerid, btlPolicy, couchDBConfig)
	defer env.Cleanup(ledgerid)
	req := require.New(t)
	s := env.TestStore

	// no pvt data with block 0
	req.NoError(s.Prepare(0, nil, nil))
	req.NoError(s.Commit())

	// construct missing data for block 1
	blk1MissingData := make(ledger.TxMissingPvtDataMap)
	// eligible missing data in tx1
	blk1MissingData.Add(1, "ns-1", "coll-1", true)
	blk1MissingData.Add(1, "ns-1", "coll-2", true)
	// ineligible missing data in tx4
	blk1MissingData.Add(4, "ns-3", "coll-1", false)
	blk1MissingData.Add(4, "ns-3", "coll-2", false)

	// write pvt data for block 1
	testDataForBlk1 := []*ledger.TxPvtData{
		produceSamplePvtdata(t, 2, []string{"ns-1:coll-1", "ns-1:coll-2", "ns-2:coll-1", "ns-2:coll-2"}),
		produceSamplePvtdata(t, 4, []string{"ns-1:coll-1", "ns-1:coll-2", "ns-2:coll-1", "ns-2:coll-2"}),
	}
	req.NoError(s.Prepare(1, testDataForBlk1, blk1MissingData))
	req.NoError(s.Commit())

	// write pvt data for block 2
	req.NoError(s.Prepare(2, nil, nil))
	req.NoError(s.Commit())
	// data for ns-1:coll-1 and ns-2:coll-2 should exist in store
	ns1Coll1 := &common.DataKey{NsCollBlk: common.NsCollBlk{Ns: "ns-1", Coll: "coll-1", BlkNum: 1}, TxNum: 2}
	ns2Coll2 := &common.DataKey{NsCollBlk: common.NsCollBlk{Ns: "ns-2", Coll: "coll-2", BlkNum: 1}, TxNum: 2}

	// eligible missingData entries for ns-1:coll-1, ns-1:coll-2 (neverExpires) should exist in store
	ns1Coll1elgMD := &common.MissingDataKey{NsCollBlk: common.NsCollBlk{Ns: "ns-1", Coll: "coll-1", BlkNum: 1}, IsEligible: true}
	ns1Coll2elgMD := &common.MissingDataKey{NsCollBlk: common.NsCollBlk{Ns: "ns-1", Coll: "coll-2", BlkNum: 1}, IsEligible: true}

	// ineligible missingData entries for ns-3:col-1, ns-3:coll-2 (neverExpires) should exist in store
	ns3Coll1inelgMD := &common.MissingDataKey{NsCollBlk: common.NsCollBlk{Ns: "ns-3", Coll: "coll-1", BlkNum: 1}, IsEligible: false}
	ns3Coll2inelgMD := &common.MissingDataKey{NsCollBlk: common.NsCollBlk{Ns: "ns-3", Coll: "coll-2", BlkNum: 1}, IsEligible: false}

	testWaitForPurgerRoutineToFinish(s)
	req.True(testDataKeyExists(t, s, ns1Coll1))
	req.True(testDataKeyExists(t, s, ns2Coll2))

	req.True(testMissingDataKeyExists(t, s, ns1Coll1elgMD))
	req.True(testMissingDataKeyExists(t, s, ns1Coll2elgMD))

	req.True(testMissingDataKeyExists(t, s, ns3Coll1inelgMD))
	req.True(testMissingDataKeyExists(t, s, ns3Coll2inelgMD))

	// write pvt data for block 3
	req.NoError(s.Prepare(3, nil, nil))
	req.NoError(s.Commit())
	// data for ns-1:coll-1 and ns-2:coll-2 should exist in store (because purger should not be launched at block 3)
	testWaitForPurgerRoutineToFinish(s)
	req.True(testDataKeyExists(t, s, ns1Coll1))
	req.True(testDataKeyExists(t, s, ns2Coll2))
	// eligible missingData entries for ns-1:coll-1, ns-1:coll-2 (neverExpires) should exist in store
	req.True(testMissingDataKeyExists(t, s, ns1Coll1elgMD))
	req.True(testMissingDataKeyExists(t, s, ns1Coll2elgMD))
	// ineligible missingData entries for ns-3:col-1, ns-3:coll-2 (neverExpires) should exist in store
	req.True(testMissingDataKeyExists(t, s, ns3Coll1inelgMD))
	req.True(testMissingDataKeyExists(t, s, ns3Coll2inelgMD))

	// write pvt data for block 4
	req.NoError(s.Prepare(4, nil, nil))
	req.NoError(s.Commit())
	// data for ns-1:coll-1 should not exist in store (because purger should be launched at block 4)
	// but ns-2:coll-2 should exist because it expires at block 5
	testWaitForPurgerRoutineToFinish(s)
	req.False(testDataKeyExists(t, s, ns1Coll1))
	req.True(testDataKeyExists(t, s, ns2Coll2))
	// eligible missingData entries for ns-1:coll-1 should have expired and ns-1:coll-2 (neverExpires) should exist in store
	req.False(testMissingDataKeyExists(t, s, ns1Coll1elgMD))
	req.True(testMissingDataKeyExists(t, s, ns1Coll2elgMD))
	// ineligible missingData entries for ns-3:col-1 should have expired and ns-3:coll-2 (neverExpires) should exist in store
	req.False(testMissingDataKeyExists(t, s, ns3Coll1inelgMD))
	req.True(testMissingDataKeyExists(t, s, ns3Coll2inelgMD))

	// write pvt data for block 5
	req.NoError(s.Prepare(5, nil, nil))
	req.NoError(s.Commit())
	// ns-2:coll-2 should exist because though the data expires at block 5 but purger is launched every second block
	testWaitForPurgerRoutineToFinish(s)
	req.False(testDataKeyExists(t, s, ns1Coll1))
	req.True(testDataKeyExists(t, s, ns2Coll2))

	// write pvt data for block 6
	req.NoError(s.Prepare(6, nil, nil))
	req.NoError(s.Commit())
	// ns-2:coll-2 should not exists now (because purger should be launched at block 6)
	testWaitForPurgerRoutineToFinish(s)
	req.False(testDataKeyExists(t, s, ns1Coll1))
	req.False(testDataKeyExists(t, s, ns2Coll2))

	// "ns-2:coll-1" should never have been purged (because, it was no btl was declared for this)
	req.True(testDataKeyExists(t, s, &common.DataKey{NsCollBlk: common.NsCollBlk{Ns: "ns-1", Coll: "coll-2", BlkNum: 1}, TxNum: 2}))

}

func testWaitForPurgerRoutineToFinish(s pvtdatastorage.Store) {
	time.Sleep(1 * time.Second)
	s.(*store).purgerLock.Lock()
	s.(*store).purgerLock.Unlock()
}

func TestEmptyStore(t *testing.T) {
	env := NewTestStoreEnv(t, "testemptystore", nil, couchDBConfig)
	defer env.Cleanup("testemptystore")
	req := require.New(t)
	store := env.TestStore
	testEmpty(true, req, store)
	testPendingBatch(false, req, store)
}

func TestEndorserRole(t *testing.T) {
	btlPolicy := btltestutil.SampleBTLPolicy(
		map[[2]string]uint64{
			{"ns-1", "coll-1"}: 0,
			{"ns-1", "coll-2"}: 0,
			{"ns-2", "coll-1"}: 0,
			{"ns-2", "coll-2"}: 0,
			{"ns-3", "coll-1"}: 0,
			{"ns-4", "coll-1"}: 0,
			{"ns-4", "coll-2"}: 0,
		},
	)
	env := NewTestStoreEnv(t, "testendorserrole", btlPolicy, couchDBConfig)
	defer env.Cleanup("testendorserrole")
	req := require.New(t)
	committerStore := env.TestStore
	testData := []*ledger.TxPvtData{
		produceSamplePvtdata(t, 2, []string{"ns-1:coll-1", "ns-1:coll-2", "ns-2:coll-1", "ns-2:coll-2"}),
		produceSamplePvtdata(t, 4, []string{"ns-1:coll-1", "ns-1:coll-2", "ns-2:coll-1", "ns-2:coll-2"}),
	}
	// no pvt data with block 0
	req.NoError(committerStore.Prepare(0, nil, nil))
	req.NoError(committerStore.Commit())

	// pvt data with block 1 - commit
	req.NoError(committerStore.Prepare(1, testData, nil))
	req.NoError(committerStore.Commit())

	// create endorser store
	rolesValue := make(map[roles.Role]struct{})
	rolesValue[roles.EndorserRole] = struct{}{}
	roles.SetRoles(rolesValue)
	defer func() { roles.SetRoles(nil) }()
	endorserStore := NewTestStoreEnv(t, "testendorserrole", btlPolicy, couchDBConfig).TestStore

	var nilFilter ledger.PvtNsCollFilter
	retrievedData, err := endorserStore.GetPvtDataByBlockNum(0, nilFilter)
	req.NoError(err)
	req.Nil(retrievedData)

	// pvt data retrieval for block 1 should return full pvtdata
	retrievedData, err = endorserStore.GetPvtDataByBlockNum(1, nilFilter)
	req.NoError(err)
	for i, data := range retrievedData {
		req.Equal(data.SeqInBlock, testData[i].SeqInBlock)
		req.True(proto.Equal(data.WriteSet, testData[i].WriteSet))
	}

}

func TestStoreBasicCommitAndRetrieval(t *testing.T) {
	btlPolicy := btltestutil.SampleBTLPolicy(
		map[[2]string]uint64{
			{"ns-1", "coll-1"}: 0,
			{"ns-1", "coll-2"}: 0,
			{"ns-2", "coll-1"}: 0,
			{"ns-2", "coll-2"}: 0,
			{"ns-3", "coll-1"}: 0,
			{"ns-4", "coll-1"}: 0,
			{"ns-4", "coll-2"}: 0,
		},
	)

	env := NewTestStoreEnv(t, "teststorebasiccommitandretrieval", btlPolicy, couchDBConfig)
	defer env.Cleanup("teststorebasiccommitandretrieval")
	req := require.New(t)
	store := env.TestStore
	testData := []*ledger.TxPvtData{
		produceSamplePvtdata(t, 2, []string{"ns-1:coll-1", "ns-1:coll-2", "ns-2:coll-1", "ns-2:coll-2"}),
		produceSamplePvtdata(t, 4, []string{"ns-1:coll-1", "ns-1:coll-2", "ns-2:coll-1", "ns-2:coll-2"}),
	}

	// construct missing data for block 1
	blk1MissingData := make(ledger.TxMissingPvtDataMap)

	// eligible missing data in tx1
	blk1MissingData.Add(1, "ns-1", "coll-1", true)
	blk1MissingData.Add(1, "ns-1", "coll-2", true)
	blk1MissingData.Add(1, "ns-2", "coll-1", true)
	blk1MissingData.Add(1, "ns-2", "coll-2", true)
	// eligible missing data in tx2
	blk1MissingData.Add(2, "ns-3", "coll-1", true)
	// ineligible missing data in tx4
	blk1MissingData.Add(4, "ns-4", "coll-1", false)
	blk1MissingData.Add(4, "ns-4", "coll-2", false)

	// construct missing data for block 2
	blk2MissingData := make(ledger.TxMissingPvtDataMap)
	// eligible missing data in tx1
	blk2MissingData.Add(1, "ns-1", "coll-1", true)
	blk2MissingData.Add(1, "ns-1", "coll-2", true)
	// eligible missing data in tx3
	blk2MissingData.Add(3, "ns-1", "coll-1", true)

	// no pvt data with block 0
	req.NoError(store.Prepare(0, nil, nil))
	req.NoError(store.Commit())

	// pvt data with block 1 - commit
	req.NoError(store.Prepare(1, testData, blk1MissingData))
	req.NoError(store.Commit())

	// pvt data with block 2 - rollback
	req.NoError(store.Prepare(2, testData, nil))
	req.NoError(store.Rollback())

	// pvt data retrieval for block 0 should return nil
	var nilFilter ledger.PvtNsCollFilter
	retrievedData, err := store.GetPvtDataByBlockNum(0, nilFilter)
	req.NoError(err)
	req.Nil(retrievedData)

	// pvt data retrieval for block 1 should return full pvtdata
	retrievedData, err = store.GetPvtDataByBlockNum(1, nilFilter)
	req.NoError(err)
	for i, data := range retrievedData {
		req.Equal(data.SeqInBlock, testData[i].SeqInBlock)
		req.True(proto.Equal(data.WriteSet, testData[i].WriteSet))
	}

	// pvt data retrieval for block 1 with filter should return filtered pvtdata
	filter := ledger.NewPvtNsCollFilter()
	filter.Add("ns-1", "coll-1")
	filter.Add("ns-2", "coll-2")
	retrievedData, err = store.GetPvtDataByBlockNum(1, filter)
	expectedRetrievedData := []*ledger.TxPvtData{
		produceSamplePvtdata(t, 2, []string{"ns-1:coll-1", "ns-2:coll-2"}),
		produceSamplePvtdata(t, 4, []string{"ns-1:coll-1", "ns-2:coll-2"}),
	}
	for i, data := range retrievedData {
		req.Equal(data.SeqInBlock, expectedRetrievedData[i].SeqInBlock)
		req.True(proto.Equal(data.WriteSet, expectedRetrievedData[i].WriteSet))
	}

	// pvt data retrieval for block 2 should return ErrOutOfRange
	retrievedData, err = store.GetPvtDataByBlockNum(2, nilFilter)
	_, ok := err.(*pvtdatastorage.ErrOutOfRange)
	req.True(ok)
	req.Nil(retrievedData)

	// pvt data with block 2 - commit
	req.NoError(store.Prepare(2, testData, blk2MissingData))
	req.NoError(store.Commit())

	// retrieve the stored missing entries using GetMissingPvtDataInfoForMostRecentBlocks
	// Only the code path of eligible entries would be covered in this unit-test. For
	// ineligible entries, the code path will be covered in FAB-11437

	expectedMissingPvtDataInfo := make(ledger.MissingPvtDataInfo)
	// missing data in block2, tx1
	expectedMissingPvtDataInfo.Add(2, 1, "ns-1", "coll-1")
	expectedMissingPvtDataInfo.Add(2, 1, "ns-1", "coll-2")
	expectedMissingPvtDataInfo.Add(2, 3, "ns-1", "coll-1")

	missingPvtDataInfo, err := store.GetMissingPvtDataInfoForMostRecentBlocks(1)
	req.NoError(err)
	req.Equal(expectedMissingPvtDataInfo, missingPvtDataInfo)

	// missing data in block1, tx1
	expectedMissingPvtDataInfo.Add(1, 1, "ns-1", "coll-1")
	expectedMissingPvtDataInfo.Add(1, 1, "ns-1", "coll-2")
	expectedMissingPvtDataInfo.Add(1, 1, "ns-2", "coll-1")
	expectedMissingPvtDataInfo.Add(1, 1, "ns-2", "coll-2")

	// missing data in block1, tx2
	expectedMissingPvtDataInfo.Add(1, 2, "ns-3", "coll-1")

	missingPvtDataInfo, err = store.GetMissingPvtDataInfoForMostRecentBlocks(2)
	req.NoError(err)
	req.Equal(expectedMissingPvtDataInfo, missingPvtDataInfo)

	missingPvtDataInfo, err = store.GetMissingPvtDataInfoForMostRecentBlocks(10)
	req.NoError(err)
	req.Equal(expectedMissingPvtDataInfo, missingPvtDataInfo)
}

func TestCommitPvtDataOfOldBlocks(t *testing.T) {
	t.Run("test error lastUpdatedOldBlocksList is set", func(t *testing.T) {
		ledgerId := "ledger"
		env := NewTestStoreEnv(t, ledgerId, nil, couchDBConfig)
		defer env.Cleanup(ledgerId)
		s := env.TestStore.(*store)
		s.isLastUpdatedOldBlocksSet = true
		err := s.CommitPvtDataOfOldBlocks(nil)
		require.Error(t, err)
		require.Contains(t, err.Error(), "lastUpdatedOldBlocksList is set")

	})

	t.Run("test error from BatchUpdateDocuments", func(t *testing.T) {
		ledgerId := "ledger"
		env := NewTestStoreEnv(t, ledgerId, nil, couchDBConfig)
		defer env.Cleanup(ledgerId)
		s := env.TestStore.(*store)
		s.isLastUpdatedOldBlocksSet = false
		s.db = mockCouchDB{batchUpdateDocumentsErr: fmt.Errorf("batchUpdateDocuments error")}
		err := s.CommitPvtDataOfOldBlocks(nil)
		require.Error(t, err)
		require.Contains(t, err.Error(), "BatchUpdateDocuments failed")

	})

	t.Run("test error from WriteBatch", func(t *testing.T) {
		ledgerId := "ledger"
		env := NewTestStoreEnv(t, ledgerId, nil, couchDBConfig)
		defer env.Cleanup(ledgerId)
		s := env.TestStore.(*store)
		s.isLastUpdatedOldBlocksSet = false
		s.db = mockCouchDB{}
		s.missingKeysIndexDB = mockDBHandler{writeBatchErr: fmt.Errorf("WriteBatch error")}
		err := s.CommitPvtDataOfOldBlocks(nil)
		require.Error(t, err)
		require.Contains(t, err.Error(), "WriteBatch failed")

	})

	viper.Set("ledger.pvtdataStore.purgeInterval", 2)
	btlPolicy := btltestutil.SampleBTLPolicy(
		map[[2]string]uint64{
			{"ns-1", "coll-1"}: 3,
			{"ns-1", "coll-2"}: 1,
			{"ns-2", "coll-1"}: 0,
			{"ns-2", "coll-2"}: 1,
			{"ns-3", "coll-1"}: 0,
			{"ns-3", "coll-2"}: 3,
			{"ns-4", "coll-1"}: 0,
			{"ns-4", "coll-2"}: 0,
		},
	)
	env := NewTestStoreEnv(t, "testcommitpvtdataofoldblocks", btlPolicy, couchDBConfig)
	defer env.Cleanup("testcommitpvtdataofoldblocks")
	req := require.New(t)
	store := env.TestStore

	testData := []*ledger.TxPvtData{
		produceSamplePvtdata(t, 2, []string{"ns-2:coll-1", "ns-2:coll-2"}),
		produceSamplePvtdata(t, 4, []string{"ns-1:coll-1", "ns-1:coll-2", "ns-2:coll-1", "ns-2:coll-2"}),
	}

	// CONSTRUCT MISSING DATA FOR BLOCK 1
	blk1MissingData := make(ledger.TxMissingPvtDataMap)

	// eligible missing data in tx1
	blk1MissingData.Add(1, "ns-1", "coll-1", true)
	blk1MissingData.Add(1, "ns-1", "coll-2", true)
	blk1MissingData.Add(1, "ns-2", "coll-1", true)
	blk1MissingData.Add(1, "ns-2", "coll-2", true)
	// eligible missing data in tx2
	blk1MissingData.Add(2, "ns-1", "coll-1", true)
	blk1MissingData.Add(2, "ns-1", "coll-2", true)
	blk1MissingData.Add(2, "ns-3", "coll-1", true)
	blk1MissingData.Add(2, "ns-3", "coll-2", true)

	// CONSTRUCT MISSING DATA FOR BLOCK 2
	blk2MissingData := make(ledger.TxMissingPvtDataMap)
	// eligible missing data in tx1
	blk2MissingData.Add(1, "ns-1", "coll-1", true)
	blk2MissingData.Add(1, "ns-1", "coll-2", true)
	// eligible missing data in tx3
	blk2MissingData.Add(3, "ns-1", "coll-1", true)

	// COMMIT BLOCK 0 WITH NO DATA
	req.NoError(store.Prepare(0, nil, nil))
	req.NoError(store.Commit())

	// COMMIT BLOCK 1 WITH PVTDATA AND MISSINGDATA
	req.NoError(store.Prepare(1, testData, blk1MissingData))
	req.NoError(store.Commit())

	// COMMIT BLOCK 2 WITH PVTDATA AND MISSINGDATA
	req.NoError(store.Prepare(2, nil, blk2MissingData))
	req.NoError(store.Commit())

	// CHECK MISSINGDATA ENTRIES ARE CORRECTLY STORED
	expectedMissingPvtDataInfo := make(ledger.MissingPvtDataInfo)
	// missing data in block1, tx1
	expectedMissingPvtDataInfo.Add(1, 1, "ns-1", "coll-1")
	expectedMissingPvtDataInfo.Add(1, 1, "ns-1", "coll-2")
	expectedMissingPvtDataInfo.Add(1, 1, "ns-2", "coll-1")
	expectedMissingPvtDataInfo.Add(1, 1, "ns-2", "coll-2")

	// missing data in block1, tx2
	expectedMissingPvtDataInfo.Add(1, 2, "ns-1", "coll-1")
	expectedMissingPvtDataInfo.Add(1, 2, "ns-1", "coll-2")
	expectedMissingPvtDataInfo.Add(1, 2, "ns-3", "coll-1")
	expectedMissingPvtDataInfo.Add(1, 2, "ns-3", "coll-2")

	// missing data in block2, tx1
	expectedMissingPvtDataInfo.Add(2, 1, "ns-1", "coll-1")
	expectedMissingPvtDataInfo.Add(2, 1, "ns-1", "coll-2")
	// missing data in block2, tx3
	expectedMissingPvtDataInfo.Add(2, 3, "ns-1", "coll-1")

	missingPvtDataInfo, err := store.GetMissingPvtDataInfoForMostRecentBlocks(2)
	req.NoError(err)
	req.Equal(expectedMissingPvtDataInfo, missingPvtDataInfo)

	// COMMIT THE MISSINGDATA IN BLOCK 1 AND BLOCK 2
	oldBlocksPvtData := make(map[uint64][]*ledger.TxPvtData)
	oldBlocksPvtData[1] = []*ledger.TxPvtData{
		produceSamplePvtdata(t, 1, []string{"ns-1:coll-1", "ns-2:coll-1"}),
		produceSamplePvtdata(t, 2, []string{"ns-1:coll-1", "ns-3:coll-1"}),
	}
	oldBlocksPvtData[2] = []*ledger.TxPvtData{
		produceSamplePvtdata(t, 3, []string{"ns-1:coll-1"}),
	}

	err = store.CommitPvtDataOfOldBlocks(oldBlocksPvtData)
	req.NoError(err)

	// ENSURE THAT THE CURRENT PVTDATA OF BLOCK 1 STILL EXIST IN THE STORE
	ns2Coll1Blk1Tx2 := &common.DataKey{NsCollBlk: common.NsCollBlk{Ns: "ns-2", Coll: "coll-1", BlkNum: 1}, TxNum: 2}
	ns2Coll2Blk1Tx2 := &common.DataKey{NsCollBlk: common.NsCollBlk{Ns: "ns-2", Coll: "coll-2", BlkNum: 1}, TxNum: 2}
	req.True(testDataKeyExists(t, store, ns2Coll1Blk1Tx2))
	req.True(testDataKeyExists(t, store, ns2Coll2Blk1Tx2))

	// ENSURE THAT THE PREVIOUSLY MISSING PVTDATA OF BLOCK 1 & 2 EXIST IN THE STORE
	ns1Coll1Blk1Tx1 := &common.DataKey{NsCollBlk: common.NsCollBlk{Ns: "ns-1", Coll: "coll-1", BlkNum: 1}, TxNum: 1}
	ns2Coll1Blk1Tx1 := &common.DataKey{NsCollBlk: common.NsCollBlk{Ns: "ns-2", Coll: "coll-1", BlkNum: 1}, TxNum: 1}
	ns1Coll1Blk1Tx2 := &common.DataKey{NsCollBlk: common.NsCollBlk{Ns: "ns-1", Coll: "coll-1", BlkNum: 1}, TxNum: 2}
	ns3Coll1Blk1Tx2 := &common.DataKey{NsCollBlk: common.NsCollBlk{Ns: "ns-3", Coll: "coll-1", BlkNum: 1}, TxNum: 2}
	ns1Coll1Blk2Tx3 := &common.DataKey{NsCollBlk: common.NsCollBlk{Ns: "ns-1", Coll: "coll-1", BlkNum: 2}, TxNum: 3}

	req.True(testDataKeyExists(t, store, ns1Coll1Blk1Tx1))
	req.True(testDataKeyExists(t, store, ns2Coll1Blk1Tx1))
	req.True(testDataKeyExists(t, store, ns1Coll1Blk1Tx2))
	req.True(testDataKeyExists(t, store, ns3Coll1Blk1Tx2))
	req.True(testDataKeyExists(t, store, ns1Coll1Blk2Tx3))

	// pvt data retrieval for block 2 should return the just committed pvtdata
	var nilFilter ledger.PvtNsCollFilter
	retrievedData, err := store.GetPvtDataByBlockNum(2, nilFilter)
	req.NoError(err)
	for i, data := range retrievedData {
		req.Equal(data.SeqInBlock, oldBlocksPvtData[2][i].SeqInBlock)
		req.True(proto.Equal(data.WriteSet, oldBlocksPvtData[2][i].WriteSet))
	}

	expectedMissingPvtDataInfo = make(ledger.MissingPvtDataInfo)
	// missing data in block1, tx1
	expectedMissingPvtDataInfo.Add(1, 1, "ns-1", "coll-2")
	expectedMissingPvtDataInfo.Add(1, 1, "ns-2", "coll-2")

	// missing data in block1, tx2
	expectedMissingPvtDataInfo.Add(1, 2, "ns-1", "coll-2")
	expectedMissingPvtDataInfo.Add(1, 2, "ns-3", "coll-2")

	// missing data in block2, tx1
	expectedMissingPvtDataInfo.Add(2, 1, "ns-1", "coll-1")
	expectedMissingPvtDataInfo.Add(2, 1, "ns-1", "coll-2")

	missingPvtDataInfo, err = store.GetMissingPvtDataInfoForMostRecentBlocks(2)
	req.NoError(err)
	req.Equal(expectedMissingPvtDataInfo, missingPvtDataInfo)

	// blksPvtData returns all the pvt data for a block for which the any pvtdata has been submitted
	// using CommitPvtDataOfOldBlocks
	blksPvtData, err := store.GetLastUpdatedOldBlocksPvtData()
	req.NoError(err)

	expectedLastupdatedPvtdata := make(map[uint64][]*ledger.TxPvtData)
	expectedLastupdatedPvtdata[1] = []*ledger.TxPvtData{
		produceSamplePvtdata(t, 1, []string{"ns-1:coll-1", "ns-2:coll-1"}),
		produceSamplePvtdata(t, 2, []string{"ns-1:coll-1", "ns-2:coll-1", "ns-2:coll-2", "ns-3:coll-1"}),
		produceSamplePvtdata(t, 4, []string{"ns-1:coll-1", "ns-1:coll-2", "ns-2:coll-1", "ns-2:coll-2"}),
	}
	expectedLastupdatedPvtdata[2] = []*ledger.TxPvtData{
		produceSamplePvtdata(t, 3, []string{"ns-1:coll-1"}),
	}

	req.Equal(expectedLastupdatedPvtdata, blksPvtData)

	err = store.ResetLastUpdatedOldBlocksList()
	req.NoError(err)

	blksPvtData, err = store.GetLastUpdatedOldBlocksPvtData()
	req.NoError(err)
	req.Nil(blksPvtData)

	// COMMIT BLOCK 3 WITH NO PVTDATA
	req.NoError(store.Prepare(3, nil, nil))
	req.NoError(store.Commit())

	// IN BLOCK 1, NS-1:COLL-2 AND NS-2:COLL-2 SHOULD HAVE EXPIRED BUT NOT PURGED
	// HENCE, THE FOLLOWING COMMIT SHOULD CREATE ENTRIES IN THE STORE
	oldBlocksPvtData = make(map[uint64][]*ledger.TxPvtData)
	oldBlocksPvtData[1] = []*ledger.TxPvtData{
		produceSamplePvtdata(t, 1, []string{"ns-1:coll-2"}), // though expired, it
		// would get committed to the store as it is not purged yet
		produceSamplePvtdata(t, 2, []string{"ns-3:coll-2"}), // never expires
	}

	err = store.CommitPvtDataOfOldBlocks(oldBlocksPvtData)
	req.NoError(err)

	ns1Coll2Blk1Tx1 := &common.DataKey{NsCollBlk: common.NsCollBlk{Ns: "ns-1", Coll: "coll-2", BlkNum: 1}, TxNum: 1}
	ns2Coll2Blk1Tx1 := &common.DataKey{NsCollBlk: common.NsCollBlk{Ns: "ns-2", Coll: "coll-2", BlkNum: 1}, TxNum: 1}
	ns1Coll2Blk1Tx2 := &common.DataKey{NsCollBlk: common.NsCollBlk{Ns: "ns-1", Coll: "coll-2", BlkNum: 1}, TxNum: 2}
	ns3Coll2Blk1Tx2 := &common.DataKey{NsCollBlk: common.NsCollBlk{Ns: "ns-3", Coll: "coll-2", BlkNum: 1}, TxNum: 2}

	// though the pvtdata are expired but not purged yet, we do
	// commit the data and hence the entries would exist in the
	// store
	req.True(testDataKeyExists(t, store, ns1Coll2Blk1Tx1))  // expired but committed
	req.False(testDataKeyExists(t, store, ns2Coll2Blk1Tx1)) // expired but still missing
	req.False(testDataKeyExists(t, store, ns1Coll2Blk1Tx2)) // expired still missing
	req.True(testDataKeyExists(t, store, ns3Coll2Blk1Tx2))  // never expires

	err = store.ResetLastUpdatedOldBlocksList()
	req.NoError(err)

	// COMMIT BLOCK 4 WITH NO PVTDATA
	req.NoError(store.Prepare(4, nil, nil))
	req.NoError(store.Commit())

	testWaitForPurgerRoutineToFinish(store)

	// IN BLOCK 1, NS-1:COLL-2 AND NS-2:COLL-2 SHOULD HAVE EXPIRED BUT NOT PURGED
	// HENCE, THE FOLLOWING COMMIT SHOULD NOT CREATE ENTRIES IN THE STORE
	oldBlocksPvtData = make(map[uint64][]*ledger.TxPvtData)
	oldBlocksPvtData[1] = []*ledger.TxPvtData{
		// both data are expired and purged. hence, it won't be
		// committed to the store
		produceSamplePvtdata(t, 1, []string{"ns-2:coll-2"}),
		produceSamplePvtdata(t, 2, []string{"ns-1:coll-2"}),
	}

	err = store.CommitPvtDataOfOldBlocks(oldBlocksPvtData)
	req.NoError(err)

	ns1Coll2Blk1Tx1 = &common.DataKey{NsCollBlk: common.NsCollBlk{Ns: "ns-1", Coll: "coll-2", BlkNum: 1}, TxNum: 1}
	ns2Coll2Blk1Tx1 = &common.DataKey{NsCollBlk: common.NsCollBlk{Ns: "ns-2", Coll: "coll-2", BlkNum: 1}, TxNum: 1}
	ns1Coll2Blk1Tx2 = &common.DataKey{NsCollBlk: common.NsCollBlk{Ns: "ns-1", Coll: "coll-2", BlkNum: 1}, TxNum: 2}
	ns3Coll2Blk1Tx2 = &common.DataKey{NsCollBlk: common.NsCollBlk{Ns: "ns-3", Coll: "coll-2", BlkNum: 1}, TxNum: 2}

	req.False(testDataKeyExists(t, store, ns1Coll2Blk1Tx1)) // purged
	req.False(testDataKeyExists(t, store, ns2Coll2Blk1Tx1)) // purged
	req.False(testDataKeyExists(t, store, ns1Coll2Blk1Tx2)) // purged
	req.True(testDataKeyExists(t, store, ns3Coll2Blk1Tx2))  // never expires
}

func TestExpiryDataNotIncluded(t *testing.T) {
	ledgerid := "testexpirydatanotincluded"
	btlPolicy := btltestutil.SampleBTLPolicy(
		map[[2]string]uint64{
			{"ns-1", "coll-1"}: 1,
			{"ns-1", "coll-2"}: 0,
			{"ns-2", "coll-1"}: 0,
			{"ns-2", "coll-2"}: 2,
			{"ns-3", "coll-1"}: 1,
			{"ns-3", "coll-2"}: 0,
		},
	)
	env := NewTestStoreEnv(t, ledgerid, btlPolicy, couchDBConfig)
	defer env.Cleanup(ledgerid)
	req := require.New(t)
	store := env.TestStore

	// construct missing data for block 1
	blk1MissingData := make(ledger.TxMissingPvtDataMap)
	// eligible missing data in tx1
	blk1MissingData.Add(1, "ns-1", "coll-1", true)
	blk1MissingData.Add(1, "ns-1", "coll-2", true)
	// ineligible missing data in tx4
	blk1MissingData.Add(4, "ns-3", "coll-1", false)
	blk1MissingData.Add(4, "ns-3", "coll-2", false)

	// construct missing data for block 2
	blk2MissingData := make(ledger.TxMissingPvtDataMap)
	// eligible missing data in tx1
	blk2MissingData.Add(1, "ns-1", "coll-1", true)
	blk2MissingData.Add(1, "ns-1", "coll-2", true)

	// no pvt data with block 0
	req.NoError(store.Prepare(0, nil, nil))
	req.NoError(store.Commit())

	// write pvt data for block 1
	testDataForBlk1 := []*ledger.TxPvtData{
		produceSamplePvtdata(t, 2, []string{"ns-1:coll-1", "ns-1:coll-2", "ns-2:coll-1", "ns-2:coll-2"}),
		produceSamplePvtdata(t, 4, []string{"ns-1:coll-1", "ns-1:coll-2", "ns-2:coll-1", "ns-2:coll-2"}),
	}
	req.NoError(store.Prepare(1, testDataForBlk1, blk1MissingData))
	req.NoError(store.Commit())

	// write pvt data for block 2
	testDataForBlk2 := []*ledger.TxPvtData{
		produceSamplePvtdata(t, 3, []string{"ns-1:coll-1", "ns-1:coll-2", "ns-2:coll-1", "ns-2:coll-2"}),
		produceSamplePvtdata(t, 5, []string{"ns-1:coll-1", "ns-1:coll-2", "ns-2:coll-1", "ns-2:coll-2"}),
	}
	req.NoError(store.Prepare(2, testDataForBlk2, blk2MissingData))
	req.NoError(store.Commit())

	retrievedData, _ := store.GetPvtDataByBlockNum(1, nil)
	// block 1 data should still be not expired
	for i, data := range retrievedData {
		req.Equal(data.SeqInBlock, testDataForBlk1[i].SeqInBlock)
		req.True(proto.Equal(data.WriteSet, testDataForBlk1[i].WriteSet))
	}

	// none of the missing data entries would have expired
	expectedMissingPvtDataInfo := make(ledger.MissingPvtDataInfo)
	// missing data in block2, tx1
	expectedMissingPvtDataInfo.Add(2, 1, "ns-1", "coll-1")
	expectedMissingPvtDataInfo.Add(2, 1, "ns-1", "coll-2")

	// missing data in block1, tx1
	expectedMissingPvtDataInfo.Add(1, 1, "ns-1", "coll-1")
	expectedMissingPvtDataInfo.Add(1, 1, "ns-1", "coll-2")

	missingPvtDataInfo, err := store.GetMissingPvtDataInfoForMostRecentBlocks(10)
	req.NoError(err)
	req.Equal(expectedMissingPvtDataInfo, missingPvtDataInfo)

	// Commit block 3 with no pvtdata
	req.NoError(store.Prepare(3, nil, nil))
	req.NoError(store.Commit())

	// After committing block 3, the data for "ns-1:coll1" of block 1 should have expired and should not be returned by the store
	expectedPvtdataFromBlock1 := []*ledger.TxPvtData{
		produceSamplePvtdata(t, 2, []string{"ns-1:coll-2", "ns-2:coll-1", "ns-2:coll-2"}),
		produceSamplePvtdata(t, 4, []string{"ns-1:coll-2", "ns-2:coll-1", "ns-2:coll-2"}),
	}
	retrievedData, _ = store.GetPvtDataByBlockNum(1, nil)
	req.Equal(expectedPvtdataFromBlock1, retrievedData)

	// After committing block 3, the missing data of "ns1-coll1" in block1-tx1 should have expired
	expectedMissingPvtDataInfo = make(ledger.MissingPvtDataInfo)
	// missing data in block2, tx1
	expectedMissingPvtDataInfo.Add(2, 1, "ns-1", "coll-1")
	expectedMissingPvtDataInfo.Add(2, 1, "ns-1", "coll-2")
	// missing data in block1, tx1
	expectedMissingPvtDataInfo.Add(1, 1, "ns-1", "coll-2")

	missingPvtDataInfo, err = store.GetMissingPvtDataInfoForMostRecentBlocks(10)
	req.NoError(err)
	req.Equal(expectedMissingPvtDataInfo, missingPvtDataInfo)

	// Commit block 4 with no pvtdata
	req.NoError(store.Prepare(4, nil, nil))
	req.NoError(store.Commit())

	// After committing block 4, the data for "ns-2:coll2" of block 1 should also have expired and should not be returned by the store
	expectedPvtdataFromBlock1 = []*ledger.TxPvtData{
		produceSamplePvtdata(t, 2, []string{"ns-1:coll-2", "ns-2:coll-1"}),
		produceSamplePvtdata(t, 4, []string{"ns-1:coll-2", "ns-2:coll-1"}),
	}
	retrievedData, _ = store.GetPvtDataByBlockNum(1, nil)
	req.Equal(expectedPvtdataFromBlock1, retrievedData)

	// Now, for block 2, "ns-1:coll1" should also have expired
	expectedPvtdataFromBlock2 := []*ledger.TxPvtData{
		produceSamplePvtdata(t, 3, []string{"ns-1:coll-2", "ns-2:coll-1", "ns-2:coll-2"}),
		produceSamplePvtdata(t, 5, []string{"ns-1:coll-2", "ns-2:coll-1", "ns-2:coll-2"}),
	}
	retrievedData, _ = store.GetPvtDataByBlockNum(2, nil)
	req.Equal(expectedPvtdataFromBlock2, retrievedData)

	// After committing block 4, the missing data of "ns1-coll1" in block2-tx1 should have expired
	expectedMissingPvtDataInfo = make(ledger.MissingPvtDataInfo)
	// missing data in block2, tx1
	expectedMissingPvtDataInfo.Add(2, 1, "ns-1", "coll-2")

	// missing data in block1, tx1
	expectedMissingPvtDataInfo.Add(1, 1, "ns-1", "coll-2")

	missingPvtDataInfo, err = store.GetMissingPvtDataInfoForMostRecentBlocks(10)
	req.NoError(err)
	req.Equal(expectedMissingPvtDataInfo, missingPvtDataInfo)

}

func TestLookupLastBlock(t *testing.T) {

	t.Run("test error from Unmarshal", func(t *testing.T) {
		_, _, err := lookupLastBlock(mockCouchDB{readDocValue: &couchdb.CouchDoc{JSONValue: []byte("wrongData")}})
		require.Error(t, err)
		require.Contains(t, err.Error(), "Unmarshal lastBlockResponse failed")
	})

	t.Run("test error from strconv ParseInt", func(t *testing.T) {
		jsonBinary, err := json.Marshal(lastCommittedBlockResponse{Data: "wrongData"})
		require.NoError(t, err)
		_, _, err = lookupLastBlock(mockCouchDB{readDocValue: &couchdb.CouchDoc{JSONValue: jsonBinary}})
		require.Error(t, err)
		require.Contains(t, err.Error(), "strconv.ParseInt lastBlockResponse.Data failed")
	})

	btlPolicy := btltestutil.SampleBTLPolicy(
		map[[2]string]uint64{
			{"ns-1", "coll-1"}: 0,
			{"ns-1", "coll-2"}: 0,
		},
	)
	env := NewTestStoreEnv(t, "teststorestate", btlPolicy, couchDBConfig)
	defer env.Cleanup("teststorestate")
	req := require.New(t)
	s := env.TestStore
	testData := []*ledger.TxPvtData{
		produceSamplePvtdata(t, 0, []string{"ns-1:coll-1", "ns-1:coll-2"}),
	}
	checkLastCommittedBlock(t, s, uint64(0))

	req.Nil(s.Prepare(0, nil, nil))
	req.NoError(s.Commit())
	checkLastCommittedBlock(t, s, uint64(0))

	req.Nil(s.Prepare(1, testData, nil))
	req.NoError(s.Commit())
	checkLastCommittedBlock(t, s, uint64(1))

	req.Nil(s.Prepare(2, nil, nil))
	req.NoError(s.Commit())
	checkLastCommittedBlock(t, s, uint64(2))

	req.Nil(s.Prepare(3, testData, nil))
	req.NoError(s.Commit())
	checkLastCommittedBlock(t, s, uint64(3))

	// Delete block num 2
	req.NoError(s.(*store).db.DeleteDoc(blockNumberToKey(2), ""))
	checkLastCommittedBlock(t, s, uint64(3))

}

func checkLastCommittedBlock(t *testing.T, s pvtdatastorage.Store, expectedLastCommittedBlock uint64) {
	lastCommitBlock, _, err := lookupLastBlock(s.(*store).db)
	require.NoError(t, err)
	require.Equal(t, expectedLastCommittedBlock, lastCommitBlock)
}

func TestStoreState(t *testing.T) {
	btlPolicy := btltestutil.SampleBTLPolicy(
		map[[2]string]uint64{
			{"ns-1", "coll-1"}: 0,
			{"ns-1", "coll-2"}: 0,
		},
	)
	env := NewTestStoreEnv(t, "teststorestate", btlPolicy, couchDBConfig)
	defer env.Cleanup("teststorestate")
	req := require.New(t)
	store := env.TestStore
	testData := []*ledger.TxPvtData{
		produceSamplePvtdata(t, 0, []string{"ns-1:coll-1", "ns-1:coll-2"}),
	}
	_, ok := store.Prepare(1, testData, nil).(*pvtdatastorage.ErrIllegalCall)
	req.True(ok)

	req.Nil(store.Prepare(0, testData, nil))
	req.NoError(store.Commit())

	req.Nil(store.Prepare(1, testData, nil))
	_, ok = store.Prepare(2, testData, nil).(*pvtdatastorage.ErrIllegalCall)
	req.True(ok)
}

func TestInitLastCommittedBlock(t *testing.T) {
	t.Run("test error from lookupLastBlock", func(t *testing.T) {
		ledgerId := "ledger"
		env := NewTestStoreEnv(t, ledgerId, nil, couchDBConfig)
		defer env.Cleanup(ledgerId)
		s := env.TestStore.(*store)
		s.db = mockCouchDB{readDocErr: fmt.Errorf("readDoc error")}
		s.isEmpty = true
		err := s.InitLastCommittedBlock(0)
		require.Error(t, err)
		require.Contains(t, err.Error(), "lookupLastBlock failed")
	})

	t.Run("test error from batchUpdateDocuments", func(t *testing.T) {
		ledgerId := "ledger"
		env := NewTestStoreEnv(t, ledgerId, nil, couchDBConfig)
		defer env.Cleanup(ledgerId)
		s := env.TestStore.(*store)
		s.db = mockCouchDB{batchUpdateDocumentsErr: fmt.Errorf("batchUpdateDocuments error")}
		s.isEmpty = true
		err := s.InitLastCommittedBlock(0)
		require.Error(t, err)
		require.Contains(t, err.Error(), "BatchUpdateDocuments failed")
	})

	env := NewTestStoreEnv(t, "teststorestate", nil, couchDBConfig)
	defer env.Cleanup("teststorestate")
	req := require.New(t)
	store := env.TestStore
	existingLastBlockNum := uint64(25)
	req.NoError(store.InitLastCommittedBlock(existingLastBlockNum))

	testEmpty(false, req, store)
	testPendingBatch(false, req, store)
	testLastCommittedBlockHeight(existingLastBlockNum+1, req, store)

	env.CloseAndReopen()
	testEmpty(false, req, store)
	testPendingBatch(false, req, store)
	testLastCommittedBlockHeight(existingLastBlockNum+1, req, store)

	err := store.InitLastCommittedBlock(30)
	_, ok := err.(*pvtdatastorage.ErrIllegalCall)
	req.True(ok)
}

func TestCollElgEnabled(t *testing.T) {
	testCollElgEnabled(t)
	defaultValBatchSize := xtestutil.TestLedgerConf().PrivateData.MaxBatchSize
	defaultValInterval := xtestutil.TestLedgerConf().PrivateData.BatchesInterval
	defer func() {
		viper.Set("ledger.pvtdataStore.collElgProcMaxDbBatchSize", defaultValBatchSize)
		viper.Set("ledger.pvtdataStore.collElgProcMaxDbBatchSize", defaultValInterval)
	}()
	viper.Set("ledger.pvtdataStore.collElgProcMaxDbBatchSize", 1)
	viper.Set("ledger.pvtdataStore.collElgProcDbBatchesInterval", 1)
	testCollElgEnabled(t)
}

func testCollElgEnabled(t *testing.T) {
	ledgerid := "testcollelgenabled"
	btlPolicy := btltestutil.SampleBTLPolicy(
		map[[2]string]uint64{
			{"ns-1", "coll-1"}: 0,
			{"ns-1", "coll-2"}: 0,
			{"ns-2", "coll-1"}: 0,
			{"ns-2", "coll-2"}: 0,
		},
	)
	env := NewTestStoreEnv(t, ledgerid, btlPolicy, couchDBConfig)
	defer env.Cleanup(ledgerid)
	req := require.New(t)
	store := env.TestStore

	// Initial state: eligible for {ns-1:coll-1 and ns-2:coll-1 }

	// no pvt data with block 0
	req.NoError(store.Prepare(0, nil, nil))
	req.NoError(store.Commit())

	// construct and commit block 1
	blk1MissingData := make(ledger.TxMissingPvtDataMap)
	blk1MissingData.Add(1, "ns-1", "coll-1", true)
	blk1MissingData.Add(1, "ns-2", "coll-1", true)
	blk1MissingData.Add(4, "ns-1", "coll-2", false)
	blk1MissingData.Add(4, "ns-2", "coll-2", false)
	testDataForBlk1 := []*ledger.TxPvtData{
		produceSamplePvtdata(t, 2, []string{"ns-1:coll-1"}),
	}
	req.NoError(store.Prepare(1, testDataForBlk1, blk1MissingData))
	req.NoError(store.Commit())

	// construct and commit block 2
	blk2MissingData := make(ledger.TxMissingPvtDataMap)
	// ineligible missing data in tx1
	blk2MissingData.Add(1, "ns-1", "coll-2", false)
	blk2MissingData.Add(1, "ns-2", "coll-2", false)
	testDataForBlk2 := []*ledger.TxPvtData{
		produceSamplePvtdata(t, 3, []string{"ns-1:coll-1"}),
	}
	req.NoError(store.Prepare(2, testDataForBlk2, blk2MissingData))
	req.NoError(store.Commit())

	// Retrieve and verify missing data reported
	// Expected missing data should be only blk1-tx1 (because, the other missing data is marked as ineliigible)
	expectedMissingPvtDataInfo := make(ledger.MissingPvtDataInfo)
	expectedMissingPvtDataInfo.Add(1, 1, "ns-1", "coll-1")
	expectedMissingPvtDataInfo.Add(1, 1, "ns-2", "coll-1")
	missingPvtDataInfo, err := store.GetMissingPvtDataInfoForMostRecentBlocks(10)
	req.NoError(err)
	req.Equal(expectedMissingPvtDataInfo, missingPvtDataInfo)

	// Enable eligibility for {ns-1:coll2}
	err = store.ProcessCollsEligibilityEnabled(
		5,
		map[string][]string{
			"ns-1": {"coll-2"},
		},
	)
	req.NoError(err)
	testutilWaitForCollElgProcToFinish(store)

	// Retrieve and verify missing data reported
	// Expected missing data should include newly eiligible collections
	expectedMissingPvtDataInfo.Add(1, 4, "ns-1", "coll-2")
	expectedMissingPvtDataInfo.Add(2, 1, "ns-1", "coll-2")
	missingPvtDataInfo, err = store.GetMissingPvtDataInfoForMostRecentBlocks(10)
	req.NoError(err)
	req.Equal(expectedMissingPvtDataInfo, missingPvtDataInfo)

	// Enable eligibility for {ns-2:coll2}
	err = store.ProcessCollsEligibilityEnabled(6,
		map[string][]string{
			"ns-2": {"coll-2"},
		},
	)
	req.NoError(err)
	testutilWaitForCollElgProcToFinish(store)

	// Retrieve and verify missing data reported
	// Expected missing data should include newly eiligible collections
	expectedMissingPvtDataInfo.Add(1, 4, "ns-2", "coll-2")
	expectedMissingPvtDataInfo.Add(2, 1, "ns-2", "coll-2")
	missingPvtDataInfo, err = store.GetMissingPvtDataInfoForMostRecentBlocks(10)
	req.Equal(expectedMissingPvtDataInfo, missingPvtDataInfo)
}

func TestRollBack(t *testing.T) {
	btlPolicy := btltestutil.SampleBTLPolicy(
		map[[2]string]uint64{
			{"ns-1", "coll-1"}: 0,
			{"ns-1", "coll-2"}: 0,
		},
	)
	env := NewTestStoreEnv(t, "testrollback", btlPolicy, couchDBConfig)
	defer env.Cleanup("testrollback")
	req := require.New(t)
	store := env.TestStore
	req.NoError(store.Prepare(0, nil, nil))
	req.NoError(store.Commit())

	pvtdata := []*ledger.TxPvtData{
		produceSamplePvtdata(t, 0, []string{"ns-1:coll-1", "ns-1:coll-2"}),
		produceSamplePvtdata(t, 5, []string{"ns-1:coll-1", "ns-1:coll-2"}),
	}
	missingData := make(ledger.TxMissingPvtDataMap)
	missingData.Add(1, "ns-1", "coll-1", true)
	missingData.Add(5, "ns-1", "coll-1", true)
	missingData.Add(5, "ns-2", "coll-2", false)

	for i := 1; i <= 9; i++ {
		req.NoError(store.Prepare(uint64(i), pvtdata, missingData))
		req.NoError(store.Commit())
	}

	datakeyTx0 := &common.DataKey{
		NsCollBlk: common.NsCollBlk{Ns: "ns-1", Coll: "coll-1"},
		TxNum:     0,
	}
	datakeyTx5 := &common.DataKey{
		NsCollBlk: common.NsCollBlk{Ns: "ns-1", Coll: "coll-1"},
		TxNum:     5,
	}
	eligibleMissingdatakey := &common.MissingDataKey{
		NsCollBlk:  common.NsCollBlk{Ns: "ns-1", Coll: "coll-1"},
		IsEligible: true,
	}

	// test store state before preparing for block 10
	testPendingBatch(false, req, store)
	testLastCommittedBlockHeight(10, req, store)

	// prepare for block 10 and test store for presence of datakeys and eligibile missingdatakeys
	req.NoError(store.Prepare(10, pvtdata, missingData))
	testPendingBatch(true, req, store)
	testLastCommittedBlockHeight(10, req, store)

	datakeyTx0.BlkNum = 10
	datakeyTx5.BlkNum = 10
	eligibleMissingdatakey.BlkNum = 10
	req.True(testPendingDataKeyExists(t, store, datakeyTx0))
	req.True(testPendingDataKeyExists(t, store, datakeyTx5))
	req.True(testPendingMissingDataKeyExists(t, store, eligibleMissingdatakey))

	// rollback last prepared block and test store for absence of datakeys and eligibile missingdatakeys
	err := store.Rollback()
	req.NoError(err)
	testPendingBatch(false, req, store)
	testLastCommittedBlockHeight(10, req, store)
	req.False(testPendingDataKeyExists(t, store, datakeyTx0))
	req.False(testPendingDataKeyExists(t, store, datakeyTx5))
	req.False(testPendingMissingDataKeyExists(t, store, eligibleMissingdatakey))

	// For previously committed blocks the datakeys and eligibile missingdatakeys should still be present
	for i := 1; i <= 9; i++ {
		datakeyTx0.BlkNum = uint64(i)
		datakeyTx5.BlkNum = uint64(i)
		eligibleMissingdatakey.BlkNum = uint64(i)
		req.True(testDataKeyExists(t, store, datakeyTx0))
		req.True(testDataKeyExists(t, store, datakeyTx5))
		req.True(testMissingDataKeyExists(t, store, eligibleMissingdatakey))
	}
}

func testMissingDataKeyExists(t *testing.T, s pvtdatastorage.Store, missingDataKey *common.MissingDataKey) bool {
	dataKeyBytes := common.EncodeMissingDataKey(missingDataKey)
	val, err := s.(*store).missingKeysIndexDB.Get(dataKeyBytes)
	require.NoError(t, err)
	return len(val) != 0
}

func testLastCommittedBlockHeight(expectedBlockHt uint64, req *require.Assertions, store pvtdatastorage.Store) {
	blkHt, err := store.LastCommittedBlockHeight()
	req.NoError(err)
	req.Equal(expectedBlockHt, blkHt)
}

func testDataKeyExists(t *testing.T, s pvtdatastorage.Store, dataKey *common.DataKey) bool {
	r, err := retrieveBlockPvtData(s.(*store).db, blockNumberToKey(dataKey.BlkNum))
	require.NoError(t, err)
	dataKeyBytes := common.EncodeDataKey(dataKey)
	_, exists := r.Data[hex.EncodeToString(dataKeyBytes)]
	return exists
}

func testPendingDataKeyExists(t *testing.T, s pvtdatastorage.Store, dataKey *common.DataKey) bool {
	var blockPvtData blockPvtDataResponse
	if s.(*store).pendingPvtData.PvtDataDoc == nil {
		return false
	}
	err := json.Unmarshal(s.(*store).pendingPvtData.PvtDataDoc.JSONValue, &blockPvtData)
	require.NoError(t, err)
	dataKeyBytes := common.EncodeDataKey(dataKey)
	_, exists := blockPvtData.Data[hex.EncodeToString(dataKeyBytes)]
	return exists
}

func testPendingMissingDataKeyExists(t *testing.T, s pvtdatastorage.Store, missingDataKey *common.MissingDataKey) bool {
	keyBytes := common.EncodeMissingDataKey(missingDataKey)
	_, exists := s.(*store).pendingPvtData.MissingDataEntries[string(keyBytes)]
	return exists
}

func testEmpty(expectedEmpty bool, req *require.Assertions, store pvtdatastorage.Store) {
	isEmpty, err := store.IsEmpty()
	req.NoError(err)
	req.Equal(expectedEmpty, isEmpty)
}

func testPendingBatch(expectedPending bool, req *require.Assertions, store pvtdatastorage.Store) {
	hasPendingBatch, err := store.HasPendingBatch()
	req.NoError(err)
	req.Equal(expectedPending, hasPendingBatch)
}

func produceSamplePvtdata(t *testing.T, txNum uint64, nsColls []string) *ledger.TxPvtData {
	builder := rwsetutil.NewRWSetBuilder()
	for _, nsColl := range nsColls {
		nsCollSplit := strings.Split(nsColl, ":")
		ns := nsCollSplit[0]
		coll := nsCollSplit[1]
		builder.AddToPvtAndHashedWriteSet(ns, coll, fmt.Sprintf("key-%s-%s", ns, coll), []byte(fmt.Sprintf("value-%s-%s", ns, coll)))
	}
	simRes, err := builder.GetTxSimulationResults()
	require.NoError(t, err)
	return &ledger.TxPvtData{SeqInBlock: txNum, WriteSet: simRes.PvtSimulationResults}
}

func testutilWaitForCollElgProcToFinish(s pvtdatastorage.Store) {
	s.(*store).collElgProc.WaitForDone()
}

// mockCouchDB
type mockCouchDB struct {
	existsWithRetryValue               bool
	existsWithRetryErr                 error
	indexDesignDocExistsWithRetryValue bool
	indexDesignDocExistsWithRetryErr   error
	createNewIndexWithRetryErr         error
	readDocValue                       *couchdb.CouchDoc
	readDocErr                         error
	batchUpdateDocumentsValue          []*couchdb.BatchUpdateResponse
	batchUpdateDocumentsErr            error
	queryDocumentsValue                []*couchdb.QueryResult
	queryDocumentsErr                  error
}

func (m mockCouchDB) ExistsWithRetry() (bool, error) {
	return m.existsWithRetryValue, m.existsWithRetryErr
}
func (m mockCouchDB) IndexDesignDocExistsWithRetry(designDocs ...string) (bool, error) {
	return m.indexDesignDocExistsWithRetryValue, m.indexDesignDocExistsWithRetryErr
}

func (m mockCouchDB) CreateNewIndexWithRetry(indexdefinition string, designDoc string) error {
	return m.createNewIndexWithRetryErr
}

func (m mockCouchDB) ReadDoc(id string) (*couchdb.CouchDoc, string, error) {
	return m.readDocValue, "", m.readDocErr
}

func (m mockCouchDB) BatchUpdateDocuments(documents []*couchdb.CouchDoc) ([]*couchdb.BatchUpdateResponse, error) {
	return m.batchUpdateDocumentsValue, m.batchUpdateDocumentsErr
}

func (m mockCouchDB) QueryDocuments(query string) ([]*couchdb.QueryResult, string, error) {
	return m.queryDocumentsValue, "", m.queryDocumentsErr
}
func (m mockCouchDB) DeleteDoc(id, rev string) error {
	return nil
}

type mockDBHandler struct {
	getFunc       func(key []byte) ([]byte, error)
	writeBatchErr error
	deleteErr     error
	putErr        error
}

func (m mockDBHandler) WriteBatch(batch *leveldbhelper.UpdateBatch, sync bool) error {
	return m.writeBatchErr
}
func (m mockDBHandler) Delete(key []byte, sync bool) error {
	return m.deleteErr
}
func (m mockDBHandler) Get(key []byte) ([]byte, error) {
	if m.getFunc != nil {
		return m.getFunc(key)
	}
	return nil, nil
}
func (m mockDBHandler) GetIterator(startKey []byte, endKey []byte) *leveldbhelper.Iterator {
	return nil
}
func (m mockDBHandler) Put(key []byte, value []byte, sync bool) error {
	return m.putErr
}
