// Code generated by counterfeiter. DO NOT EDIT.
package mocks

import (
	"sync"
)

type DCASConfig struct {
	GetDCASMaxLinksPerBlockStub        func() int
	getDCASMaxLinksPerBlockMutex       sync.RWMutex
	getDCASMaxLinksPerBlockArgsForCall []struct{}
	getDCASMaxLinksPerBlockReturns     struct {
		result1 int
	}
	getDCASMaxLinksPerBlockReturnsOnCall map[int]struct {
		result1 int
	}
	IsDCASRawLeavesStub        func() bool
	isDCASRawLeavesMutex       sync.RWMutex
	isDCASRawLeavesArgsForCall []struct{}
	isDCASRawLeavesReturns     struct {
		result1 bool
	}
	isDCASRawLeavesReturnsOnCall map[int]struct {
		result1 bool
	}
	GetDCASMaxBlockSizeStub        func() int64
	getDCASMaxBlockSizeMutex       sync.RWMutex
	getDCASMaxBlockSizeArgsForCall []struct{}
	getDCASMaxBlockSizeReturns     struct {
		result1 int64
	}
	getDCASMaxBlockSizeReturnsOnCall map[int]struct {
		result1 int64
	}
	GetDCASBlockLayoutStub        func() string
	getDCASBlockLayoutMutex       sync.RWMutex
	getDCASBlockLayoutArgsForCall []struct{}
	getDCASBlockLayoutReturns     struct {
		result1 string
	}
	getDCASBlockLayoutReturnsOnCall map[int]struct {
		result1 string
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *DCASConfig) GetDCASMaxLinksPerBlock() int {
	fake.getDCASMaxLinksPerBlockMutex.Lock()
	ret, specificReturn := fake.getDCASMaxLinksPerBlockReturnsOnCall[len(fake.getDCASMaxLinksPerBlockArgsForCall)]
	fake.getDCASMaxLinksPerBlockArgsForCall = append(fake.getDCASMaxLinksPerBlockArgsForCall, struct{}{})
	fake.recordInvocation("GetDCASMaxLinksPerBlock", []interface{}{})
	fake.getDCASMaxLinksPerBlockMutex.Unlock()
	if fake.GetDCASMaxLinksPerBlockStub != nil {
		return fake.GetDCASMaxLinksPerBlockStub()
	}
	if specificReturn {
		return ret.result1
	}
	return fake.getDCASMaxLinksPerBlockReturns.result1
}

func (fake *DCASConfig) GetDCASMaxLinksPerBlockCallCount() int {
	fake.getDCASMaxLinksPerBlockMutex.RLock()
	defer fake.getDCASMaxLinksPerBlockMutex.RUnlock()
	return len(fake.getDCASMaxLinksPerBlockArgsForCall)
}

func (fake *DCASConfig) GetDCASMaxLinksPerBlockReturns(result1 int) {
	fake.GetDCASMaxLinksPerBlockStub = nil
	fake.getDCASMaxLinksPerBlockReturns = struct {
		result1 int
	}{result1}
}

func (fake *DCASConfig) GetDCASMaxLinksPerBlockReturnsOnCall(i int, result1 int) {
	fake.GetDCASMaxLinksPerBlockStub = nil
	if fake.getDCASMaxLinksPerBlockReturnsOnCall == nil {
		fake.getDCASMaxLinksPerBlockReturnsOnCall = make(map[int]struct {
			result1 int
		})
	}
	fake.getDCASMaxLinksPerBlockReturnsOnCall[i] = struct {
		result1 int
	}{result1}
}

func (fake *DCASConfig) IsDCASRawLeaves() bool {
	fake.isDCASRawLeavesMutex.Lock()
	ret, specificReturn := fake.isDCASRawLeavesReturnsOnCall[len(fake.isDCASRawLeavesArgsForCall)]
	fake.isDCASRawLeavesArgsForCall = append(fake.isDCASRawLeavesArgsForCall, struct{}{})
	fake.recordInvocation("IsDCASRawLeaves", []interface{}{})
	fake.isDCASRawLeavesMutex.Unlock()
	if fake.IsDCASRawLeavesStub != nil {
		return fake.IsDCASRawLeavesStub()
	}
	if specificReturn {
		return ret.result1
	}
	return fake.isDCASRawLeavesReturns.result1
}

func (fake *DCASConfig) IsDCASRawLeavesCallCount() int {
	fake.isDCASRawLeavesMutex.RLock()
	defer fake.isDCASRawLeavesMutex.RUnlock()
	return len(fake.isDCASRawLeavesArgsForCall)
}

func (fake *DCASConfig) IsDCASRawLeavesReturns(result1 bool) {
	fake.IsDCASRawLeavesStub = nil
	fake.isDCASRawLeavesReturns = struct {
		result1 bool
	}{result1}
}

func (fake *DCASConfig) IsDCASRawLeavesReturnsOnCall(i int, result1 bool) {
	fake.IsDCASRawLeavesStub = nil
	if fake.isDCASRawLeavesReturnsOnCall == nil {
		fake.isDCASRawLeavesReturnsOnCall = make(map[int]struct {
			result1 bool
		})
	}
	fake.isDCASRawLeavesReturnsOnCall[i] = struct {
		result1 bool
	}{result1}
}

func (fake *DCASConfig) GetDCASMaxBlockSize() int64 {
	fake.getDCASMaxBlockSizeMutex.Lock()
	ret, specificReturn := fake.getDCASMaxBlockSizeReturnsOnCall[len(fake.getDCASMaxBlockSizeArgsForCall)]
	fake.getDCASMaxBlockSizeArgsForCall = append(fake.getDCASMaxBlockSizeArgsForCall, struct{}{})
	fake.recordInvocation("GetDCASMaxBlockSize", []interface{}{})
	fake.getDCASMaxBlockSizeMutex.Unlock()
	if fake.GetDCASMaxBlockSizeStub != nil {
		return fake.GetDCASMaxBlockSizeStub()
	}
	if specificReturn {
		return ret.result1
	}
	return fake.getDCASMaxBlockSizeReturns.result1
}

func (fake *DCASConfig) GetDCASMaxBlockSizeCallCount() int {
	fake.getDCASMaxBlockSizeMutex.RLock()
	defer fake.getDCASMaxBlockSizeMutex.RUnlock()
	return len(fake.getDCASMaxBlockSizeArgsForCall)
}

func (fake *DCASConfig) GetDCASMaxBlockSizeReturns(result1 int64) {
	fake.GetDCASMaxBlockSizeStub = nil
	fake.getDCASMaxBlockSizeReturns = struct {
		result1 int64
	}{result1}
}

func (fake *DCASConfig) GetDCASMaxBlockSizeReturnsOnCall(i int, result1 int64) {
	fake.GetDCASMaxBlockSizeStub = nil
	if fake.getDCASMaxBlockSizeReturnsOnCall == nil {
		fake.getDCASMaxBlockSizeReturnsOnCall = make(map[int]struct {
			result1 int64
		})
	}
	fake.getDCASMaxBlockSizeReturnsOnCall[i] = struct {
		result1 int64
	}{result1}
}

func (fake *DCASConfig) GetDCASBlockLayout() string {
	fake.getDCASBlockLayoutMutex.Lock()
	ret, specificReturn := fake.getDCASBlockLayoutReturnsOnCall[len(fake.getDCASBlockLayoutArgsForCall)]
	fake.getDCASBlockLayoutArgsForCall = append(fake.getDCASBlockLayoutArgsForCall, struct{}{})
	fake.recordInvocation("GetDCASBlockLayout", []interface{}{})
	fake.getDCASBlockLayoutMutex.Unlock()
	if fake.GetDCASBlockLayoutStub != nil {
		return fake.GetDCASBlockLayoutStub()
	}
	if specificReturn {
		return ret.result1
	}
	return fake.getDCASBlockLayoutReturns.result1
}

func (fake *DCASConfig) GetDCASBlockLayoutCallCount() int {
	fake.getDCASBlockLayoutMutex.RLock()
	defer fake.getDCASBlockLayoutMutex.RUnlock()
	return len(fake.getDCASBlockLayoutArgsForCall)
}

func (fake *DCASConfig) GetDCASBlockLayoutReturns(result1 string) {
	fake.GetDCASBlockLayoutStub = nil
	fake.getDCASBlockLayoutReturns = struct {
		result1 string
	}{result1}
}

func (fake *DCASConfig) GetDCASBlockLayoutReturnsOnCall(i int, result1 string) {
	fake.GetDCASBlockLayoutStub = nil
	if fake.getDCASBlockLayoutReturnsOnCall == nil {
		fake.getDCASBlockLayoutReturnsOnCall = make(map[int]struct {
			result1 string
		})
	}
	fake.getDCASBlockLayoutReturnsOnCall[i] = struct {
		result1 string
	}{result1}
}

func (fake *DCASConfig) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.getDCASMaxLinksPerBlockMutex.RLock()
	defer fake.getDCASMaxLinksPerBlockMutex.RUnlock()
	fake.isDCASRawLeavesMutex.RLock()
	defer fake.isDCASRawLeavesMutex.RUnlock()
	fake.getDCASMaxBlockSizeMutex.RLock()
	defer fake.getDCASMaxBlockSizeMutex.RUnlock()
	fake.getDCASBlockLayoutMutex.RLock()
	defer fake.getDCASBlockLayoutMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *DCASConfig) recordInvocation(key string, args []interface{}) {
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
