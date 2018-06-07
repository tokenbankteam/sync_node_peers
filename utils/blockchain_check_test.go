package utils

import (
	"testing"
	"github.com/tokenbankteam/sync_node_peers/model/blockchain"
	"encoding/json"
	"fmt"
)

func TestGetAdminPeers(t *testing.T) {
	blockChainModel, _ := blockchain.NewModel(nil)
	address := "https://chain3.mytokenpocket.vip"
	getAdminPeersResult, _ := blockChainModel.GetAdminPeers(address)
	marshal, _ := json.Marshal(getAdminPeersResult)
	fmt.Println(string(marshal))
	peer := getAdminPeersResult.Result[0]
	nodePeer := "enode://" + peer.Id + "@" + peer.Network.RemoteAddress
	ret, _ := blockChainModel.AddAdminPeer(address, nodePeer)
	bytes, _ := json.Marshal(ret)
	fmt.Println(bytes)
}
