package db

type DbInstConfig struct {
	Driver string
	Url    string
}

type DBConfig struct {
	Instances map[string]DbInstConfig
}
