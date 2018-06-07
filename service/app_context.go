package service

import "github.com/tokenbankteam/sync_node_peers/config"

type AppContext struct {
	Config *config.AppConfig

	Services map[string]interface{}
	Models   map[string]interface{}
}
