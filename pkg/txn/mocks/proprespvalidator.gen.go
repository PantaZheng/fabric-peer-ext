// Code generated by counterfeiter. DO NOT EDIT.
package mocks

import (
	"sync"

	pb "github.com/hyperledger/fabric-protos-go/peer"
	"github.com/trustbloc/fabric-peer-ext/pkg/txn/api"
)

type ProposalResponseValidator struct {
	ValidateStub        func(proposal *pb.SignedProposal, proposalResponses []*pb.ProposalResponse) (pb.TxValidationCode, error)
	validateMutex       sync.RWMutex
	validateArgsForCall []struct {
		proposal          *pb.SignedProposal
		proposalResponses []*pb.ProposalResponse
	}
	validateReturns struct {
		result1 pb.TxValidationCode
		result2 error
	}
	validateReturnsOnCall map[int]struct {
		result1 pb.TxValidationCode
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *ProposalResponseValidator) Validate(proposal *pb.SignedProposal, proposalResponses []*pb.ProposalResponse) (pb.TxValidationCode, error) {
	var proposalResponsesCopy []*pb.ProposalResponse
	if proposalResponses != nil {
		proposalResponsesCopy = make([]*pb.ProposalResponse, len(proposalResponses))
		copy(proposalResponsesCopy, proposalResponses)
	}
	fake.validateMutex.Lock()
	ret, specificReturn := fake.validateReturnsOnCall[len(fake.validateArgsForCall)]
	fake.validateArgsForCall = append(fake.validateArgsForCall, struct {
		proposal          *pb.SignedProposal
		proposalResponses []*pb.ProposalResponse
	}{proposal, proposalResponsesCopy})
	fake.recordInvocation("Validate", []interface{}{proposal, proposalResponsesCopy})
	fake.validateMutex.Unlock()
	if fake.ValidateStub != nil {
		return fake.ValidateStub(proposal, proposalResponses)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.validateReturns.result1, fake.validateReturns.result2
}

func (fake *ProposalResponseValidator) ValidateCallCount() int {
	fake.validateMutex.RLock()
	defer fake.validateMutex.RUnlock()
	return len(fake.validateArgsForCall)
}

func (fake *ProposalResponseValidator) ValidateArgsForCall(i int) (*pb.SignedProposal, []*pb.ProposalResponse) {
	fake.validateMutex.RLock()
	defer fake.validateMutex.RUnlock()
	return fake.validateArgsForCall[i].proposal, fake.validateArgsForCall[i].proposalResponses
}

func (fake *ProposalResponseValidator) ValidateReturns(result1 pb.TxValidationCode, result2 error) {
	fake.ValidateStub = nil
	fake.validateReturns = struct {
		result1 pb.TxValidationCode
		result2 error
	}{result1, result2}
}

func (fake *ProposalResponseValidator) ValidateReturnsOnCall(i int, result1 pb.TxValidationCode, result2 error) {
	fake.ValidateStub = nil
	if fake.validateReturnsOnCall == nil {
		fake.validateReturnsOnCall = make(map[int]struct {
			result1 pb.TxValidationCode
			result2 error
		})
	}
	fake.validateReturnsOnCall[i] = struct {
		result1 pb.TxValidationCode
		result2 error
	}{result1, result2}
}

func (fake *ProposalResponseValidator) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.validateMutex.RLock()
	defer fake.validateMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *ProposalResponseValidator) recordInvocation(key string, args []interface{}) {
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

var _ api.ProposalResponseValidator = new(ProposalResponseValidator)
