/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package config

import (
	"testing"
	"time"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestGetRoles(t *testing.T) {
	oldVal := viper.Get(confRoles)
	defer viper.Set(confRoles, oldVal)

	roles := "endorser,committer"
	viper.Set(confRoles, roles)
	assert.Equal(t, roles, GetRoles())
}

func TestGetPvtDataCacheSize(t *testing.T) {
	oldVal := viper.Get(confPvtDataCacheSize)
	defer viper.Set(confPvtDataCacheSize, oldVal)

	val := GetPvtDataCacheSize()
	assert.Equal(t, val, 10)

	viper.Set(confPvtDataCacheSize, 99)
	val = GetPvtDataCacheSize()
	assert.Equal(t, val, 99)

}

func TestGetTransientDataLevelDBPath(t *testing.T) {
	oldVal := viper.Get("peer.fileSystemPath")
	defer viper.Set("peer.fileSystemPath", oldVal)

	viper.Set("peer.fileSystemPath", "/tmp123")

	assert.Equal(t, "/tmp123/transientDataLeveldb", GetTransientDataLevelDBPath())
}

func TestGetTransientDataExpiredIntervalTime(t *testing.T) {
	oldVal := viper.Get(confTransientDataCleanupIntervalTime)
	defer viper.Set(confTransientDataCleanupIntervalTime, oldVal)

	viper.Set(confTransientDataCleanupIntervalTime, "")
	assert.Equal(t, defaultTransientDataCleanupIntervalTime, GetTransientDataExpiredIntervalTime())

	viper.Set(confTransientDataCleanupIntervalTime, 111*time.Second)
	assert.Equal(t, 111*time.Second, GetTransientDataExpiredIntervalTime())
}

func TestGetTransientDataCacheSize(t *testing.T) {
	oldVal := viper.Get(confTransientDataCacheSize)
	defer viper.Set(confTransientDataCacheSize, oldVal)

	viper.Set(confTransientDataCacheSize, 0)
	assert.Equal(t, defaultTransientDataCacheSize, GetTransientDataCacheSize())

	viper.Set(confTransientDataCacheSize, 10)
	assert.Equal(t, 10, GetTransientDataCacheSize())
}

func TestGetOLLevelDBPath(t *testing.T) {
	oldVal := viper.Get("peer.fileSystemPath")
	defer viper.Set("peer.fileSystemPath", oldVal)

	viper.Set("peer.fileSystemPath", "/tmp123")

	assert.Equal(t, "/tmp123/ledgersData/offLedgerLeveldb", GetOLCollLevelDBPath())
}

func TestGetOLCollExpiredIntervalTime(t *testing.T) {
	oldVal := viper.Get(confOLCollCleanupIntervalTime)
	defer viper.Set(confOLCollCleanupIntervalTime, oldVal)

	viper.Set(confOLCollCleanupIntervalTime, "")
	assert.Equal(t, defaultOLCollCleanupIntervalTime, GetOLCollExpirationCheckInterval())

	viper.Set(confOLCollCleanupIntervalTime, 111*time.Second)
	assert.Equal(t, 111*time.Second, GetOLCollExpirationCheckInterval())
}

func TestGetOLCacheSize(t *testing.T) {
	oldVal := viper.Get(confOLCollCacheSize)
	defer viper.Set(confOLCollCacheSize, oldVal)

	viper.Set(confOLCollCacheSize, 0)
	assert.Equal(t, defaultOLCollCacheSize, GetOLCollCacheSize())

	viper.Set(confOLCollCacheSize, 10)
	assert.Equal(t, 10, GetOLCollCacheSize())
}

func TestGetOLCacheEnabled(t *testing.T) {
	oldVal := viper.Get(confOLCollCacheEnabled)
	defer viper.Set(confOLCollCacheEnabled, oldVal)

	assert.False(t, GetOLCollCacheEnabled())

	viper.Set(confOLCollCacheEnabled, true)
	assert.True(t, GetOLCollCacheEnabled())

	viper.Set(confOLCollCacheEnabled, false)
	assert.False(t, GetOLCollCacheEnabled())
}

func TestGetTransientDataPullTimeout(t *testing.T) {
	oldVal := viper.Get(confTransientDataPullTimeout)
	defer viper.Set(confTransientDataPullTimeout, oldVal)

	viper.Set(confTransientDataPullTimeout, "")
	assert.Equal(t, defaultTransientDataPullTimeout, GetTransientDataPullTimeout())

	viper.Set(confTransientDataPullTimeout, 111*time.Second)
	assert.Equal(t, 111*time.Second, GetTransientDataPullTimeout())
}

func TestGetBlockPublisherBufferSize(t *testing.T) {
	oldVal := viper.Get(confBlockPublisherBufferSize)
	defer viper.Set(confBlockPublisherBufferSize, oldVal)

	viper.Set(confBlockPublisherBufferSize, "")
	assert.Equal(t, defaultBlockPublisherBufferSize, GetBlockPublisherBufferSize())

	viper.Set(confBlockPublisherBufferSize, 1234)
	assert.Equal(t, 1234, GetBlockPublisherBufferSize())
}

func TestGetOLCollMaxPeersForRetrieval(t *testing.T) {
	oldVal := viper.Get(confOLCollMaxPeersForRetrieval)
	defer viper.Set(confOLCollMaxPeersForRetrieval, oldVal)

	viper.Set(confOLCollMaxPeersForRetrieval, "")
	assert.Equal(t, defaultOLCollMaxPeersForRetrieval, GetOLCollMaxPeersForRetrieval())

	viper.Set(confOLCollMaxPeersForRetrieval, 7)
	assert.Equal(t, 7, GetOLCollMaxPeersForRetrieval())
}

func TestGetOLCollPullTimeout(t *testing.T) {
	oldVal := viper.Get(confOLCollPullTimeout)
	defer viper.Set(confOLCollPullTimeout, oldVal)

	viper.Set(confOLCollPullTimeout, "")
	assert.Equal(t, defaultOLCollPullTimeout, GetOLCollPullTimeout())

	viper.Set(confOLCollPullTimeout, 111*time.Second)
	assert.Equal(t, 111*time.Second, GetOLCollPullTimeout())
}

func TestGetConfigUpdatePublisherBufferSize(t *testing.T) {
	oldVal := viper.Get(confConfigUpdatePublisherBufferSize)
	defer viper.Set(confConfigUpdatePublisherBufferSize, oldVal)

	viper.Set(confConfigUpdatePublisherBufferSize, "")
	assert.Equal(t, defaultConfigUpdatePublisherBufferSize, GetConfigUpdatePublisherBufferSize())

	viper.Set(confConfigUpdatePublisherBufferSize, 1234)
	assert.Equal(t, 1234, GetConfigUpdatePublisherBufferSize())
}
