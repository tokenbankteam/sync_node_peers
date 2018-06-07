package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/kataras/iris/core/errors"
	"strconv"
)

type BaseController struct {
}

func (s *BaseController) ParamInt64(c *gin.Context, key string) (int64, error) {
	paramValue := c.Param(key)
	if paramValue == "" {
		return -1, errors.New(key + " is empty")
	}
	return strconv.ParseInt(paramValue, 10, 64)
}

func (s *BaseController) DefaultQueryInt64(c *gin.Context, key string) (int64, error) {
	paramValue := c.DefaultQuery(key, "")
	if paramValue == "" {
		return -1, errors.New(key + " is empty")
	}
	return strconv.ParseInt(paramValue, 10, 64)
}

func (s *BaseController) DefaultQueryInt641(c *gin.Context, key string, value string) (int64, error) {
	paramValue := c.DefaultQuery(key, value)
	if paramValue == "" {
		return -1, errors.New(key + " is empty")
	}
	return strconv.ParseInt(paramValue, 10, 64)
}

func (s *BaseController) PostFormInt64(c *gin.Context, key string) (int64, error) {
	param := c.PostForm(key)
	if param == "" {
		return -1, errors.New(key + " is empty")
	}
	ret, err := strconv.ParseInt(param, 10, 64)
	if err != nil {
		return -1, err
	}
	return ret, nil
}
