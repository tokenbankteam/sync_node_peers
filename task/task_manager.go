package task

import (
	log "github.com/cihub/seelog"
	"github.com/tokenbankteam/sync_node_peers/service"
	"github.com/tokenbankteam/sync_node_peers/config"
)

type Manager struct {
	AppContext *service.AppContext
	Config     *config.AppConfig

	blockChainCheckTask *BlockChainCheckTask
}

func NewTaskManager(config *config.AppConfig, appContext *service.AppContext) (*Manager, error) {
	manager := &Manager{
		AppContext: appContext,
		Config:     config,
	}
	return manager, nil
}

func (s *Manager) init() {
}

//启动定时任务管理器
func (s *Manager) Start() {
	s.init()

	blockChainCheckTask, err := NewBlockChainCheckTask(s)
	if err != nil {
		log.Errorf("new blockChainCheckTask task error, %v", err)
		return
	}
	s.blockChainCheckTask = blockChainCheckTask
	s.blockChainCheckTask.Start()
}
