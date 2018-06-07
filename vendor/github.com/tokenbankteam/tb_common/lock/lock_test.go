package lock

import (
	"database/sql"
	log "github.com/cihub/seelog"
	"testing"
	"time"
)

func TestTime(t *testing.T) {
	t.Logf("%v must equal", time.Now().UnixNano()/1000000)
}

func TestLock(t *testing.T) {
	//init start
	database, err := sql.Open("mysql", "test:test@tcp(10.10.62.27:3306)/showapp?timeout=3s&strict=true&readTimeout=3s&writeTimeout=3s&parseTime=true")
	if err != nil {
		log.Errorf("open database error %v", err)
		return
	}

	if err := database.Ping(); err != nil {
		log.Errorf("ping database error %v", err)
		return
	}
	tableName := "t_resource_lock"
	retryCount := 3
	where := " lock_code='ready_room_lock'"
	resourceLockAccessor, err := NewDbResourceLockAccessor(database, tableName, retryCount, where)
	if err != nil {
		log.Errorf("new resource lock accessor error, %v", err)
		return
	}
	_, err = resourceLockAccessor.GetResourceLockByInterval(TEN_SECOND_INTERVAL)
	if err != nil {
		log.Errorf("get lock error, %v", err)
	}
}
