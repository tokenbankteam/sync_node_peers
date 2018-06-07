package config

import (
	cp "github.com/tokenbankteam/tb_common/cache"
	dbp "github.com/tokenbankteam/tb_common/db"
)

type AppConfig struct {
	Profile string `default:"dev"`

	Monitor      bool     `default:"true"`
	MonitorAddrs []string `default:"localhost:7272"`

	PprofAddrs []string `default:"localhost:7271"`

	Addr string `default:"0.0.0.0:8080"`

	DB dbp.DBConfig

	Redis cp.RedisConfig

	Logger string `default:"conf/logger.xml"`

	MongoAddr string `default:"mongodb://192.168.10.21:27017"`

	WorkerQueueSize int64  `default:"1000"`
	EtherEndpoint   string `default:"http://web3.tokenmaster.pro"`
	MonitorPeriod   int64  `default:"10"`
	ServerEndpoint  string `default:"http://tokenbank.skyfromwell.com"`

	EtherNodeAddrs []string `default:"localhost:8545"`
	DelayThreshold int64    `default:"10"`
	AlertPeriod    int64    `default:"5"`

	BloomFilterMaxN        uint64  `default:"100000"`
	BloomFilterProbCollide float64 `default:"0.0000001"`
}
