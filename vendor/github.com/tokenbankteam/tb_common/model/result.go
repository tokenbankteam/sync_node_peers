package model

import "github.com/gin-gonic/gin"

type Result struct {
	Result  int64       `json:"result"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func (s *Result) ToGinHPointer() *gin.H {
	return &gin.H{
		"result":  s.Result,
		"message": s.Message,
		"data":    s.Data,
	}
}
