/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package retriever

import (
	"context"
	"testing"
	"time"

	storeapi "github.com/hyperledger/fabric/extensions/collections/api/store"
	supportapi "github.com/hyperledger/fabric/extensions/collections/api/support"
	gossipapi "github.com/hyperledger/fabric/extensions/gossip/api"
	gcommon "github.com/hyperledger/fabric/gossip/common"
	cb "github.com/hyperledger/fabric/protos/common"
	gproto "github.com/hyperledger/fabric/protos/gossip"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	olapi "github.com/trustbloc/fabric-peer-ext/pkg/collections/offledger/api"
	"github.com/trustbloc/fabric-peer-ext/pkg/collections/offledger/dcas"
	spmocks "github.com/trustbloc/fabric-peer-ext/pkg/collections/storeprovider/mocks"
	"github.com/trustbloc/fabric-peer-ext/pkg/common/requestmgr"
	"github.com/trustbloc/fabric-peer-ext/pkg/mocks"
	ledgerconfig "github.com/trustbloc/fabric-peer-ext/pkg/roles"
)

const (
	channelID = "testchannel"
	ns1       = "chaincode1"
	coll1     = "collection1"
	coll2     = "collection2"

	org1MSPID = "Org1MSP"
	org2MSPID = "Org2MSP"
	org3MSPID = "Org3MSP"

	p1Org1Endpoint = "p1.org1.com"
	p2Org1Endpoint = "p2.org1.com"
	p3Org1Endpoint = "p3.org1.com"
	p1Org2Endpoint = "p1.org2.com"
	p2Org2Endpoint = "p2.org2.com"
	p3Org2Endpoint = "p3.org2.com"
	p1Org3Endpoint = "p1.org3.com"
	p2Org3Endpoint = "p2.org3.com"
	p3Org3Endpoint = "p3.org3.com"

	key1 = "key1"
	key2 = "key2"
	key3 = "key3"
	key4 = "key4"

	txID  = "tx1"
	txID1 = "tx2"
)

var (
	p1Org1PKIID = gcommon.PKIidType("pkiid_P1O1")
	p2Org1PKIID = gcommon.PKIidType("pkiid_P2O1")
	p3Org1PKIID = gcommon.PKIidType("pkiid_P3O1")

	p1Org2PKIID = gcommon.PKIidType("pkiid_P1O2")
	p2Org2PKIID = gcommon.PKIidType("pkiid_P2O2")
	p3Org2PKIID = gcommon.PKIidType("pkiid_P3O2")

	p1Org3PKIID = gcommon.PKIidType("pkiid_P1O3")
	p2Org3PKIID = gcommon.PKIidType("pkiid_P2O3")
	p3Org3PKIID = gcommon.PKIidType("pkiid_P3O3")

	committerRole = string(ledgerconfig.CommitterRole)
	endorserRole  = string(ledgerconfig.EndorserRole)

	respTimeout = 100 * time.Millisecond

	value1 = &storeapi.ExpiringValue{Value: []byte("value1")}
	value2 = &storeapi.ExpiringValue{Value: []byte("value2")}
	value3 = &storeapi.ExpiringValue{Value: []byte("value3")}
	value4 = &storeapi.ExpiringValue{Value: []byte("value4")}
)

