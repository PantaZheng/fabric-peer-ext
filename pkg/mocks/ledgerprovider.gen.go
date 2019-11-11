// Code generated by counterfeiter. DO NOT EDIT.
package mocks

import (
	"sync"

	"github.com/hyperledger/fabric/core/ledger"
)

type LedgerProvider struct {
	GetLedgerStub        func(cid string) ledger.PeerLedger
	getLedgerMutex       sync.RWMutex
	getLedgerArgsForCall []struct {
		cid string
	}
	getLedgerReturns struct {
		result1 ledger.PeerLedger
	}
	getLedgerReturnsOnCall map[int]struct {
		result1 ledger.PeerLedger
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *LedgerProvider) GetLedger(cid string) ledger.PeerLedger {
	fake.getLedgerMutex.Lock()
	ret, specificReturn := fake.getLedgerReturnsOnCall[len(fake.getLedgerArgsForCall)]
	fake.getLedgerArgsForCall = append(fake.getLedgerArgsForCall, struct {
		cid string
	}{cid})
	fake.recordInvocation("GetLedger", []interface{}{cid})
	fake.getLedgerMutex.Unlock()
	if fake.GetLedgerStub != nil {
		return fake.GetLedgerStub(cid)
	}
	if specificReturn {
		return ret.result1
	}
	return fake.getLedgerReturns.result1
}

func (fake *LedgerProvider) GetLedgerCallCount() int {
	fake.getLedgerMutex.RLock()
	defer fake.getLedgerMutex.RUnlock()
	return len(fake.getLedgerArgsForCall)
}

func (fake *LedgerProvider) GetLedgerArgsForCall(i int) string {
	fake.getLedgerMutex.RLock()
	defer fake.getLedgerMutex.RUnlock()
	return fake.getLedgerArgsForCall[i].cid
}

func (fake *LedgerProvider) GetLedgerReturns(result1 ledger.PeerLedger) {
	fake.GetLedgerStub = nil
	fake.getLedgerReturns = struct {
		result1 ledger.PeerLedger
	}{result1}
}

func (fake *LedgerProvider) GetLedgerReturnsOnCall(i int, result1 ledger.PeerLedger) {
	fake.GetLedgerStub = nil
	if fake.getLedgerReturnsOnCall == nil {
		fake.getLedgerReturnsOnCall = make(map[int]struct {
			result1 ledger.PeerLedger
		})
	}
	fake.getLedgerReturnsOnCall[i] = struct {
		result1 ledger.PeerLedger
	}{result1}
}

func (fake *LedgerProvider) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.getLedgerMutex.RLock()
	defer fake.getLedgerMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *LedgerProvider) recordInvocation(key string, args []interface{}) {
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
