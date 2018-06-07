package model

type ResultHeader struct {
	JsonRpc string `json:"jsonrpc"`
	Id      int64  `json:"id"`
}
