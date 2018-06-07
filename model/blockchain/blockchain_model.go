package blockchain

import (
	log "github.com/cihub/seelog"
	"github.com/tokenbankteam/sync_node_peers/config"
	"encoding/json"
	"time"
	"github.com/parnurzeal/gorequest"
	"github.com/tokenbankteam/sync_node_peers/model"
)

type Model struct {
	AppConfig *config.AppConfig
}

const (
	methodAdminPeers   = "admin_peers"
	methodAddAdminPeer = "admin_addPeer"
)

var timeoutDuration = time.Second * 10
var id int

type Request struct {
	JsonRpc string        `json:"jsonrpc"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
	Id      int           `json:"id"`
}

// cache request
func NewReq(jsonRpc string, method string, param ...interface{}) *Request {
	r := new(Request)
	r.JsonRpc = jsonRpc
	r.Method = method
	r.Params = append([]interface{}{}, param...)
	id++
	r.Id = id
	return r
}

func NewModel(config *config.AppConfig) (*Model, error) {
	return &Model{
		AppConfig: config,
	}, nil
}

func (s *Model) GetAdminPeers(address string) (*model.GetAdminPeersResult, error) {
	marshal, _ := json.Marshal(NewReq("2.0", methodAdminPeers))
	resp, body, errs := gorequest.New().Timeout(timeoutDuration).Post(address).Send(string(marshal)).End()
	if errs != nil && len(errs) > 0 && errs[0] != nil {
		var err error
		for _, err = range errs {
			log.Errorf("request %v params %v error, %v", address, string(marshal), err)
		}
		return nil, err
	} else if resp.StatusCode != 200 {
		log.Errorf("request %v result invalid, params %v, status %v, body %v", address, string(marshal), resp.StatusCode, body)
		return nil, nil
	}
	defer resp.Body.Close()
	var ret = &model.GetAdminPeersResult{}
	err := json.Unmarshal([]byte(body), ret)
	return ret, err
}

func (s *Model) AddAdminPeer(address string, peer string) (*model.AddAdminPeerResult, error) {
	marshal, _ := json.Marshal(NewReq("2.0", methodAddAdminPeer, peer))
	resp, body, errs := gorequest.New().Timeout(timeoutDuration).Post(address).Send(string(marshal)).End()
	if errs != nil && len(errs) > 0 && errs[0] != nil {
		var err error
		for _, err = range errs {
			log.Errorf("request %v params %v error, %v", address, string(marshal), err)
		}
		return nil, err
	} else if resp.StatusCode != 200 {
		log.Errorf("request %v result invalid, params %v, status %v, body %v", address, string(marshal), resp.StatusCode, body)
		return nil, nil
	}
	defer resp.Body.Close()
	var ret = &model.AddAdminPeerResult{}
	err := json.Unmarshal([]byte(body), ret)
	return ret, err
}
