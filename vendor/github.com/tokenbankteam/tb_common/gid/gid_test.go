package gid

import (
	"testing"
	"fmt"
)

var S = &Server{
	UrlPrefix: "http://sonyflake.live.xunlei.com/",
}

func TestGetId(t *testing.T) {
	result, err := S.GetId()
	if err != nil {
		t.Fatalf("get userRoomStat error, %v", err)
	}
	fmt.Println(result)
}