func TestRetriever(t *testing.T) {
	support := mocks.NewMockSupport().
		CollectionPolicy(&mocks.MockAccessPolicy{
			MaxPeerCount: 2,
			Orgs:         []string{org1MSPID, org2MSPID, org3MSPID},
		}).
		CollectionConfig(&cb.StaticCollectionConfig{
			Type: cb.CollectionType_COL_OFFLEDGER,
			Name: coll1,
		}).
		CollectionConfig(&cb.StaticCollectionConfig{
			Type: cb.CollectionType_COL_DCAS,
			Name: coll2,
		})

	getLocalMSPIDRestore := getLocalMSPID
	getLocalMSPID = func() (string, error) { return org1MSPID, nil }
	defer func() {
		getLocalMSPID = getLocalMSPIDRestore
	}()

	gossip := mocks.NewMockGossipAdapter()
	gossip.Self(org1MSPID, mocks.NewMember(p1Org1Endpoint, p1Org1PKIID)).
		Member(org1MSPID, mocks.NewMember(p2Org1Endpoint, p2Org1PKIID, committerRole)).
		Member(org1MSPID, mocks.NewMember(p3Org1Endpoint, p3Org1PKIID, committerRole)).
		Member(org2MSPID, mocks.NewMember(p1Org2Endpoint, p1Org2PKIID, endorserRole)).
		Member(org2MSPID, mocks.NewMember(p2Org2Endpoint, p2Org2PKIID, committerRole)).
		Member(org2MSPID, mocks.NewMember(p3Org2Endpoint, p3Org2PKIID, endorserRole)).
		Member(org3MSPID, mocks.NewMember(p1Org3Endpoint, p1Org3PKIID, endorserRole)).
		Member(org3MSPID, mocks.NewMember(p2Org3Endpoint, p2Org3PKIID, committerRole)).
		Member(org3MSPID, mocks.NewMember(p3Org3Endpoint, p3Org3PKIID, endorserRole))

	casKey1, casValue1, err := dcas.GetCASKeyAndValueBase58([]byte(`{"id":"id1","value":"value1"}`))
	require.NoError(t, err)

	localStore := spmocks.NewStore().
		Data(storeapi.NewKey(txID, ns1, coll1, key1), value1).
		Data(storeapi.NewKey(txID, ns1, coll2, key1), value1).                                      // Invalid CAS key
		Data(storeapi.NewKey(txID, ns1, coll2, casKey1), &storeapi.ExpiringValue{Value: casValue1}) // Valid CAS key

	storeProvider := func(channelID string) olapi.Store { return localStore }
	gossipProvider := func() supportapi.GossipAdapter { return gossip }

	p := NewProvider(storeProvider, support, gossipProvider,
		WithValidator(cb.CollectionType_COL_DCAS, dcas.Validator),
		WithDecorator(cb.CollectionType_COL_DCAS, dcas.Decorator),
	)

	retriever := p.RetrieverForChannel(channelID)
	require.NotNil(t, retriever)

	t.Run("GetData - From local peer", func(t *testing.T) {
		ctx, _ := context.WithTimeout(context.Background(), respTimeout)
		value, err := retriever.GetData(ctx, storeapi.NewKey(txID, ns1, coll1, key1))
		require.NoError(t, err)
		require.NotNil(t, value)
		require.Equal(t, value1, value)
	})

	t.Run("GetData - From remote peer", func(t *testing.T) {
		gossip.MessageHandler(
			newMockGossipMsgHandler(channelID).
				Value(key2, value2).
				Handle)

		ctx, _ := context.WithTimeout(context.Background(), respTimeout)
		value, err := retriever.GetData(ctx, storeapi.NewKey(txID, ns1, coll1, key2))
		require.NoError(t, err)
		require.NotNil(t, value)
		require.Equal(t, value2, value)
	})

	t.Run("GetData (DCAS) - From local peer", func(t *testing.T) {
		ctx, _ := context.WithTimeout(context.Background(), respTimeout)
		value, err := retriever.GetData(ctx, storeapi.NewKey(txID, ns1, coll2, key1))
		require.Error(t, err)
		require.Contains(t, err.Error(), "the key should be the hash of the value")
		require.Nil(t, value)

		value, err = retriever.GetData(ctx, storeapi.NewKey(txID, ns1, coll2, casKey1))
		require.NoError(t, err)
		require.NotNil(t, value)
		require.Equal(t, casValue1, value.Value)
	})

	t.Run("GetData - No response from remote peer", func(t *testing.T) {
		gossip.MessageHandler(func(msg *gproto.GossipMessage) {})

		ctx, _ := context.WithTimeout(context.Background(), respTimeout)
		value, err := retriever.GetData(ctx, storeapi.NewKey(txID, ns1, coll1, key4))
		require.NoError(t, err)
		assert.Nil(t, value)
	})

	t.Run("GetData - Cancel request from remote peer", func(t *testing.T) {
		gossip.MessageHandler(func(msg *gproto.GossipMessage) {})

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		go func() {
			time.Sleep(50 * time.Millisecond)
			cancel()
		}()

		value, err := retriever.GetData(ctx, storeapi.NewKey(txID, ns1, coll1, key4))
		require.NoError(t, err)
		assert.Nil(t, value)
	})

	t.Run("GetDataMultipleKeys -> success", func(t *testing.T) {
		gossip.MessageHandler(
			newMockGossipMsgHandler(channelID).
				Value(key2, value2).
				Value(key3, value3).
				Handle)

		ctx, _ := context.WithTimeout(context.Background(), respTimeout)
		values, err := retriever.GetDataMultipleKeys(ctx, storeapi.NewMultiKey(txID, ns1, coll1, key1, key2, key3))
		require.NoError(t, err)
		require.Equal(t, 3, len(values))
		assert.Equal(t, value1, values[0])
		assert.Equal(t, value2, values[1])
		assert.Equal(t, value3, values[2])
	})

	t.Run("GetDataMultipleKeys - Key not found -> success", func(t *testing.T) {
		gossip.MessageHandler(
			newMockGossipMsgHandler(channelID).
				Value(key1, value1).
				Value(key2, value2).
				Value(key3, value3).
				Handle)

		ctx, _ := context.WithTimeout(context.Background(), respTimeout)
		values, err := retriever.GetDataMultipleKeys(ctx, storeapi.NewMultiKey(txID, ns1, coll1, "xxx", key2, key3))
		require.NoError(t, err)
		require.Equal(t, 3, len(values))
		assert.Nil(t, values[0])
		assert.Equal(t, value2, values[1])
		assert.Equal(t, value3, values[2])
	})

	t.Run("GetDataMultipleKeys - Timeout -> fail", func(t *testing.T) {
		gossip.MessageHandler(
			newMockGossipMsgHandler(channelID).
				Value(key1, value1).
				Value(key2, value2).
				Value(key3, value3).
				Handle)

		ctx, _ := context.WithTimeout(context.Background(), time.Microsecond)
		values, err := retriever.GetDataMultipleKeys(ctx, storeapi.NewMultiKey(txID, ns1, coll1, key1, key2, key3))
		assert.NoError(t, err)
		assert.Empty(t, values.Values().IsEmpty())
	})

	t.Run("GetDataMultipleKeys with store provider error -> fail", func(t *testing.T) {
		expectedErr := errors.New("store provider error")
		localStore.Error(expectedErr)
		defer localStore.Error(nil)

		gossip.MessageHandler(
			newMockGossipMsgHandler(channelID).
				Value(key2, value2).
				Value(key3, value3).
				Handle)

		ctx, _ := context.WithTimeout(context.Background(), respTimeout)
		values, err := retriever.GetDataMultipleKeys(ctx, storeapi.NewMultiKey(txID, ns1, coll1, key1, key2, key3))
		require.Error(t, err)
		require.Contains(t, err.Error(), expectedErr.Error())
		require.Nil(t, values)
	})

	t.Run("Persist Missing -> success", func(t *testing.T) {
		gossip.MessageHandler(
			newMockGossipMsgHandler(channelID).Handle)

		// Should not exist in local store
		ctx, _ := context.WithTimeout(context.Background(), respTimeout)
		value, err := retriever.GetData(ctx, storeapi.NewKey(txID, ns1, coll1, key4))
		require.NoError(t, err)
		require.Nil(t, value)

		gossip.MessageHandler(
			newMockGossipMsgHandler(channelID).
				Value(key4, value4).
				Handle)

		// Should retrieve from other peer and then persist to local store
		value, err = retriever.GetData(ctx, storeapi.NewKey(txID, ns1, coll1, key4))
		require.NoError(t, err)
		require.NotNil(t, value)

		gossip.MessageHandler(
			newMockGossipMsgHandler(channelID).Handle)

		// Should retrieve from local store
		value, err = retriever.GetData(ctx, storeapi.NewKey(txID, ns1, coll1, key4))
		require.NoError(t, err)
		require.NotNil(t, value)
	})

	t.Run("GetDataMultipleKeys -> success", func(t *testing.T) {
		gossip.MessageHandler(
			newMockGossipMsgHandler(channelID).
				Value(key2, value2).
				Value(key3, value3).
				Handle)

		ctx, _ := context.WithTimeout(context.Background(), respTimeout)
		values, err := retriever.GetDataMultipleKeys(ctx, storeapi.NewMultiKey(txID, ns1, coll1, key1, key2, key3))
		require.NoError(t, err)
		require.Equal(t, 3, len(values))
		assert.Equal(t, value1, values[0])
		assert.Equal(t, value2, values[1])
		assert.Equal(t, value3, values[2])
	})
}

