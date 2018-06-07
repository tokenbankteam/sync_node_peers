package cache

import (
	log "github.com/cihub/seelog"
	"github.com/go-redis/redis"
	"sync"
)

type Cache struct {
	Config *RedisConfig
	*redis.Client
}

func NewCache(dbName string, config *RedisConfig) (*Cache, error) {
	instConfig := config.Instances[dbName]
	client := redis.NewClient(&redis.Options{
		Addr:     instConfig.Url,
		Password: instConfig.Password, // no password set
		DB:       0,                   // use default DB
	})

	if pong, err := client.Ping().Result(); err != nil {
		log.Errorf("ping cache error %v", err)
		return nil, err
	} else {
		log.Infof("ping cache %v", pong)
	}

	return &Cache{
		Config: config,
		Client: client,
	}, nil
}

type Caches struct {
	config    *RedisConfig
	instances map[string]*Cache
	lock      sync.Mutex
}

//获取数据库连接池, map[string]*Caches
func GetCaches(config *RedisConfig) (*Caches, error) {
	return &Caches{config: config, instances: make(map[string]*Cache), lock: sync.Mutex{}}, nil
}

// TODO 出错需返回 error
func (ds *Caches) GetCache(dbName string) *Cache {
	ds.lock.Lock()
	defer ds.lock.Unlock()
	if oldDb, ok := ds.instances[dbName]; !ok {
		database, err := NewCache(dbName, ds.config)
		if err != nil {
			log.Errorf("init cache %v error: %v", dbName, err)
			panic("init database error")
		}
		ds.instances[dbName] = database
		return database
	} else {
		return oldDb
	}
}
