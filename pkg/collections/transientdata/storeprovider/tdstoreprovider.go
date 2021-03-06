/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package storeprovider

import (
	"sync"

	"github.com/pkg/errors"
	collcommon "github.com/trustbloc/fabric-peer-ext/pkg/collections/common"
	"github.com/trustbloc/fabric-peer-ext/pkg/collections/transientdata/api"
	"github.com/trustbloc/fabric-peer-ext/pkg/collections/transientdata/storeprovider/store/dbstore"
	"github.com/trustbloc/fabric-peer-ext/pkg/config"
)

// New returns a new transient data store provider
func New(gossipProvider collcommon.GossipProvider, idProvider collcommon.IdentityDeserializerProvider) *StoreProvider {
	logger.Infof("Creating new transient data store provider")
	dbp, err := dbstore.NewDBProvider()
	if err != nil {
		panic(err)
	}
	return &StoreProvider{
		stores:         make(map[string]*store),
		dbProvider:     dbp,
		gossipProvider: gossipProvider,
		idProvider:     idProvider,
	}
}

// StoreProvider is a transient data store provider
type StoreProvider struct {
	stores         map[string]*store
	dbProvider     *dbstore.LevelDBProvider
	gossipProvider collcommon.GossipProvider
	idProvider     collcommon.IdentityDeserializerProvider
	sync.RWMutex
}

// StoreForChannel returns the transient data store for the given channel
func (sp *StoreProvider) StoreForChannel(channelID string) api.Store {
	sp.RLock()
	defer sp.RUnlock()
	return sp.stores[channelID]
}

// OpenStore opens the transient data store for the given channel
func (sp *StoreProvider) OpenStore(channelID string) (api.Store, error) {
	sp.Lock()
	defer sp.Unlock()

	_, ok := sp.stores[channelID]
	if ok {
		return nil, errors.Errorf("a store for channel [%s] already exists", channelID)
	}

	db, err := sp.dbProvider.OpenDBStore(channelID)
	if err != nil {
		return nil, err
	}

	store := newStore(channelID, config.GetTransientDataCacheSize(), config.GetTransientDataAlwaysPersist(), db, sp.gossipProvider.GetGossipService(), sp.idProvider.GetIdentityDeserializer(channelID))
	sp.stores[channelID] = store

	return store, nil
}

// Close shuts down all of the stores
func (sp *StoreProvider) Close() {
	for _, s := range sp.stores {
		s.Close()
	}
	sp.dbProvider.Close()
}