func TestRetriever_Query(t *testing.T) {
	keyX, valueX, err := dcas.GetCASKeyAndValueBase58([]byte(`{"id":"id3","value":"valueX"}`))
	require.NoError(t, err)
	keyY, valueY, err := dcas.GetCASKeyAndValueBase58([]byte(`{"id":"id4","value":"valueY"}`))
	require.NoError(t, err)

	offLedgerQueryKey := storeapi.NewQueryKey(txID, ns1, coll1, "off-ledger query")
	offLedgerQueryResults := []*storeapi.QueryResult{
		{
			Key:           storeapi.NewKey(txID1, ns1, coll1, key1),
			ExpiringValue: value1,
		},
		{
			Key:           storeapi.NewKey(txID1, ns1, coll1, key2),
			ExpiringValue: value2,
		},
	}

	dcasQueryKey := storeapi.NewQueryKey(txID, ns1, coll2, "DCAS query")
	dcasQueryResults := []*storeapi.QueryResult{
		{
			Key:           storeapi.NewKey(txID1, ns1, coll2, keyX),
			ExpiringValue: &storeapi.ExpiringValue{Value: valueX},
		},
		{
			Key:           storeapi.NewKey(txID1, ns1, coll2, keyY),
			ExpiringValue: &storeapi.ExpiringValue{Value: valueY},
		},
	}

	support := mocks.NewMockSupport().
		CollectionPolicy(&mocks.MockAccessPolicy{
			MaxPeerCount: 2,
			Orgs:         []string{org1MSPID, org2MSPID},
		}).
		CollectionConfig(&cb.StaticCollectionConfig{
			Type: cb.CollectionType_COL_OFFLEDGER,
			Name: coll1,
		}).
		CollectionConfig(&cb.StaticCollectionConfig{
			Type: cb.CollectionType_COL_DCAS,
			Name: coll2,
		})

	getLocalMSPIDRestore := getLocalMSPID
	getLocalMSPID = func() (string, error) { return org1MSPID, nil }
	defer func() {
		getLocalMSPID = getLocalMSPIDRestore
	}()

	gossip := mocks.NewMockGossipAdapter()

	localStore := spmocks.NewStore().
		WithQueryResults(dcasQueryKey, dcasQueryResults).
		WithQueryResults(offLedgerQueryKey, offLedgerQueryResults)

	storeProvider := func(channelID string) olapi.Store { return localStore }
	gossipProvider := func() supportapi.GossipAdapter { return gossip }

	p := NewProvider(storeProvider, support, gossipProvider,
		WithValidator(cb.CollectionType_COL_DCAS, dcas.Validator),
		WithDecorator(cb.CollectionType_COL_DCAS, dcas.Decorator),
	)

	retriever := p.RetrieverForChannel(channelID)
	require.NotNil(t, retriever)

	t.Run("Query off-ledger -> success", func(t *testing.T) {
		ctx, _ := context.WithTimeout(context.Background(), respTimeout)
		it, err := retriever.Query(ctx, offLedgerQueryKey)
		require.NoError(t, err)
		require.NotNil(t, it)
		defer it.Close()

		next, err := it.Next()
		require.NoError(t, err)
		require.NotNil(t, next)
		require.Equal(t, value1.Value, next.Value)
		require.Equal(t, key1, next.Key.Key)

		next, err = it.Next()
		require.NoError(t, err)
		require.NotNil(t, next)
		require.Equal(t, value2.Value, next.Value)
		require.Equal(t, key2, next.Key.Key)

		next, err = it.Next()
		require.NoError(t, err)
		require.Nil(t, next)
	})

	t.Run("Query DCAS -> success", func(t *testing.T) {
		ctx, _ := context.WithTimeout(context.Background(), respTimeout)
		it, err := retriever.Query(ctx, dcasQueryKey)
		require.NoError(t, err)
		require.NotNil(t, it)
		defer it.Close()

		next, err := it.Next()
		require.NoError(t, err)
		require.NotNil(t, next)
		require.Equal(t, valueX, next.Value)
		require.Equal(t, keyX, next.Key.Key)

		next, err = it.Next()
		require.NoError(t, err)
		require.NotNil(t, next)
		require.Equal(t, valueY, next.Value)
		require.Equal(t, keyY, next.Key.Key)

		next, err = it.Next()
		require.NoError(t, err)
		require.Nil(t, next)
	})

	t.Run("Query iterator error", func(t *testing.T) {
		expectedErr := errors.New("injected error")
		localStore.ResultsIteratorError(expectedErr)
		ctx, _ := context.WithTimeout(context.Background(), respTimeout)
		it, err := retriever.Query(ctx, dcasQueryKey)
		require.NoError(t, err)
		require.NotNil(t, it)
		defer it.Close()

		next, err := it.Next()
		require.EqualError(t, err, expectedErr.Error())
		require.Nil(t, next)
	})

	t.Run("Query access denied -> empty", func(t *testing.T) {
		restore := getLocalMSPID
		getLocalMSPID = func() (string, error) { return org3MSPID, nil }
		defer func() {
			getLocalMSPID = restore
		}()

		ctx, _ := context.WithTimeout(context.Background(), respTimeout)
		it, err := retriever.Query(ctx, offLedgerQueryKey)
		require.NoError(t, err)
		require.NotNil(t, it)

		next, err := it.Next()
		require.NoError(t, err)
		require.Nil(t, next)
	})

	t.Run("Is authorized error -> fail", func(t *testing.T) {
		support.Err = errors.New("injected error")
		defer func() { support.Err = nil }()

		ctx, _ := context.WithTimeout(context.Background(), respTimeout)
		it, err := retriever.Query(ctx, offLedgerQueryKey)
		require.Error(t, err)
		require.Contains(t, err.Error(), support.Err.Error())
		require.Nil(t, it)
	})
}

