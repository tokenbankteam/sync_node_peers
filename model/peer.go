package model

type Network struct {
	LocalAddress  string `json:"localAddress"`
	RemoteAddress string `json:"remoteAddress"`
}

type Mc struct {
	Version    string `json:"version"`
	Difficulty string `json:"difficulty"`
	Head       string `json:"head"`
}

type Protocols struct {
	Mc Mc `json:"mc"`
}

type Peer struct {
	Id      string   `json:"id"`
	Name    string   `json:"name"`
	Caps    []string `json:"caps"`
	Network Network  `json:"network"`
}

type GetAdminPeersResult struct {
	ResultHeader
	Result []Peer `json:"result"`
}

type AddAdminPeerResult struct {
	ResultHeader
	Result bool `json:"result"`
}
