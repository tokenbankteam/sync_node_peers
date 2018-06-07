package blockchain

import (
	log "github.com/cihub/seelog"
	"fmt"
	"testing"
	"strings"
	"encoding/json"
	"net/http"
	"io/ioutil"
)

//var addr = "https://api.myetherapi.com/eth"
var addr = "http://web3.skyfromwell.com"

func TestTransfer(t *testing.T) {
	hex := fmt.Sprintf("0x%x", 4953967)
	marshal, _ := json.Marshal(NewReq("2.0", methodGetBlockByNumber, hex, true))
	resp, err := http.Post(addr, "application/json", strings.NewReader(string(marshal)))
	if err != nil {
		return
	}
	defer resp.Body.Close()
	bytes, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(bytes))
}