func TestRetriever_AccessDenied(t *testing.T) {
	support := mocks.NewMockSupport().
		CollectionPolicy(&mocks.MockAccessPolicy{
			MaxPeerCount: 2,
			Orgs:         []string{org1MSPID, org2MSPID},
		}).
		CollectionConfig(&cb.StaticCollectionConfig{
			Type: cb.CollectionType_COL_OFFLEDGER,
			Name: coll1,
		})

	getLocalMSPIDRestore := getLocalMSPID
	getLocalMSPID = func() (string, error) { return org3MSPID, nil }
	defer func() {
		getLocalMSPID = getLocalMSPIDRestore
	}()

	gossip := mocks.NewMockGossipAdapter()
	gossip.Self(org3MSPID, mocks.NewMember(p1Org3Endpoint, p1Org3PKIID)).
		Member(org1MSPID, mocks.NewMember(p2Org1Endpoint, p2Org1PKIID, committerRole)).
		Member(org1MSPID, mocks.NewMember(p3Org1Endpoint, p3Org1PKIID, committerRole)).
		Member(org2MSPID, mocks.NewMember(p1Org2Endpoint, p1Org2PKIID, endorserRole)).
		Member(org2MSPID, mocks.NewMember(p2Org2Endpoint, p2Org2PKIID, committerRole)).
		Member(org2MSPID, mocks.NewMember(p3Org2Endpoint, p3Org2PKIID, endorserRole))

	localStore := spmocks.NewStore().
		Data(storeapi.NewKey(txID, ns1, coll1, key1), value1)

	storeProvider := func(channelID string) olapi.Store { return localStore }
	gossipProvider := func() supportapi.GossipAdapter { return gossip }

	p := NewProvider(storeProvider, support, gossipProvider)

	retriever := p.RetrieverForChannel(channelID)
	require.NotNil(t, retriever)

	gossip.MessageHandler(
		newMockGossipMsgHandler(channelID).
			Value(key2, value2).
			Handle)

	t.Run("GetData - From remote peer -> nil", func(t *testing.T) {
		ctx, _ := context.WithTimeout(context.Background(), respTimeout)
		value, err := retriever.GetData(ctx, storeapi.NewKey(txID, ns1, coll1, key2))
		assert.NoError(t, err)
		assert.Nil(t, value)
	})

	t.Run("GetDataMultipleKeys - From remote peer -> nil", func(t *testing.T) {
		ctx, _ := context.WithTimeout(context.Background(), respTimeout)
		values, err := retriever.GetDataMultipleKeys(ctx, storeapi.NewMultiKey(txID, ns1, coll1, "xxx", key2, key3))
		assert.NoError(t, err)
		assert.Nil(t, values)
	})

	t.Run("Missing from local", func(t *testing.T) {
		gossip.MessageHandler(
			newMockGossipMsgHandler(channelID).
				Value(key4, value4).
				Handle)

		// Shouldn't be authorized to retrieve from other peer
		ctx, _ := context.WithTimeout(context.Background(), respTimeout)
		value, err := retriever.GetData(ctx, storeapi.NewKey(txID, ns1, coll1, key4))
		require.NoError(t, err)
		require.Nil(t, value)
	})

	t.Run("GetData - From remote peer after CC upgrade -> success", func(t *testing.T) {
		gossip.MessageHandler(
			newMockGossipMsgHandler(channelID).
				Value(key2, value4).
				Handle)

		support.CollectionPolicy(&mocks.MockAccessPolicy{
			MaxPeerCount: 2,
			Orgs:         []string{org1MSPID, org2MSPID, org3MSPID},
		})

		require.NoError(t, support.Publisher.HandleUpgrade(gossipapi.TxMetadata{BlockNum: 1001, TxID: txID}, ns1))
		ctx, _ := context.WithTimeout(context.Background(), respTimeout)
		value, err := retriever.GetData(ctx, storeapi.NewKey(txID, ns1, coll1, key2))
		assert.NoError(t, err)
		assert.NotNil(t, value)
	})
}

