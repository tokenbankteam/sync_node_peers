package controller

import (
	log "github.com/cihub/seelog"
	"github.com/gin-gonic/gin"
	"github.com/tokenbankteam/sync_node_peers/service"
)

type VersionController struct {
	BaseController
}

func NewVersionController(context *service.AppContext) (*VersionController, error) {
	return &VersionController{
	}, nil
}

func (s *VersionController) GetVersion(c *gin.Context) {
	version := c.DefaultQuery("version", "0.0.1")
	if version == "" {
		log.Errorf("version is empty")
		c.JSON(400, gin.H{
			"result":  1,
			"message": "version is empty",
		})
		return
	}
	c.JSON(200, gin.H{
		"result":  0,
		"message": "success",
		"data":    nil,
	})
}
