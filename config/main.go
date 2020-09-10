package config

import (
	"fmt"
	"runtime"

	log "github.com/sirupsen/logrus"

	_ "net/http/pprof"

	"github.com/spf13/viper"

	_ "github.com/ClickHouse/clickhouse-go"
	"github.com/go-redis/redis/v8"
	"github.com/jmoiron/sqlx"
)

// TableRule struct incapsulates a specific table config - it's title ( fully qualified ), slice of key column and version column name
type TableRule struct {
	Title   string   `json:"title"`
	Key     []string `json:"key"`
	Version string   `json:"version"`
}

type Config struct {
	TestFlag      bool   `json:"TestFlag"`
	ListenPort    int    `json:"ListenPort"`
	Clickhouse    string `json:"Clickhouse"`
	Redis         string `json:"Redis"`
	FlushInterval int    `json:"FlushInterval"`
	FlushCount    int    `json:"FlushCount"`
	TableRules    struct {
		Main      TableRule   `json:"main"`
		Secondary []TableRule `json:"secondary"`
	} `json:"TableRules"`
}

type Env struct {
	Db       *sqlx.DB
	Redis    *redis.Client
	Config   *Config
	loglevel log.Level
}

func NewEnv(path string) *Env {
	viper.SetConfigType("json")
	viper.SetConfigName(path)
	viper.AddConfigPath("config")

	if err := viper.ReadInConfig(); err != nil {
		checkErr(err)
	}

	var cfg Config
	err := viper.Unmarshal(&cfg)
	if err != nil {
		checkErr(err)
	}

	loglevel := log.WarnLevel

	log.Printf("%v\n", cfg)

	return &Env{
		Config:   &cfg,
		loglevel: loglevel,
	}
}

func (e *Env) Initredis() {
	e.Redis = redis.NewClient(&redis.Options{
		Addr: e.Config.Redis,
	})
}

func (e *Env) InitLog() {
	if e.Config.TestFlag {
		e.loglevel = log.DebugLevel
	}

	log.SetLevel(e.loglevel)
	log.SetFormatter(&log.JSONFormatter{})
}

// func (e *Env) InitDb() {
// 	dsn := initDbDsn(e.Config)
// 	log.Debug(dsn)
// 	db, err := sqlx.Connect("clickhouse", dsn)
// 	checkErr(err)
// 	e.Db = db
// }

// func initDbDsn(cfg *Config) string {
// 	return fmt.Sprintf("tcp://%s:%d?username=%s&password=%s&database=%s", cfg.DbHost, cfg.DbPort, cfg.DbUser, cfg.DbPass, cfg.DbName)
// }

func checkErr(err error) {
	if err != nil {
		_, filename, lineno, ok := runtime.Caller(1)
		message := ""
		if ok {
			message = fmt.Sprintf("%v:%v: %v\n", filename, lineno, err)
		}
		log.Panic(message, err)
	}
}
