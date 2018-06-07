package db

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"sync"

	log "github.com/cihub/seelog"
)

type Database struct {
	Config  *DBConfig // 配置信息
	*sql.DB           // 数据库连接
}

// NewDatabase 创建数据库对象.
func NewDatabase(dbName string, config *DBConfig) (*Database, error) {
	instConfig := config.Instances[dbName]

	database, err := sql.Open(instConfig.Driver, instConfig.Url)
	if err != nil {
		log.Errorf("open database error %v", err)
		return nil, err
	}

	if err := database.Ping(); err != nil {
		log.Errorf("ping database error %v", err)
		return nil, err
	}
	return &Database{
		Config: config,
		DB:     database,
	}, nil
}

type Databases struct {
	config    *DBConfig
	instances map[string]*Database
	lock      sync.Mutex
}

//获取数据库连接池, map[string]*Database
func GetDatabases(config *DBConfig) (*Databases, error) {
	return &Databases{config: config, instances: make(map[string]*Database), lock: sync.Mutex{}}, nil
}

// TODO 出错需返回 error
func (ds *Databases) GetDatabase(dbName string) *Database {
	ds.lock.Lock()
	defer ds.lock.Unlock()
	if oldDb, ok := ds.instances[dbName]; !ok {
		database, err := NewDatabase(dbName, ds.config)
		if err != nil {
			log.Errorf("init database error %v", err)
			panic(err)
		}
		ds.instances[dbName] = database
		return database
	} else {
		return oldDb
	}
}
