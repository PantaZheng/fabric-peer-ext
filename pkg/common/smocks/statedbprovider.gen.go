// Code generated by counterfeiter. DO NOT EDIT.
package smocks

import (
	"sync"

	extstatedb "github.com/trustbloc/fabric-peer-ext/pkg/statedb"
)

type StateDBProvider struct {
	StateDBForChannelStub        func(channelID string) extstatedb.StateDB
	stateDBForChannelMutex       sync.RWMutex
	stateDBForChannelArgsForCall []struct {
		channelID string
	}
	stateDBForChannelReturns struct {
		result1 extstatedb.StateDB
	}
	stateDBForChannelReturnsOnCall map[int]struct {
		result1 extstatedb.StateDB
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *StateDBProvider) StateDBForChannel(channelID string) extstatedb.StateDB {
	fake.stateDBForChannelMutex.Lock()
	ret, specificReturn := fake.stateDBForChannelReturnsOnCall[len(fake.stateDBForChannelArgsForCall)]
	fake.stateDBForChannelArgsForCall = append(fake.stateDBForChannelArgsForCall, struct {
		channelID string
	}{channelID})
	fake.recordInvocation("StateDBForChannel", []interface{}{channelID})
	fake.stateDBForChannelMutex.Unlock()
	if fake.StateDBForChannelStub != nil {
		return fake.StateDBForChannelStub(channelID)
	}
	if specificReturn {
		return ret.result1
	}
	return fake.stateDBForChannelReturns.result1
}

func (fake *StateDBProvider) StateDBForChannelCallCount() int {
	fake.stateDBForChannelMutex.RLock()
	defer fake.stateDBForChannelMutex.RUnlock()
	return len(fake.stateDBForChannelArgsForCall)
}

func (fake *StateDBProvider) StateDBForChannelArgsForCall(i int) string {
	fake.stateDBForChannelMutex.RLock()
	defer fake.stateDBForChannelMutex.RUnlock()
	return fake.stateDBForChannelArgsForCall[i].channelID
}

func (fake *StateDBProvider) StateDBForChannelReturns(result1 extstatedb.StateDB) {
	fake.StateDBForChannelStub = nil
	fake.stateDBForChannelReturns = struct {
		result1 extstatedb.StateDB
	}{result1}
}

func (fake *StateDBProvider) StateDBForChannelReturnsOnCall(i int, result1 extstatedb.StateDB) {
	fake.StateDBForChannelStub = nil
	if fake.stateDBForChannelReturnsOnCall == nil {
		fake.stateDBForChannelReturnsOnCall = make(map[int]struct {
			result1 extstatedb.StateDB
		})
	}
	fake.stateDBForChannelReturnsOnCall[i] = struct {
		result1 extstatedb.StateDB
	}{result1}
}

func (fake *StateDBProvider) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.stateDBForChannelMutex.RLock()
	defer fake.stateDBForChannelMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *StateDBProvider) recordInvocation(key string, args []interface{}) {
	fake.invocationsMutex.Lock()
	defer fake.invocationsMutex.Unlock()
	if fake.invocations == nil {
		fake.invocations = map[string][][]interface{}{}
	}
	if fake.invocations[key] == nil {
		fake.invocations[key] = [][]interface{}{}
	}
	fake.invocations[key] = append(fake.invocations[key], args)
}
