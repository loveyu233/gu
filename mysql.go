package gu

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	mysqlClient *gorm.DB
)

type Setting struct {
	Debug bool
}

type MysqlConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Database string
	Debug    bool
}

var DefaultMysqlConfig = MysqlConfig{
	Host:     "127.0.0.1",
	Port:     3306,
	User:     "root",
	Password: "mysql",
	Database: "mysql",
	Debug:    true,
}

func MustInitMysqlClient(config ...MysqlConfig) *gorm.DB {
	var (
		err     error
		dConfig MysqlConfig
	)
	if len(config) > 0 {
		dConfig = config[0]
	} else {
		dConfig = DefaultMysqlConfig
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dConfig.User,
		dConfig.Password,
		dConfig.Host,
		dConfig.Port,
		dConfig.Database,
	)

	if mysqlClient, err = gorm.Open(mysql.Open(dsn), &gorm.Config{}); err != nil {
		logrus.Panicf("mysql connect err: %v\n", err)
	}

	if dConfig.Debug {
		mysqlClient = mysqlClient.Debug()
	}

	logrus.Info("successfully init mysql client")

	MySqlClient = func(cfg ...UseConfig) *gorm.DB {
		var dUseConfig UseConfig
		if len(cfg) == 0 {
			dUseConfig = DefaultUseConfig
		} else {
			dUseConfig = cfg[0]
		}
		tx := mysqlClient.Session(&gorm.Session{}).WithContext(dUseConfig.Ctx)

		if dUseConfig.Debug {
			tx = tx.Debug()
		}

		return tx
	}

	return mysqlClient
}
