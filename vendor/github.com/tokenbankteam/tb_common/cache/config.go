package cache

type RedisInstConfig struct {
	Url      string
	Password string
}

type RedisConfig struct {
	Instances map[string]RedisInstConfig
}
