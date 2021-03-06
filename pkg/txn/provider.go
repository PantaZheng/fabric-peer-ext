/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package txn

import (
	"github.com/bluele/gcache"
	"github.com/hyperledger/fabric/common/flogging"
	gossipapi "github.com/hyperledger/fabric/extensions/gossip/api"
	"github.com/trustbloc/fabric-peer-ext/pkg/config/ledgerconfig/config"
	"github.com/trustbloc/fabric-peer-ext/pkg/txn/api"
)

var logger = flogging.MustGetLogger("ext_txn")

// Provider is a transaction service provider
type Provider struct {
	services gcache.Cache
}

type configServiceProvider interface {
	ForChannel(channelID string) config.Service
}

type configValidatorRegistry interface {
	Register(v config.Validator)
}

type gossipProvider interface {
	GetGossipService() gossipapi.GossipService
}

type proposalResponseValidatorProvider interface {
	ValidatorForChannel(channelID string) api.ProposalResponseValidator
}

// NewProvider returns a new transaction service provider
func NewProvider(configProvider configServiceProvider, peerConfig api.PeerConfig, validatorRegistry configValidatorRegistry, gossipProvider gossipProvider, validatorProvider proposalResponseValidatorProvider) *Provider {
	validatorRegistry.Register(newConfigValidator())

	return newProvider(configProvider, peerConfig, gossipProvider, validatorProvider, &defaultClientProvider{})
}

func newProvider(configProvider configServiceProvider, peerConfig api.PeerConfig, gossipProvider gossipProvider, validatorProvider proposalResponseValidatorProvider, clientProvider clientProvider) *Provider {
	logger.Info("Creating transaction service provider")

	return &Provider{
		services: gcache.New(0).LoaderFunc(func(chID interface{}) (interface{}, error) {
			channelID := chID.(string)

			return newService(channelID,
				&providers{
					peerConfig:                peerConfig,
					configService:             configProvider.ForChannel(channelID),
					clientProvider:            clientProvider,
					gossip:                    gossipProvider.GetGossipService(),
					proposalResponseValidator: validatorProvider.ValidatorForChannel(channelID),
				})
		}).Build(),
	}
}

// ForChannel returns the transaction service for the given channel
func (p *Provider) ForChannel(channelID string) (api.Service, error) {
	svc, err := p.services.Get(channelID)
	if err != nil {
		return nil, err
	}

	return svc.(api.Service), nil
}

// Close closes all of the channel services
func (p *Provider) Close() {
	logger.Debug("Closing transaction services...")

	for _, channelID := range p.services.Keys() {
		svc, err := p.services.Get(channelID.(string))
		if err != nil {
			// This shouldn't happen since all of the services should already be cached
			logger.Warnf("Unable to close service for channel [%s]", channelID)
			continue
		}

		logger.Debugf("... closing service for channel [%s]", channelID)
		svc.(*Service).Close()
	}
}
