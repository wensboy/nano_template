package config

import (
	"context"
	"log"
	"os"
	"time"

	"example.com/nano_template/pkg/util"
	"github.com/glebarez/sqlite"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type (
	// DatabaseConfig holds the configuration options for the database.
	DatabaseConfig struct {
		Enable        bool         `yaml:"enable"`
		Type          string       `yaml:"type"`
		EnableLog     bool         `yaml:"enableLog"`
		LogLevel      int          `yaml:"logLevel"`
		SlowThreshold int          `yaml:"slowThreshold"`
		AutoMigrate   bool         `yaml:"autoMigrate"`
		MysqlConfig   MysqlConfig  `yaml:"mysql"`
		SqliteConfig  SqliteConfig `yaml:"sqlite"`
	}
	MysqlConfig struct {
		DSN             string `yaml:"dsn"`
		MaxIdleConns    int    `yaml:"maxIdleConns"`
		MaxOpenConns    int    `yaml:"maxOpenConns"`
		ConnMaxLifetime int    `yaml:"connMaxLifetime"`
	}
	SqliteConfig struct {
		Path string `yaml:"path"`
	}
	ValkeyConfig struct {
		Enable   bool   `yaml:"enable"`
		Address  string `yaml:"address"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
		Database int    `yaml:"database"`
	}
)

// DefaultDatabaseConfig provides a default configuration for the database.
func DefaultDatabaseConfig() DatabaseConfig {
	return DatabaseConfig{
		Type:          "mysql",
		EnableLog:     true,
		LogLevel:      1,
		SlowThreshold: 500,
		AutoMigrate:   false,
	}
}

func DefaultValkeyConfig() ValkeyConfig {
	return ValkeyConfig{
		Enable:   false,
		Address:  "localhost:6379",
		Username: "",
		Password: "",
		Database: 0,
	}
}

var GDB *gorm.DB
var VDB *redis.Client

func CloseDB() {
	util.Info("close all database...")
	db, _ := GDB.DB()
	db.Close()
	VDB.Close()
	util.Info("close all database successfully")
}

func InitDB(cfg *DatabaseConfig) {
	util.Info("load database...")
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Duration(cfg.SlowThreshold) * time.Millisecond,
			LogLevel:      logger.LogLevel(cfg.LogLevel),
			Colorful:      true,
		},
	)
	var db *gorm.DB
	switch cfg.Type {
	case "mysql":
		db = initMysqlDB(&cfg.MysqlConfig, &gorm.Config{
			Logger: newLogger,
		})
	case "sqlite":
		db = initSqliteDB(&cfg.SqliteConfig, &gorm.Config{
			Logger: newLogger,
		})
	default:
		util.Warn("不支持的数据库类型: " + cfg.Type)
	}
	GDB = db
	util.Info("database connected.")
}

func initMysqlDB(cfg *MysqlConfig, opts *gorm.Config) *gorm.DB {
	gdb, err := gorm.Open(mysql.Open(cfg.DSN), opts)
	if err != nil {
		panic("数据库连接失败: " + err.Error())
	}
	db, _ := gdb.DB()
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetConnMaxLifetime(time.Duration(cfg.ConnMaxLifetime) * time.Second)
	return gdb
}

func initSqliteDB(cfg *SqliteConfig, opts *gorm.Config) *gorm.DB {
	path := cfg.Path
	if path == "" {
		path = "dev.db"
	}

	gdb, err := gorm.Open(sqlite.Open(path), opts)
	if err != nil {
		panic("sqlite 数据库连接失败: " + err.Error())
	}

	db, err := gdb.DB()
	if err == nil {
		db.SetMaxIdleConns(1)
		db.SetMaxOpenConns(1)
		db.SetConnMaxLifetime(0)
	}

	return gdb
}

func InitValkey(cfg *ValkeyConfig) {
	util.Info("load valkey...")
	VDB = redis.NewClient(&redis.Options{
		Addr:            cfg.Address,
		Username:        cfg.Username,
		Password:        cfg.Password,
		DB:              cfg.Database,
		PoolSize:        20,
		MinIdleConns:    5,
		MaxActiveConns:  20,
		DialTimeout:     5 * time.Second,
		ReadTimeout:     3 * time.Second,
		WriteTimeout:    3 * time.Second,
		PoolTimeout:     4 * time.Second,
		ConnMaxIdleTime: 30 * time.Minute,
		ConnMaxLifetime: 1 * time.Hour,
	})
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := VDB.Ping(ctx).Err(); err != nil {
		util.Warn("load valkey fail")
		return
	}
	util.Info("load valkey successfully")
}
