package task

import (
	log "github.com/cihub/seelog"
	"github.com/tokenbankteam/sync_node_peers/service"
	"github.com/tokenbankteam/sync_node_peers/config"
	"time"
	"github.com/parnurzeal/gorequest"
	"encoding/json"
	"fmt"
)

type BlockChainCheckTask struct {
	BlockChainService *service.BlockChainService
	AppConfig         *config.AppConfig
	AppContext        *service.AppContext
	LastSendEmailTime time.Time
}

func NewBlockChainCheckTask(taskManager *Manager) (*BlockChainCheckTask, error) {
	context := taskManager.AppContext
	blockChainSyncTask := &BlockChainCheckTask{
		AppConfig:         context.Config,
		BlockChainService: context.Services["blockchainService"].(*service.BlockChainService),
	}
	return blockChainSyncTask, nil
}

func (s *BlockChainCheckTask) Start() {
	go func() {
		for {
			s.doCheck()
			time.Sleep(time.Minute * 1)
		}
	}()
}

func (s *BlockChainCheckTask) doCheck() {
	peerInfoList := []string{}
	nodeAddrList := s.AppConfig.EtherNodeAddrs
	for _, nodeAddr := range nodeAddrList {
		result, err := s.BlockChainService.GetAdminPeers(nodeAddr)
		if err != nil {
			log.Errorf("get admin peers of %v error %v", nodeAddr, err)
			continue
		}
		peerList := result.Result
		if peerList == nil || len(peerList) == 0 {
			s.sendSyncPeerEmail(nodeAddr, len(peerList))
		}
		for _, peer := range peerList {
			peerInfoList = append(peerInfoList, "enode://"+peer.Id+"@"+peer.Network.RemoteAddress)
		}
	}

	for _, nodeAddr := range nodeAddrList {
		for _, peer := range peerInfoList {
			ret, err := s.BlockChainService.AddAdminPeer(nodeAddr, peer)
			if err != nil {
				log.Errorf("add admin peer %v to %v error %v", peer, nodeAddr, err)
				break
			}
			bytes, _ := json.Marshal(ret)
			log.Infof("add peer %v to %v success %v", peer, nodeAddr, string(bytes))
		}
	}
}

func (s *BlockChainCheckTask) sendSyncPeerEmail(nodeAddr string, peerNum int) {
	to := "tokenpocket@163.com"
	subject := nodeAddr + "节点peer数" + fmt.Sprintf("%v", peerNum) + "报警"
	body := fmt.Sprintf("peer num report, node: %v, number: %v", nodeAddr, peerNum)
	log.Warn(body)
	if time.Now().Sub(s.LastSendEmailTime).Minutes() > float64(s.AppConfig.AlertPeriod) {
		_, err := s.sendEmail(to, subject, body)
		if err != nil {
			log.Errorf("send email to %v, subject %v, body %v error, %v", to, subject, body, err)
		}
		s.LastSendEmailTime = time.Now()
	}
}

func (s *BlockChainCheckTask) sendEmail(to, subject, body string) (int64, error) {
	addr := "http://emailinternal.mytokenpocket.vip/v1/email"
	resp, body, errs := gorequest.New().Timeout(time.Second * 10).Post(addr).Type("form").SendMap(
		map[string]string{
			"to":      to,
			"subject": subject,
			"body":    body,
		}).End()
	if errs != nil && len(errs) > 0 && errs[0] != nil {
		var err error
		for _, err = range errs {
			log.Errorf("request %v params to %v, subject %v error, %v", addr, to, subject, err)
		}
		return -1, err
	} else if resp.StatusCode != 200 {
		log.Errorf("request %v result invalid, params to %v, subject %v, status %v, body %v", addr, to, subject, resp.StatusCode, body)
		return -1, nil
	}
	defer resp.Body.Close()
	type Result struct {
		Result  int64  `json:"result"`
		Message string `json:"message"`
		Data    int64  `json:"data"`
	}
	ret := Result{}
	if err := json.Unmarshal([]byte(body), ret); err != nil {
		return -1, err
	}
	return ret.Data, nil
}