func TestGetLocalMSPID(t *testing.T) {
	m, err := getLocalMSPID()
	require.NoError(t, err)
	require.NotNil(t, m)
}

type mockGossipMsgHandler struct {
	channelID string
	values    map[string]*storeapi.ExpiringValue
}

func newMockGossipMsgHandler(channelID string) *mockGossipMsgHandler {
	return &mockGossipMsgHandler{
		channelID: channelID,
		values:    make(map[string]*storeapi.ExpiringValue),
	}
}

func (m *mockGossipMsgHandler) Value(key string, value *storeapi.ExpiringValue) *mockGossipMsgHandler {
	m.values[key] = value
	return m
}

func (m *mockGossipMsgHandler) Handle(msg *gproto.GossipMessage) {
	req := msg.GetCollDataReq()

	res := &requestmgr.Response{
		Endpoint:  "p1.org1",
		MSPID:     "org1",
		Identity:  []byte("p1.org1"),
		Signature: []byte("sig"),
	}

	for _, d := range req.Digests {
		e := &requestmgr.Element{
			Namespace:  d.Namespace,
			Collection: d.Collection,
			Key:        d.Key,
		}

		v := m.values[d.Key]
		if v != nil {
			e.Value = v.Value
			e.Expiry = v.Expiry
		}

		res.Data = append(res.Data, e)
	}

	requestmgr.Get(m.channelID).Respond(req.Nonce, res)
}
