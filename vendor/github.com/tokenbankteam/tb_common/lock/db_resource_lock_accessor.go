package lock

import (
	"database/sql"
	log "github.com/cihub/seelog"
	"strings"
	"time"
)

// 基于数据库实现的通用行记录锁机制 需要在目标表上添加一个lock_time字段, 还必须有一个唯一键id
type DbResourceLockAccessor struct {
	//数据库访问对象
	dbTemplate        *sql.DB
	emptyResourceLock *ResourceLock
	retryCount        int
	tableName         string
	where             string
}

func NewDbResourceLockAccessor(dbTemplate *sql.DB, tableName string, retryCount int, where string) (ResourceLockAccessor, error) {
	if retryCount <= 0 {
		retryCount = 3
	}
	resourceLockAccessor := &DbResourceLockAccessor{
		dbTemplate:        dbTemplate,
		emptyResourceLock: nil,
		retryCount:        retryCount,
		tableName:         tableName,
		where:             where,
	}
	return resourceLockAccessor, nil
}

func (s *DbResourceLockAccessor) GetResourceLock(IntervalTime int64, lockTime int64) (*ResourceLock, error) {
	interval := &Interval{
		IntervalTime: IntervalTime,
		LockTime:     lockTime,
	}
	return s.GetResourceLockByInterval(interval)
}

func (s *DbResourceLockAccessor) GetResourceLockByInterval(interval *Interval) (*ResourceLock, error) {
	retryCount := s.retryCount
	for ; retryCount > 0; retryCount-- {
		resourceLock, err := s.getResourceLock(s.tableName, interval)
		if err != nil {
			log.Errorf("get resource lock error, %v", err)
			return nil, err
		}
		if resourceLock != nil {
			ret, err := s.update(resourceLock, interval.IntervalTime, interval.LockTime, interval.Unit)
			if err != nil {
				log.Errorf("update resource lock error, %v", err)
				return nil, err
			}
			if ret {
				return resourceLock, nil
			}
		}
	}
	return s.emptyResourceLock, nil
}

func (s *DbResourceLockAccessor) GetResourceLockList(IntervalTime int64, lockTime int64, max int) ([]*ResourceLock, error) {
	interval := &Interval{
		IntervalTime: IntervalTime,
		LockTime:     lockTime,
	}
	return s.GetResourceLockListByInterval(interval, max)
}

func (s *DbResourceLockAccessor) GetResourceLockListByInterval(interval *Interval, max int) ([]*ResourceLock, error) {
	resourceLockList, err := s.getResourceLockList(s.tableName, interval, max)
	if err != nil {
		return nil, err
	}
	if resourceLockList == nil || len(resourceLockList) == 0 {
		return nil, nil
	}
	result := []*ResourceLock{}
	//todo 添加批量更新操作如果resourceLockList不等于1
	for _, resourceLock := range resourceLockList {
		effect, err := s.update(resourceLock, interval.IntervalTime, interval.LockTime, interval.Unit)
		if err != nil {
			log.Errorf("update error, %v", err)
			return result, err
		}
		if effect {
			result = append(result, resourceLock)
		}
	}
	return result, nil
}

func (s *DbResourceLockAccessor) DeleteResourceLock(resourceLock *ResourceLock) (int64, error) {
	updateSql := "DELETE FROM " + s.tableName + " WHERE id=?"
	stmt, err := s.dbTemplate.Prepare(updateSql)
	if err != nil {
		log.Errorf("delete %v error, %v", updateSql, err)
		return -1, err
	}
	affect, err := stmt.Exec(resourceLock.Id)
	if err != nil {
		log.Errorf("exec delete %v error, %v", updateSql, err)
		return -1, err
	}
	effects, err := affect.RowsAffected()
	if err != nil {
		log.Errorf("delete %v rows affected error, %v", updateSql, err)
		return -1, err
	}
	return effects, nil
}

func (s *DbResourceLockAccessor) AddResourceLock(resourceLock *ResourceLock) (int64, error) {
	insertSql := s.getInsertSql()
	stmt, err := s.dbTemplate.Prepare(insertSql)
	if err != nil {
		log.Errorf("add %v error, %v", insertSql, err)
		return -1, err
	}
	affect, err := stmt.Exec(resourceLock.Id, resourceLock.LockTime)
	if err != nil {
		log.Errorf("exec add %v error, %v", insertSql, err)
		return -1, err
	}
	effects, err := affect.RowsAffected()
	if err != nil {
		log.Errorf("add %v rows affected error, %v", insertSql, err)
		return -1, err
	}
	return effects, nil
}

func (s *DbResourceLockAccessor) AddResourceLockList(resourceLockList []*ResourceLock) ([]int64, error) {
	//todo 添加批量操作
	result := []int64{}
	for _, resourceLock := range resourceLockList {
		effect, err := s.AddResourceLock(resourceLock)
		if err != nil {
			return result, err
		}
		result = append(result, effect)
	}
	return result, nil
}

func (s *DbResourceLockAccessor) getInsertSql() string {
	return "insert ignore into " + s.tableName + " (id, lock_time) values (?, ?)"
}

func (s *DbResourceLockAccessor) getResourceLock(tableName string, interval *Interval) (*ResourceLock, error) {
	resourceLockList, err := s.getResourceLockList(s.tableName, interval, 1)
	if err == nil && len(resourceLockList) > 0 {
		return resourceLockList[0], nil
	}
	return s.emptyResourceLock, err
}

//获取资源列表，
func (s *DbResourceLockAccessor) getResourceLockList(tableName string, interval *Interval, num int) ([]*ResourceLock, error) {
	querySql := "SELECT id, lock_time FROM " + tableName + " WHERE "
	if s.where != "" && strings.Trim(s.where, " ") != "" {
		querySql += strings.Trim(s.where, " ") + " AND "
	}
	querySql += "lock_time < ? ORDER BY id ASC LIMIT ?"
	lastLockTime := time.Now().UnixNano()/1000000 - (interval.IntervalTime-interval.LockTime)*interval.Unit + interval.Unit/2
	rows, err := s.dbTemplate.Query(querySql, lastLockTime, num)
	defer rows.Close()
	if err != nil {
		log.Errorf("get resource lock list from database error, %v", err)
		return nil, err
	}
	resourceLockList := []*ResourceLock{}
	for rows.Next() {
		resourceLock := &ResourceLock{}
		if err = rows.Scan(&resourceLock.Id, &resourceLock.LockTime); err != nil {
			log.Errorf("scan resource lock error: %v", err)
			return nil, err
		}
		resourceLockList = append(resourceLockList, resourceLock)
	}
	return resourceLockList, nil
}

func (s *DbResourceLockAccessor) update(resourceLock *ResourceLock, IntervalTime int64, lockTime int64, unit int64) (bool, error) {
	updateSql := "UPDATE " + s.tableName + " SET lock_time=? WHERE id=? AND lock_time=?"
	stmt, err := s.dbTemplate.Prepare(updateSql)
	if err != nil {
		log.Errorf("update %v error, %v", updateSql, err)
		return false, err
	}
	newLockTime := time.Now().UnixNano()/1000000 + lockTime*unit
	affect, err := stmt.Exec(newLockTime, resourceLock.Id, resourceLock.LockTime)
	if err != nil {
		log.Errorf("exec update %v error, %v", updateSql, err)
		return false, err
	}
	effects, err := affect.RowsAffected()
	if err != nil {
		log.Errorf("update %v rows affected error, %v", updateSql, err)
		return false, err
	}
	return effects > 0, nil
}
