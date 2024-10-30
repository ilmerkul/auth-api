package mysql

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type StorageConfig struct {
	ConnMaxLifetime int    `yaml:"conn_max_lifetime" env:"DB_CONN_MAX_LIFETIME" env-default:"0"`
	MaxIdleConns    int    `yaml:"max_idle_conns" env:"DB_MAX_IDLE_CONNS" env-default:"50"`
	MaxOpenConns    int    `yaml:"max_open_conns" env:"DB_MAX_OPEN_CONNS" env-default:"50"`
	StoragePath     string `yaml:"storage_path" env:"STORAGE_PATH"`
	ConfigDSN       `yaml:"dsn"`
}

type ConfigDSN struct {
	User     string `yaml:"user" env:"STORAGE_DSN_USER"`
	Password string `yaml:"password" env:"STORAGE_DSN_PASSWORD"`
	Protocol string `yaml:"protocol" env:"STORAGE_DSN_PROTOCOL"`
	Host     string `yaml:"host" env:"STORAGE_DSN_HOST"`
	Port     string `yaml:"port" env:"STORAGE_DSN_PORT"`
	NameDB   string `yaml:"name_db" env:"STORAGE_DSN_NAME_DB"`
}

func NewClient(storageCfg *StorageConfig) (*sql.DB, error) {
	var (
		DB  *sql.DB
		err error
	)
	if len(storageCfg.StoragePath) == 0 {
		DB, err = sql.Open("mysql", fmt.Sprintf("%s:%s@%s(%s:%s)/%s", storageCfg.ConfigDSN.User, storageCfg.ConfigDSN.Password, storageCfg.ConfigDSN.Protocol, storageCfg.ConfigDSN.Host, storageCfg.ConfigDSN.Port, storageCfg.ConfigDSN.NameDB))
	} else {
		DB, err = sql.Open("mysql", storageCfg.StoragePath)
	}

	if err != nil {
		return nil, err
	}

	DB.SetConnMaxLifetime(time.Duration(storageCfg.ConnMaxLifetime))
	DB.SetMaxIdleConns(storageCfg.MaxIdleConns)
	DB.SetMaxOpenConns(storageCfg.MaxOpenConns)

	if err := DB.Ping(); err != nil {
		return nil, err
	}

	return DB, nil
}
