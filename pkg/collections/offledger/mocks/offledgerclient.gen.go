// Code generated by counterfeiter. DO NOT EDIT.
package mocks

import (
	"sync"

	commonledger "github.com/hyperledger/fabric/common/ledger"
	"github.com/trustbloc/fabric-peer-ext/pkg/collections/client"
)

type OffLedgerClient struct {
	PutStub        func(ns, coll, key string, value []byte) error
	putMutex       sync.RWMutex
	putArgsForCall []struct {
		ns    string
		coll  string
		key   string
		value []byte
	}
	putReturns struct {
		result1 error
	}
	putReturnsOnCall map[int]struct {
		result1 error
	}
	PutMultipleValuesStub        func(ns, coll string, kvs []*client.KeyValue) error
	putMultipleValuesMutex       sync.RWMutex
	putMultipleValuesArgsForCall []struct {
		ns   string
		coll string
		kvs  []*client.KeyValue
	}
	putMultipleValuesReturns struct {
		result1 error
	}
	putMultipleValuesReturnsOnCall map[int]struct {
		result1 error
	}
	DeleteStub        func(ns, coll string, keys ...string) error
	deleteMutex       sync.RWMutex
	deleteArgsForCall []struct {
		ns   string
		coll string
		keys []string
	}
	deleteReturns struct {
		result1 error
	}
	deleteReturnsOnCall map[int]struct {
		result1 error
	}
	GetStub        func(ns, coll, key string) ([]byte, error)
	getMutex       sync.RWMutex
	getArgsForCall []struct {
		ns   string
		coll string
		key  string
	}
	getReturns struct {
		result1 []byte
		result2 error
	}
	getReturnsOnCall map[int]struct {
		result1 []byte
		result2 error
	}
	GetMultipleKeysStub        func(ns, coll string, keys ...string) ([][]byte, error)
	getMultipleKeysMutex       sync.RWMutex
	getMultipleKeysArgsForCall []struct {
		ns   string
		coll string
		keys []string
	}
	getMultipleKeysReturns struct {
		result1 [][]byte
		result2 error
	}
	getMultipleKeysReturnsOnCall map[int]struct {
		result1 [][]byte
		result2 error
	}
	QueryStub        func(ns, coll, query string) (commonledger.ResultsIterator, error)
	queryMutex       sync.RWMutex
	queryArgsForCall []struct {
		ns    string
		coll  string
		query string
	}
	queryReturns struct {
		result1 commonledger.ResultsIterator
		result2 error
	}
	queryReturnsOnCall map[int]struct {
		result1 commonledger.ResultsIterator
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *OffLedgerClient) Put(ns string, coll string, key string, value []byte) error {
	var valueCopy []byte
	if value != nil {
		valueCopy = make([]byte, len(value))
		copy(valueCopy, value)
	}
	fake.putMutex.Lock()
	ret, specificReturn := fake.putReturnsOnCall[len(fake.putArgsForCall)]
	fake.putArgsForCall = append(fake.putArgsForCall, struct {
		ns    string
		coll  string
		key   string
		value []byte
	}{ns, coll, key, valueCopy})
	fake.recordInvocation("Put", []interface{}{ns, coll, key, valueCopy})
	fake.putMutex.Unlock()
	if fake.PutStub != nil {
		return fake.PutStub(ns, coll, key, value)
	}
	if specificReturn {
		return ret.result1
	}
	return fake.putReturns.result1
}

func (fake *OffLedgerClient) PutCallCount() int {
	fake.putMutex.RLock()
	defer fake.putMutex.RUnlock()
	return len(fake.putArgsForCall)
}

func (fake *OffLedgerClient) PutArgsForCall(i int) (string, string, string, []byte) {
	fake.putMutex.RLock()
	defer fake.putMutex.RUnlock()
	return fake.putArgsForCall[i].ns, fake.putArgsForCall[i].coll, fake.putArgsForCall[i].key, fake.putArgsForCall[i].value
}

func (fake *OffLedgerClient) PutReturns(result1 error) {
	fake.PutStub = nil
	fake.putReturns = struct {
		result1 error
	}{result1}
}

func (fake *OffLedgerClient) PutReturnsOnCall(i int, result1 error) {
	fake.PutStub = nil
	if fake.putReturnsOnCall == nil {
		fake.putReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.putReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *OffLedgerClient) PutMultipleValues(ns string, coll string, kvs []*client.KeyValue) error {
	var kvsCopy []*client.KeyValue
	if kvs != nil {
		kvsCopy = make([]*client.KeyValue, len(kvs))
		copy(kvsCopy, kvs)
	}
	fake.putMultipleValuesMutex.Lock()
	ret, specificReturn := fake.putMultipleValuesReturnsOnCall[len(fake.putMultipleValuesArgsForCall)]
	fake.putMultipleValuesArgsForCall = append(fake.putMultipleValuesArgsForCall, struct {
		ns   string
		coll string
		kvs  []*client.KeyValue
	}{ns, coll, kvsCopy})
	fake.recordInvocation("PutMultipleValues", []interface{}{ns, coll, kvsCopy})
	fake.putMultipleValuesMutex.Unlock()
	if fake.PutMultipleValuesStub != nil {
		return fake.PutMultipleValuesStub(ns, coll, kvs)
	}
	if specificReturn {
		return ret.result1
	}
	return fake.putMultipleValuesReturns.result1
}

func (fake *OffLedgerClient) PutMultipleValuesCallCount() int {
	fake.putMultipleValuesMutex.RLock()
	defer fake.putMultipleValuesMutex.RUnlock()
	return len(fake.putMultipleValuesArgsForCall)
}

func (fake *OffLedgerClient) PutMultipleValuesArgsForCall(i int) (string, string, []*client.KeyValue) {
	fake.putMultipleValuesMutex.RLock()
	defer fake.putMultipleValuesMutex.RUnlock()
	return fake.putMultipleValuesArgsForCall[i].ns, fake.putMultipleValuesArgsForCall[i].coll, fake.putMultipleValuesArgsForCall[i].kvs
}

func (fake *OffLedgerClient) PutMultipleValuesReturns(result1 error) {
	fake.PutMultipleValuesStub = nil
	fake.putMultipleValuesReturns = struct {
		result1 error
	}{result1}
}

func (fake *OffLedgerClient) PutMultipleValuesReturnsOnCall(i int, result1 error) {
	fake.PutMultipleValuesStub = nil
	if fake.putMultipleValuesReturnsOnCall == nil {
		fake.putMultipleValuesReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.putMultipleValuesReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *OffLedgerClient) Delete(ns string, coll string, keys ...string) error {
	fake.deleteMutex.Lock()
	ret, specificReturn := fake.deleteReturnsOnCall[len(fake.deleteArgsForCall)]
	fake.deleteArgsForCall = append(fake.deleteArgsForCall, struct {
		ns   string
		coll string
		keys []string
	}{ns, coll, keys})
	fake.recordInvocation("Delete", []interface{}{ns, coll, keys})
	fake.deleteMutex.Unlock()
	if fake.DeleteStub != nil {
		return fake.DeleteStub(ns, coll, keys...)
	}
	if specificReturn {
		return ret.result1
	}
	return fake.deleteReturns.result1
}

func (fake *OffLedgerClient) DeleteCallCount() int {
	fake.deleteMutex.RLock()
	defer fake.deleteMutex.RUnlock()
	return len(fake.deleteArgsForCall)
}

func (fake *OffLedgerClient) DeleteArgsForCall(i int) (string, string, []string) {
	fake.deleteMutex.RLock()
	defer fake.deleteMutex.RUnlock()
	return fake.deleteArgsForCall[i].ns, fake.deleteArgsForCall[i].coll, fake.deleteArgsForCall[i].keys
}

func (fake *OffLedgerClient) DeleteReturns(result1 error) {
	fake.DeleteStub = nil
	fake.deleteReturns = struct {
		result1 error
	}{result1}
}

func (fake *OffLedgerClient) DeleteReturnsOnCall(i int, result1 error) {
	fake.DeleteStub = nil
	if fake.deleteReturnsOnCall == nil {
		fake.deleteReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.deleteReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *OffLedgerClient) Get(ns string, coll string, key string) ([]byte, error) {
	fake.getMutex.Lock()
	ret, specificReturn := fake.getReturnsOnCall[len(fake.getArgsForCall)]
	fake.getArgsForCall = append(fake.getArgsForCall, struct {
		ns   string
		coll string
		key  string
	}{ns, coll, key})
	fake.recordInvocation("Get", []interface{}{ns, coll, key})
	fake.getMutex.Unlock()
	if fake.GetStub != nil {
		return fake.GetStub(ns, coll, key)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.getReturns.result1, fake.getReturns.result2
}

func (fake *OffLedgerClient) GetCallCount() int {
	fake.getMutex.RLock()
	defer fake.getMutex.RUnlock()
	return len(fake.getArgsForCall)
}

func (fake *OffLedgerClient) GetArgsForCall(i int) (string, string, string) {
	fake.getMutex.RLock()
	defer fake.getMutex.RUnlock()
	return fake.getArgsForCall[i].ns, fake.getArgsForCall[i].coll, fake.getArgsForCall[i].key
}

func (fake *OffLedgerClient) GetReturns(result1 []byte, result2 error) {
	fake.GetStub = nil
	fake.getReturns = struct {
		result1 []byte
		result2 error
	}{result1, result2}
}

func (fake *OffLedgerClient) GetReturnsOnCall(i int, result1 []byte, result2 error) {
	fake.GetStub = nil
	if fake.getReturnsOnCall == nil {
		fake.getReturnsOnCall = make(map[int]struct {
			result1 []byte
			result2 error
		})
	}
	fake.getReturnsOnCall[i] = struct {
		result1 []byte
		result2 error
	}{result1, result2}
}

func (fake *OffLedgerClient) GetMultipleKeys(ns string, coll string, keys ...string) ([][]byte, error) {
	fake.getMultipleKeysMutex.Lock()
	ret, specificReturn := fake.getMultipleKeysReturnsOnCall[len(fake.getMultipleKeysArgsForCall)]
	fake.getMultipleKeysArgsForCall = append(fake.getMultipleKeysArgsForCall, struct {
		ns   string
		coll string
		keys []string
	}{ns, coll, keys})
	fake.recordInvocation("GetMultipleKeys", []interface{}{ns, coll, keys})
	fake.getMultipleKeysMutex.Unlock()
	if fake.GetMultipleKeysStub != nil {
		return fake.GetMultipleKeysStub(ns, coll, keys...)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.getMultipleKeysReturns.result1, fake.getMultipleKeysReturns.result2
}

func (fake *OffLedgerClient) GetMultipleKeysCallCount() int {
	fake.getMultipleKeysMutex.RLock()
	defer fake.getMultipleKeysMutex.RUnlock()
	return len(fake.getMultipleKeysArgsForCall)
}

func (fake *OffLedgerClient) GetMultipleKeysArgsForCall(i int) (string, string, []string) {
	fake.getMultipleKeysMutex.RLock()
	defer fake.getMultipleKeysMutex.RUnlock()
	return fake.getMultipleKeysArgsForCall[i].ns, fake.getMultipleKeysArgsForCall[i].coll, fake.getMultipleKeysArgsForCall[i].keys
}

func (fake *OffLedgerClient) GetMultipleKeysReturns(result1 [][]byte, result2 error) {
	fake.GetMultipleKeysStub = nil
	fake.getMultipleKeysReturns = struct {
		result1 [][]byte
		result2 error
	}{result1, result2}
}

func (fake *OffLedgerClient) GetMultipleKeysReturnsOnCall(i int, result1 [][]byte, result2 error) {
	fake.GetMultipleKeysStub = nil
	if fake.getMultipleKeysReturnsOnCall == nil {
		fake.getMultipleKeysReturnsOnCall = make(map[int]struct {
			result1 [][]byte
			result2 error
		})
	}
	fake.getMultipleKeysReturnsOnCall[i] = struct {
		result1 [][]byte
		result2 error
	}{result1, result2}
}

func (fake *OffLedgerClient) Query(ns string, coll string, query string) (commonledger.ResultsIterator, error) {
	fake.queryMutex.Lock()
	ret, specificReturn := fake.queryReturnsOnCall[len(fake.queryArgsForCall)]
	fake.queryArgsForCall = append(fake.queryArgsForCall, struct {
		ns    string
		coll  string
		query string
	}{ns, coll, query})
	fake.recordInvocation("Query", []interface{}{ns, coll, query})
	fake.queryMutex.Unlock()
	if fake.QueryStub != nil {
		return fake.QueryStub(ns, coll, query)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.queryReturns.result1, fake.queryReturns.result2
}

func (fake *OffLedgerClient) QueryCallCount() int {
	fake.queryMutex.RLock()
	defer fake.queryMutex.RUnlock()
	return len(fake.queryArgsForCall)
}

func (fake *OffLedgerClient) QueryArgsForCall(i int) (string, string, string) {
	fake.queryMutex.RLock()
	defer fake.queryMutex.RUnlock()
	return fake.queryArgsForCall[i].ns, fake.queryArgsForCall[i].coll, fake.queryArgsForCall[i].query
}

func (fake *OffLedgerClient) QueryReturns(result1 commonledger.ResultsIterator, result2 error) {
	fake.QueryStub = nil
	fake.queryReturns = struct {
		result1 commonledger.ResultsIterator
		result2 error
	}{result1, result2}
}

func (fake *OffLedgerClient) QueryReturnsOnCall(i int, result1 commonledger.ResultsIterator, result2 error) {
	fake.QueryStub = nil
	if fake.queryReturnsOnCall == nil {
		fake.queryReturnsOnCall = make(map[int]struct {
			result1 commonledger.ResultsIterator
			result2 error
		})
	}
	fake.queryReturnsOnCall[i] = struct {
		result1 commonledger.ResultsIterator
		result2 error
	}{result1, result2}
}

func (fake *OffLedgerClient) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.putMutex.RLock()
	defer fake.putMutex.RUnlock()
	fake.putMultipleValuesMutex.RLock()
	defer fake.putMultipleValuesMutex.RUnlock()
	fake.deleteMutex.RLock()
	defer fake.deleteMutex.RUnlock()
	fake.getMutex.RLock()
	defer fake.getMutex.RUnlock()
	fake.getMultipleKeysMutex.RLock()
	defer fake.getMultipleKeysMutex.RUnlock()
	fake.queryMutex.RLock()
	defer fake.queryMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *OffLedgerClient) recordInvocation(key string, args []interface{}) {
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

var _ client.OffLedger = new(OffLedgerClient)