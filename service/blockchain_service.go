package service

import (
	"github.com/tokenbankteam/sync_node_peers/model/blockchain"
	"time"
	log "github.com/cihub/seelog"
	"github.com/tokenbankteam/sync_node_peers/model"
)

type BlockChainService struct {
	BlockChainModel *blockchain.Model
}

func NewBlockChainService(context *AppContext) (*BlockChainService, error) {
	return &BlockChainService{
		BlockChainModel: context.Models["blockchainModel"].(*blockchain.Model),
	}, nil
}

func (s *BlockChainService) GetAdminPeers(addr string) (*model.GetAdminPeersResult, error) {
	start := time.Now().Nanosecond()
	defer func() {
		elapsed := (time.Now().Nanosecond() - start) / (1000 * 1000)
		if elapsed > 500 {
			log.Infof("GetAdminPeers elapsed %vms", elapsed)
		}
	}()
	return s.BlockChainModel.GetAdminPeers(addr)
}

func (s *BlockChainService) AddAdminPeer(addr string, peer string) (*model.AddAdminPeerResult, error) {
	start := time.Now().Nanosecond()
	defer func() {
		elapsed := (time.Now().Nanosecond() - start) / (1000 * 1000)
		if elapsed > 500 {
			log.Infof("AddAdminPeers elapsed %vms", elapsed)
		}
	}()
	return s.BlockChainModel.AddAdminPeer(addr, peer)
}
